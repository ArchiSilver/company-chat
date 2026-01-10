package repository

import (
	"context"
	"fmt"
	"time"

	"company-chat/internal/domain"
)

type UploadRepository struct{ db *DB }

// UploadSaver abstracts saving upload metadata so handlers can be unit-tested without DB
type UploadSaver interface {
    Save(ctx context.Context, u *domain.Upload) error
}

func NewUploadRepository(db *DB) *UploadRepository { return &UploadRepository{db: db} }

func (r *UploadRepository) Save(ctx context.Context, u *domain.Upload) error {
    // let Postgres set created_at using DEFAULT to ensure DB timezone consistency
    query := `INSERT INTO uploads (id, user_id, path, mime, size, created_at) VALUES (gen_random_uuid(), $1, $2, $3, $4, DEFAULT) RETURNING id, created_at`
    var id string
    var created time.Time
    err := r.db.Pool.QueryRow(ctx, query, u.UserID, u.Path, u.MIME, u.Size).Scan(&id, &created)
    if err != nil {
        return fmt.Errorf("insert upload: %w", err)
    }
    u.ID = id
    u.CreatedAt = created
    return nil
}
