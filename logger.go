package l4g

import (
	"fmt"
	"io"
	"runtime"
	"time"
)

// Logger defines the logging interface.
type Logger interface {
	Output() io.Writer
	SetOutput(w io.Writer)
	Level() Level
	SetLevel(lvl Level)
	StacktraceLevel() Level
	SetStacktraceLevel(lvl Level)
	Enabled(lvl Level) bool
	StacktraceEnabled(lvl Level) bool
	Trace(i ...any)
	Tracef(format string, args ...any)
	Tracej(j map[string]any)
	Debug(i ...any)
	Debugf(format string, args ...any)
	Debugj(j map[string]any)
	Info(i ...any)
	Infof(format string, args ...any)
	Infoj(j map[string]any)
	Warn(i ...any)
	Warnf(format string, args ...any)
	Warnj(j map[string]any)
	Error(i ...any)
	Errorf(format string, args ...any)
	Errorj(j map[string]any)
	Panic(i ...any)
	Panicj(j map[string]any)
	Panicf(format string, args ...any)
	Fatal(i ...any)
	Fatalj(j map[string]any)
	Fatalf(format string, args ...any)
}

type Options struct {
	Level           Level
	StacktraceLevel Level
	NewHandlerFunc  func(opts HandlerOptions) Handler
	Handler         Handler
	ReplacePart     func(kind PartKind, r *Record, last bool) (string, bool)
}

type Option func(*Options)

func WithLevel(lvl Level) Option {
	return func(opts *Options) {
		opts.Level = lvl
	}
}

func WithStacktraceLevel(lvl Level) Option {
	return func(opts *Options) {
		opts.StacktraceLevel = lvl
	}
}

func WithNewHandlerFunc(f func(opts HandlerOptions) Handler) Option {
	return func(opts *Options) {
		opts.NewHandlerFunc = f
	}
}

func WithHandler(h Handler) Option {
	return func(opts *Options) {
		opts.NewHandlerFunc = func(_ HandlerOptions) Handler {
			return h
		}
	}
}

func WithReplacePart(f func(PartKind, *Record, bool) (string, bool)) Option {
	return func(opts *Options) {
		opts.ReplacePart = f
	}
}

func New(out io.Writer, options ...Option) Logger {
	opts := Options{
		Level:           LevelInfo,
		StacktraceLevel: LevelPanic,
		NewHandlerFunc:  NewSimpleHandler,
	}
	for _, option := range options {
		option(&opts)
	}
	l := &defaultLogger{
		level:           NewLevelVar(opts.Level),
		stacktraceLevel: NewLevelVar(opts.StacktraceLevel),
		output:          NewOutputVar(out),
		handler:         opts.Handler,
	}
	if opts.Handler == nil {
		l.handler = opts.NewHandlerFunc(HandlerOptions{
			Level:           l.level,
			StacktraceLevel: l.stacktraceLevel,
			Output:          l.output,
			ReplacePart:     opts.ReplacePart,
		})
	}
	return l
}

type defaultLogger struct {
	level           *LevelVar
	stacktraceLevel *LevelVar
	output          *OutputVar
	handler         Handler
}

func (l *defaultLogger) Output() io.Writer {
	return l.output.Output()
}

func (l *defaultLogger) SetOutput(w io.Writer) {
	l.output.Set(w)
}

func (l *defaultLogger) Level() Level {
	return l.level.Level()
}

func (l *defaultLogger) SetLevel(lvl Level) {
	l.level.Set(lvl)
}

func (l *defaultLogger) StacktraceLevel() Level {
	return l.stacktraceLevel.Level()
}

func (l *defaultLogger) SetStacktraceLevel(lvl Level) {
	l.stacktraceLevel.Set(lvl)
}

func (l *defaultLogger) Enabled(level Level) bool {
	return l.handler.Enabled(level)
}

func (l *defaultLogger) StacktraceEnabled(level Level) bool {
	return l.handler.StacktraceEnabled(level)
}

// Write writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is used to recover the PC and is
// provided for generality, although at the moment on all pre-defined
// paths it will be 2.
func (l *defaultLogger) Write(calldepth int, level Leveler, s string) error {
	calldepth++ // +1 for this frame.
	return l.write(calldepth, level.Level(), func() string { return s })
}

// write can take either a calldepth or a pc to get source line information.
// It uses the pc if it is non-zero.
func (l *defaultLogger) write(calldepth int, level Level, msg func() string) error {
	if l.output.Discard() || !l.Enabled(level) {
		return nil
	}
	var pcs [32]uintptr
	runtime.Callers(calldepth, pcs[:])
	r := Record{time.Now(), msg(), level, pcs[0], pcs[1:]}
	return l.handler.Handle(r, l.StacktraceEnabled(level))
}

func (l *defaultLogger) log(calldepth int, level Level, appendOutput func() string) {
	err := l.write(calldepth+1, level, appendOutput)
	if err != nil {
		FallbackErrorf("unable to write log message: %v", err)
	}
}

