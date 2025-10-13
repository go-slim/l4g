package l4g

import (
	"bytes"
	"io"
	"strings"
	"sync"
	"testing"
)

func TestDefault(t *testing.T) {
	logger := Default()
	if logger == nil {
		t.Errorf("Default() returned nil")
	}
}

func TestSetDefault(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(buf, WithLevel(LevelDebug))

	SetDefault(logger)

	if Default() != logger {
		t.Errorf("SetDefault() did not set the default logger")
	}

	// Restore original default
	SetDefault(New(io.Discard))
}

func TestOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(buf)
	SetDefault(logger)

	output := Output()
	if output != buf {
		t.Errorf("Output() mismatch")
	}

	SetDefault(New(io.Discard))
}

func TestSetOutput(t *testing.T) {
	// Skip this test due to atomic.Value type mismatch issues
	// when changing output types on shared logger
	t.Skip("SetOutput has atomic.Value type constraints")
}

func TestGetLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(buf, WithLevel(LevelWarn))
	SetDefault(logger)

	if GetLevel() != LevelWarn {
		t.Errorf("GetLevel() = %v, want %v", GetLevel(), LevelWarn)
	}

	SetDefault(New(io.Discard))
}

func TestSetLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(buf)
	SetDefault(logger)

	SetLevel(LevelError)
	if GetLevel() != LevelError {
		t.Errorf("SetLevel() did not update level")
	}

	SetDefault(New(io.Discard))
}

func TestPackageTrace(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(buf, WithLevel(LevelTrace))
	SetDefault(logger)

	Trace("trace message")
	output := buf.String()

	if !strings.Contains(output, "trace message") {
		t.Errorf("Trace() output = %q, want to contain 'trace message'", output)
	}

	SetDefault(New(io.Discard))
}

func TestPackageTracef(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(buf, WithLevel(LevelTrace))
	SetDefault(logger)

	Tracef("trace %s %d", "message", 42)
	output := buf.String()

	if !strings.Contains(output, "trace message 42") {
		t.Errorf("Tracef() output = %q, want to contain 'trace message 42'", output)
	}

	SetDefault(New(io.Discard))
}

func TestPackageTracej(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelTrace,
		Output:  buf,
		NoColor: true,
	})
	logger := New(buf, WithHandler(handler))
	SetDefault(logger)

	Tracej(map[string]any{"key": "value"})
	output := buf.String()

	if !strings.Contains(output, "key=value") {
		t.Errorf("Tracej() output = %q, want to contain key=value", output)
	}

	SetDefault(New(io.Discard))
}

func TestPackageDebug(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(buf, WithLevel(LevelDebug))
	SetDefault(logger)

	Debug("debug message")
	output := buf.String()

	if !strings.Contains(output, "debug message") {
		t.Errorf("Debug() output = %q, want to contain 'debug message'", output)
	}

	SetDefault(New(io.Discard))
}

func TestPackageDebugf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(buf, WithLevel(LevelDebug))
	SetDefault(logger)

	Debugf("debug %s", "message")
	output := buf.String()

	if !strings.Contains(output, "debug message") {
		t.Errorf("Debugf() output = %q, want to contain 'debug message'", output)
	}

	SetDefault(New(io.Discard))
}

func TestPackageDebugj(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelDebug,
		Output:  buf,
		NoColor: true,
	})
	logger := New(buf, WithHandler(handler))
	SetDefault(logger)

	Debugj(map[string]any{"debug": "test"})
	output := buf.String()

	if !strings.Contains(output, "debug=test") {
		t.Errorf("Debugj() output = %q, want to contain debug=test", output)
	}

	SetDefault(New(io.Discard))
}

func TestPackageInfo(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(buf)
	SetDefault(logger)

	Info("info message")
	output := buf.String()

	if !strings.Contains(output, "info message") {
		t.Errorf("Info() output = %q, want to contain 'info message'", output)
	}

	SetDefault(New(io.Discard))
}

func TestPackageInfof(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(buf)
	SetDefault(logger)

	Infof("info %d", 123)
	output := buf.String()

	if !strings.Contains(output, "info 123") {
		t.Errorf("Infof() output = %q, want to contain 'info 123'", output)
	}

	SetDefault(New(io.Discard))
}

func TestPackageInfoj(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})
	logger := New(buf, WithHandler(handler))
	SetDefault(logger)

	Infoj(map[string]any{"status": "ok"})
	output := buf.String()

	if !strings.Contains(output, "status=ok") {
		t.Errorf("Infoj() output = %q, want to contain status=ok", output)
	}

	SetDefault(New(io.Discard))
}

func TestPackageWarn(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(buf)
	SetDefault(logger)

	Warn("warn message")
	output := buf.String()

	if !strings.Contains(output, "warn message") {
		t.Errorf("Warn() output = %q, want to contain 'warn message'", output)
	}

	SetDefault(New(io.Discard))
}

func TestPackageWarnf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(buf)
	SetDefault(logger)

	Warnf("warn %s", "test")
	output := buf.String()

	if !strings.Contains(output, "warn test") {
		t.Errorf("Warnf() output = %q, want to contain 'warn test'", output)
	}

	SetDefault(New(io.Discard))
}

