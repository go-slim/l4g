package l4g

import (
	"fmt"
	"io"
	"os"
	"sync"
)

var (
	std *defaultLogger

	mu sync.Mutex
	ls map[string]Logger

	OsExiter func(code int)
	NewFunc  func(name string) Logger
)

func init() {
	std = New(os.Stderr).(*defaultLogger)
	OsExiter = os.Exit
	NewFunc = func(_ string) Logger { return New(os.Stderr) }
}

// FallbackErrorf is the last chance to show an error if the logger has internal errors
func FallbackErrorf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func Channel(name string) Logger {
	mu.Lock()
	defer mu.Unlock()

	if ls == nil {
		ls = make(map[string]Logger)
	}

	if l, ok := ls[name]; ok {
		return l
	}

	l := NewFunc(name)
	ls[name] = l
	return l
}

// Default returns the default logger used by the package-level output functions.
func Default() Logger {
	return std
}

// SetDefault sets the default logger used by the package-level output functions.
func SetDefault(l Logger) bool {
	if dl, ok := l.(*defaultLogger); ok {
		mu.Lock()
		defer mu.Unlock()
		std = dl
		return true
	}
	return false
}

// Output returns the output destination for the standard logger.
func Output() io.Writer {
	return std.Output()
}

// SetOutput sets the output destination for the standard logger.
func SetOutput(w io.Writer) {
	std.SetOutput(w)
}

func GetLevel() Level {
	return Level(std.Level())
}

func SetLevel(level Level) {
	std.SetLevel(level)
}

func StacktraceLevel() Level {
	return std.StacktraceLevel()
}

func SetStacktraceLevel(level Level) {
	std.SetStacktraceLevel(level)
}

func Trace(i ...any) {
	std.log(2, LevelTrace, func() string { return fmt.Sprint(i...) })
}

func Tracef(format string, args ...any) {
	std.log(2, LevelTrace, func() string { return fmt.Sprintf(format, args...) })
}

func Tracej(j map[string]any) {
	std.log(2, LevelTrace, func() string { return stringify(j) })
}

// Debug calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Print].
func Debug(v ...any) {
	std.log(2, LevelDebug, func() string { return fmt.Sprint(v...) })
}

// Debugf calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Printf].
func Debugf(format string, v ...any) {
	std.log(2, LevelDebug, func() string { return fmt.Sprintf(format, v...) })
}

// Debugj calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Println].
func Debugj(j map[string]any) {
	std.log(2, LevelDebug, func() string { return stringify(j) })
}

// Info calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Print].
func Info(v ...any) {
	std.log(2, LevelInfo, func() string { return fmt.Sprint(v...) })
}

// Infof calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Printf].
func Infof(format string, v ...any) {
	std.log(2, LevelInfo, func() string { return fmt.Sprintf(format, v...) })
}

// Infoj calls Write to print to the standard logger.
// Arguments are handled in the manner of [fmt.Println].
func Infoj(j map[string]any) {
	std.log(2, LevelInfo, func() string { return stringify(j) })
}

// Warn calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Print].
func Warn(v ...any) {
	std.log(2, LevelWarn, func() string { return fmt.Sprint(v...) })
}

// Warnf calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Printf].
func Warnf(format string, v ...any) {
	std.log(2, LevelWarn, func() string { return fmt.Sprintf(format, v...) })
}

// Warnj calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Println].
func Warnj(j map[string]any) {
	std.log(2, LevelWarn, func() string { return stringify(j) })
}

// Error calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Print].
func Error(v ...any) {
	std.log(2, LevelError, func() string { return fmt.Sprint(v...) })
}

// Errorf calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Printf].
func Errorf(format string, v ...any) {
	std.log(2, LevelError, func() string { return fmt.Sprintf(format, v...) })
}

// Errorj calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Println].
func Errorj(j map[string]any) {
	std.log(2, LevelError, func() string { return stringify(j) })
}

// Panic is equivalent to [Print] followed by a call to panic().
func Panic(v ...any) {
	s := fmt.Sprint(v...)
	std.log(2, LevelPanic, func() string { return s })
	//panic(s)
}

// Panicf is equivalent to [Printf] followed by a call to panic().
func Panicf(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	std.log(2, LevelPanic, func() string { return s })
	panic(s)
}

// Panicj is equivalent to [Println] followed by a call to panic().
func Panicj(j map[string]any) {
	s := stringify(j)
	std.log(2, LevelPanic, func() string { return s })
	panic(s)
}

// Fatal is equivalent to [Print] followed by a call to [os.Exit](1).
func Fatal(v ...any) {
	std.log(2, LevelFatal, func() string { return fmt.Sprint(v...) })
	//OsExiter(1)
}

// Fatalf is equivalent to [Printf] followed by a call to [os.Exit](1).
func Fatalf(format string, v ...any) {
	std.log(2, LevelFatal, func() string { return fmt.Sprintf(format, v...) })
	OsExiter(1)
}

// Fatalj is equivalent to [Println] followed by a call to [os.Exit](1).
func Fatalj(j map[string]any) {
	std.log(2, LevelFatal, func() string { return stringify(j) })
	OsExiter(1)
}

// Write writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is the count of the number of
// frames to skip when computing the file name and line number
func Write(calldepth int, level Leveler, s string) error {
	return std.Write(calldepth+1, level, s) // +1 for this frame.
}
