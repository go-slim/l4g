package l4g

import (
	"errors"
	"log/slog"
	"testing"
	"time"
)

func TestString(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
		want  string
	}{
		{"simple", "key", "value", "value"},
		{"empty", "key", "", ""},
		{"unicode", "key", "你好", "你好"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := String(tt.key, tt.value)
			if attr.Key != tt.key {
				t.Errorf("String() key = %v, want %v", attr.Key, tt.key)
			}
			if attr.Value.String() != tt.want {
				t.Errorf("String() value = %v, want %v", attr.Value.String(), tt.want)
			}
		})
	}
}

func TestInt(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value int
		want  int64
	}{
		{"positive", "count", 42, 42},
		{"negative", "offset", -10, -10},
		{"zero", "zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := Int(tt.key, tt.value)
			if attr.Key != tt.key {
				t.Errorf("Int() key = %v, want %v", attr.Key, tt.key)
			}
			if attr.Value.Int64() != tt.want {
				t.Errorf("Int() value = %v, want %v", attr.Value.Int64(), tt.want)
			}
		})
	}
}

func TestInt_CustomTypes(t *testing.T) {
	type MyInt int
	type MyInt8 int8
	type MyInt16 int16
	type MyInt32 int32
	type MyInt64 int64

	t.Run("int", func(t *testing.T) {
		attr := Int("key", MyInt(42))
		if attr.Value.Int64() != 42 {
			t.Errorf("Int(MyInt) value = %v, want 42", attr.Value.Int64())
		}
	})

	t.Run("int8", func(t *testing.T) {
		attr := Int("key", MyInt8(8))
		if attr.Value.Int64() != 8 {
			t.Errorf("Int(MyInt8) value = %v, want 8", attr.Value.Int64())
		}
	})

	t.Run("int16", func(t *testing.T) {
		attr := Int("key", MyInt16(16))
		if attr.Value.Int64() != 16 {
			t.Errorf("Int(MyInt16) value = %v, want 16", attr.Value.Int64())
		}
	})

	t.Run("int32", func(t *testing.T) {
		attr := Int("key", MyInt32(32))
		if attr.Value.Int64() != 32 {
			t.Errorf("Int(MyInt32) value = %v, want 32", attr.Value.Int64())
		}
	})

	t.Run("int64", func(t *testing.T) {
		attr := Int("key", MyInt64(64))
		if attr.Value.Int64() != 64 {
			t.Errorf("Int(MyInt64) value = %v, want 64", attr.Value.Int64())
		}
	})
}

func TestUint(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value uint
		want  uint64
	}{
		{"positive", "count", 42, 42},
		{"zero", "zero", 0, 0},
		{"large", "large", 18446744073709551615, 18446744073709551615},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := Uint(tt.key, tt.value)
			if attr.Key != tt.key {
				t.Errorf("Uint() key = %v, want %v", attr.Key, tt.key)
			}
			if attr.Value.Uint64() != tt.want {
				t.Errorf("Uint() value = %v, want %v", attr.Value.Uint64(), tt.want)
			}
		})
	}
}

func TestUint_CustomTypes(t *testing.T) {
	type MyUint uint
	type MyUint8 uint8
	type MyUint16 uint16
	type MyUint32 uint32
	type MyUint64 uint64

	t.Run("uint", func(t *testing.T) {
		attr := Uint("key", MyUint(42))
		if attr.Value.Uint64() != 42 {
			t.Errorf("Uint(MyUint) value = %v, want 42", attr.Value.Uint64())
		}
	})

	t.Run("uint8", func(t *testing.T) {
		attr := Uint("key", MyUint8(8))
		if attr.Value.Uint64() != 8 {
			t.Errorf("Uint(MyUint8) value = %v, want 8", attr.Value.Uint64())
		}
	})

	t.Run("uint16", func(t *testing.T) {
		attr := Uint("key", MyUint16(16))
		if attr.Value.Uint64() != 16 {
			t.Errorf("Uint(MyUint16) value = %v, want 16", attr.Value.Uint64())
		}
	})

	t.Run("uint32", func(t *testing.T) {
		attr := Uint("key", MyUint32(32))
		if attr.Value.Uint64() != 32 {
			t.Errorf("Uint(MyUint32) value = %v, want 32", attr.Value.Uint64())
		}
	})

	t.Run("uint64", func(t *testing.T) {
		attr := Uint("key", MyUint64(64))
		if attr.Value.Uint64() != 64 {
			t.Errorf("Uint(MyUint64) value = %v, want 64", attr.Value.Uint64())
		}
	})
}

