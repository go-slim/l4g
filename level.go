package l4g

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
)

// A Level is the importance or severity of a log event.
type Level int

// Names for common levels.
//
// Level numbers are inherently arbitrary,
// but we picked them to satisfy three constraints.
// Any system can map them to another numbering scheme if it wishes.
const (
	LevelTrace Level = iota + 1
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal
)

// Int returns the integer value of the level.
func (l Level) Int() int {
	return int(l)
}

// Real returns a real level of the level.
func (l Level) Real() Level {
	return max(min(l, LevelFatal), LevelTrace)
}

// String returns a name in lowercase for the level.
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelPanic:
		return "panic"
	default:
		if l <= LevelTrace {
			return "trace"
		}
		return "fatal"
	}
}

// MarshalJSON implements [encoding/json.Marshaler]
// by quoting the output of [Level.String].
func (l Level) MarshalJSON() ([]byte, error) {
	// AppendQuote is sufficient for JSON-encoding all Level strings.
	// They don't contain any runes that would produce invalid JSON
	// when escaped.
	return strconv.AppendQuote(nil, l.String()), nil
}

// UnmarshalJSON implements [encoding/json.Unmarshaler]
// It accepts any string produced by [Level.MarshalJSON],
// ignoring case.
// It also accepts numeric offsets that would result in a different string on
// output. For example, "Error-8" would marshal as "LevelInfo".
func (l *Level) UnmarshalJSON(data []byte) error {
	s, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	return l.parse(s)
}

// AppendText implements [encoding.TextAppender]
// by calling [Level.String].
func (l Level) AppendText(b []byte) ([]byte, error) {
	return append(b, l.String()...), nil
}

// MarshalText implements [encoding.TextMarshaler]
// by calling [Level.AppendText].
func (l Level) MarshalText() ([]byte, error) {
	return l.AppendText(nil)
}

// UnmarshalText implements [encoding.TextUnmarshaler].
// It accepts any string produced by [Level.MarshalText],
// ignoring case.
// It also accepts numeric offsets that would result in a different string on
// output. For example, "Error-8" would marshal as "LevelInfo".
func (l *Level) UnmarshalText(data []byte) error {
	return l.parse(string(data))
}

func (l *Level) parse(s string) (err error) {
	switch strings.ToLower(s) {
	case "trace":
		*l = LevelTrace
	case "debug":
		*l = LevelDebug
	case "info":
		*l = LevelInfo
	case "warn":
		*l = LevelWarn
	case "error":
		*l = LevelError
	case "panic":
		*l = LevelPanic
	case "fatal":
		*l = LevelFatal
	default:
		return fmt.Errorf("l4g: level string %q: unknown name", s)
	}
	return nil
}

// Level returns the receiver.
// It implements [Leveler].
func (l Level) Level() Level { return l }

// A LevelVar is a [Level] variable, to allow a [Handler] level to change
// dynamically.
// It implements [Leveler] as well as a Set method,
// and it is safe for use by multiple goroutines.
// The zero LevelVar corresponds to [LevelInfo].
type LevelVar struct {
	val atomic.Int64
}

// NewLevelVar creates a new LevelVar from a Leveler.
// If the provided Leveler is already a *LevelVar, it is returned as-is.
// Otherwise, a new LevelVar is created with the level from the provided Leveler.
func NewLevelVar(lvl Leveler) *LevelVar {
	if l, ok := lvl.(*LevelVar); ok {
		return l
	}
	l := &LevelVar{}
	l.Set(lvl.Level())
	return l
}

// Int returns the integer representation of the level.
func (v *LevelVar) Int() int {
	return v.Level().Int()
}

// Level returns v's level.
func (v *LevelVar) Level() Level {
	return Level(int(v.val.Load()))
}

// Set sets v's level to l.
func (v *LevelVar) Set(l Level) {
	v.val.Store(int64(l))
}

// String returns a string representation of the LevelVar in the form "LevelVar(level)".
func (v *LevelVar) String() string {
	return fmt.Sprintf("LevelVar(%s)", v.Level())
}

// AppendText implements [encoding.TextAppender]
// by calling [Level.AppendText].
func (v *LevelVar) AppendText(b []byte) ([]byte, error) {
	return v.Level().AppendText(b)
}

// MarshalText implements [encoding.TextMarshaler]
// by calling [LevelVar.AppendText].
func (v *LevelVar) MarshalText() ([]byte, error) {
	return v.AppendText(nil)
}

// UnmarshalText implements [encoding.TextUnmarshaler]
// by calling [Level.UnmarshalText].
func (v *LevelVar) UnmarshalText(data []byte) error {
	var l Level
	if err := l.UnmarshalText(data); err != nil {
		return err
	}
	v.Set(l)
	return nil
}

// A Leveler provides a [Level] value.
//
// As Level itself implements Leveler, clients typically supply
// a Level value wherever a Leveler is needed, such as in [HandlerOptions].
// Clients who need to vary the level dynamically can provide a more complex
// Leveler implementation such as *[LevelVar].
type Leveler interface {
	Int() int
	Level() Level
}
