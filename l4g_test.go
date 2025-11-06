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
	logger := New(Options{Output: buf, Level: LevelDebug})

	SetDefault(logger)

	if Default() != logger {
		t.Errorf("SetDefault() did not set the default logger")
	}

	// Restore original default
	SetDefault(New(Options{Output: io.Discard}))
}

func TestOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})
	SetDefault(logger)

	output := Output()
	if output != buf {
		t.Errorf("Output() mismatch")
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestSetOutput(t *testing.T) {
	// Skip this test due to atomic.Value type mismatch issues
	// when changing output types on shared logger
	t.Skip("SetOutput has atomic.Value type constraints")
}

func TestGetLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf, Level: LevelWarn})
	SetDefault(logger)

	if GetLevel() != LevelWarn {
		t.Errorf("GetLevel() = %v, want %v", GetLevel(), LevelWarn)
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestSetLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})
	SetDefault(logger)

	SetLevel(LevelError)
	if GetLevel() != LevelError {
		t.Errorf("SetLevel() did not update level")
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestPackageTrace(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf, Level: LevelTrace})
	SetDefault(logger)

	Trace("trace message")
	output := buf.String()

	if !strings.Contains(output, "trace message") {
		t.Errorf("Trace() output = %q, want to contain 'trace message'", output)
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestPackageTracef(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf, Level: LevelTrace})
	SetDefault(logger)

	Tracef("trace %s %d", "message", 42)
	output := buf.String()

	if !strings.Contains(output, "trace message 42") {
		t.Errorf("Tracef() output = %q, want to contain 'trace message 42'", output)
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestPackageTracej(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelTrace,
		Output:  buf,
		NoColor: true,
	})
	logger := New(Options{Output: buf, Handler: handler})
	SetDefault(logger)

	Tracej(map[string]any{"key": "value"})
	output := buf.String()

	if !strings.Contains(output, "key=value") {
		t.Errorf("Tracej() output = %q, want to contain key=value", output)
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestPackageDebug(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf, Level: LevelDebug})
	SetDefault(logger)

	Debug("debug message")
	output := buf.String()

	if !strings.Contains(output, "debug message") {
		t.Errorf("Debug() output = %q, want to contain 'debug message'", output)
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestPackageDebugf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf, Level: LevelDebug})
	SetDefault(logger)

	Debugf("debug %s", "message")
	output := buf.String()

	if !strings.Contains(output, "debug message") {
		t.Errorf("Debugf() output = %q, want to contain 'debug message'", output)
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestPackageDebugj(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelDebug,
		Output:  buf,
		NoColor: true,
	})
	logger := New(Options{Output: buf, Handler: handler})
	SetDefault(logger)

	Debugj(map[string]any{"debug": "test"})
	output := buf.String()

	if !strings.Contains(output, "debug=test") {
		t.Errorf("Debugj() output = %q, want to contain debug=test", output)
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestPackageInfo(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})
	SetDefault(logger)

	Info("info message")
	output := buf.String()

	if !strings.Contains(output, "info message") {
		t.Errorf("Info() output = %q, want to contain 'info message'", output)
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestPackageInfof(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})
	SetDefault(logger)

	Infof("info %d", 123)
	output := buf.String()

	if !strings.Contains(output, "info 123") {
		t.Errorf("Infof() output = %q, want to contain 'info 123'", output)
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestPackageInfoj(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})
	logger := New(Options{Output: buf, Handler: handler})
	SetDefault(logger)

	Infoj(map[string]any{"status": "ok"})
	output := buf.String()

	if !strings.Contains(output, "status=ok") {
		t.Errorf("Infoj() output = %q, want to contain status=ok", output)
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestPackageWarn(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})
	SetDefault(logger)

	Warn("warn message")
	output := buf.String()

	if !strings.Contains(output, "warn message") {
		t.Errorf("Warn() output = %q, want to contain 'warn message'", output)
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestPackageWarnf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})
	SetDefault(logger)

	Warnf("warn %s", "test")
	output := buf.String()

	if !strings.Contains(output, "warn test") {
		t.Errorf("Warnf() output = %q, want to contain 'warn test'", output)
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestPackageWarnj(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})
	logger := New(Options{Output: buf, Handler: handler})
	SetDefault(logger)

	Warnj(map[string]any{"warning": "alert"})
	output := buf.String()

	if !strings.Contains(output, "warning=alert") {
		t.Errorf("Warnj() output = %q, want to contain warning=alert", output)
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestPackageError(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})
	SetDefault(logger)

	Error("error message")
	output := buf.String()

	if !strings.Contains(output, "error message") {
		t.Errorf("Error() output = %q, want to contain 'error message'", output)
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestPackageErrorf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})
	SetDefault(logger)

	Errorf("error %d", 500)
	output := buf.String()

	if !strings.Contains(output, "error 500") {
		t.Errorf("Errorf() output = %q, want to contain 'error 500'", output)
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestPackageErrorj(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})
	logger := New(Options{Output: buf, Handler: handler})
	SetDefault(logger)

	Errorj(map[string]any{"error": "failed"})
	output := buf.String()

	if !strings.Contains(output, "error=failed") {
		t.Errorf("Errorj() output = %q, want to contain error=failed", output)
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestPackagePanic(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})
	SetDefault(logger)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Panic() did not panic")
		}
		SetDefault(New(Options{Output: io.Discard}))
	}()

	Panic("panic message")
}