func TestFloat(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value float64
		want  float64
	}{
		{"positive", "value", 3.14, 3.14},
		{"negative", "value", -2.5, -2.5},
		{"zero", "value", 0.0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := Float(tt.key, tt.value)
			if attr.Key != tt.key {
				t.Errorf("Float() key = %v, want %v", attr.Key, tt.key)
			}
			if attr.Value.Float64() != tt.want {
				t.Errorf("Float() value = %v, want %v", attr.Value.Float64(), tt.want)
			}
		})
	}
}

func TestFloat_CustomTypes(t *testing.T) {
	type MyFloat32 float32
	type MyFloat64 float64

	t.Run("float32", func(t *testing.T) {
		attr := Float("key", MyFloat32(3.14))
		if attr.Value.Float64()-3.14 > 0.01 {
			t.Errorf("Float(MyFloat32) value = %v, want ~3.14", attr.Value.Float64())
		}
	})

	t.Run("float64", func(t *testing.T) {
		attr := Float("key", MyFloat64(2.71))
		if attr.Value.Float64() != 2.71 {
			t.Errorf("Float(MyFloat64) value = %v, want 2.71", attr.Value.Float64())
		}
	})
}

func TestBool(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value bool
		want  bool
	}{
		{"true", "enabled", true, true},
		{"false", "enabled", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := Bool(tt.key, tt.value)
			if attr.Key != tt.key {
				t.Errorf("Bool() key = %v, want %v", attr.Key, tt.key)
			}
			if attr.Value.Bool() != tt.want {
				t.Errorf("Bool() value = %v, want %v", attr.Value.Bool(), tt.want)
			}
		})
	}
}

func TestTime(t *testing.T) {
	now := time.Now()
	attr := Time("time", now)

	if attr.Key != "time" {
		t.Errorf("Time() key = %v, want 'time'", attr.Key)
	}
	if !attr.Value.Time().Equal(now) {
		t.Errorf("Time() value = %v, want %v", attr.Value.Time(), now)
	}
}

func TestDuration(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		duration time.Duration
	}{
		{"seconds", "duration", 5 * time.Second},
		{"milliseconds", "duration", 100 * time.Millisecond},
		{"nanoseconds", "duration", 1 * time.Nanosecond},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := Duration(tt.key, tt.duration)
			if attr.Key != tt.key {
				t.Errorf("Duration() key = %v, want %v", attr.Key, tt.key)
			}
			if attr.Value.Duration() != tt.duration {
				t.Errorf("Duration() value = %v, want %v", attr.Value.Duration(), tt.duration)
			}
		})
	}
}

func TestGroup(t *testing.T) {
	attr := Group("group", String("a", "1"), Int("b", 2))

	if attr.Key != "group" {
		t.Errorf("Group() key = %v, want 'group'", attr.Key)
	}
	if attr.Value.Kind() != slog.KindGroup {
		t.Errorf("Group() kind = %v, want KindGroup", attr.Value.Kind())
	}

	group := attr.Value.Group()
	if len(group) != 2 {
		t.Errorf("Group() length = %v, want 2", len(group))
	}
}

func TestAny(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value any
	}{
		{"struct", "key", struct{ Name string }{Name: "test"}},
		{"nil", "key", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := Any(tt.key, tt.value)
			if attr.Key != tt.key {
				t.Errorf("Any() key = %v, want %v", attr.Key, tt.key)
			}
			// Note: slog.Any() optimizes some types, so we just check the key
		})
	}
}

