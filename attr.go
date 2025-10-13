package l4g

import (
	"log/slog"
	"time"
)

// Attr is an alias for slog.Attr, representing a key-value pair for structured logging.
type Attr = slog.Attr

// String returns an Attr for a string value.
// It supports any type with an underlying string type.
func String[T ~string](key string, value T) slog.Attr {
	return slog.String(key, string(value))
}

// Int64 returns an Attr for an int64 value.
func Int64(key string, value int64) Attr {
	return slog.Int64(key, value)
}

// Int returns an Attr for an integer value.
// It supports any type with an underlying int, int8, int16, int32, or int64 type.
func Int[T ~int | ~int8 | ~int16 | ~int32 | ~int64](key string, value T) slog.Attr {
	return slog.Int(key, int(value))
}

// Uint returns an Attr for an unsigned integer value.
// It supports any type with an underlying uint, uint8, uint16, uint32, or uint64 type.
func Uint[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](key string, value T) slog.Attr {
	return slog.Uint64(key, uint64(value))
}

// Float returns an Attr for a floating-point value.
// It supports any type with an underlying float32 or float64 type.
func Float[T ~float32 | ~float64](key string, value T) slog.Attr {
	return slog.Float64(key, float64(value))
}

// Bool returns an Attr for a boolean value.
// It supports any type with an underlying bool type.
func Bool[T ~bool](key string, v T) slog.Attr {
	return slog.Bool(key, bool(v))
}

// Time returns an Attr for a time.Time value.
func Time(key string, v time.Time) slog.Attr {
	return slog.Time(key, v)
}

// Duration returns an Attr for a time.Duration value.
func Duration(key string, value time.Duration) slog.Attr {
	return slog.Duration(key, value)
}

// Group returns an Attr for a group of attributes.
// The args can be Attr values or alternating key-value pairs (string, any, string, any, ...).
func Group(key string, args ...any) slog.Attr {
	return slog.Group(key, args...)
}

// Any returns an Attr for any value type.
// The value is stored as-is and formatted according to its type at output time.
func Any(key string, value any) slog.Attr {
	return slog.Any(key, value)
}

// colorValue wraps a slog.Value with a color code for colorized output.
// It implements slog.LogValuer to transparently pass through the underlying value
// while preserving the color information for handlers that support it.
type colorValue struct {
	slog.Value
	Color uint8
}

// LogValue implements the [slog.LogValuer] interface.
func (v colorValue) LogValue() slog.Value {
	return v.Value
}

// ColorAttr returns a tinted (colorized) [slog.Attr] that will be written in the
// specified color by the [tint.Handler]. When used with any other [slog.Handler], it behaves as a
// plain [slog.Attr].
//
// Use the uint8 color value to specify the color of the attribute:
//
//   - 0-7: standard ANSI colors
//   - 8-15: high intensity ANSI colors
//   - 16-231: 216 colors (6×6×6 cube)
//   - 232-255: grayscale from dark to light in 24 steps
//
// See https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit
func ColorAttr(color uint8, attr slog.Attr) slog.Attr {
	attr.Value = slog.AnyValue(colorValue{attr.Value, color})
	return attr
}

func argsToAttrSlice(args []any) []slog.Attr {
	if len(args) == 0 {
		return nil
	}
	var (
		attr slog.Attr
		// Pre-allocate with estimated capacity to reduce allocations
		// argsToAttr typically consumes 1-2 args per iteration
		attrs = make([]slog.Attr, 0, len(args)/2+1)
	)
	for len(args) > 0 {
		attr, args = argsToAttr(args)
		attrs = append(attrs, attr)
	}
	return attrs
}

func splitAttrs(args []any) ([]slog.Attr, []any) {
	if len(args) == 0 {
		return nil, nil
	}
	// Pre-allocate with full capacity to avoid reallocation
	// In worst case, all items go into one slice
	attrs := make([]slog.Attr, 0, len(args))
	remaining := make([]any, 0, len(args))

	for _, arg := range args {
		if attr, ok := arg.(slog.Attr); ok {
			attrs = append(attrs, attr)
		} else {
			remaining = append(remaining, arg)
		}
	}

	return attrs, remaining
}
