package repository

import (
	"context"
	"fmt"
	"time"

	"company-chat/internal/domain"

	"github.com/jackc/pgx/v5"
)

// UserRepository работает с пользователями
type UserRepository struct {
	db *DB
}

// NewUserRepository создаёт репозиторий
func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create добавляет пользователя и возвращает его ID
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (email, username, password_hash, role, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	var id string
	var created time.Time
	err := r.db.Pool.QueryRow(ctx, query, user.Email, user.Username, user.PasswordHash, user.Role, time.Now()).Scan(&id, &created)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	user.ID = id
	user.CreatedAt = created
	return nil
}

// FindByEmail возвращает пользователя по email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, username, password_hash, role, created_at FROM users WHERE email=$1`
	row := r.db.Pool.QueryRow(ctx, query, email)
	var u domain.User
	if err := row.Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.Role, &u.CreatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("select user: %w", err)
	}
	return &u, nil
}

// FindByID возвращает пользователя по id
func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, email, username, password_hash, role, created_at FROM users WHERE id=$1`
	row := r.db.Pool.QueryRow(ctx, query, id)
	var u domain.User
	if err := row.Scan(&u.ID, &u.Email, &u.Username, &u.PasswordHash, &u.Role, &u.CreatedAt); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("select user by id: %w", err)
	}
	return &u, nil
}
