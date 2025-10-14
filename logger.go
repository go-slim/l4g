// Package l4g provides a high-performance, structured logging library for Go
// that is compatible with the standard log/slog package. It offers fast logging
// with zero allocations for disabled log levels, buffer pooling, and multiple
// output formats including text, formatted strings, and structured JSON-like logging.
package l4g

import (
	"fmt"
	"io"
	"time"
)

// Options holds configuration options for creating a new Logger.
type Options struct {
	// Prefix is the prefix to use for all log messages.
	Prefix string
	// Level minimum log level to output
	Level Level
	// NewHandlerFunc factory function to create a handler
	NewHandlerFunc func(opts HandlerOptions) Handler
	// Handler custom handler to use (overrides NewHandlerFunc)
	Handler Handler
	// ReplaceAttr function to rewrite attributes before logging
	ReplaceAttr func(groups []string, attr Attr) Attr
	// TimeFormat time format string (default: time.StampMilli)
	TimeFormat string
	// Output destination (default: os.Stderr)
	Output io.Writer
	// NoColor disable color output (default: false)
	NoColor bool
}

// New creates a new Logger that writes to the given io.Writer.
// By default, it uses LevelInfo as the minimum log level and SimpleHandler for output formatting.
// The behavior can be customized using Option functions.
func New(opts Options) *Logger {
	if opts.Level == 0 {
		opts.Level = LevelInfo
	}
	if opts.NewHandlerFunc == nil {
		opts.NewHandlerFunc = NewSimpleHandler
	}
	l := &Logger{
		level:   NewLevelVar(opts.Level.Real()),
		output:  NewOutputVar(opts.Output),
		handler: opts.Handler,
	}
	if opts.Handler == nil {
		l.handler = opts.NewHandlerFunc(HandlerOptions{
			Prefix:      opts.Prefix,
			Level:       l.level,
			Output:      l.output,
			ReplaceAttr: opts.ReplaceAttr,
			TimeFormat:  opts.TimeFormat,
			NoColor:     opts.NoColor,
		})
	}
	return l
}

// Logger represents a logger instance that outputs log messages through a handler.
// It is safe for concurrent use by multiple goroutines.
type Logger struct {
	level   *LevelVar  // Minimum log level, can be changed dynamically
	output  *OutputVar // Output destination, can be changed dynamically
	handler Handler    // Handler for processing and formatting log records
}

// Output returns the current output destination for the logger.
func (l *Logger) Output() io.Writer {
	return l.output.Output()
}

// SetOutput sets the output destination for the logger.
// This can be called at runtime to redirect log output.
func (l *Logger) SetOutput(w io.Writer) {
	l.output.Set(w)
}

// Level returns the current minimum log level of the logger.
func (l *Logger) Level() Level {
	return l.level.Level()
}

// SetLevel sets the minimum log level of the logger.
// Log messages below this level will not be output.
func (l *Logger) SetLevel(lvl Level) {
	l.level.Set(lvl)
}

// Enabled reports whether the logger is enabled for the given log level.
// It returns true if a log message at the given level would be output.
func (l *Logger) Enabled(level Level) bool {
	return l.handler.Enabled(level)
}

// WithAttrs returns a new Logger that includes the given attributes in all subsequent log output.
// The attributes are added to every log record produced by the returned logger.
// args can be key-value pairs (string, any, string, any, ...) or Attr values.
func (l *Logger) WithAttrs(args ...any) *Logger {
	if len(args) == 0 {
		return l
	}
	return &Logger{
		level:   l.level,
		output:  l.output,
		handler: l.handler.WithAttrs(argsToAttrSlice(args)),
	}
}

// WithPrefix returns a new Logger that includes the given prefix in all subsequent log output.
// The prefix is prepended to the logger's existing prefix (if any).
func (l *Logger) WithPrefix(prefix string) *Logger {
	if prefix == "" {
		return l
	}
	return &Logger{
		level:   l.level,
		output:  l.output,
		handler: l.handler.WithPrefix(prefix),
	}
}

