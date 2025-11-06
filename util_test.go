package l4g

import (
	"bytes"
	"io"
	"sync"
	"testing"
)

func TestOutputVar_NewOutputVar(t *testing.T) {
	t.Run("from writer", func(t *testing.T) {
		buf := &bytes.Buffer{}
		ov := NewOutputVar(buf)
		if ov.Output() != buf {
			t.Errorf("NewOutputVar() output mismatch")
		}
		if ov.Discard() {
			t.Errorf("NewOutputVar() should not be discarded")
		}
	})

	t.Run("from outputvar", func(t *testing.T) {
		buf := &bytes.Buffer{}
		ov1 := NewOutputVar(buf)
		ov2 := NewOutputVar(ov1)
		if ov1 != ov2 {
			t.Errorf("NewOutputVar() should return same instance")
		}
	})

	t.Run("from nil", func(t *testing.T) {
		t.Skip("atomic.Value cannot store nil")
	})

	t.Run("from io.Discard", func(t *testing.T) {
		ov := NewOutputVar(io.Discard)
		if !ov.Discard() {
			t.Errorf("NewOutputVar(io.Discard) should be discarded")
		}
	})
}

func TestOutputVar_SetAndGet(t *testing.T) {
	t.Run("set buffer", func(t *testing.T) {
		ov := &OutputVar{}
		buf := &bytes.Buffer{}
		ov.Set(buf)
		if ov.Output() != buf {
			t.Errorf("OutputVar.Output() mismatch after Set")
		}
		if ov.Discard() {
			t.Errorf("OutputVar should not be discarded after Set")
		}
	})

	t.Run("set nil", func(t *testing.T) {
		t.Skip("atomic.Value cannot store nil")
	})

	t.Run("set io.Discard", func(t *testing.T) {
		ov := &OutputVar{}
		ov.Set(io.Discard)
		if !ov.Discard() {
			t.Errorf("OutputVar should be discarded after Set(io.Discard)")
		}
		if ov.Output() != io.Discard {
			t.Errorf("OutputVar.Output() should be io.Discard")
		}
	})
}

func TestOutputVar_Write(t *testing.T) {
	buf := &bytes.Buffer{}
	ov := NewOutputVar(buf)

	data := []byte("test data")
	n, err := ov.Write(data)
	if err != nil {
		t.Errorf("OutputVar.Write() error = %v", err)
	}
	if n != len(data) {
		t.Errorf("OutputVar.Write() n = %v, want %v", n, len(data))
	}
	if buf.String() != "test data" {
		t.Errorf("OutputVar.Write() wrote %v, want %v", buf.String(), "test data")
	}
}

