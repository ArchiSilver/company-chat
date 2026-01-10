package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"company-chat/internal/config"
	"company-chat/internal/domain"
)

func TestUploadRepositoryIntegration(t *testing.T) {
	cfg := config.LoadConfig()
	db, err := NewDB(cfg.GetDBConnectionString())
	if err != nil {
		t.Skipf("no db available: %v", err)
	}
	defer db.Close()

	// check uploads table exists
	var exists bool
	err = db.Pool.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name='uploads')").Scan(&exists)
	if err != nil {
		t.Skipf("could not query information_schema: %v", err)
	}
	if !exists {
		t.Skip("migrations not applied, skipping integration test")
	}

	repo := NewUploadRepository(db)
	ctx := context.Background()

	// create a user for FK constraint
	userID := uuid.New().String()
	if _, err := db.Pool.Exec(ctx, "INSERT INTO users (id, email, username, password_hash) VALUES ($1,$2,$3,$4)", userID, "integration+user@example.com", "integration_user", "hash"); err != nil {
		t.Fatalf("failed to insert user for integration test: %v", err)
	}
	defer func() {
		_, _ = db.Pool.Exec(ctx, "DELETE FROM users WHERE id=$1", userID)
	}()

	u := &domain.Upload{UserID: userID, Path: "/uploads/integration-test.bin", MIME: "application/octet-stream", Size: 42}
	if err := repo.Save(ctx, u); err != nil {
		t.Fatalf("save: %v", err)
	}
	if u.ID == "" {
		t.Fatalf("expected id to be set")
	}
	if time.Since(u.CreatedAt) > time.Minute*5 {
		t.Fatalf("created_at seems incorrect: %v", u.CreatedAt)
	}

	// cleanup
	_, _ = db.Pool.Exec(ctx, "DELETE FROM uploads WHERE id=$1", u.ID)
}
