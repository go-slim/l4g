package l4g

import (
	"errors"
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
	VERBOSE Level = iota*4 - 8
	DEBUG
	INFO
	WARN
	ERROR
	PANIC
	FATAL
)

// Int returns the integer value of the level.
func (l Level) Int() int {
	return int(l)
}

// String returns a name for the level.
// If the level has a name, then that name
// in uppercase is returned.
// If the level is between named values, then
// an integer is appended to the uppercased name.
// Examples:
//
//	LevelWarn.String() => "WARN"
//	(LevelInfo+2).String() => "INFO+2"
func (l Level) String() string {
	str := func(base string, val Level) string {
		if val == 0 {
			return base
		}
		return fmt.Sprintf("%s%+d", base, val)
	}
	switch {
	case l < VERBOSE:
		return str("VERBOSE", l-VERBOSE)
	case l < DEBUG:
		return str("DEBUG", l-DEBUG)
	case l < INFO:
		return str("INFO", l-INFO)
	case l < WARN:
		return str("WARN", l-WARN)
	case l < ERROR:
		return str("ERROR", l-ERROR)
	case l < PANIC:
		return str("PANIC", l-PANIC)
	default:
		return str("FATAL", l-FATAL)
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
// output. For example, "Error-8" would marshal as "INFO".
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
// output. For example, "Error-8" would marshal as "INFO".
func (l *Level) UnmarshalText(data []byte) error {
	return l.parse(string(data))
}

func (l *Level) parse(s string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("slog: level string %q: %w", s, err)
		}
	}()

	name := s
	offset := 0
	if i := strings.IndexAny(s, "+-"); i >= 0 {
		name = s[:i]
		offset, err = strconv.Atoi(s[i:])
		if err != nil {
			return err
		}
	}
	switch strings.ToUpper(name) {
	case "VERBOSE":
		*l = VERBOSE
	case "DEBUG":
		*l = DEBUG
	case "INFO":
		*l = INFO
	case "WARN":
		*l = WARN
	case "ERROR":
		*l = ERROR
	case "PANIC":
		*l = PANIC
	case "FATAL":
		*l = FATAL
	default:
		return errors.New("unknown name")
	}
	*l += Level(offset)
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

// Level returns v's level.
func (v *LevelVar) Level() Level {
	return Level(int(v.val.Load()))
}

// Set sets v's level to l.
func (v *LevelVar) Set(l Level) {
	v.val.Store(int64(l))
}

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
