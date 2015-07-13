package mypackage

import (
	"testing"
)

func TestCount(t *testing.T) {
	if Count() != 1 {
		t.Error("expected 1")
	}
	if count != 1 {
		t.Error("expected 1 for count too")
	}
}
