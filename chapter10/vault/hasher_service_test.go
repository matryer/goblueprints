package main

import "testing"

func TestHasherService(t *testing.T) {
	srv := hasherService{}
	h, err := srv.Hash("password")
	if err != nil {
		t.Errorf("Hash: %s", err)
	}
	ok := srv.Valid("password", h)
	if !ok {
		t.Error("expected true from Valid")
	}
	ok = srv.Valid("wrong password", h)
	if ok {
		t.Error("expected false from Valid")
	}
}
