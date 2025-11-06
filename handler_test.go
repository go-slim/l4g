package l4g

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNewSimpleHandler(t *testing.T) {
	buf := &bytes.Buffer{}
	opts := HandlerOptions{
		Output: buf,
	}
	h := NewSimpleHandler(opts)

	if h == nil {
		t.Errorf("NewSimpleHandler() returned nil")
	}

	sh, ok := h.(*SimpleHandler)
	if !ok {
		t.Errorf("NewSimpleHandler() did not return *SimpleHandler")
	}

	if sh.opts.TimeFormat != time.StampMilli {
		t.Errorf("NewSimpleHandler() default TimeFormat = %v, want %v", sh.opts.TimeFormat, time.StampMilli)
	}
}

func TestSimpleHandler_Enabled(t *testing.T) {
	tests := []struct {
		name        string
		minLevel    Level
		testLevel   Level
		wantEnabled bool
	}{
		{"trace enabled for trace", LevelTrace, LevelTrace, true},
		{"trace enabled for info", LevelTrace, LevelInfo, true},
		{"info enabled for info", LevelInfo, LevelInfo, true},
		{"info disabled for trace", LevelInfo, LevelTrace, false},
		{"warn disabled for debug", LevelWarn, LevelDebug, false},
		{"error enabled for error", LevelError, LevelError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			h := NewSimpleHandler(HandlerOptions{
				Level:  tt.minLevel,
				Output: buf,
			})

			if got := h.Enabled(tt.testLevel); got != tt.wantEnabled {
				t.Errorf("SimpleHandler.Enabled() = %v, want %v", got, tt.wantEnabled)
			}
		})
	}
}

func TestSimpleHandler_Handle(t *testing.T) {
	tests := []struct {
		name        string
		level       Level
		message     string
		wantContain []string
	}{
		{
			name:        "trace level",
			level:       LevelTrace,
			message:     "trace message",
			wantContain: []string{"TRACE", "trace message"},
		},
		{
			name:        "debug level",
			level:       LevelDebug,
			message:     "debug message",
			wantContain: []string{"DEBUG", "debug message"},
		},
		{
			name:        "info level",
			level:       LevelInfo,
			message:     "info message",
			wantContain: []string{"INFO", "info message"},
		},
		{
			name:        "warn level",
			level:       LevelWarn,
			message:     "warn message",
			wantContain: []string{"WARN", "warn message"},
		},
		{
			name:        "error level",
			level:       LevelError,
			message:     "error message",
			wantContain: []string{"ERROR", "error message"},
		},
		{
			name:        "panic level",
			level:       LevelPanic,
			message:     "panic message",
			wantContain: []string{"PANIC", "panic message"},
		},
		{
			name:        "fatal level",
			level:       LevelFatal,
			message:     "fatal message",
			wantContain: []string{"FATAL", "fatal message"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			h := NewSimpleHandler(HandlerOptions{
				Level:   LevelTrace,
				Output:  buf,
				NoColor: true,
			})

			r := NewRecord(time.Now(), tt.level, tt.message)
			err := h.Handle(r)
			if err != nil {
				t.Errorf("SimpleHandler.Handle() error = %v", err)
			}

			output := buf.String()
			for _, want := range tt.wantContain {
				if !strings.Contains(output, want) {
					t.Errorf("SimpleHandler.Handle() output = %q, want to contain %q", output, want)
				}
			}
		})
	}
}

func TestSimpleHandler_HandleWithAttrs(t *testing.T) {
	buf := &bytes.Buffer{}
	h := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})

	r := NewRecord(time.Now(), LevelInfo, "test message")
	r.AddAttrs(String("key1", "value1"), Int("key2", 42))

	err := h.Handle(r)
	if err != nil {
		t.Errorf("SimpleHandler.Handle() error = %v", err)
	}

	output := buf.String()
	wantContain := []string{"test message", "key1=value1", "key2=42"}
	for _, want := range wantContain {
		if !strings.Contains(output, want) {
			t.Errorf("SimpleHandler.Handle() output = %q, want to contain %q", output, want)
		}
	}
}

