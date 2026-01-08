package auth

import (
	"testing"
	"time"
)

func TestGenerateAndParseAccessToken(t *testing.T) {
	secret := "test-secret"
	uid := "user-123"
	tok, err := GenerateAccessToken(secret, uid, time.Minute)
	if err != nil {
		t.Fatalf("GenerateAccessToken error: %v", err)
	}
	got, err := ParseAndValidate(secret, tok)
	if err != nil {
		t.Fatalf("ParseAndValidate error: %v", err)
	}
	if got != uid {
		t.Fatalf("expected user id %s, got %s", uid, got)
	}
}

func TestExpiredToken(t *testing.T) {
	secret := "test-secret"
	uid := "u"
	tok, err := GenerateAccessToken(secret, uid, -time.Minute)
	if err != nil {
		t.Fatalf("GenerateAccessToken error: %v", err)
	}
	if _, err := ParseAndValidate(secret, tok); err == nil {
		t.Fatalf("expected error for expired token")
	}
}

func TestGenerateRandomToken(t *testing.T) {
	a, err := GenerateRandomToken(16)
	if err != nil {
		t.Fatalf("GenerateRandomToken error: %v", err)
	}
	b, err := GenerateRandomToken(16)
	if err != nil {
		t.Fatalf("GenerateRandomToken error: %v", err)
	}
	if a == b {
		t.Fatalf("expected two tokens to differ")
	}
	if len(a) != 32 {
		t.Fatalf("expected hex length 32, got %d", len(a))
	}
}
