package l4g

import (
	"log/slog"
	"slices"
	"time"
)

// nAttrsInline is the number of attributes to store inline in a Record
// to avoid heap allocation for small log calls. This value is tuned based
// on examination of typical log usage patterns.
const nAttrsInline = 5

// A Record holds information about a log event.
// Do not modify a Record after handing out a copy to it.
type Record struct {
	// The time at which the output method (Log, Info, etc.) was called.
	Time time.Time

	// The log prefix.
	Prefix string

	// The log message.
	Message string

	// The level of the event.
	Level Level

	// Allocation optimization: an inline array sized to hold
	// the majority of log calls (based on examination of open-source
	// code). It holds the start of the list of Attrs.
	front [nAttrsInline]Attr

	// The number of Attrs in front.
	nFront int

	// The list of Attrs except for those in front.
	// Invariants:
	//   - len(back) > 0 iff nFront == len(front)
	//   - Unused array elements are zero. Used to detect mistakes.
	back []Attr
}

// NewRecord creates a [Record] from the given arguments.
// Use [Record.AddAttrs] to add attributes to the Record.
//
// NewRecord is intended for logging APIs that want to support a [Handler] as
// a backend.
func NewRecord(t time.Time, level Level, msg string) Record {
	return Record{
		Time:    t,
		Message: msg,
		Level:   level,
	}
}

// Clone returns a copy of the record with no shared state.
// The original record and the clone can both be modified
// without interfering with each other.
func (r Record) Clone() Record {
	r.back = slices.Clip(r.back) // prevent append from mutating shared array
	return r
}

// NumAttrs returns the number of attributes in the [Record].
func (r Record) NumAttrs() int {
	return r.nFront + len(r.back)
}

// Attrs calls f on each Attr in the [Record].
// Iteration stops if f returns false.
func (r Record) Attrs(f func(Attr) bool) {
	for i := 0; i < r.nFront; i++ {
		if !f(r.front[i]) {
			return
		}
	}
	for _, a := range r.back {
		if !f(a) {
			return
		}
	}
}

// AddAttrs appends the given Attrs to the [Record]'s list of Attrs.
// It omits empty groups.
func (r *Record) AddAttrs(attrs ...Attr) {
	var i int
	for i = 0; i < len(attrs) && r.nFront < len(r.front); i++ {
		a := attrs[i]
		if isEmptyGroup(a.Value) {
			continue
		}
		r.front[r.nFront] = a
		r.nFront++
	}
	// Check if a copy was modified by slicing past the end
	// and seeing if the Attr there is non-zero.
	if cap(r.back) > len(r.back) {
		end := r.back[:len(r.back)+1][len(r.back)]
		if !isEmptyAttr(end) {
			// Don't panic; copy and muddle through.
			r.back = slices.Clip(r.back)
			r.back = append(r.back, String("!BUG", "AddAttrs unsafely called on copy of Record made without using Record.Clone"))
		}
	}
	ne := countEmptyGroups(attrs[i:])
	r.back = slices.Grow(r.back, len(attrs[i:])-ne)
	for _, a := range attrs[i:] {
		if !isEmptyGroup(a.Value) {
			r.back = append(r.back, a)
		}
	}
}

// Add converts the args to Attrs as described in [Logger.Log],
// then appends the Attrs to the [Record]'s list of Attrs.
// It omits empty groups.
func (r *Record) Add(args ...any) {
	var a Attr
	for len(args) > 0 {
		a, args = argsToAttr(args)
		if isEmptyGroup(a.Value) {
			continue
		}
		if r.nFront < len(r.front) {
			r.front[r.nFront] = a
			r.nFront++
		} else {
			if r.back == nil {
				r.back = make([]Attr, 0, countAttrs(args)+1)
			}
			r.back = append(r.back, a)
		}
	}
}

// countAttrs returns the number of Attrs that would be created from args.
func countAttrs(args []any) int {
	n := 0
	for i := 0; i < len(args); i++ {
		n++
		if _, ok := args[i].(string); ok {
			i++
		}
	}
	return n
}

// countEmptyGroups returns the number of empty group values in its argument.
func countEmptyGroups(as []Attr) int {
	n := 0
	for _, a := range as {
		if isEmptyGroup(a.Value) {
			n++
		}
	}
	return n
}

// isEmptyAttr reports whether a has an empty key and a nil value.
// That can be written as Attr{} or Any("", nil).
func isEmptyAttr(a Attr) bool {
	//return a.Key == "" && a.Value.num == 0 && a.Value.any == nil
	return a.Key == "" && (a.Value.Equal(slog.Value{}) || a.Value.Equal(slog.AnyValue(nil)))
}

// isEmptyGroup reports whether v is a group that has no attributes.
func isEmptyGroup(v slog.Value) bool {
	if v.Kind() != slog.KindGroup {
		return false
	}
	//// We do not need to recursively examine the group's Attrs for emptiness,
	//// because GroupValue removed them when the group was constructed, and
	//// groups are immutable.
	//return len(v.group()) == 0
	return len(v.Group()) == 0
}

// badKey is used as the key for attributes that have a value but no key.
// This happens when argsToAttr receives a non-string, non-Attr argument.
const badKey = "!BADKEY"

// argsToAttr turns a prefix of the nonempty args slice into an Attr
// and returns the unconsumed portion of the slice.
// If args[0] is an Attr, it returns it.
// If args[0] is a string, it treats the first two elements as
// a key-value pair.
// Otherwise, it treats args[0] as a value with a missing key.
func argsToAttr(args []any) (Attr, []any) {
	switch x := args[0].(type) {
	case string:
		if len(args) == 1 {
			return String(badKey, x), nil
		}
		return Any(x, args[1]), args[2:]

	case Attr:
		return x, args[1:]

	default:
		return Any(badKey, x), args[1:]
	}
}
