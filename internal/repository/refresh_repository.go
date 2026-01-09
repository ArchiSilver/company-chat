package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type RefreshTokenRepository struct {
	db *DB
}

func NewRefreshTokenRepository(db *DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func (r *RefreshTokenRepository) Save(ctx context.Context, userID, token string, expiresAt time.Time) error {
	q := `INSERT INTO refresh_tokens (user_id, token_hash, expires_at, revoked, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Pool.Exec(ctx, q, userID, hashToken(token), expiresAt, false, time.Now())
	if err != nil {
		return fmt.Errorf("save refresh token: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) RevokeByHash(ctx context.Context, token string) error {
	q := `UPDATE refresh_tokens SET revoked=true WHERE token_hash=$1`
	_, err := r.db.Pool.Exec(ctx, q, hashToken(token))
	if err != nil {
		return fmt.Errorf("revoke token: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) IsValid(ctx context.Context, token string) (bool, error) {
	q := `SELECT revoked, expires_at FROM refresh_tokens WHERE token_hash=$1 LIMIT 1`
	var revoked bool
	var expires time.Time
	err := r.db.Pool.QueryRow(ctx, q, hashToken(token)).Scan(&revoked, &expires)
	if err != nil {
		return false, err
	}
	if revoked {
		return false, nil
	}
	if time.Now().After(expires) {
		return false, nil
	}
	return true, nil
}

// GetUserIDByToken возвращает user_id, если токен действителен (не revoked и не истёк)
func (r *RefreshTokenRepository) GetUserIDByToken(ctx context.Context, token string) (string, error) {
	q := `SELECT user_id FROM refresh_tokens WHERE token_hash=$1 AND revoked=false AND expires_at > now() LIMIT 1`
	var userID string
	err := r.db.Pool.QueryRow(ctx, q, hashToken(token)).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}
