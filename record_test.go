package l4g

import (
	"log/slog"
	"testing"
	"time"
)

func TestNewRecord(t *testing.T) {
	now := time.Now()
	msg := "test message"
	level := LevelInfo

	r := NewRecord(now, level, msg)

	if !r.Time.Equal(now) {
		t.Errorf("NewRecord() time = %v, want %v", r.Time, now)
	}
	if r.Level != level {
		t.Errorf("NewRecord() level = %v, want %v", r.Level, level)
	}
	if r.Message != msg {
		t.Errorf("NewRecord() message = %v, want %v", r.Message, msg)
	}
	if r.NumAttrs() != 0 {
		t.Errorf("NewRecord() attrs count = %v, want 0", r.NumAttrs())
	}
}

func TestRecord_AddAttrs(t *testing.T) {
	t.Run("single attr", func(t *testing.T) {
		r := NewRecord(time.Now(), LevelInfo, "test")
		r.AddAttrs(String("key", "value"))
		if r.NumAttrs() != 1 {
			t.Errorf("Record.AddAttrs() count = %v, want 1", r.NumAttrs())
		}
	})

	t.Run("multiple attrs", func(t *testing.T) {
		r := NewRecord(time.Now(), LevelInfo, "test")
		r.AddAttrs(String("a", "1"), Int("b", 2), Bool("c", true))
		if r.NumAttrs() != 3 {
			t.Errorf("Record.AddAttrs() count = %v, want 3", r.NumAttrs())
		}
	})

	t.Run("many attrs", func(t *testing.T) {
		r := NewRecord(time.Now(), LevelInfo, "test")
		attrs := make([]Attr, 10)
		for i := range 10 {
			attrs[i] = Int("key", i)
		}
		r.AddAttrs(attrs...)
		if r.NumAttrs() != 10 {
			t.Errorf("Record.AddAttrs() count = %v, want 10", r.NumAttrs())
		}
	})

	t.Run("empty group", func(t *testing.T) {
		r := NewRecord(time.Now(), LevelInfo, "test")
		r.AddAttrs(Group("empty"), String("key", "value"))
		if r.NumAttrs() != 1 {
			t.Errorf("Record.AddAttrs() should skip empty groups, count = %v", r.NumAttrs())
		}
	})
}

func TestRecord_Add(t *testing.T) {
	t.Run("key-value pairs", func(t *testing.T) {
		r := NewRecord(time.Now(), LevelInfo, "test")
		r.Add("key1", "value1", "key2", 42)
		if r.NumAttrs() != 2 {
			t.Errorf("Record.Add() count = %v, want 2", r.NumAttrs())
		}
	})

	t.Run("with attrs", func(t *testing.T) {
		r := NewRecord(time.Now(), LevelInfo, "test")
		r.Add(String("x", "y"), "a", 1)
		if r.NumAttrs() != 2 {
			t.Errorf("Record.Add() count = %v, want 2", r.NumAttrs())
		}
	})

	t.Run("single string", func(t *testing.T) {
		r := NewRecord(time.Now(), LevelInfo, "test")
		r.Add("orphan")
		if r.NumAttrs() != 1 {
			t.Errorf("Record.Add() count = %v, want 1", r.NumAttrs())
		}
	})

	t.Run("many pairs", func(t *testing.T) {
		r := NewRecord(time.Now(), LevelInfo, "test")
		args := make([]any, 20)
		for i := range 10 {
			args[i*2] = "key"
			args[i*2+1] = i
		}
		r.Add(args...)
		if r.NumAttrs() != 10 {
			t.Errorf("Record.Add() count = %v, want 10", r.NumAttrs())
		}
	})

	t.Run("empty group", func(t *testing.T) {
		r := NewRecord(time.Now(), LevelInfo, "test")
		r.Add(Group("empty"), "key", "value")
		if r.NumAttrs() != 1 {
			t.Errorf("Record.Add() should skip empty groups, count = %v", r.NumAttrs())
		}
	})
}

