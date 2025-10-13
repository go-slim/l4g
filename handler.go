package l4g

import (
	"encoding"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// PartKind represents the different parts of a log record that can be customized.
type PartKind int8

const (
	// PartTime represents the timestamp part of a log record.
	PartTime PartKind = iota
	// PartLevel represents the log level part of a log record.
	PartLevel
	// PartMessage represents the message part of a log record.
	PartMessage
	// PartAttrs represents the attributes part of a log record.
	PartAttrs
	// PartPrefix represents the prefix part of a log record.
	PartPrefix
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

	// Handle handles the Record.
	// It will only be called when Enabled returns true.
	Handle(r Record) error

	// WithAttrs returns a new Handler whose attributes consist of
	// both the receiver's attributes and the arguments.
	// The Handler owns the slice: it may retain, modify or discard it.
	WithAttrs(attrs []Attr) Handler

	// WithGroup returns a new Handler with the given group appended to
	// the receiver's existing groups.
	// The keys of all subsequent attributes, whether added by With or in a
	// Record, should be qualified by the sequence of group names.
	//
	// How this qualification happens is up to the Handler, so long as
	// this Handler's attribute keys differ from those of another Handler
	// with a different sequence of group names.
	//
	// A Handler should treat WithGroup as starting a Group of Attrs that ends
	// at the end of the log event. That is,
	//
	//     logger.WithGroup("s").LogAttrs(ctx, level, msg, slog.Int("a", 1), slog.Int("b", 2))
	//
	// should behave like
	//
	//     logger.LogAttrs(ctx, level, msg, slog.Group("s", slog.Int("a", 1), slog.Int("b", 2)))
	//
	// If the name is empty, WithGroup returns the receiver.
	WithGroup(name string) Handler

	// WithPrefix returns a new Handler with the given prefix prepended to
	// the receiver's existing prefix.
	WithPrefix(prefix string) Handler
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

	// ReplaceAttr is called to rewrite each non-group attribute before it is logged.
	// See https://pkg.go.dev/log/slog#HandlerOptions for details.
	ReplaceAttr func(groups []string, attr Attr) Attr

	// Time format (Default: time.StampMilli)
	TimeFormat string

	// Disable color (Default: false)
	NoColor bool

	// Output is a destination to which log data will be written.
	Output io.Writer
}

const (
	// ANSI modes
	ansiEsc          = '\u001b'
	ansiReset        = "\u001b[0m"
	ansiFaint        = "\u001b[2m"
	ansiResetFaint   = "\u001b[22m"
	ansiBrightRed    = "\u001b[91m"
	ansiBrightGreen  = "\u001b[92m"
	ansiBrightYellow = "\u001b[93m"
	ansiBrightCyan   = "\u001b[96m"
	ansiGray         = "\u001b[90m"
	ansiWhite        = "\u001b[97m"
)

// Keys for "built-in" attributes.
const (
	// TimeKey is the key used by the built-in handlers for the time
	// when the log method is called. The associated Value is a [time.Time].
	TimeKey = "time"
	// LevelKey is the key used by the built-in handlers for the level
	// of the log call. The associated value is a [Level].
	LevelKey = "level"
	// MessageKey is the key used by the built-in handlers for the
	// message of the log call. The associated value is a string.
	MessageKey = "msg"
	// PrefixKey is the key used by the built-in handlers for the
	// prefix of the log call. The associated value is a string.
	PrefixKey = "prefix"
)

// NewSimpleHandler creates a [SimpleHandler] that writes to w,
// using the given options.
// If opts is nil, the default options are used.
func NewSimpleHandler(opts HandlerOptions) Handler {
	if opts.TimeFormat == "" {
		opts.TimeFormat = time.StampMilli
	}

	return &SimpleHandler{
		opts: &opts,
	}
}

var _ Handler = (*SimpleHandler)(nil)

// SimpleHandler is a Handler that writes formatted log records to an io.Writer.
// It formats records as single lines with space-separated fields, optionally
// colorized for terminal output. It supports structured logging with attributes,
// groups, and prefixes.
type SimpleHandler struct {
	attrsPrefix string          // Pre-formatted attributes from WithAttrs
	groupPrefix string          // Dot-separated group names for attributes
	groups      []string        // Stack of group names
	prefix      string          // Log prefix from WithPrefix
	opts        *HandlerOptions // Configuration options
}

// clone creates a shallow copy of the handler with a new groups slice.
// This is used by WithAttrs, WithGroup, and WithPrefix to create derived handlers.
func (h *SimpleHandler) clone() *SimpleHandler {
	return &SimpleHandler{
		attrsPrefix: h.attrsPrefix,
		groupPrefix: h.groupPrefix,
		groups:      h.groups,
		prefix:      h.prefix,
		opts:        h.opts,
	}
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
func (h *SimpleHandler) Enabled(level Level) bool {
	minLevel := LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

// Handle formats its argument [Record] as a single line of space-separated
// fields.
func (h *SimpleHandler) Handle(r Record) error {
	// get a buffer from the sync pool
	buf := newBuffer()
	defer buf.Free()

	rep := h.opts.ReplaceAttr

	// write time
	if !r.Time.IsZero() {
		val := r.Time.Round(0) // strip monotonic to match Attr behavior
		if rep == nil {
			h.appendTintTime(buf, r.Time, -1)
			buf.WriteByte(' ')
		} else if a := rep(nil /* groups */, slog.Time(TimeKey, val)); a.Key != "" {
			val, color := h.resolve(a.Value)
			if val.Kind() == slog.KindTime {
				h.appendTintTime(buf, val.Time(), color)
			} else {
				h.appendTintValue(buf, val, false, color, true)
			}
			buf.WriteByte(' ')
		}
	}

	// write level
	if rep == nil {
		h.appendTintLevel(buf, r.Level, -1)
		buf.WriteByte(' ')
	} else if a := rep(nil /* groups */, slog.Any(LevelKey, r.Level)); a.Key != "" {
		val, color := h.resolve(a.Value)
		if val.Kind() == slog.KindAny {
			if lvlVal, ok := val.Any().(Level); ok {
				h.appendTintLevel(buf, lvlVal, color)
			} else {
				h.appendTintValue(buf, val, false, color, false)
			}
		} else {
			h.appendTintValue(buf, val, false, color, false)
		}
		buf.WriteByte(' ')
	}

	//write prefix
	if r.Prefix != "" {
		if rep == nil {
			buf.WriteString("[" + r.Prefix + "]")
			buf.WriteByte(' ')
		} else if a := rep(nil /* groups */, slog.String(PrefixKey, r.Prefix)); a.Key != "" {
			val, color := h.resolve(a.Value)
			h.appendTintValue(buf, val, false, color, true)
			buf.WriteByte(' ')
		}
	}

	// write message
	if rep == nil {
		buf.WriteString(r.Message)
		buf.WriteByte(' ')
	} else if a := rep(nil /* groups */, slog.String(MessageKey, r.Message)); a.Key != "" {
		val, color := h.resolve(a.Value)
		h.appendTintValue(buf, val, false, color, false)
		buf.WriteByte(' ')
	}

	// write handler attributes
	if len(h.attrsPrefix) > 0 {
		buf.WriteString(h.attrsPrefix)
	}

	// write attributes
	r.Attrs(func(attr slog.Attr) bool {
		h.appendAttr(buf, attr, h.groupPrefix, h.groups)
		return true
	})

	if len(*buf) == 0 {
		buf.WriteByte('\n')
	} else {
		(*buf)[len(*buf)-1] = '\n' // replace last space with newline
	}

	_, err := h.opts.Output.Write(*buf)
	return err
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
func (h *SimpleHandler) WithAttrs(attrs []Attr) Handler {
	if len(attrs) == 0 {
		return h
	}

	buf := newBuffer()
	defer buf.Free()

	// write attributes to buffer
	for _, attr := range attrs {
		h.appendAttr(buf, attr, h.groupPrefix, h.groups)
	}

	h2 := h.clone()
	h2.attrsPrefix = h.attrsPrefix + string(*buf)
	return h2
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
func (h *SimpleHandler) WithGroup(name string) Handler {
	if name == "" {
		return h
	}

	h2 := h.clone()
	h2.groupPrefix += name + "."
	h2.groups = append(h2.groups, name)
	return h2
}

// WithPrefix returns a new Handler with the given prefix prepended to
// the receiver's existing prefix.
func (h *SimpleHandler) WithPrefix(prefix string) Handler {
	if prefix == "" {
		return h
	}

	h2 := h.clone()
	if h2.prefix == "" {
		h2.prefix = prefix
	} else {
		h2.prefix = prefix + h2.prefix
	}
	return h2
}

func (h *SimpleHandler) appendTintTime(buf *buffer, t time.Time, color int16) {
	if h.opts.NoColor {
		*buf = t.AppendFormat(*buf, h.opts.TimeFormat)
	} else {
		if color >= 0 {
			appendAnsi(buf, uint8(color), true)
		} else {
			buf.WriteString(ansiFaint)
		}
		*buf = t.AppendFormat(*buf, h.opts.TimeFormat)
		buf.WriteString(ansiReset)
	}
}

func (h *SimpleHandler) appendTintLevel(buf *buffer, level Level, color int16) {
	if !h.opts.NoColor {
		if color >= 0 {
			appendAnsi(buf, uint8(color), false)
		} else {
			switch level {
			case LevelTrace:
				buf.WriteString(ansiGray)
			case LevelDebug:
				buf.WriteString(ansiBrightCyan)
			case LevelInfo:
				buf.WriteString(ansiWhite)
			case LevelWarn:
				buf.WriteString(ansiBrightGreen)
			case LevelError:
				buf.WriteString(ansiBrightYellow)
			case LevelPanic:
				buf.WriteString(ansiBrightRed)
			default:
				if level <= LevelTrace {
					buf.WriteString(ansiGray)
				} else {
					buf.WriteString(ansiBrightRed)
				}
			}
		}
	}

	switch level {
	case LevelTrace:
		buf.WriteString("TRC")
	case LevelDebug:
		buf.WriteString("DBG")
	case LevelInfo:
		buf.WriteString("INF")
	case LevelWarn:
		buf.WriteString("WRN")
	case LevelError:
		buf.WriteString("ERR")
	case LevelPanic:
		buf.WriteString("PNL")
	default:
		// LevelFatal or higher
		buf.WriteString("FTL")
	}

	if !h.opts.NoColor {
		buf.WriteString(ansiReset)
	}
}

func appendSource(buf *buffer, src *slog.Source) {
	dir, file := filepath.Split(src.File)

	buf.WriteString(filepath.Join(filepath.Base(dir), file))
	buf.WriteByte(':')
	*buf = strconv.AppendInt(*buf, int64(src.Line), 10)
}

func (h *SimpleHandler) resolve(val slog.Value) (resolvedVal slog.Value, color int16) {
	if !h.opts.NoColor && val.Kind() == slog.KindLogValuer {
		if tintVal, ok := val.Any().(colorValue); ok {
			return tintVal.Value.Resolve(), int16(tintVal.Color)
		}
	}
	return val.Resolve(), -1
}

func (h *SimpleHandler) appendAttr(buf *buffer, attr slog.Attr, groupsPrefix string, groups []string) {
	var color int16 // -1 if no color
	attr.Value, color = h.resolve(attr.Value)
	if rep := h.opts.ReplaceAttr; rep != nil && attr.Value.Kind() != slog.KindGroup {
		attr = rep(groups, attr)
		var colorRep int16
		attr.Value, colorRep = h.resolve(attr.Value)
		if colorRep >= 0 {
			color = colorRep
		}
	}

	if attr.Equal(slog.Attr{}) {
		return
	}

	if attr.Value.Kind() == slog.KindGroup {
		if attr.Key != "" {
			groupsPrefix += attr.Key + "."
			groups = append(groups, attr.Key)
		}
		for _, groupAttr := range attr.Value.Group() {
			h.appendAttr(buf, groupAttr, groupsPrefix, groups)
		}
		return
	}

	if h.opts.NoColor {
		h.appendKey(buf, attr.Key, groupsPrefix)
		h.appendValue(buf, attr.Value, true)
	} else {
		if color >= 0 {
			appendAnsi(buf, uint8(color), true)
			h.appendKey(buf, attr.Key, groupsPrefix)
			buf.WriteString(ansiResetFaint)
			h.appendValue(buf, attr.Value, true)
			buf.WriteString(ansiReset)
		} else {
			buf.WriteString(ansiFaint)
			h.appendKey(buf, attr.Key, groupsPrefix)
			buf.WriteString(ansiReset)
			h.appendValue(buf, attr.Value, true)
		}
	}
	buf.WriteByte(' ')
}

func (h *SimpleHandler) appendKey(buf *buffer, key, groups string) {
	appendString(buf, groups+key, true, !h.opts.NoColor)
	buf.WriteByte('=')
}

func (h *SimpleHandler) appendValue(buf *buffer, v slog.Value, quote bool) {
	switch v.Kind() {
	case slog.KindString:
		appendString(buf, v.String(), quote, !h.opts.NoColor)
	case slog.KindInt64:
		*buf = strconv.AppendInt(*buf, v.Int64(), 10)
	case slog.KindUint64:
		*buf = strconv.AppendUint(*buf, v.Uint64(), 10)
	case slog.KindFloat64:
		*buf = strconv.AppendFloat(*buf, v.Float64(), 'g', -1, 64)
	case slog.KindBool:
		*buf = strconv.AppendBool(*buf, v.Bool())
	case slog.KindDuration:
		appendString(buf, v.Duration().String(), quote, !h.opts.NoColor)
	case slog.KindTime:
		*buf = appendRFC3339Millis(*buf, v.Time())
	case slog.KindAny:
		defer func() {
			// Copied from log/slog/handler.go.
			if r := recover(); r != nil {
				// If it panics with a nil pointer, the most likely cases are
				// an encoding.TextMarshaler or error fails to guard against nil,
				// in which case "<nil>" seems to be the feasible choice.
				//
				// Adapted from the code in fmt/print.go.
				if v := reflect.ValueOf(v.Any()); v.Kind() == reflect.Pointer && v.IsNil() {
					buf.WriteString("<nil>")
					return
				}

				// Otherwise just print the original panic message.
				appendString(buf, fmt.Sprintf("!PANIC: %v", r), true, !h.opts.NoColor)
			}
		}()

		switch cv := v.Any().(type) {
		case encoding.TextMarshaler:
			data, err := cv.MarshalText()
			if err != nil {
				break
			}
			appendString(buf, string(data), quote, !h.opts.NoColor)
		case *slog.Source:
			appendSource(buf, cv)
		default:
			appendString(buf, fmt.Sprintf("%+v", cv), quote, !h.opts.NoColor)
		}
	default:
		// Handle unknown kinds (e.g., KindGroup, KindLogValuer, or future kinds)
		// KindGroup is typically handled in appendAttr, but this provides a fallback
		// KindLogValuer should be resolved before reaching here, but this is defensive
		if v.Kind() == slog.KindGroup {
			// Format group as inline attributes
			attrs := v.Group()
			buf.WriteByte('{')
			for i, attr := range attrs {
				if i > 0 {
					buf.WriteByte(' ')
				}
				buf.WriteString(attr.Key)
				buf.WriteByte(':')
				h.appendValue(buf, attr.Value, true)
			}
			buf.WriteByte('}')
		} else {
			// For any other unknown kind, format as string using %+v
			appendString(buf, fmt.Sprintf("%+v(%v)", v.Kind(), v.Any()), quote, !h.opts.NoColor)
		}
	}
}

func (h *SimpleHandler) appendTintValue(buf *buffer, val slog.Value, quote bool, color int16, faint bool) {
	if h.opts.NoColor {
		h.appendValue(buf, val, quote)
	} else {
		if color >= 0 {
			appendAnsi(buf, uint8(color), faint)
		} else if faint {
			buf.WriteString(ansiFaint)
		}
		h.appendValue(buf, val, quote)
		if color >= 0 || faint {
			buf.WriteString(ansiReset)
		}
	}
}

// Copied from log/slog/handler.go.
func appendRFC3339Millis(b []byte, t time.Time) []byte {
	// Format according to time.RFC3339Nano since it is highly optimized,
	// but truncate it to use millisecond resolution.
	// Unfortunately, that format trims trailing 0s, so add 1/10 millisecond
	// to guarantee that there are exactly 4 digits after the period.
	const prefixLen = len("2006-01-02T15:04:05.000")
	n := len(b)
	t = t.Truncate(time.Millisecond).Add(time.Millisecond / 10)
	b = t.AppendFormat(b, time.RFC3339Nano)
	b = append(b[:n+prefixLen], b[n+prefixLen+1:]...) // drop the 4th digit
	return b
}

func appendAnsi(buf *buffer, color uint8, faint bool) {
	buf.WriteString("\u001b[")
	if faint {
		buf.WriteString("2;")
	}
	if color < 8 {
		*buf = strconv.AppendUint(*buf, uint64(color)+30, 10)
	} else if color < 16 {
		*buf = strconv.AppendUint(*buf, uint64(color)+82, 10)
	} else {
		buf.WriteString("38;5;")
		*buf = strconv.AppendUint(*buf, uint64(color), 10)
	}
	buf.WriteByte('m')
}

func appendString(buf *buffer, s string, quote, color bool) {
	if quote && !color {
		// trim ANSI escape sequences
		var inEscape bool
		s = cut(s, func(r rune) bool {
			if r == ansiEsc {
				inEscape = true
			} else if inEscape && unicode.IsLetter(r) {
				inEscape = false
				return true
			}

			return inEscape
		})
	}

	quote = quote && needsQuoting(s)
	switch {
	case color && quote:
		s = strconv.Quote(s)
		s = strings.ReplaceAll(s, `\x1b`, string(ansiEsc))
		buf.WriteString(s)
	case !color && quote:
		*buf = strconv.AppendQuote(*buf, s)
	default:
		buf.WriteString(s)
	}
}

func cut(s string, f func(r rune) bool) string {
	var res []rune
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError {
			break
		}
		if !f(r) {
			res = append(res, r)
		}
		i += size
	}
	return string(res)
}

// Copied from log/slog/text_handler.go.
func needsQuoting(s string) bool {
	if len(s) == 0 {
		return true
	}
	for i := 0; i < len(s); {
		b := s[i]
		if b < utf8.RuneSelf {
			// Quote anything except a backslash that would need quoting in a
			// JSON string, as well as space and '='
			if b != '\\' && (b == ' ' || b == '=' || !safeSet[b]) {
				return true
			}
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError || unicode.IsSpace(r) || !unicode.IsPrint(r) {
			return true
		}
		i += size
	}
	return false
}

// Copied from log/slog/json_handler.go.
//
// safeSet is extended by the ANSI escape code "\u001b".
var safeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      true,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      true,
	'=':      true,
	'>':      true,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
	'\u001b': true,
}
