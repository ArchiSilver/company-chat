package repository

import (
	"context"
	"testing"

	"company-chat/internal/config"
	"company-chat/internal/domain"
)

func TestUserRepositoryIntegration(t *testing.T) {
	cfg := config.LoadConfig()
	db, err := NewDB(cfg.GetDBConnectionString())
	if err != nil {
		t.Skipf("no db available: %v", err)
	}
	defer db.Close()

	// ensure migrations applied (table exists)
	var exists bool
	err = db.Pool.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name='users')").Scan(&exists)
	if err != nil {
		t.Skipf("could not query information_schema: %v", err)
	}
	if !exists {
		t.Skip("migrations not applied, skipping integration test")
	}

	// ensure role column exists (schema as expected)
	var colExists bool
	err = db.Pool.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='role')").Scan(&colExists)
	if err != nil {
		t.Skipf("could not query information_schema.columns: %v", err)
	}
	if !colExists {
		t.Skip("users.role column missing, skipping integration test")
	}

	repo := NewUserRepository(db)
	ctx := context.Background()
	u := &domain.User{
		Email:        "integ+user@example.com",
		Username:     "integuser",
		PasswordHash: "hash",
		Role:         "user",
	}

	if err := repo.Create(ctx, u); err != nil {
		t.Fatalf("create user: %v", err)
	}
	if u.ID == "" {
		t.Fatalf("expected user id to be set")
	}

	got, err := repo.FindByEmail(ctx, u.Email)
	if err != nil {
		t.Fatalf("find by email error: %v", err)
	}
	if got == nil || got.Email != u.Email {
		t.Fatalf("expected to find user by email")
	}

	got2, err := repo.FindByID(ctx, u.ID)
	if err != nil {
		t.Fatalf("find by id error: %v", err)
	}
	if got2 == nil || got2.ID != u.ID {
		t.Fatalf("expected to find user by id")
	}

	// cleanup: remove created test user
	_, _ = db.Pool.Exec(ctx, "DELETE FROM users WHERE id=$1", u.ID)
}
