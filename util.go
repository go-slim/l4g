package l4g

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/MatusOllah/stripansi"
)

type OutputVar struct {
	ready  atomic.Bool
	writer atomic.Value
}

func NewOutputVar(w io.Writer) *OutputVar {
	if o, ok := w.(*OutputVar); ok {
		return o
	}
	v := &OutputVar{}
	v.Set(w)
	return v
}

func (v *OutputVar) Set(w io.Writer) {
	v.ready.Store(w != nil && w != io.Discard)
	v.writer.Store(w)
}

func (v *OutputVar) Discard() bool {
	return !v.ready.Load()
}

func (v *OutputVar) Output() io.Writer {
	if v.Discard() {
		return io.Discard
	}
	return v.writer.Load().(io.Writer)
}

func (v *OutputVar) Write(p []byte) (int, error) {
	return v.Output().Write(p)
}

func stringify(j map[string]any) string {
	bts, _ := json.Marshal(j)
	return *(*string)(unsafe.Pointer(&bts))
}

// Having an initial size gives a dramatic speedup.
var rwPool = sync.Pool{
	New: func() any {
		return &recordWriter{}
	},
}

func newRecordWriter(l *SimpleHandler) *recordWriter {
	rw := rwPool.Get().(*recordWriter)
	rw.replace = l.options.ReplacePart
	return rw
}

type recordWriter struct {
	replace func(PartKind, *Record, bool) (string, bool)
	buf     []byte
	sep     bool
}

func (w *recordWriter) Reset() []byte {
	defer rwPool.Put(w)
	buf := w.buf[:]
	w.replace = nil
	w.buf = w.buf[:0]
	w.sep = false
	return buf
}

func (w *recordWriter) Write(kind PartKind, r *Record, last bool) {
	if w.replace != nil {
		s, ok := w.replace(kind, r, last)
		if ok {
			w.write(s)
			return
		}
	}

	switch kind {
	case PartLevel:
		w.write("%-5s", r.Level.String())
	case PartTime:
		w.write(r.Time.Format("15:04:05.000"))
	case PartMessage:
		lines := strings.Split(r.Message, "\n")
		for i, line := range lines {
			if i > 0 {
				w.write("    ")
			}
			w.write(line)
			w.write("\n")
		}
	case PartLocation:
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		w.write(" -> %s\n", f.Function)
		w.write(" ->   %s:%d\n", f.File, f.Line)
	case PartStacktrace:
		fs := runtime.CallersFrames(r.Frames)
		for {
			f, more := fs.Next()
			w.write("    %s\n", f.Function)
			w.write("      %s:%d\n", f.File, f.Line)
			if !more {
				break
			}
		}
	}
}

func (w *recordWriter) write(s string, args ...any) {
	if len(args) > 0 {
		s = fmt.Sprintf(s, args...)
	}
	raw := stripansi.String(s) // removes ANSI escape sequences
	if len(raw) == 0 {
		return
	}
	if w.sep {
		w.buf = append(w.buf, ' ')
	}
	w.buf = append(w.buf, s...)
	w.sep = raw[len(raw)-1] != '\n'
}

func (w *recordWriter) FlushTo(out io.Writer) error {
	_, err := out.Write(w.Reset())
	return err
}