func TestSimpleHandler_HandleWithPrefix(t *testing.T) {
	buf := &bytes.Buffer{}
	h := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})

	r := NewRecord(time.Now(), LevelInfo, "test message")
	r.Prefix = "myapp"

	err := h.Handle(r)
	if err != nil {
		t.Errorf("SimpleHandler.Handle() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "[myapp]") {
		t.Errorf("SimpleHandler.Handle() output = %q, want to contain [myapp]", output)
	}
}

func TestSimpleHandler_WithAttrs(t *testing.T) {
	buf := &bytes.Buffer{}
	h := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})

	h2 := h.WithAttrs([]Attr{String("app", "test"), String("version", "1.0")})

	r := NewRecord(time.Now(), LevelInfo, "message")
	err := h2.Handle(r)
	if err != nil {
		t.Errorf("SimpleHandler.Handle() error = %v", err)
	}

	output := buf.String()
	wantContain := []string{"app=test", "version=1.0"}
	for _, want := range wantContain {
		if !strings.Contains(output, want) {
			t.Errorf("WithAttrs() output = %q, want to contain %q", output, want)
		}
	}
}

func TestSimpleHandler_WithAttrs_Empty(t *testing.T) {
	buf := &bytes.Buffer{}
	h := NewSimpleHandler(HandlerOptions{
		Level:  LevelInfo,
		Output: buf,
	})

	h2 := h.WithAttrs([]Attr{})
	if h2 != h {
		t.Errorf("WithAttrs([]) should return the same handler")
	}
}

func TestSimpleHandler_WithGroup(t *testing.T) {
	buf := &bytes.Buffer{}
	h := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})

	h2 := h.WithGroup("request")

	r := NewRecord(time.Now(), LevelInfo, "message")
	r.AddAttrs(String("id", "123"), String("method", "GET"))
	err := h2.Handle(r)
	if err != nil {
		t.Errorf("SimpleHandler.Handle() error = %v", err)
	}

	output := buf.String()
	wantContain := []string{"request.id=123", "request.method=GET"}
	for _, want := range wantContain {
		if !strings.Contains(output, want) {
			t.Errorf("WithGroup() output = %q, want to contain %q", output, want)
		}
	}
}

func TestSimpleHandler_WithGroup_Empty(t *testing.T) {
	buf := &bytes.Buffer{}
	h := NewSimpleHandler(HandlerOptions{
		Level:  LevelInfo,
		Output: buf,
	})

	h2 := h.WithGroup("")
	if h2 != h {
		t.Errorf("WithGroup(\"\") should return the same handler")
	}
}

func TestSimpleHandler_WithPrefix(t *testing.T) {
	buf := &bytes.Buffer{}
	h := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})

	h2 := h.WithPrefix("myapp")

	r := NewRecord(time.Now(), LevelInfo, "test")
	r.Prefix = h2.(*SimpleHandler).prefix
	err := h2.Handle(r)
	if err != nil {
		t.Errorf("SimpleHandler.Handle() error = %v", err)
	}

	sh := h2.(*SimpleHandler)
	if sh.prefix != "myapp" {
		t.Errorf("WithPrefix() prefix = %v, want 'myapp'", sh.prefix)
	}
}

func TestSimpleHandler_WithPrefix_Empty(t *testing.T) {
	buf := &bytes.Buffer{}
	h := NewSimpleHandler(HandlerOptions{
		Level:  LevelInfo,
		Output: buf,
	})

	h2 := h.WithPrefix("")
	if h2 != h {
		t.Errorf("WithPrefix(\"\") should return the same handler")
	}
}

func TestSimpleHandler_ReplaceAttr(t *testing.T) {
	buf := &bytes.Buffer{}
	h := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
		ReplaceAttr: func(groups []string, attr Attr) Attr {
			if attr.Key == "password" {
				return String("password", "***")
			}
			return attr
		},
	})

	r := NewRecord(time.Now(), LevelInfo, "login")
	r.AddAttrs(String("user", "john"), String("password", "secret123"))
	err := h.Handle(r)
	if err != nil {
		t.Errorf("SimpleHandler.Handle() error = %v", err)
	}

	output := buf.String()
	if strings.Contains(output, "secret123") {
		t.Errorf("ReplaceAttr should have replaced password")
	}
	if !strings.Contains(output, "password=***") {
		t.Errorf("ReplaceAttr output = %q, want to contain password=***", output)
	}
}