// WithGroup returns a new Logger that starts a group for all subsequent log output.
// All attributes added by the returned logger will be nested under the given group name.
// If the name is empty, WithGroup returns the receiver unchanged.
func (l *Logger) WithGroup(name string) *Logger {
	if name == "" {
		return l
	}
	return &Logger{
		level:   l.level,
		output:  l.output,
		handler: l.handler.WithGroup(name),
	}
}

// Log outputs a log record at the specified level with the given message and optional attributes.
// args can be key-value pairs (string, any, string, any, ...) or Attr values.
// If the log level is disabled, this function returns immediately without allocating.
func (l *Logger) Log(level Leveler, msg string, args ...any) {
	l.log(level.Level(), msg, args)
}

// Logf outputs a formatted log record at the specified level.
// It supports both [fmt.Printf]-style formatting and optional structured attributes.
// args can mix format arguments with Attr values for structured logging.
func (l *Logger) Logf(level Level, format string, args ...any) {
	l.logf(level, format, args)
}

// Logj outputs a log record at the specified level with structured key-value pairs from a map.
// The map is converted to structured attributes in the log output.
func (l *Logger) Logj(level Level, j map[string]any) {
	l.logj(level, j)
}

// Trace logs a message at trace level with optional structured attributes.
// args can be key-value pairs (string, any, string, any, ...) or Attr values.
func (l *Logger) Trace(msg string, args ...any) {
	l.Log(LevelTrace, msg, args...)
}

// Tracef logs a formatted message at trace level.
// It supports [fmt.Printf]-style formatting and optional structured attributes.
func (l *Logger) Tracef(format string, args ...any) {
	l.logf(LevelTrace, format, args)
}

// Tracej logs a message at trace level with structured key-value pairs from a map.
func (l *Logger) Tracej(j map[string]any) {
	l.logj(LevelTrace, j)
}

// Debug logs a message at debug level with optional structured attributes.
// args can be key-value pairs (string, any, string, any, ...) or Attr values.
func (l *Logger) Debug(msg string, args ...any) {
	l.log(LevelDebug, msg, args)
}

// Debugf logs a formatted message at debug level.
// It supports [fmt.Printf]-style formatting and optional structured attributes.
func (l *Logger) Debugf(format string, args ...any) {
	l.logf(LevelDebug, format, args)
}

// Debugj logs a message at debug level with structured key-value pairs from a map.
func (l *Logger) Debugj(j map[string]any) {
	l.logj(LevelDebug, j)
}

// Info logs a message at info level with optional structured attributes.
// args can be key-value pairs (string, any, string, any, ...) or Attr values.
func (l *Logger) Info(msg string, args ...any) {
	l.log(LevelInfo, msg, args)
}

// Infof logs a formatted message at info level.
// It supports [fmt.Printf]-style formatting and optional structured attributes.
func (l *Logger) Infof(format string, args ...any) {
	l.logf(LevelInfo, format, args)
}

// Infoj logs a message at info level with structured key-value pairs from a map.
func (l *Logger) Infoj(j map[string]any) {
	l.logj(LevelInfo, j)
}

// Warn logs a message at warn level with optional structured attributes.
// args can be key-value pairs (string, any, string, any, ...) or Attr values.
func (l *Logger) Warn(msg string, args ...any) {
	l.log(LevelWarn, msg, args)
}

// Warnf logs a formatted message at warn level.
// It supports [fmt.Printf]-style formatting and optional structured attributes.
func (l *Logger) Warnf(format string, args ...any) {
	l.logf(LevelWarn, format, args)
}

// Warnj logs a message at warn level with structured key-value pairs from a map.
func (l *Logger) Warnj(j map[string]any) {
	l.logj(LevelWarn, j)
}

// Error logs a message at error level with optional structured attributes.
// args can be key-value pairs (string, any, string, any, ...) or Attr values.
func (l *Logger) Error(msg string, args ...any) {
	l.log(LevelError, msg, args)
}

