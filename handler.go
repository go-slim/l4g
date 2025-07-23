package l4g

import (
	"io"
	"time"
)

// A Record holds information about a log event.
// Do not modify a Record after handing out a copy to it.
type Record struct {
	// The time at which the output method (Log, Info, etc.) was called.
	Time time.Time

	// The log message.
	Message string

	// The level of the event.
	Level Level

	// The program counter at the time the record was constructed, as determined
	// by runtime.Callers. If zero, no program counter is available.
	//
	// The only valid use for this value is as an argument to
	// [runtime.CallersFrames]. In particular, it must not be passed to
	// [runtime.FuncForPC].
	PC uintptr

	// Frames is a stack trace of the program counter at the time the record was
	// constructed, as determined by runtime.Callers. If empty, no stack trace is
	// available.
	Frames []uintptr
}

type PartKind int8

const (
	PartTime PartKind = iota
	PartLevel
	PartMessage
	PartLocation
	PartStacktrace
)

// A Handler handles log records produced by a Logger.
//
// Any of the Handler's methods may be called concurrently with itself
// or with other methods. It is the responsibility of the Handler to
// manage this concurrency.
//
// Users of the l4g package should not invoke Handler methods directly.
// They should use the methods of [Logger] instead.
type Handler interface {
	// Enabled reports whether the handler handles records at the given level.
	// The handler ignores records whose level is lower.
	// It is called early, before any arguments are processed,
	// to save effort if the log event should be discarded.
	Enabled(Level) bool

	// StacktraceEnabled reports whether the handler handles records at the given level.
	// The handler ignores stacktrace whose level is lower.
	// It is called early, before any arguments are processed,
	// to save effort if the log event should be discarded.
	StacktraceEnabled(Level) bool

	// Handle handles the Record.
	// It will only be called when Enabled returns true.
	Handle(r Record, stacktrace bool) error
}

// HandlerOptions are options for a [SimpleHandler].
// A zero HandlerOptions consists entirely of default values.
type HandlerOptions struct {
	// Level reports the minimum record level that will be logged.
	// The handler discards records with lower levels.
	// If Level is nil, the handler assumes LevelInfo.
	// The handler calls Level.Level for each record processed;
	// to adjust the minimum level dynamically, use a LevelVar.
	Level Leveler

	// StacktraceLevel reports the minimum record level that will be logged.
	// The handler discards records with lower levels.
	// If Level is nil, the handler assumes LevelInfo.
	// The handler calls Level.Level for each record processed;
	// to adjust the minimum level dynamically, use a LevelVar.
	StacktraceLevel Leveler

	// ReplacePart is called to write each record part on it is logged.
	//
	// The record parts with kinds "time", "level", "message", "location" and
	// "stacktrace" are passed to this function.
	ReplacePart func(kind PartKind, r *Record, last bool) (string, bool)

	// Output is a destination to which log data will be written.
	Output io.Writer
}

// NewSimpleHandler creates a [SimpleHandler] that writes to w,
// using the given options.
// If opts is nil, the default options are used.
func NewSimpleHandler(opts HandlerOptions) Handler {
	return &SimpleHandler{
		options: &opts,
	}
}

type SimpleHandler struct {
	options *HandlerOptions
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
func (l *SimpleHandler) Enabled(level Level) bool {
	minLevel := LevelInfo
	if l.options.Level != nil {
		minLevel = l.options.Level.Level()
	}
	return level >= minLevel
}

// StacktraceEnabled reports whether the handler handles stacktrace at the given level.
// The handler ignores records whose level is lower.
func (l *SimpleHandler) StacktraceEnabled(level Level) bool {
	minLevel := LevelPanic
	if l.options.StacktraceLevel != nil {
		minLevel = l.options.StacktraceLevel.Level()
	}
	return level >= minLevel
}

// Handle formats its argument [Record] as a single line of space-separated
// fields.
func (l *SimpleHandler) Handle(r Record, stacktrace bool) error {
	rw := newRecordWriter(l)

	rw.Write(PartLevel, &r, false)
	rw.Write(PartTime, &r, false)
	rw.Write(PartMessage, &r, !stacktrace)

	if stacktrace {
		rw.Write(PartLocation, &r, false)
		rw.Write(PartStacktrace, &r, true)
	}

	return rw.FlushTo(l.options.Output)
}