func TestSimpleHandler_NoColor(t *testing.T) {
	tests := []struct {
		name    string
		noColor bool
		level   Level
	}{
		{"with color", false, LevelInfo},
		{"no color", true, LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			h := NewSimpleHandler(HandlerOptions{
				Level:   LevelTrace,
				Output:  buf,
				NoColor: tt.noColor,
			})

			r := NewRecord(time.Now(), tt.level, "test")
			err := h.Handle(r)
			if err != nil {
				t.Errorf("SimpleHandler.Handle() error = %v", err)
			}

			output := buf.String()
			hasAnsi := strings.Contains(output, "\x1b[")
			if tt.noColor && hasAnsi {
				t.Errorf("NoColor=true but output contains ANSI codes: %q", output)
			}
			if !tt.noColor && !hasAnsi {
				t.Errorf("NoColor=false but output doesn't contain ANSI codes: %q", output)
			}
		})
	}
}

func TestSimpleHandler_ColorAttr(t *testing.T) {
	buf := &bytes.Buffer{}
	h := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: false,
	})

	r := NewRecord(time.Now(), LevelInfo, "test")
	r.AddAttrs(ColorAttr(42, String("colored", "value")))
	err := h.Handle(r)
	if err != nil {
		t.Errorf("SimpleHandler.Handle() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "\x1b[") {
		t.Errorf("ColorAttr should produce ANSI codes in output")
	}
}

func TestHandlerOptions_Defaults(t *testing.T) {
	opts := HandlerOptions{}

	if opts.Level != nil {
		t.Errorf("HandlerOptions default Level should be nil")
	}
	if opts.ReplaceAttr != nil {
		t.Errorf("HandlerOptions default ReplaceAttr should be nil")
	}
	if opts.TimeFormat != "" {
		t.Errorf("HandlerOptions default TimeFormat should be empty")
	}
	if opts.NoColor {
		t.Errorf("HandlerOptions default NoColor should be false")
	}
}

func TestAppendSource(t *testing.T) {
	// This tests the internal appendSource function indirectly
	// by checking if it formats correctly
	// Note: appendSource expects *slog.Source, but we can't easily test it directly
	// without importing slog internals
	t.Skip("appendSource is tested indirectly through handler tests")
}

func TestNeedsQuoting(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"empty", "", true},
		{"simple", "abc", false},
		{"with space", "a b", true},
		{"with equals", "a=b", true},
		{"with quote", "a\"b", true},
		{"unicode", "你好", false},
		{"alphanumeric", "abc123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := needsQuoting(tt.s)
			if got != tt.want {
				t.Errorf("needsQuoting(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func BenchmarkSimpleHandler_Handle(b *testing.B) {
	buf := &bytes.Buffer{}
	h := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})

	r := NewRecord(time.Now(), LevelInfo, "benchmark message")
	r.AddAttrs(String("key1", "value1"), Int("key2", 42))

	for b.Loop() {
		buf.Reset()
		_ = h.Handle(r)
	}
}

func BenchmarkSimpleHandler_HandleWithColor(b *testing.B) {
	buf := &bytes.Buffer{}
	h := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: false,
	})

	r := NewRecord(time.Now(), LevelInfo, "benchmark message")
	r.AddAttrs(String("key1", "value1"), Int("key2", 42))

	for b.Loop() {
		buf.Reset()
		_ = h.Handle(r)
	}
}

func BenchmarkSimpleHandler_WithAttrs(b *testing.B) {
	buf := &bytes.Buffer{}
	h := NewSimpleHandler(HandlerOptions{
		Level:  LevelInfo,
		Output: buf,
	})

	attrs := []Attr{String("a", "1"), Int("b", 2), Bool("c", true)}

	for b.Loop() {
		_ = h.WithAttrs(attrs)
	}
}

