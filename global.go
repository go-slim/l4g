package l4g

import (
	"fmt"
	"io"
	"os"
	"sync"
)

var (
	std      *Logger
	OsExiter func(code int)

	pollMu sync.Mutex
	poll   map[string]*Logger

	NewLoggerFunc func(name string) *Logger
)

func init() {
	std = New(os.Stderr)
	poll = make(map[string]*Logger)
	NewLoggerFunc = func(name string) *Logger {
		return New(os.Stderr, WithPrefix(name))
	}
}

// FallbackErrorf is the last chance to show an error if the logger has internal errors
func FallbackErrorf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func Channel(name string) *Logger {
	pollMu.Lock()
	defer pollMu.Unlock()
	logger := poll[name]
	if logger == nil {
		logger = NewLoggerFunc(name)
		poll[name] = logger
	}
	return logger
}

// Default returns the default logger used by the package-level output functions.
func Default() *Logger { return std }

// SetOutput sets the output destination for the standard logger.
func SetOutput(w io.Writer) {
	std.SetOutput(w)
}

// Flags returns the output flags for the standard logger.
// The flag bits are [Ldate], [Ltime], and so on.
func Flags() int {
	return std.Flags()
}

// SetFlags sets the output flags for the standard logger.
// The flag bits are [Ldate], [Ltime], and so on.
func SetFlags(flag int) {
	std.SetFlags(flag)
}

// Prefix returns the output prefix for the standard logger.
func Prefix() string {
	return std.Prefix()
}

// SetPrefix sets the output prefix for the standard logger.
func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}

func GetLevel() Level {
	return std.Level()
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

// Writer returns the output destination for the standard logger.
func Writer() io.Writer {
	return std.Writer()
}

// Debug calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Print].
func Debug(v ...any) {
	std.log(2, DEBUG, func(b []byte) []byte {
		return fmt.Append(b, v...)
	})
}

// Debugf calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Printf].
func Debugf(format string, v ...any) {
	std.log(2, DEBUG, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	})
}

// Debugln calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Println].
func Debugln(v ...any) {
	std.log(2, DEBUG, func(b []byte) []byte {
		return fmt.Appendln(b, v...)
	})
}

// Info calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Print].
func Info(v ...any) {
	std.log(2, INFO, func(b []byte) []byte {
		return fmt.Append(b, v...)
	})
}

// Infof calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Printf].
func Infof(format string, v ...any) {
	std.log(2, INFO, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	})
}

// Infoln calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Println].
func Infoln(v ...any) {
	std.log(2, INFO, func(b []byte) []byte {
		return fmt.Appendln(b, v...)
	})
}

// Warn calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Print].
func Warn(v ...any) {
	std.log(2, WARN, func(b []byte) []byte {
		return fmt.Append(b, v...)
	})
}

// Warnf calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Printf].
func Warnf(format string, v ...any) {
	std.log(2, WARN, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	})
}

// Warnln calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Println].
func Warnln(v ...any) {
	std.log(2, WARN, func(b []byte) []byte {
		return fmt.Appendln(b, v...)
	})
}

// Error calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Print].
func Error(v ...any) {
	std.log(2, ERROR, func(b []byte) []byte {
		return fmt.Append(b, v...)
	})
}

// Errorf calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Printf].
func Errorf(format string, v ...any) {
	std.log(2, ERROR, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	})
}

// Errorln calls Output to print to the standard logger.
// Arguments are handled in the manner of [fmt.Println].
func Errorln(v ...any) {
	std.log(2, ERROR, func(b []byte) []byte {
		return fmt.Appendln(b, v...)
	})
}

// Panic is equivalent to [Print] followed by a call to panic().
func Panic(v ...any) {
	s := fmt.Sprint(v...)
	std.log(2, PANIC, func(b []byte) []byte {
		return append(b, s...)
	})
	panic(s)
}

// Panicf is equivalent to [Printf] followed by a call to panic().
func Panicf(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	std.log(2, PANIC, func(b []byte) []byte {
		return append(b, s...)
	})
	panic(s)
}

// Panicln is equivalent to [Println] followed by a call to panic().
func Panicln(v ...any) {
	s := fmt.Sprintln(v...)
	std.log(2, PANIC, func(b []byte) []byte {
		return append(b, s...)
	})
	panic(s)
}

// Fatal is equivalent to [Print] followed by a call to [os.Exit](1).
func Fatal(v ...any) {
	std.log(2, FATAL, func(b []byte) []byte {
		return fmt.Append(b, v...)
	})
	OsExiter(1)
}

// Fatalf is equivalent to [Printf] followed by a call to [os.Exit](1).
func Fatalf(format string, v ...any) {
	std.log(2, FATAL, func(b []byte) []byte {
		return fmt.Appendf(b, format, v...)
	})
	OsExiter(1)
}

// Fatalln is equivalent to [Println] followed by a call to [os.Exit](1).
func Fatalln(v ...any) {
	std.log(2, FATAL, func(b []byte) []byte {
		return fmt.Appendln(b, v...)
	})
	OsExiter(1)
}

// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is the count of the number of
// frames to skip when computing the file name and line number
// if [Llongfile] or [Lshortfile] is set; a value of 1 will print the details
// for the caller of Output.
func Output(calldepth int, level Level, s string) error {
	return std.Output(calldepth+1, level, s) // +1 for this frame.
}