func TestRecord_NumAttrs(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*Record)
		want  int
	}{
		{
			name:  "empty",
			setup: func(r *Record) {},
			want:  0,
		},
		{
			name: "inline only",
			setup: func(r *Record) {
				r.AddAttrs(String("a", "1"), String("b", "2"))
			},
			want: 2,
		},
		{
			name: "overflow to back",
			setup: func(r *Record) {
				for i := range 10 {
					r.AddAttrs(Int("key", i))
				}
			},
			want: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRecord(time.Now(), LevelInfo, "test")
			tt.setup(&r)
			if got := r.NumAttrs(); got != tt.want {
				t.Errorf("Record.NumAttrs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecord_Attrs(t *testing.T) {
	r := NewRecord(time.Now(), LevelInfo, "test")
	r.AddAttrs(String("a", "1"), Int("b", 2), Bool("c", true))

	count := 0
	r.Attrs(func(attr Attr) bool {
		count++
		return true
	})

	if count != 3 {
		t.Errorf("Record.Attrs() iterated %v times, want 3", count)
	}
}

func TestRecord_Attrs_EarlyExit(t *testing.T) {
	r := NewRecord(time.Now(), LevelInfo, "test")
	r.AddAttrs(String("a", "1"), Int("b", 2), Bool("c", true))

	count := 0
	r.Attrs(func(attr Attr) bool {
		count++
		return count < 2 // Stop after 2 iterations
	})

	if count != 2 {
		t.Errorf("Record.Attrs() iterated %v times, want 2", count)
	}
}

func TestRecord_Clone(t *testing.T) {
	original := NewRecord(time.Now(), LevelInfo, "test")
	original.AddAttrs(String("a", "1"), Int("b", 2))

	cloned := original.Clone()

	// Verify fields are copied
	if !cloned.Time.Equal(original.Time) {
		t.Errorf("Clone() time mismatch")
	}
	if cloned.Level != original.Level {
		t.Errorf("Clone() level mismatch")
	}
	if cloned.Message != original.Message {
		t.Errorf("Clone() message mismatch")
	}
	if cloned.NumAttrs() != original.NumAttrs() {
		t.Errorf("Clone() attrs count mismatch")
	}

	// Verify modifications don't affect original
	cloned.AddAttrs(String("c", "3"))
	if original.NumAttrs() == cloned.NumAttrs() {
		t.Errorf("Clone() should create independent copy")
	}
}

func TestRecord_Prefix(t *testing.T) {
	r := NewRecord(time.Now(), LevelInfo, "test")
	r.Prefix = "myapp"

	if r.Prefix != "myapp" {
		t.Errorf("Record.Prefix = %v, want 'myapp'", r.Prefix)
	}
}

func TestCountAttrs(t *testing.T) {
	tests := []struct {
		name string
		args []any
		want int
	}{
		{"empty", []any{}, 0},
		{"key-value", []any{"key", "value"}, 1},
		{"multiple pairs", []any{"a", 1, "b", 2}, 2},
		{"with attr", []any{String("x", "y")}, 1},
		{"single string", []any{"orphan"}, 1},
		{"single value", []any{42}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countAttrs(tt.args)
			if got != tt.want {
				t.Errorf("countAttrs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCountEmptyGroups(t *testing.T) {
	tests := []struct {
		name  string
		attrs []Attr
		want  int
	}{
		{"no groups", []Attr{String("a", "1"), Int("b", 2)}, 0},
		{"one empty group", []Attr{Group("empty"), String("a", "1")}, 1},
		{"multiple empty groups", []Attr{Group("e1"), Group("e2"), String("a", "1")}, 2},
		{"non-empty group", []Attr{Group("g", String("x", "y"))}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countEmptyGroups(tt.attrs)
			if got != tt.want {
				t.Errorf("countEmptyGroups() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEmptyAttr(t *testing.T) {
	tests := []struct {
		name string
		attr Attr
		want bool
	}{
		{"empty attr", Attr{}, true},
		{"any nil", Any("", nil), true},
		{"string attr", String("key", "value"), false},
		{"empty key with value", slog.Attr{Key: "", Value: slog.IntValue(42)}, false},
		{"key with nil", slog.Attr{Key: "key", Value: slog.AnyValue(nil)}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isEmptyAttr(tt.attr)
			if got != tt.want {
				t.Errorf("isEmptyAttr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEmptyGroup(t *testing.T) {
	tests := []struct {
		name  string
		value slog.Value
		want  bool
	}{
		{"empty group", slog.GroupValue(), true},
		{"non-empty group", slog.GroupValue(String("a", "1")), false},
		{"string value", slog.StringValue("test"), false},
		{"int value", slog.IntValue(42), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isEmptyGroup(tt.value)
			if got != tt.want {
				t.Errorf("isEmptyGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArgsToAttr(t *testing.T) {
	tests := []struct {
		name       string
		args       []any
		wantKey    string
		wantRest   int
		wantBadKey bool
	}{
		{
			name:       "string key-value",
			args:       []any{"key", "value", "extra"},
			wantKey:    "key",
			wantRest:   1,
			wantBadKey: false,
		},
		{
			name:       "attr",
			args:       []any{String("x", "y"), "extra"},
			wantKey:    "x",
			wantRest:   1,
			wantBadKey: false,
		},
		{
			name:       "single string",
			args:       []any{"orphan"},
			wantKey:    badKey,
			wantRest:   0,
			wantBadKey: true,
		},
		{
			name:       "non-string non-attr",
			args:       []any{42, "extra"},
			wantKey:    badKey,
			wantRest:   1,
			wantBadKey: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr, rest := argsToAttr(tt.args)
			if attr.Key != tt.wantKey {
				t.Errorf("argsToAttr() key = %v, want %v", attr.Key, tt.wantKey)
			}
			if len(rest) != tt.wantRest {
				t.Errorf("argsToAttr() rest length = %v, want %v", len(rest), tt.wantRest)
			}
		})
	}
}

func TestRecord_FrontBackSplit(t *testing.T) {
	r := NewRecord(time.Now(), LevelInfo, "test")

	// Add exactly nAttrsInline attributes
	for i := range nAttrsInline {
		r.AddAttrs(Int("key", i))
	}

	if r.nFront != nAttrsInline {
		t.Errorf("nFront = %v, want %v", r.nFront, nAttrsInline)
	}
	if len(r.back) != 0 {
		t.Errorf("back length = %v, want 0", len(r.back))
	}

	// Add one more to overflow to back
	r.AddAttrs(String("overflow", "value"))

	if r.nFront != nAttrsInline {
		t.Errorf("nFront = %v, want %v", r.nFront, nAttrsInline)
	}
	if len(r.back) != 1 {
		t.Errorf("back length = %v, want 1", len(r.back))
	}
}

func BenchmarkRecord_AddAttrs(b *testing.B) {
	attrs := []Attr{
		String("key1", "value1"),
		Int("key2", 42),
		Bool("key3", true),
	}

	for b.Loop() {
		r := NewRecord(time.Now(), LevelInfo, "test")
		r.AddAttrs(attrs...)
	}
}

func BenchmarkRecord_Add(b *testing.B) {
	for b.Loop() {
		r := NewRecord(time.Now(), LevelInfo, "test")
		r.Add("key1", "value1", "key2", 42, "key3", true)
	}
}

func BenchmarkRecord_Clone(b *testing.B) {
	r := NewRecord(time.Now(), LevelInfo, "test")
	r.AddAttrs(String("a", "1"), Int("b", 2), Bool("c", true))

	for b.Loop() {
		_ = r.Clone()
	}
}
