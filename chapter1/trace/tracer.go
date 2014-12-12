package trace

import (
	"fmt"
	"io"
)

// Tracer is the interface that describes an object capable of
// tracing events throughout code.
type Tracer interface {
	Trace(...interface{})
}

// New creates a new Tracer that will write the output to
// the specified io.Writer.
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

// tracer is a Tracer that writes to an
// io.Writer.
type tracer struct {
	out io.Writer
}

// Trace writes the arguments to this Tracers io.Writer.
func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

// nilTracer
type nilTracer struct{}

// Trace for a nil tracer does nothing.
func (t *nilTracer) Trace(a ...interface{}) {}

// Off creates a Tracer that will ignore calls to Trace.
func Off() Tracer {
	return &nilTracer{}
}