// Errorf logs a formatted message at error level.
// It supports [fmt.Printf]-style formatting and optional structured attributes.
func (l *Logger) Errorf(format string, args ...any) {
	l.logf(LevelError, format, args)
}

// Errorj logs a message at error level with structured key-value pairs from a map.
func (l *Logger) Errorj(j map[string]any) {
	l.logj(LevelError, j)
}

// Panic logs a message at panic level with optional structured attributes, then panics.
// args can be key-value pairs (string, any, string, any, ...) or Attr values.
func (l *Logger) Panic(msg string, args ...any) {
	l.log(LevelPanic, msg, args)
	panic(msg)
}

// Panicf logs a formatted message at panic level, then panics.
// It supports [fmt.Printf]-style formatting and optional structured attributes.
func (l *Logger) Panicf(format string, args ...any) {
	l.logf(LevelPanic, format, args)

	_, anies := splitAttrs(args)
	msg := format
	if len(anies) > 0 {
		msg = fmt.Sprintf(format, anies...)
	}
	panic(msg)
}

// Panicj logs a message at panic level with structured key-value pairs from a map, then panics.
func (l *Logger) Panicj(j map[string]any) {
	l.logj(LevelPanic, j)
	panic(j)
}

// Fatal logs a message at fatal level with optional structured attributes, then calls os.Exit(1).
// args can be key-value pairs (string, any, string, any, ...) or Attr values.
func (l *Logger) Fatal(msg string, args ...any) {
	l.log(LevelFatal, msg, args)
	OsExiter(1)
}

// Fatalf logs a formatted message at fatal level, then calls os.Exit(1).
// It supports [fmt.Printf]-style formatting and optional structured attributes.
func (l *Logger) Fatalf(format string, args ...any) {
	l.logf(LevelFatal, format, args)
	OsExiter(1)
}

// Fatalj logs a message at fatal level with structured key-value pairs from a map, then calls os.Exit(1).
func (l *Logger) Fatalj(j map[string]any) {
	l.logj(LevelFatal, j)
	OsExiter(1)
}

// log is the internal implementation for logging with optional structured attributes.
// It returns early without allocating if the output is disabled or the level is not enabled.
func (l *Logger) log(level Level, msg string, args []any) {
	if l.output.Discard() || !l.Enabled(level) {
		return
	}
	r := NewRecord(time.Now(), level, msg)
	if len(args) > 0 {
		r.AddAttrs(argsToAttrSlice(args)...)
	}
	if err := l.handler.Handle(r); err != nil {
		FallbackErrorf("unable to write log message: %v", err)
	}
}

// logf is the internal implementation for formatted logging with optional structured attributes.
// It returns early without allocating if the output is disabled or the level is not enabled.
// args are split into Attr values for structured logging and regular values for fmt.Sprintf formatting.
func (l *Logger) logf(level Level, format string, args []any) {
	if l.output.Discard() || !l.Enabled(level) {
		return
	}
	attrs, anies := splitAttrs(args)
	msg := format
	if len(anies) > 0 {
		msg = fmt.Sprintf(format, anies...)
	}
	r := NewRecord(time.Now(), level, msg)
	if len(attrs) > 0 {
		r.AddAttrs(attrs...)
	}
	if err := l.handler.Handle(r); err != nil {
		FallbackErrorf("unable to write log message: %v", err)
	}
}

// logj is the internal implementation for logging with structured key-value pairs from a map.
// It returns early without allocating if the output is disabled or the level is not enabled.
func (l *Logger) logj(level Level, j map[string]any) {
	if l.output.Discard() || !l.Enabled(level) {
		return
	}
	r := NewRecord(time.Now(), level, "")
	for key, value := range j {
		r.Add(key, value)
	}
	if err := l.handler.Handle(r); err != nil {
		FallbackErrorf("unable to write log message: %v", err)
	}
}
