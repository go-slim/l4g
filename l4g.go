package l4g

import (
	"fmt"
	"io"
	"os"
	"sync"
)

var (
	// std is the default logger instance used by package-level functions.
	std *Logger

	// mu protects the std logger during SetDefault operations.
	mu sync.Mutex

	// ls stores named channel loggers, keyed by channel name.
	ls *sync.Map // map[string]*Logger

	// OsExiter is the function called by Fatal and Fatalf to exit the program.
	// It is set to os.Exit by default but can be overridden for testing.
	OsExiter func(code int)

	// NewFunc is the factory function used by Channel to create new loggers.
	// It can be overridden to customize logger creation for channels.
	NewFunc func(name string) *Logger
)

func init() {
	std = New(Options{Output: os.Stderr})
	ls = new(sync.Map)
	OsExiter = os.Exit
	NewFunc = func(_ string) *Logger { return New(Options{Output: os.Stderr}) }
}

// FallbackErrorf is the last-resort error reporting function used when the logger
// itself encounters an error. It writes directly to stderr, bypassing all logging handlers.
func FallbackErrorf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
}

// Channel returns a named logger instance. Multiple calls with the same name return
// the same logger instance. This allows different parts of an application to share
// a logger without explicitly passing it around.
// The returned logger is created using NewFunc, which can be customized.
func Channel(name string) *Logger {
	// Fast path: check if logger already exists
	if l, ok := ls.Load(name); ok {
		return l.(*Logger)
	}

	// Slow path: create new logger
	// Note: NewFunc is called without holding any locks
	newLogger := NewFunc(name)

	// Store the logger, or return existing one if another goroutine created it first
	actual, _ := ls.LoadOrStore(name, newLogger)
	return actual.(*Logger)
}

// Default returns the default logger used by the package-level output functions.
func Default() *Logger {
	return std
}

// SetDefault sets the default logger used by the package-level output functions.
func SetDefault(l *Logger) {
	mu.Lock()
	defer mu.Unlock()
	std = l
}

// Output returns the output destination for the standard logger.
func Output() io.Writer {
	return std.Output()
}

// SetOutput sets the output destination for the standard logger.
func SetOutput(w io.Writer) {
	std.SetOutput(w)
}

// GetLevel returns the current minimum log level of the standard logger.
func GetLevel() Level {
	return std.Level()
}

// SetLevel sets the minimum log level for the standard logger.
// Log messages below this level will not be output.
func SetLevel(level Level) {
	std.SetLevel(level)
}

// WithAttrs returns a new Logger based on the standard logger that includes the given attributes
// in all subsequent log output. The attributes are added to every log record.
// args can be key-value pairs (string, any, string, any, ...) or Attr values.
func WithAttrs(args ...any) *Logger {
	return std.WithAttrs(args...)
}

// WithPrefix returns a new Logger based on the standard logger that includes the given prefix
// in all subsequent log output. The prefix is prepended to the logger's existing prefix (if any).
func WithPrefix(prefix string) *Logger {
	return std.WithPrefix(prefix)
}

// WithGroup returns a new Logger based on the standard logger that starts a group for all
// subsequent log output. All attributes added by the returned logger will be nested under
// the given group name. If the name is empty, WithGroup returns the standard logger.
func WithGroup(name string) *Logger {
	return std.WithGroup(name)
}

// Trace logs a message at trace level using the standard logger.
// args can be key-value pairs (string, any, string, any, ...) or Attr values.
func Trace(msg string, args ...any) {
	std.Trace(msg, args...)
}

// Tracef logs a formatted message at trace level using the standard logger.
// It supports [fmt.Printf]-style formatting and optional structured attributes.
func Tracef(format string, args ...any) {
	std.Tracef(format, args...)
}

// Tracej logs a message at trace level with structured key-value pairs from a map using the standard logger.
func Tracej(j map[string]any) {
	std.Tracej(j)
}

// Debug logs a message at debug level using the standard logger.
// args can be key-value pairs (string, any, string, any, ...) or Attr values.
func Debug(msg string, args ...any) {
	std.Debug(msg, args...)
}

// Debugf logs a formatted message at debug level using the standard logger.
// It supports [fmt.Printf]-style formatting and optional structured attributes.
func Debugf(format string, args ...any) {
	std.Debugf(format, args...)
}

// Debugj logs a message at debug level with structured key-value pairs from a map using the standard logger.
func Debugj(j map[string]any) {
	std.Debugj(j)
}

// Info logs a message at info level using the standard logger.
// args can be key-value pairs (string, any, string, any, ...) or Attr values.
func Info(msg string, args ...any) {
	std.Info(msg, args...)
}

// Infof logs a formatted message at info level using the standard logger.
// It supports [fmt.Printf]-style formatting and optional structured attributes.
func Infof(format string, args ...any) {
	std.Infof(format, args...)
}

// Infoj logs a message at info level with structured key-value pairs from a map using the standard logger.
func Infoj(j map[string]any) {
	std.Infoj(j)
}

// Warn logs a message at warn level using the standard logger.
// args can be key-value pairs (string, any, string, any, ...) or Attr values.
func Warn(msg string, args ...any) {
	std.Warn(msg, args...)
}

// Warnf logs a formatted message at warn level using the standard logger.
// It supports [fmt.Printf]-style formatting and optional structured attributes.
func Warnf(format string, args ...any) {
	std.Warnf(format, args...)
}

// Warnj logs a message at warn level with structured key-value pairs from a map using the standard logger.
func Warnj(j map[string]any) {
	std.Warnj(j)
}

// Error logs a message at error level using the standard logger.
// args can be key-value pairs (string, any, string, any, ...) or Attr values.
func Error(msg string, args ...any) {
	std.Error(msg, args...)
}

// Errorf logs a formatted message at error level using the standard logger.
// It supports [fmt.Printf]-style formatting and optional structured attributes.
func Errorf(format string, args ...any) {
	std.Errorf(format, args...)
}

// Errorj logs a message at error level with structured key-value pairs from a map using the standard logger.
func Errorj(j map[string]any) {
	std.Errorj(j)
}

// Panic logs a message at panic level using the standard logger, then panics.
// args can be key-value pairs (string, any, string, any, ...) or Attr values.
func Panic(msg string, args ...any) {
	std.Panic(msg, args...)
}

// Panicf logs a formatted message at panic level using the standard logger, then panics.
// It supports [fmt.Printf]-style formatting and optional structured attributes.
func Panicf(format string, args ...any) {
	std.Panicf(format, args...)
}

// Panicj logs a message at panic level with structured key-value pairs from a map using the standard logger, then panics.
func Panicj(j map[string]any) {
	std.Panicj(j)
}

// Fatal logs a message at fatal level using the standard logger, then calls os.Exit(1).
// args can be key-value pairs (string, any, string, any, ...) or Attr values.
func Fatal(msg string, args ...any) {
	std.Fatal(msg, args...)
}

// Fatalf logs a formatted message at fatal level using the standard logger, then calls os.Exit(1).
// It supports [fmt.Printf]-style formatting and optional structured attributes.
func Fatalf(format string, v ...any) {
	std.Fatalf(format, v...)
}

// Fatalj logs a message at fatal level with structured key-value pairs from a map using the standard logger, then calls os.Exit(1).
func Fatalj(j map[string]any) {
	std.Fatalj(j)
}