func TestColorAttr(t *testing.T) {
	attr := String("key", "value")
	coloredAttr := ColorAttr(42, attr)

	if coloredAttr.Key != "key" {
		t.Errorf("ColorAttr() key = %v, want 'key'", coloredAttr.Key)
	}

	// Verify that the value is wrapped in colorValue (which is a LogValuer)
	if coloredAttr.Value.Kind() != slog.KindLogValuer {
		t.Errorf("ColorAttr() kind = %v, want KindLogValuer", coloredAttr.Value.Kind())
	}

	// Verify LogValue() returns the original value
	if cv, ok := coloredAttr.Value.Any().(colorValue); ok {
		if cv.Color != 42 {
			t.Errorf("ColorAttr() color = %v, want 42", cv.Color)
		}
		if cv.Value.String() != "value" {
			t.Errorf("ColorAttr() original value = %v, want 'value'", cv.Value.String())
		}
	} else {
		t.Errorf("ColorAttr() value is not colorValue")
	}
}

func TestArgsToAttrSlice(t *testing.T) {
	tests := []struct {
		name string
		args []any
		want int
	}{
		{"empty", []any{}, 0},
		{"key-value", []any{"key", "value"}, 1},
		{"multiple", []any{"a", 1, "b", 2, "c", 3}, 3},
		{"with attr", []any{String("x", "y"), "a", 1}, 2},
		{"single string", []any{"orphan"}, 1},
		{"single value", []any{42}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := argsToAttrSlice(tt.args)
			if len(attrs) != tt.want {
				t.Errorf("argsToAttrSlice() length = %v, want %v", len(attrs), tt.want)
			}
		})
	}
}

func TestSplitAttrs(t *testing.T) {
	tests := []struct {
		name      string
		args      []any
		wantAttrs int
		wantAnies int
	}{
		{"empty", []any{}, 0, 0},
		{"only attrs", []any{String("a", "1"), Int("b", 2)}, 2, 0},
		{"only anies", []any{"a", 1, 2.5}, 0, 3},
		{"mixed", []any{String("x", "y"), "a", 1, Int("z", 3)}, 2, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs, anies := splitAttrs(tt.args)
			if len(attrs) != tt.wantAttrs {
				t.Errorf("splitAttrs() attrs length = %v, want %v", len(attrs), tt.wantAttrs)
			}
			if len(anies) != tt.wantAnies {
				t.Errorf("splitAttrs() anies length = %v, want %v", len(anies), tt.wantAnies)
			}
		})
	}
}

func TestColorValue_LogValue(t *testing.T) {
	cv := colorValue{
		Value: slog.StringValue("test"),
		Color: 42,
	}

	lv := cv.LogValue()
	if lv.String() != "test" {
		t.Errorf("colorValue.LogValue() = %v, want 'test'", lv.String())
	}
}

func TestErr(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"simple error", errors.New("simple error")},
		{"custom error", &customError{"custom error message"}},
		{"nil error", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := Err(tt.err)

			// Verify the key is "error"
			if attr.Key != "error" {
				t.Errorf("Err() key = %v, want 'error'", attr.Key)
			}

			// Verify that the value is wrapped in colorValue
			if attr.Value.Kind() != slog.KindLogValuer {
				t.Errorf("Err() kind = %v, want KindLogValuer", attr.Value.Kind())
			}

			// Verify the color is 9 (bright red)
			if cv, ok := attr.Value.Any().(colorValue); ok {
				if cv.Color != 9 {
					t.Errorf("Err() color = %v, want 9 (bright red)", cv.Color)
				}
				// Verify the error value is preserved
				if tt.err != nil {
					actualErr := cv.Value.Any()
					if actualErr != tt.err {
						t.Errorf("Err() error value = %v, want %v", actualErr, tt.err)
					}
				}
			} else {
				t.Errorf("Err() value is not colorValue")
			}
		})
	}
}

func TestErr_Integration(t *testing.T) {
	// Test that Err behaves equivalently to ColorAttr(9, Any("error", err))
	testErr := &customError{"test error"}

	errAttr := Err(testErr)
	manualAttr := ColorAttr(9, Any("error", testErr))

	if errAttr.Key != manualAttr.Key {
		t.Errorf("Err() key = %v, manual key = %v", errAttr.Key, manualAttr.Key)
	}

	errCV, ok1 := errAttr.Value.Any().(colorValue)
	manualCV, ok2 := manualAttr.Value.Any().(colorValue)

	if !ok1 || !ok2 {
		t.Fatal("Both should be colorValue")
	}

	if errCV.Color != manualCV.Color {
		t.Errorf("Err() color = %v, manual color = %v", errCV.Color, manualCV.Color)
	}
}

// customError is a test helper for error testing
type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}