// Trace calls l.Write to print to the logger.
// Arguments are handled in the manner of [fmt.Println].
func (l *defaultLogger) Trace(v ...any) {
	l.log(2, LevelTrace, func() string { return fmt.Sprint(v...) })
}

// Tracef calls l.Write to print to the logger.
// Arguments are handled in the manner of [fmt.Printf].
func (l *defaultLogger) Tracef(format string, args ...any) {
	l.log(2, LevelTrace, func() string { return fmt.Sprintf(format, args...) })
}

// Tracej calls l.Write to print to the logger.
// Arguments are handled in the manner of [fmt.Printf] an [json.Marshal].
func (l *defaultLogger) Tracej(j map[string]any) {
	l.log(2, LevelTrace, func() string { return stringify(j) })
}

// Debug calls l.Write to print to the logger.
// Arguments are handled in the manner of [fmt.Print].
func (l *defaultLogger) Debug(v ...any) {
	l.log(2, LevelDebug, func() string { return fmt.Sprint(v...) })
}

// Debugf calls l.Write to print to the logger.
// Arguments are handled in the manner of [fmt.Printf].
func (l *defaultLogger) Debugf(format string, v ...any) {
	l.log(2, LevelDebug, func() string { return fmt.Sprintf(format, v...) })
}

// Debugj calls l.Write to print to the logger.
// Arguments are handled in the manner of [fmt.Println].
func (l *defaultLogger) Debugj(j map[string]any) {
	l.log(2, LevelDebug, func() string { return stringify(j) })
}

// Info calls l.Write to print to the logger.
// Arguments are handled in the manner of [fmt.Print].
func (l *defaultLogger) Info(v ...any) {
	l.log(2, LevelInfo, func() string { return fmt.Sprint(v...) })
}

// Infof calls l.Write to print to the logger.
// Arguments are handled in the manner of [fmt.Printf].
func (l *defaultLogger) Infof(format string, v ...any) {
	l.log(2, LevelInfo, func() string { return fmt.Sprintf(format, v...) })
}

// Infoj calls l.Write to print to the logger.
// Arguments are handled in the manner of [fmt.Println].
func (l *defaultLogger) Infoj(j map[string]any) {
	l.log(2, LevelInfo, func() string { return stringify(j) })
}

// Warn calls l.Write to print to the logger.
// Arguments are handled in the manner of [fmt.Print].
func (l *defaultLogger) Warn(v ...any) {
	l.log(2, LevelWarn, func() string { return fmt.Sprint(v...) })
}

// Warnf calls l.Write to print to the logger.
// Arguments are handled in the manner of [fmt.Printf].
func (l *defaultLogger) Warnf(format string, v ...any) {
	l.log(2, LevelWarn, func() string { return fmt.Sprintf(format, v...) })
}

// Warnj calls l.Write to print to the logger.
// Arguments are handled in the manner of [fmt.Println].
func (l *defaultLogger) Warnj(j map[string]any) {
	l.log(2, LevelWarn, func() string { return stringify(j) })
}

// Error calls l.Write to print to the logger.
// Arguments are handled in the manner of [fmt.Print].
func (l *defaultLogger) Error(v ...any) {
	l.log(2, LevelError, func() string { return fmt.Sprint(v...) })
}

// Errorf calls l.Write to print to the logger.
// Arguments are handled in the manner of [fmt.Printf].
func (l *defaultLogger) Errorf(format string, v ...any) {
	l.log(2, LevelError, func() string { return fmt.Sprintf(format, v...) })
}

// Errorj calls l.Write to print to the logger.
// Arguments are handled in the manner of [fmt.Println].
func (l *defaultLogger) Errorj(j map[string]any) {
	l.log(2, LevelError, func() string { return stringify(j) })
}

// Panic is equivalent to l.Print() followed by a call to panic().
func (l *defaultLogger) Panic(v ...any) {
	s := fmt.Sprint(v...)
	l.log(2, LevelPanic, func() string { return s })
	panic(s)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *defaultLogger) Panicf(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	l.log(2, LevelPanic, func() string { return s })
	panic(s)
}

// Panicj is equivalent to l.Println() followed by a call to panic().
func (l *defaultLogger) Panicj(j map[string]any) {
	s := stringify(j)
	l.log(2, LevelPanic, func() string { return s })
	panic(s)
}

// Fatal is equivalent to l.Print() followed by a call to [os.Exit](1).
func (l *defaultLogger) Fatal(v ...any) {
	l.log(2, LevelFatal, func() string { return fmt.Sprint(v...) })
	OsExiter(1)
}

// Fatalf is equivalent to l.Printf() followed by a call to [os.Exit](1).
func (l *defaultLogger) Fatalf(format string, v ...any) {
	l.log(2, LevelFatal, func() string { return fmt.Sprintf(format, v...) })
	OsExiter(1)
}

// Fatalj is equivalent to l.Println() followed by a call to [os.Exit](1).
func (l *defaultLogger) Fatalj(j map[string]any) {
	l.log(2, LevelFatal, func() string { return stringify(j) })
	OsExiter(1)
}
