package vault

import (
	"testing"

	"golang.org/x/net/context"
)

func TestHasherService(t *testing.T) {
	srv := NewService()
	ctx := context.Background()
	h, err := srv.Hash(ctx, "password")
	if err != nil {
		t.Errorf("Hash: %s", err)
	}
	ok, err := srv.Validate(ctx, "password", h)
	if err != nil {
		t.Errorf("Valid: %s", err)
	}
	if !ok {
		t.Error("expected true from Valid")
	}
	ok, err = srv.Validate(ctx, "wrong password", h)
	if err != nil {
		t.Errorf("Valid: %s", err)
	}
	if ok {
		t.Error("expected false from Valid")
	}
}
