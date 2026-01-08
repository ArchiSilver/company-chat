package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestHashTokenKnown(t *testing.T) {
	input := "abc"
	expected := "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	got := hashToken(input)
	if got != expected {
		t.Fatalf("hash mismatch: expected %s got %s", expected, got)
	}
}

func TestHashIsSHA256Len(t *testing.T) {
	s := "some-random-string"
	got := hashToken(s)
	// ensure it's valid hex of length 64
	if len(got) != 64 {
		t.Fatalf("expected length 64, got %d", len(got))
	}
	if _, err := hex.DecodeString(got); err != nil {
		t.Fatalf("not valid hex: %v", err)
	}
	// cross-check with standard library
	h := sha256.Sum256([]byte(s))
	if hex.EncodeToString(h[:]) != got {
		t.Fatalf("mismatch with stdlib sha256")
	}
}
