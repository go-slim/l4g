package l4g

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})

	if logger == nil {
		t.Errorf("New() returned nil")
	}
	if logger.Level() != LevelInfo {
		t.Errorf("New() default level = %v, want %v", logger.Level(), LevelInfo)
	}
}

func TestNew_WithOptions(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{
		Output: buf,
		Level:  LevelDebug,
	})

	if logger.Level() != LevelDebug {
		t.Errorf("New() with level option = %v, want %v", logger.Level(), LevelDebug)
	}
}

func TestLogger_SetLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})

	tests := []Level{LevelTrace, LevelDebug, LevelInfo, LevelWarn, LevelError, LevelPanic, LevelFatal}
	for _, level := range tests {
		t.Run(level.String(), func(t *testing.T) {
			logger.SetLevel(level)
			if logger.Level() != level {
				t.Errorf("Logger.SetLevel() level = %v, want %v", logger.Level(), level)
			}
		})
	}
}

func TestLogger_SetOutput(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	logger := New(Options{Output: buf1})
	logger.Info("test1")

	logger.SetOutput(buf2)
	logger.Info("test2")

	if strings.Contains(buf1.String(), "test2") {
		t.Errorf("Logger should not write to old output after SetOutput")
	}
	if !strings.Contains(buf2.String(), "test2") {
		t.Errorf("Logger should write to new output after SetOutput")
	}
}

func TestLogger_Output(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})

	if logger.Output() != buf {
		t.Errorf("Logger.Output() mismatch")
	}
}

func TestLogger_Enabled(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf, Level: LevelWarn})

	tests := []struct {
		level   Level
		enabled bool
	}{
		{LevelTrace, false},
		{LevelDebug, false},
		{LevelInfo, false},
		{LevelWarn, true},
		{LevelError, true},
		{LevelPanic, true},
		{LevelFatal, true},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			if got := logger.Enabled(tt.level); got != tt.enabled {
				t.Errorf("Logger.Enabled(%v) = %v, want %v", tt.level, got, tt.enabled)
			}
		})
	}
}

func TestLogger_Trace(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf, Level: LevelTrace})

	logger.Trace("trace message")
	output := buf.String()

	if !strings.Contains(output, "trace message") {
		t.Errorf("Logger.Trace() output = %q, want to contain 'trace message'", output)
	}
}

func TestLogger_Tracef(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf, Level: LevelTrace})

	logger.Tracef("trace %s %d", "message", 42)
	output := buf.String()

	if !strings.Contains(output, "trace message 42") {
		t.Errorf("Logger.Tracef() output = %q, want to contain 'trace message 42'", output)
	}
}

func TestLogger_Tracej(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelTrace,
		Output:  buf,
		NoColor: true,
	})
	logger := New(Options{Output: buf, Handler: handler})

	logger.Tracej(map[string]any{"key": "value", "count": 42})
	output := buf.String()

	if !strings.Contains(output, "key=value") || !strings.Contains(output, "count=42") {
		t.Errorf("Logger.Tracej() output = %q, want to contain key=value and count=42", output)
	}
}

func TestLogger_Debug(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf, Level: LevelDebug})

	logger.Debug("debug message")
	output := buf.String()

	if !strings.Contains(output, "debug message") {
		t.Errorf("Logger.Debug() output = %q, want to contain 'debug message'", output)
	}
}

func TestLogger_Debugf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf, Level: LevelDebug})

	logger.Debugf("debug %s %d", "message", 42)
	output := buf.String()

	if !strings.Contains(output, "debug message 42") {
		t.Errorf("Logger.Debugf() output = %q, want to contain 'debug message 42'", output)
	}
}

func TestLogger_Debugj(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelDebug,
		Output:  buf,
		NoColor: true,
	})
	logger := New(Options{Output: buf, Handler: handler})

	logger.Debugj(map[string]any{"key": "value"})
	output := buf.String()

	if !strings.Contains(output, "key=value") {
		t.Errorf("Logger.Debugj() output = %q, want to contain key=value", output)
	}
}

func TestLogger_Info(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})

	logger.Info("info message")
	output := buf.String()

	if !strings.Contains(output, "info message") {
		t.Errorf("Logger.Info() output = %q, want to contain 'info message'", output)
	}
}

func TestLogger_Infof(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})

	logger.Infof("info %s %d", "message", 42)
	output := buf.String()

	if !strings.Contains(output, "info message 42") {
		t.Errorf("Logger.Infof() output = %q, want to contain 'info message 42'", output)
	}
}

func TestLogger_Infoj(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})
	logger := New(Options{Output: buf, Handler: handler})

	logger.Infoj(map[string]any{"status": "ok", "code": 200})
	output := buf.String()

	if !strings.Contains(output, "status=ok") || !strings.Contains(output, "code=200") {
		t.Errorf("Logger.Infoj() output = %q, want to contain status and code", output)
	}
}