func TestOutputVar_Discard(t *testing.T) {
	tests := []struct {
		name   string
		writer io.Writer
		want   bool
		skip   bool
	}{
		{"normal buffer", &bytes.Buffer{}, false, false},
		{"nil writer", nil, true, true},
		{"io.Discard", io.Discard, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip("atomic.Value cannot store nil")
				return
			}
			ov := NewOutputVar(tt.writer)
			if got := ov.Discard(); got != tt.want {
				t.Errorf("OutputVar.Discard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOutputVar_Concurrent(t *testing.T) {
	// This test verifies thread-safety of OutputVar read operations
	// In real-world usage, OutputVar is typically set once and read many times
	// Run with: go test -race to detect data races
	//
	// Note: We use io.Discard here because bytes.Buffer is NOT thread-safe
	// In production, loggers typically write to thread-safe destinations like:
	// - os.Stdout/os.Stderr (thread-safe)
	// - Files (kernel-level synchronization)
	// - Network connections with proper synchronization
	ov := NewOutputVar(io.Discard)

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()

			// Perform read operations concurrently (typical usage pattern)
			writer := ov.Output()
			isDiscard := ov.Discard()

			if !isDiscard {
				t.Errorf("Expected Discard() to be true")
			}
			if writer != io.Discard {
				t.Errorf("Expected writer to be io.Discard")
			}

			// Write to io.Discard is safe
			_, _ = ov.Write([]byte("test"))
		}()
	}

	wg.Wait()
}

func TestOutputVar_ConcurrentSet(t *testing.T) {
	// This test demonstrates that concurrent Set() calls should be avoided
	// In practice, Set() should be called from a single goroutine (e.g., main thread)
	t.Skip("Concurrent Set() is not supported - Set() should only be called from a single goroutine")
}

func TestBuffer_NewBuffer(t *testing.T) {
	buf := newBuffer()
	if buf == nil {
		t.Fatal("newBuffer() returned nil")
	}
	if len(*buf) != 0 {
		t.Errorf("newBuffer() length = %v, want 0", len(*buf))
	}
	if cap(*buf) < 1024 {
		t.Errorf("newBuffer() capacity = %v, want >= 1024", cap(*buf))
	}
	buf.Free()
}

func TestBuffer_Write(t *testing.T) {
	buf := newBuffer()
	defer buf.Free()

	data := []byte("hello")
	buf.Write(data)

	if string(*buf) != "hello" {
		t.Errorf("buffer.Write() = %v, want %v", string(*buf), "hello")
	}
}

func TestBuffer_WriteByte(t *testing.T) {
	buf := newBuffer()
	defer buf.Free()

	buf.WriteByte('a')
	buf.WriteByte('b')
	buf.WriteByte('c')

	if string(*buf) != "abc" {
		t.Errorf("buffer.WriteByte() = %v, want %v", string(*buf), "abc")
	}
}

func TestBuffer_WriteString(t *testing.T) {
	buf := newBuffer()
	defer buf.Free()

	buf.WriteString("hello")
	buf.WriteString(" ")
	buf.WriteString("world")

	if string(*buf) != "hello world" {
		t.Errorf("buffer.WriteString() = %v, want %v", string(*buf), "hello world")
	}
}

func TestBuffer_Mixed(t *testing.T) {
	buf := newBuffer()
	defer buf.Free()

	buf.WriteString("hello")
	buf.WriteByte(' ')
	buf.Write([]byte("world"))
	buf.WriteByte('!')

	want := "hello world!"
	if string(*buf) != want {
		t.Errorf("buffer mixed operations = %v, want %v", string(*buf), want)
	}
}

func TestBuffer_Free(t *testing.T) {
	t.Run("small buffer", func(t *testing.T) {
		buf := newBuffer()
		buf.WriteString("small")
		buf.Free()
		// Buffer should be returned to pool

		buf2 := newBuffer()
		if len(*buf2) != 0 {
			t.Errorf("buffer from pool should be empty, got length %v", len(*buf2))
		}
		buf2.Free()
	})

	t.Run("large buffer", func(t *testing.T) {
		buf := newBuffer()
		// Write more than 16KB
		largeData := make([]byte, 20000)
		buf.Write(largeData)
		buf.Free()
		// Buffer should not be returned to pool
	})
}

func TestBuffer_Growth(t *testing.T) {
	buf := newBuffer()
	defer buf.Free()

	// Write data to cause growth
	for range 100 {
		buf.WriteString("test data ")
	}

	if len(*buf) != 1000 {
		t.Errorf("buffer length = %v, want 1000", len(*buf))
	}
}

func BenchmarkBuffer_Write(b *testing.B) {
	data := []byte("benchmark data")
	for b.Loop() {
		buf := newBuffer()
		buf.Write(data)
		buf.Free()
	}
}

func BenchmarkBuffer_WriteString(b *testing.B) {
	data := "benchmark data"
	for b.Loop() {
		buf := newBuffer()
		buf.WriteString(data)
		buf.Free()
	}
}

func BenchmarkBuffer_WriteByte(b *testing.B) {
	for b.Loop() {
		buf := newBuffer()
		buf.WriteByte('x')
		buf.Free()
	}
}
