package repository

import (
    "context"
    "fmt"
    "log"

    "github.com/jackc/pgx/v5/pgxpool"
)

// DB представляет пул соединений к базе данных
type DB struct {
    Pool *pgxpool.Pool
}

// NewDB создает новое подключение к базе данных
func NewDB(connectionString string) (*DB, error) {
    pool, err := pgxpool.New(context.Background(), connectionString)
    if err != nil {
        return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
    }

    // Проверим подключение
    if err := pool.Ping(context.Background()); err != nil {
        return nil, fmt.Errorf("не удалось проверить подключение: %w", err)
    }

    log.Println("✅ Успешное подключение к PostgreSQL")
    return &DB{Pool: pool}, nil
}

// Close закрывает соединение с БД
func (db *DB) Close() {
    db.Pool.Close()
}