func TestLogger_Warn(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})

	logger.Warn("warn message")
	output := buf.String()

	if !strings.Contains(output, "warn message") {
		t.Errorf("Logger.Warn() output = %q, want to contain 'warn message'", output)
	}
}

func TestLogger_Warnf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})

	logger.Warnf("warn %s %d", "message", 42)
	output := buf.String()

	if !strings.Contains(output, "warn message 42") {
		t.Errorf("Logger.Warnf() output = %q, want to contain 'warn message 42'", output)
	}
}

func TestLogger_Warnj(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})
	logger := New(Options{Output: buf, Handler: handler})

	logger.Warnj(map[string]any{"warning": "test"})
	output := buf.String()

	if !strings.Contains(output, "warning=test") {
		t.Errorf("Logger.Warnj() output = %q, want to contain warning=test", output)
	}
}

func TestLogger_Error(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})

	logger.Error("error message")
	output := buf.String()

	if !strings.Contains(output, "error message") {
		t.Errorf("Logger.Error() output = %q, want to contain 'error message'", output)
	}
}

func TestLogger_Errorf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})

	logger.Errorf("error %s %d", "message", 42)
	output := buf.String()

	if !strings.Contains(output, "error message 42") {
		t.Errorf("Logger.Errorf() output = %q, want to contain 'error message 42'", output)
	}
}

func TestLogger_Errorj(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})
	logger := New(Options{Output: buf, Handler: handler})

	logger.Errorj(map[string]any{"error": "failed", "code": 500})
	output := buf.String()

	if !strings.Contains(output, "error=failed") || !strings.Contains(output, "code=500") {
		t.Errorf("Logger.Errorj() output = %q, want to contain error and code", output)
	}
}

func TestLogger_Panic(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Logger.Panic() did not panic")
		}
	}()

	logger.Panic("panic message")
}

func TestLogger_Panicf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Logger.Panicf() did not panic")
		}
	}()

	logger.Panicf("panic %s", "message")
}

func TestLogger_Panicj(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Logger.Panicj() did not panic")
		}
	}()

	logger.Panicj(map[string]any{"panic": "test"})
}

func TestLogger_Fatal(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})

	exitCalled := false
	exitCode := 0

	// Replace OsExiter with a test function
	oldExiter := OsExiter
	OsExiter = func(code int) {
		exitCalled = true
		exitCode = code
	}
	defer func() { OsExiter = oldExiter }()

	logger.Fatal("fatal message")

	if !exitCalled {
		t.Errorf("Logger.Fatal() did not call os.Exit")
	}
	if exitCode != 1 {
		t.Errorf("Logger.Fatal() exit code = %v, want 1", exitCode)
	}
}

func TestLogger_Fatalf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})

	exitCalled := false

	oldExiter := OsExiter
	OsExiter = func(code int) {
		exitCalled = true
	}
	defer func() { OsExiter = oldExiter }()

	logger.Fatalf("fatal %s", "message")

	if !exitCalled {
		t.Errorf("Logger.Fatalf() did not call os.Exit")
	}
}

func TestLogger_Fatalj(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})

	exitCalled := false

	oldExiter := OsExiter
	OsExiter = func(code int) {
		exitCalled = true
	}
	defer func() { OsExiter = oldExiter }()

	logger.Fatalj(map[string]any{"fatal": "error"})

	if !exitCalled {
		t.Errorf("Logger.Fatalj() did not call os.Exit")
	}
}

func TestLogger_Log(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})

	logger.Log(LevelWarn, "log message")
	output := buf.String()

	if !strings.Contains(output, "log message") {
		t.Errorf("Logger.Log() output = %q, want to contain 'log message'", output)
	}
}

func TestLogger_Logf(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})

	logger.Logf(LevelError, "log %s %d", "message", 42)
	output := buf.String()

	if !strings.Contains(output, "log message 42") {
		t.Errorf("Logger.Logf() output = %q, want to contain 'log message 42'", output)
	}
}

func TestLogger_Logj(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})
	logger := New(Options{Output: buf, Handler: handler})

	logger.Logj(LevelInfo, map[string]any{"test": "data"})
	output := buf.String()

	if !strings.Contains(output, "test=data") {
		t.Errorf("Logger.Logj() output = %q, want to contain test=data", output)
	}
}