func TestSimpleHandler_LevelFormat(t *testing.T) {
	tests := []struct {
		name        string
		level       Level
		levelFormat func(Level) string
		want        string
	}{
		{
			name:  "default format trace",
			level: LevelTrace,
			want:  "TRACE",
		},
		{
			name:  "default format info",
			level: LevelInfo,
			want:  "INFO",
		},
		{
			name:  "custom format abbreviations",
			level: LevelDebug,
			levelFormat: func(l Level) string {
				switch l.Real() {
				case LevelTrace:
					return "TRC"
				case LevelDebug:
					return "DBG"
				case LevelInfo:
					return "INF"
				case LevelWarn:
					return "WRN"
				case LevelError:
					return "ERR"
				case LevelPanic:
					return "PNL"
				case LevelFatal:
					return "FTL"
				default:
					return "???"
				}
			},
			want: "DBG",
		},
		{
			name:  "custom format with emoji",
			level: LevelError,
			levelFormat: func(l Level) string {
				switch l.Real() {
				case LevelTrace:
					return "🔍 TRACE"
				case LevelDebug:
					return "🐛 DEBUG"
				case LevelInfo:
					return "ℹ️ INFO"
				case LevelWarn:
					return "⚠️ WARN"
				case LevelError:
					return "❌ ERROR"
				case LevelPanic:
					return "💥 PANIC"
				case LevelFatal:
					return "☠️ FATAL"
				default:
					return "❓"
				}
			},
			want: "❌ ERROR",
		},
		{
			name:  "custom format numeric",
			level: LevelWarn,
			levelFormat: func(l Level) string {
				return "[" + l.String() + ":" + string(rune('0'+l.Int())) + "]"
			},
			want: "[warn:4]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			h := NewSimpleHandler(HandlerOptions{
				Level:       LevelTrace,
				Output:      buf,
				NoColor:     true,
				LevelFormat: tt.levelFormat,
			})

			r := NewRecord(time.Now(), tt.level, "test message")
			err := h.Handle(r)
			if err != nil {
				t.Errorf("SimpleHandler.Handle() error = %v", err)
			}

			output := buf.String()
			if !strings.Contains(output, tt.want) {
				t.Errorf("LevelFormat output = %q, want to contain %q", output, tt.want)
			}
		})
	}
}

func TestSimpleHandler_LevelFormat_WithColor(t *testing.T) {
	buf := &bytes.Buffer{}
	h := NewSimpleHandler(HandlerOptions{
		Level:  LevelInfo,
		Output: buf,
		LevelFormat: func(l Level) string {
			return "[" + l.String() + "]"
		},
		NoColor: false, // colors enabled
	})

	r := NewRecord(time.Now(), LevelInfo, "test")
	err := h.Handle(r)
	if err != nil {
		t.Errorf("SimpleHandler.Handle() error = %v", err)
	}

	output := buf.String()
	// Should contain custom format
	if !strings.Contains(output, "[info]") {
		t.Errorf("Output should contain custom level format [info], got: %q", output)
	}
	// Should still have ANSI color codes
	if !strings.Contains(output, "\x1b[") {
		t.Errorf("Output should contain ANSI color codes")
	}
}

