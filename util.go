package l4g

import (
	"io"
	"sync"
	"sync/atomic"
)

// OutputVar is an atomically updatable io.Writer variable.
// It is safe for concurrent use by multiple goroutines.
// It optimizes for the case where the writer is nil or io.Discard
// by storing a ready flag to avoid unnecessary Write operations.
type OutputVar struct {
	ready  atomic.Bool  // true if writer is not nil and not io.Discard
	writer atomic.Value // holds the io.Writer
}

// NewOutputVar creates a new OutputVar from an io.Writer.
// If the provided writer is already an *OutputVar, it is returned as-is.
// Otherwise, a new OutputVar is created wrapping the writer.
func NewOutputVar(w io.Writer) *OutputVar {
	if o, ok := w.(*OutputVar); ok {
		return o
	}
	v := &OutputVar{}
	v.Set(w)
	return v
}

// Set atomically sets the output writer.
// If w is nil or io.Discard, the OutputVar is marked as disabled for optimization.
func (v *OutputVar) Set(w io.Writer) {
	v.ready.Store(w != nil && w != io.Discard)
	v.writer.Store(w)
}

// Discard reports whether writes to this OutputVar should be discarded.
// It returns true if the writer is nil, io.Discard, or not set.
func (v *OutputVar) Discard() bool {
	return !v.ready.Load()
}

// Output returns the current io.Writer.
// If the writer is nil or marked for discard, it returns io.Discard.
func (v *OutputVar) Output() io.Writer {
	if v.Discard() {
		return io.Discard
	}
	return v.writer.Load().(io.Writer)
}

// Write implements io.Writer by writing to the current output writer.
func (v *OutputVar) Write(p []byte) (int, error) {
	return v.Output().Write(p)
}

// buffer is a byte slice used for building log output.
// It implements efficient Write, WriteByte, and WriteString methods.
type buffer []byte

// bufPool is a sync.Pool for reusing buffer instances to reduce allocations.
// Buffers are initially allocated with 1KB capacity.
var bufPool = sync.Pool{
	New: func() any {
		b := make(buffer, 0, 1024)
		return &b
	},
}

// newBuffer gets a buffer from the pool.
func newBuffer() *buffer {
	return bufPool.Get().(*buffer)
}

// Free returns the buffer to the pool for reuse if it's not too large.
// Buffers larger than 16KB are discarded to avoid keeping large allocations.
func (b *buffer) Free() {
	// To reduce peak allocation, return only
	// smaller buffers to the pool.
	const maxBufferSize = 16 << 10
	if cap(*b) <= maxBufferSize {
		*b = (*b)[:0]
		bufPool.Put(b)
	}
}

// Write appends bytes to the buffer.
func (b *buffer) Write(bytes []byte) {
	*b = append(*b, bytes...)
}

// WriteByte appends a single byte to the buffer.
func (b *buffer) WriteByte(char byte) {
	*b = append(*b, char)
}

// WriteString appends a string to the buffer.
func (b *buffer) WriteString(str string) {
	*b = append(*b, str...)
}
