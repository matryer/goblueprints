package trace

import (
	"bytes"
	"testing"
)

// TestNew tests the tracing behaviour.
func TestNew(t *testing.T) {

	var buf bytes.Buffer
	tracer := New(&buf)

	if tracer == nil {
		t.Error("Return from New should not be nil")
	}
	tracer.Trace("Hello trace package.")
	if buf.String() != "Hello trace package.\n" {
		t.Errorf("Trace should not write '%s'.", buf.String())
	}

}

func TestOff(t *testing.T) {
	silentTracer := Off()
	silentTracer.Trace("something")
}
