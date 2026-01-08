package config

import (
    "fmt"
    "os"
)

// Config хранит все настройки приложения
type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    ServerPort string
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() *Config {
    return &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "app"),
        DBPassword: getEnv("DB_PASSWORD", "password"),
        DBName:     getEnv("DB_NAME", "app"),
        ServerPort: getEnv("SERVER_PORT", "8080"),
    }
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}

// GetDBConnectionString возвращает строку подключения к PostgreSQL
func (c *Config) GetDBConnectionString() string {
    return fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=disable",
        c.DBUser,
        c.DBPassword,
        c.DBHost,
        c.DBPort,
        c.DBName,
    )
}