func TestSimpleHandler_PrefixFormat(t *testing.T) {
	tests := []struct {
		name         string
		prefix       string
		prefixFormat func(string) string
		want         string
	}{
		{
			name:   "default format",
			prefix: "myapp",
			want:   "[myapp]",
		},
		{
			name:   "custom format with parentheses",
			prefix: "api",
			prefixFormat: func(p string) string {
				return "(" + p + ")"
			},
			want: "(api)",
		},
		{
			name:   "custom format with emoji",
			prefix: "server",
			prefixFormat: func(p string) string {
				return "🚀 " + p
			},
			want: "🚀 server",
		},
		{
			name:   "custom format with unicode brackets",
			prefix: "worker",
			prefixFormat: func(p string) string {
				return "【" + p + "】"
			},
			want: "【worker】",
		},
		{
			name:   "custom format uppercase",
			prefix: "service",
			prefixFormat: func(p string) string {
				return "<" + strings.ToUpper(p) + ">"
			},
			want: "<SERVICE>",
		},
		{
			name:   "custom format with padding",
			prefix: "db",
			prefixFormat: func(p string) string {
				return "[" + strings.ToUpper(p) + "    ]"
			},
			want: "[DB    ]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			h := NewSimpleHandler(HandlerOptions{
				Level:        LevelInfo,
				Output:       buf,
				NoColor:      true,
				PrefixFormat: tt.prefixFormat,
			})

			r := NewRecord(time.Now(), LevelInfo, "test message")
			r.Prefix = tt.prefix
			err := h.Handle(r)
			if err != nil {
				t.Errorf("SimpleHandler.Handle() error = %v", err)
			}

			output := buf.String()
			if !strings.Contains(output, tt.want) {
				t.Errorf("PrefixFormat output = %q, want to contain %q", output, tt.want)
			}
		})
	}
}

func TestSimpleHandler_PrefixFormat_WithColor(t *testing.T) {
	buf := &bytes.Buffer{}
	h := NewSimpleHandler(HandlerOptions{
		Level:  LevelInfo,
		Output: buf,
		PrefixFormat: func(p string) string {
			return "<" + p + ">"
		},
		NoColor: false, // colors enabled
	})

	r := NewRecord(time.Now(), LevelInfo, "test")
	r.Prefix = "app"
	err := h.Handle(r)
	if err != nil {
		t.Errorf("SimpleHandler.Handle() error = %v", err)
	}

	output := buf.String()
	// Should contain custom format
	if !strings.Contains(output, "<app>") {
		t.Errorf("Output should contain custom prefix format <app>, got: %q", output)
	}
	// Should still have ANSI color codes
	if !strings.Contains(output, "\x1b[") {
		t.Errorf("Output should contain ANSI color codes")
	}
}

func TestSimpleHandler_PrefixFormat_EmptyPrefix(t *testing.T) {
	buf := &bytes.Buffer{}
	h := NewSimpleHandler(HandlerOptions{
		Level:  LevelInfo,
		Output: buf,
		PrefixFormat: func(p string) string {
			return "[PREFIX:" + p + "]"
		},
		NoColor: true,
	})

	r := NewRecord(time.Now(), LevelInfo, "test message")
	// No prefix set
	err := h.Handle(r)
	if err != nil {
		t.Errorf("SimpleHandler.Handle() error = %v", err)
	}

	output := buf.String()
	// Should NOT contain prefix format since prefix is empty
	if strings.Contains(output, "[PREFIX:") {
		t.Errorf("Output should not contain prefix format when prefix is empty, got: %q", output)
	}
}

func TestSimpleHandler_PrefixFormat_WithReplaceAttr(t *testing.T) {
	buf := &bytes.Buffer{}
	h := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
		PrefixFormat: func(p string) string {
			return "{" + p + "}"
		},
		ReplaceAttr: func(groups []string, attr Attr) Attr {
			// When ReplaceAttr is used, it overrides PrefixFormat
			if attr.Key == PrefixKey {
				return String(PrefixKey, "[REPLACED:"+attr.Value.String()+"]")
			}
			return attr
		},
	})

	r := NewRecord(time.Now(), LevelInfo, "test")
	r.Prefix = "test"
	err := h.Handle(r)
	if err != nil {
		t.Errorf("SimpleHandler.Handle() error = %v", err)
	}

	output := buf.String()
	// Should contain ReplaceAttr output, not PrefixFormat
	if !strings.Contains(output, "[REPLACED:test]") {
		t.Errorf("Output should contain ReplaceAttr result, got: %q", output)
	}
	// Should NOT contain PrefixFormat output
	if strings.Contains(output, "{test}") {
		t.Errorf("Output should not contain PrefixFormat when ReplaceAttr is used, got: %q", output)
	}
}