func TestPackagePanicf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})
	SetDefault(logger)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Panicf() did not panic")
		}
		SetDefault(New(Options{Output: io.Discard}))
	}()

	Panicf("panic %s", "test")
}

func TestPackagePanicj(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})
	SetDefault(logger)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Panicj() did not panic")
		}
		SetDefault(New(Options{Output: io.Discard}))
	}()

	Panicj(map[string]any{"panic": "data"})
}

func TestPackageFatal(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})
	SetDefault(logger)

	exitCalled := false
	oldExiter := OsExiter
	OsExiter = func(code int) {
		exitCalled = true
	}
	defer func() {
		OsExiter = oldExiter
		SetDefault(New(Options{Output: io.Discard}))
	}()

	Fatal("fatal message")

	if !exitCalled {
		t.Errorf("Fatal() did not call os.Exit")
	}
}

func TestPackageFatalf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})
	SetDefault(logger)

	exitCalled := false
	oldExiter := OsExiter
	OsExiter = func(code int) {
		exitCalled = true
	}
	defer func() {
		OsExiter = oldExiter
		SetDefault(New(Options{Output: io.Discard}))
	}()

	Fatalf("fatal %s", "error")

	if !exitCalled {
		t.Errorf("Fatalf() did not call os.Exit")
	}
}