func TestLogger_WithAttrs(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})
	logger := New(Options{Output: buf, Handler: handler})

	// Create a new logger with preset attributes
	loggerWithAttrs := logger.WithAttrs("service", "api", "version", "v1.0")
	loggerWithAttrs.Info("test message")
	output := buf.String()

	if !strings.Contains(output, "test message") {
		t.Errorf("Logger.WithAttrs() output = %q, want to contain 'test message'", output)
	}
	if !strings.Contains(output, "service=api") {
		t.Errorf("Logger.WithAttrs() output = %q, want to contain 'service=api'", output)
	}
	if !strings.Contains(output, "version=v1.0") {
		t.Errorf("Logger.WithAttrs() output = %q, want to contain 'version=v1.0'", output)
	}

	// Original logger should not have the attributes
	buf.Reset()
	logger.Info("original message")
	output2 := buf.String()
	if strings.Contains(output2, "service=api") {
		t.Errorf("Original logger should not have attributes from derived logger")
	}
}

func TestLogger_WithPrefix(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})
	logger := New(Options{Output: buf, Handler: handler})

	// Create a new logger with prefix
	loggerWithPrefix := logger.WithPrefix("HTTP")
	loggerWithPrefix.Info("request received")
	output := buf.String()

	if !strings.Contains(output, "[HTTP]") {
		t.Errorf("Logger.WithPrefix() output = %q, want to contain '[HTTP]'", output)
	}
	if !strings.Contains(output, "request received") {
		t.Errorf("Logger.WithPrefix() output = %q, want to contain 'request received'", output)
	}

	// Test empty prefix returns same logger
	samLogger := logger.WithPrefix("")
	if samLogger != logger {
		t.Errorf("Logger.WithPrefix(\"\") should return the same logger")
	}
}

func TestLogger_WithGroup(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})
	logger := New(Options{Output: buf, Handler: handler})

	// Create a new logger with group
	loggerWithGroup := logger.WithGroup("request")
	loggerWithGroup.Info("processing", "method", "GET", "path", "/api")
	output := buf.String()

	if !strings.Contains(output, "request.method=GET") {
		t.Errorf("Logger.WithGroup() output = %q, want to contain 'request.method=GET'", output)
	}
	if !strings.Contains(output, "request.path=/api") {
		t.Errorf("Logger.WithGroup() output = %q, want to contain 'request.path=/api'", output)
	}

	// Test empty group returns same logger
	sameLogger := logger.WithGroup("")
	if sameLogger != logger {
		t.Errorf("Logger.WithGroup(\"\") should return the same logger")
	}
}

func TestLogger_WithAttrs_Chaining(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := NewSimpleHandler(HandlerOptions{
		Level:   LevelInfo,
		Output:  buf,
		NoColor: true,
	})
	logger := New(Options{Output: buf, Handler: handler})

	// Test chaining WithAttrs, WithPrefix, and WithGroup
	derived := logger.WithPrefix("API").WithAttrs("version", "v1").WithGroup("metrics")
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
}

func TestLogger_DiscardOutput(t *testing.T) {
	logger := New(Options{Output: io.Discard})

	// Should not panic and should not write anything
	logger.Info("test message")
	logger.Error("error message")
}

func TestLogger_LevelFiltering(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf, Level: LevelWarn})

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")

	output := buf.String()

	if strings.Contains(output, "debug message") {
		t.Errorf("Logger should not log below configured level")
	}
	if strings.Contains(output, "info message") {
		t.Errorf("Logger should not log below configured level")
	}
	if !strings.Contains(output, "warn message") {
		t.Errorf("Logger should log at configured level")
	}
}

func TestWithLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf, Level: LevelError})

	if logger.Level() != LevelError {
		t.Errorf("WithLevel() level = %v, want %v", logger.Level(), LevelError)
	}
}

func TestWithHandler(t *testing.T) {
	buf := &bytes.Buffer{}
	customHandler := NewSimpleHandler(HandlerOptions{
		Level:   LevelDebug,
		Output:  buf,
		NoColor: true,
	})

	logger := New(Options{Output: buf, Handler: customHandler})

	logger.Info("test")
	if buf.Len() == 0 {
		t.Errorf("WithHandler() should use custom handler")
	}
}

func BenchmarkLogger_Info(b *testing.B) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})

	for b.Loop() {
		logger.Info("benchmark message")
	}
}

func BenchmarkLogger_InfoWithAttrs(b *testing.B) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})

	for b.Loop() {
		logger.Info("benchmark message", String("key1", "value1"), Int("key2", 42))
	}
}

func BenchmarkLogger_Infof(b *testing.B) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf})

	for b.Loop() {
		logger.Infof("benchmark %s %d", "message", 42)
	}
}

func BenchmarkLogger_Disabled(b *testing.B) {
	buf := &bytes.Buffer{}
	logger := New(Options{Output: buf, Level: LevelError})

	for b.Loop() {
		logger.Debug("this should be skipped")
	}
}
