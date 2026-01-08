-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

COMMENT ON TABLE users IS 'Таблица пользователей мессенджера';
COMMENT ON COLUMN users.email IS 'Электронная почта (уникальная)';
COMMENT ON COLUMN users.username IS 'Имя пользователя (уникальное)';
COMMENT ON COLUMN users.password_hash IS 'Хэш пароля';
