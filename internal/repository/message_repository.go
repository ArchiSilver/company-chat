package repository

import (
	"context"
	"fmt"
	"time"

	"company-chat/internal/domain"
)

type MessageRepository struct {
	db *DB
}

func NewMessageRepository(db *DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(ctx context.Context, m *domain.Message) error {
	query := `INSERT INTO messages (chat_id, sender_id, content, created_at) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	var id string
	var created time.Time
	err := r.db.Pool.QueryRow(ctx, query, m.ChatID, m.SenderID, m.Content, time.Now()).Scan(&id, &created)
	if err != nil {
		return fmt.Errorf("insert message: %w", err)
	}
	m.ID = id
	m.CreatedAt = created
	return nil
}

func (r *MessageRepository) ListByChat(ctx context.Context, chatID string, limit, offset int) ([]domain.Message, error) {
	query := `SELECT id, chat_id, sender_id, content, created_at FROM messages WHERE chat_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Pool.Query(ctx, query, chatID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query messages: %w", err)
	}
	defer rows.Close()
	var res []domain.Message
	for rows.Next() {
		var m domain.Message
		if err := rows.Scan(&m.ID, &m.ChatID, &m.SenderID, &m.Content, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}
		res = append(res, m)
	}
	return res, nil
}
