# Company Chat — проект для обучения Backend (Go)

Коротко: это учебный репозиторий для разработки корпоративного чата на Go. В нём реализуются шаги из вашего Roadmap: HTTP-сервер, PostgreSQL, аутентификация, WebSocket и т.д.


Как запустить локально (быстрый старт)

1. Сборка/запуск:

```zsh
make run
```

2. Полезные Make цели:

```zsh
make run    # запустить приложение
make docker-build  # собрать образ локально
make test   # прогнать тесты
```

3. Быстрый dev с docker-compose:

```bash
docker-compose up -d
make migrate
make run
```

Обязательные переменные окружения (примеры):

```bash
export DATABASE_URL='postgres://app:password@localhost:5432/app?sslmode=disable'
export JWT_SECRET='dev-secret'
# Включить per-chat WS метрику (по умолчанию выключена):
export METRICS_WS_BY_CHAT=0
```

Production notes:
- Храните `JWT_SECRET` безопасно (Vault or secrets manager). `JWT_SECRET_FILE` поддерживается для контейнеров.
- Для prod используйте внешние rate limiters / API gateway; встроенный in-memory rate limiter служит только для dev/small deployments.
- Настройте TLS / reverse proxy (nginx, Caddy) перед приложением.


3. Структура проекта:

- `cmd/` — исполняемые пакеты (app, migrate и т.д.)
- `internal/config` — конфигурация
- `internal/domain` — доменные модели
- `internal/repository` — слой доступа к БД
- `migrations/` — SQL миграции

Следующие шаги: создайте ветку `feature/bootstrap-001`, сделайте коммит и пуш. После этого перейдём к Day 2 (вынесение конфига и graceful shutdown).
