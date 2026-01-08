package main

import (
    "context"
    "fmt"
    "log"
    "net/http"

    "company-chat/internal/config"
    "company-chat/internal/repository"
)

func main() {
    // Загружаем конфигурацию
    cfg := config.LoadConfig()
    
    // Подключаемся к базе данных
    db, err := repository.NewDB(cfg.GetDBConnectionString())
    if err != nil {
        log.Fatalf("Ошибка подключения к БД: %v", err)
    }
    defer db.Close()

    // Обработчик для корневого пути
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "🚀 Company Chat API v2.0\n")
        fmt.Fprintf(w, "✅ Сервер работает\n")
        fmt.Fprintf(w, "📊 База данных: подключено\n")
    })

    // Обработчик для health check
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        // Проверяем подключение к БД
        if err := db.Pool.Ping(context.Background()); err != nil {
            w.WriteHeader(http.StatusServiceUnavailable)
            fmt.Fprintf(w, "❌ Database: DOWN\n")
            return
        }
        
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, "✅ Все системы работают\n")
        fmt.Fprintf(w, "✅ Database: CONNECTED\n")
    })

    // Эндпоинт для проверки таблицы users
    http.HandleFunc("/api/db-check", func(w http.ResponseWriter, r *http.Request) {
        var count int
        err := db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users").Scan(&count)
        
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "❌ Ошибка запроса к БД: %v", err)
            return
        }
        
        fmt.Fprintf(w, "✅ Таблица 'users' существует\n")
        fmt.Fprintf(w, "📊 Количество пользователей: %d\n", count)
    })

    // Запуск сервера
    port := ":" + cfg.ServerPort
    fmt.Printf("🚀 Сервер запущен на http://localhost%s\n", port)
    fmt.Printf("�� База данных: %s:%s\n", cfg.DBHost, cfg.DBPort)
    log.Fatal(http.ListenAndServe(port, nil))
}
