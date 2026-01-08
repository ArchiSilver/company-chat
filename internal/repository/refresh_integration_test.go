package repository

import (
	"context"
	"testing"
	"time"

	"company-chat/internal/config"
)

func TestRefreshTokenRepositoryIntegration(t *testing.T) {
	cfg := config.LoadConfig()
	db, err := NewDB(cfg.GetDBConnectionString())
	if err != nil {
		t.Skipf("no db available: %v", err)
	}
	defer db.Close()

	// check migrations applied (table exists) and expected columns present
	var exists bool
	err = db.Pool.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name='refresh_tokens')").Scan(&exists)
	if err != nil {
		t.Skipf("could not query information_schema: %v", err)
	}
	if !exists {
		t.Skip("migrations not applied, skipping integration test")
	}

	// ensure token_hash column exists
	var colExists bool
	err = db.Pool.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='refresh_tokens' AND column_name='token_hash')").Scan(&colExists)
	if err != nil {
		t.Skipf("could not query information_schema.columns: %v", err)
	}
	if !colExists {
		t.Skip("refresh_tokens.token_hash column missing, skipping integration test")
	}

	repo := NewRefreshTokenRepository(db)
	ctx := context.Background()
	token := "integration-test-token"
	userID := "integration-user-1"
	expires := time.Now().Add(1 * time.Hour)

	if err := repo.Save(ctx, userID, token, expires); err != nil {
		t.Fatalf("save: %v", err)
	}

	valid, err := repo.IsValid(ctx, token)
	if err != nil {
		t.Fatalf("isvalid error: %v", err)
	}
	if !valid {
		t.Fatalf("expected token to be valid")
	}

	gotUser, err := repo.GetUserIDByToken(ctx, token)
	if err != nil {
		t.Fatalf("get user by token error: %v", err)
	}
	if gotUser != userID {
		t.Fatalf("expected user %s, got %s", userID, gotUser)
	}

	if err := repo.RevokeByHash(ctx, token); err != nil {
		t.Fatalf("revoke error: %v", err)
	}

	validAfter, err := repo.IsValid(ctx, token)
	if err != nil {
		t.Fatalf("isvalid after revoke error: %v", err)
	}
	if validAfter {
		t.Fatalf("expected token to be invalid after revoke")
	}

	// cleanup: remove test token
	_, _ = db.Pool.Exec(ctx, "DELETE FROM refresh_tokens WHERE token_hash=$1", hashToken(token))
}