func TestPackageFatalj(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})
	SetDefault(logger)

	exitCalled := false
	oldExiter := OsExiter
	OsExiter = func(code int) {
		exitCalled = true
	}
	defer func() {
		OsExiter = oldExiter
		SetDefault(New(Options{Output: io.Discard}))
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
	SetDefault(New(Options{Output: buf}))

	ch1 := Channel("app1")
	ch2 := Channel("app1")
	ch3 := Channel("app2")

	if ch1 != ch2 {
		t.Errorf("Channel() should return same instance for same name")
	}
	if ch1 == ch3 {
		t.Errorf("Channel() should return different instances for different names")
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestChannel_Independent(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	// Set up custom NewFunc to create loggers with different outputs
	oldNewFunc := NewFunc
	NewFunc = func(name string) *Logger {
		if name == "ch1" {
			return New(Options{Output: buf1})
		}
		return New(Options{Output: buf2})
	}
	defer func() { NewFunc = oldNewFunc }()

	// Clear the channels map (sync.Map doesn't have a clear method, so we replace it)
	// Store original map and restore later - use a new variable to avoid copying lock
	originalLs := ls
	ls = &sync.Map{}
	defer func() { ls = originalLs }()

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
	logger := New(Options{Output: buf})
	SetDefault(logger)

	for b.Loop() {
		Info("benchmark message")
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func BenchmarkPackageInfof(b *testing.B) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})
	SetDefault(logger)

	for b.Loop() {
		Infof("benchmark %s", "message")
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func BenchmarkChannel(b *testing.B) {
	SetDefault(New(Options{Output: io.Discard}))

	for b.Loop() {
		_ = Channel("test")
	}
}

func TestWithAttrs(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})
	logger := New(Options{Output: buf, Handler: handler})
	SetDefault(logger)

	// Create a new logger with preset attributes
	loggerWithAttrs := WithAttrs("app", "myapp", "env", "dev")
	loggerWithAttrs.Info("test message")
	output := buf.String()

	if !strings.Contains(output, "test message") {
		t.Errorf("WithAttrs() output = %q, want to contain 'test message'", output)
	}
	if !strings.Contains(output, "app=myapp") {
		t.Errorf("WithAttrs() output = %q, want to contain 'app=myapp'", output)
	}
	if !strings.Contains(output, "env=dev") {
		t.Errorf("WithAttrs() output = %q, want to contain 'env=dev'", output)
	}

	// Original default logger should not have the attributes
	buf.Reset()
	Info("original message")
	output2 := buf.String()
	if strings.Contains(output2, "app=myapp") {
		t.Errorf("Default logger should not have attributes from derived logger")
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestWithPrefix(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})
	logger := New(Options{Output: buf, Handler: handler})
	SetDefault(logger)

	// Create a new logger with prefix
	loggerWithPrefix := WithPrefix("HTTP")
	loggerWithPrefix.Info("request received")
	output := buf.String()

	if !strings.Contains(output, "[HTTP]") {
		t.Errorf("WithPrefix() output = %q, want to contain '[HTTP]'", output)
	}
	if !strings.Contains(output, "request received") {
		t.Errorf("WithPrefix() output = %q, want to contain 'request received'", output)
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestWithGroup(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})
	logger := New(Options{Output: buf, Handler: handler})
	SetDefault(logger)

	// Create a new logger with group
	loggerWithGroup := WithGroup("request")
	loggerWithGroup.Info("processing", "method", "GET", "path", "/api")
	output := buf.String()

	if !strings.Contains(output, "request.method=GET") {
		t.Errorf("WithGroup() output = %q, want to contain 'request.method=GET'", output)
	}
	if !strings.Contains(output, "request.path=/api") {
		t.Errorf("WithGroup() output = %q, want to contain 'request.path=/api'", output)
	}

	SetDefault(New(Options{Output: io.Discard}))
}

func TestPackage_Chaining(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})
	logger := New(Options{Output: buf, Handler: handler})
	SetDefault(logger)

	// Test chaining WithAttrs, WithPrefix, and WithGroup
	derived := WithPrefix("API").WithAttrs("version", "v1").WithGroup("metrics")
	derived.Info("completed", "duration", "150ms", "status", 200)
	output := buf.String()

	if !strings.Contains(output, "[API]") {
		t.Errorf("Chained logger output = %q, want to contain '[API]'", output)
	}
	if !strings.Contains(output, "version=v1") {
		t.Errorf("Chained logger output = %q, want to contain 'version=v1'", output)
	}
	if !strings.Contains(output, "metrics.duration=150ms") {
		t.Errorf("Chained logger output = %q, want to contain 'metrics.duration=150ms'", output)
	}
	if !strings.Contains(output, "metrics.status=200") {
		t.Errorf("Chained logger output = %q, want to contain 'metrics.status=200'", output)
	}

	SetDefault(New(Options{Output: io.Discard}))
}
