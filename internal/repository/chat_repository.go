package repository

import (
	"context"
	"fmt"
	"time"

	"company-chat/internal/domain"
)

type ChatRepository struct {
	db *DB
}

func NewChatRepository(db *DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) Create(ctx context.Context, chat *domain.Chat) error {
	query := `INSERT INTO chats (name, type, created_by, created_at) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	var id string
	var created time.Time
	err := r.db.Pool.QueryRow(ctx, query, chat.Name, chat.Type, chat.CreatedBy, time.Now()).Scan(&id, &created)
	if err != nil {
		return fmt.Errorf("insert chat: %w", err)
	}
	chat.ID = id
	chat.CreatedAt = created
	return nil
}

func (r *ChatRepository) FindByID(ctx context.Context, id string) (*domain.Chat, error) {
	query := `SELECT id, name, type, created_by, created_at FROM chats WHERE id=$1`
	row := r.db.Pool.QueryRow(ctx, query, id)
	var c domain.Chat
	if err := row.Scan(&c.ID, &c.Name, &c.Type, &c.CreatedBy, &c.CreatedAt); err != nil {
		return nil, fmt.Errorf("select chat: %w", err)
	}
	return &c, nil
}
