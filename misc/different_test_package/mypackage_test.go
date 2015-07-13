package mypackage_test

import (
	"testing"

	"github.com/matryer/goblueprints/misc/different_test_package"
)

func TestCount(t *testing.T) {
	if mypackage.Count() != 1 {
		t.Error("expected 1")
	}
	if mypackage.Count() != 2 {
		t.Error("expected 2")
	}
	if mypackage.Count() != 3 {
		t.Error("expected 3")
	}
}
