package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB представляет пул соединений
type DB struct {
	Pool *pgxpool.Pool
}

// NewDB инициализирует пул соединений
func NewDB(conn string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), conn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return &DB{Pool: pool}, nil
}

// Close закрывает пул
func (d *DB) Close() {
	if d.Pool != nil {
		d.Pool.Close()
	}
}