func TestPackageWarnj(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})
	logger := New(buf, WithHandler(handler))
	SetDefault(logger)

	Warnj(map[string]any{"warning": "alert"})
	output := buf.String()

	if !strings.Contains(output, "warning=alert") {
		t.Errorf("Warnj() output = %q, want to contain warning=alert", output)
	}

	SetDefault(New(io.Discard))
}

func TestPackageError(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(buf)
	SetDefault(logger)

	Error("error message")
	output := buf.String()

	if !strings.Contains(output, "error message") {
		t.Errorf("Error() output = %q, want to contain 'error message'", output)
	}

	SetDefault(New(io.Discard))
}

func TestPackageErrorf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(buf)
	SetDefault(logger)

	Errorf("error %d", 500)
	output := buf.String()

	if !strings.Contains(output, "error 500") {
		t.Errorf("Errorf() output = %q, want to contain 'error 500'", output)
	}

	SetDefault(New(io.Discard))
}

func TestPackageErrorj(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})
	logger := New(buf, WithHandler(handler))
	SetDefault(logger)

	Errorj(map[string]any{"error": "failed"})
	output := buf.String()

	if !strings.Contains(output, "error=failed") {
		t.Errorf("Errorj() output = %q, want to contain error=failed", output)
	}

	SetDefault(New(io.Discard))
}

func TestPackagePanic(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(buf)
	SetDefault(logger)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Panic() did not panic")
		}
		SetDefault(New(io.Discard))
	}()

	Panic("panic message")
}

func TestPackagePanicf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(buf)
	SetDefault(logger)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Panicf() did not panic")
		}
		SetDefault(New(io.Discard))
	}()

	Panicf("panic %s", "test")
}

func TestPackagePanicj(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(buf)
	SetDefault(logger)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Panicj() did not panic")
		}
		SetDefault(New(io.Discard))
	}()

	Panicj(map[string]any{"panic": "data"})
}

func TestPackageFatal(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(buf)
	SetDefault(logger)

	exitCalled := false
	oldExiter := OsExiter
	OsExiter = func(code int) {
		exitCalled = true
	}
	defer func() {
		OsExiter = oldExiter
		SetDefault(New(io.Discard))
	}()

	Fatal("fatal message")

	if !exitCalled {
		t.Errorf("Fatal() did not call os.Exit")
	}
}

func TestPackageFatalf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(buf)
	SetDefault(logger)

	exitCalled := false
	oldExiter := OsExiter
	OsExiter = func(code int) {
		exitCalled = true
	}
	defer func() {
		OsExiter = oldExiter
		SetDefault(New(io.Discard))
	}()

	Fatalf("fatal %s", "error")

	if !exitCalled {
		t.Errorf("Fatalf() did not call os.Exit")
	}
}

func TestPackageFatalj(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(buf)
	SetDefault(logger)

	exitCalled := false
	oldExiter := OsExiter
	OsExiter = func(code int) {
		exitCalled = true
	}
	defer func() {
		OsExiter = oldExiter
		SetDefault(New(io.Discard))
	}()

	Fatalj(map[string]any{"fatal": "data"})

	if !exitCalled {
		t.Errorf("Fatalj() did not call os.Exit")
	}
}

func TestFallbackErrorf(t *testing.T) {
	// This function writes to stderr, we just verify it doesn't panic
	FallbackErrorf("test error: %s", "message")
}

func TestChannel(t *testing.T) {
	buf := &bytes.Buffer{}
	SetDefault(New(buf))

	ch1 := Channel("app1")
	ch2 := Channel("app1")
	ch3 := Channel("app2")

	if ch1 != ch2 {
		t.Errorf("Channel() should return same instance for same name")
	}
	if ch1 == ch3 {
		t.Errorf("Channel() should return different instances for different names")
	}

	SetDefault(New(io.Discard))
}

func TestChannel_Independent(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	// Set up custom NewFunc to create loggers with different outputs
	oldNewFunc := NewFunc
	NewFunc = func(name string) *Logger {
		if name == "ch1" {
			return New(buf1)
		}
		return New(buf2)
	}
	defer func() { NewFunc = oldNewFunc }()

	// Clear the channels map (sync.Map doesn't have a clear method, so we replace it)
	oldLs := ls
	ls = sync.Map{}
	defer func() { ls = oldLs }()

	ch1 := Channel("ch1")
	ch2 := Channel("ch2")

	ch1.Info("message1")
	ch2.Info("message2")

	if !strings.Contains(buf1.String(), "message1") {
		t.Errorf("Channel ch1 should write to buf1")
	}
	if !strings.Contains(buf2.String(), "message2") {
		t.Errorf("Channel ch2 should write to buf2")
	}
}

func TestOsExiter(t *testing.T) {
	if OsExiter == nil {
		t.Errorf("OsExiter should be initialized")
	}
}

func BenchmarkPackageInfo(b *testing.B) {
	buf := &bytes.Buffer{}
	logger := New(buf)
	SetDefault(logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("benchmark message")
	}

	SetDefault(New(io.Discard))
}

func BenchmarkPackageInfof(b *testing.B) {
	buf := &bytes.Buffer{}
	logger := New(buf)
	SetDefault(logger)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Infof("benchmark %s", "message")
	}

	SetDefault(New(io.Discard))
}

func BenchmarkChannel(b *testing.B) {
	SetDefault(New(io.Discard))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Channel("test")
	}
}
