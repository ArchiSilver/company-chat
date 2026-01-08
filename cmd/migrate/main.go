package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
    // Подключение к базе данных
    db, err := sql.Open("pgx", "postgres://app:password@localhost:5432/app?sslmode=disable")
    if err != nil {
        log.Fatal("Ошибка подключения к БД:", err)
    }
    defer db.Close()

    // Проверка подключения
    if err := db.Ping(); err != nil {
        log.Fatal("Не удалось подключиться к БД:", err)
    }

    // Чтение SQL файла
    sqlBytes, err := os.ReadFile("migrations/001_create_users.up.sql")
    if err != nil {
        log.Fatal("Не удалось прочитать файл миграции:", err)
    }

    sqlString := string(sqlBytes)

    // Выполнение миграции
    _, err = db.ExecContext(context.Background(), sqlString)
    if err != nil {
        log.Fatal("Ошибка выполнения миграции:", err)
    }

    fmt.Println("✅ Миграция успешно выполнена!")
    fmt.Println("✅ Таблица 'users' создана")
}
