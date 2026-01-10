```markdown
# Company Chat — обучающий проект (Backend на Go)

Небольшое описание: это учебный репозиторий для разработки корпоративного чат‑сервиса на Go. В проекте присутствуют: HTTP‑сервер, PostgreSQL, аутентификация, WebSocket, загрузки файлов и другие компоненты, которые входят в ваш Roadmap.

Как запустить локально (быстрый старт)

1. Быстрый запуск через Make (если настроено):

```zsh
make run
```

2. Полезные цели Makefile:

```zsh
make run            # запустить приложение локально
make docker-build   # собрать Docker-образ
make test           # прогнать тесты
```

3. Быстрая разработка с Docker Compose:

```bash
docker compose up -d
make migrate
make run
```

Обязательные переменные окружения (пример):

```bash
export DATABASE_URL='postgres://app:password@localhost:5432/app?sslmode=disable'
export JWT_SECRET='dev-secret'
# Включить per-chat WS метрику (по умолчанию выключена):
export METRICS_WS_BY_CHAT=0
```

Production — заметки безопасности
- Храните `JWT_SECRET` в безопасном хранилище (Vault, Secrets Manager). Для контейнеров поддерживается `JWT_SECRET_FILE`.
- Для production используйте внешний rate limiter / API gateway; встроенный in-memory rate limiter подходит только для разработки или мелких развёртываний.
- Обязательно настройте TLS и reverse proxy (nginx, Caddy) перед приложением.

Структура проекта
- `cmd/` — исполняемые пакеты (app, migrate и т.д.)
- `internal/config` — конфигурация
- `internal/domain` — доменные модели
- `internal/repository` — слой доступа к базе данных
- `migrations/` — SQL‑миграции

Рекомендации по локальному запуску (рекомендуется Docker)
-----------------------------------------------------

Требования (один из вариантов):
- Docker & Docker Compose (рекомендуется)
- Локальная среда: Go 1.20+, PostgreSQL, Node.js/npm (если вы хотите запускать Electron‑оболочку)

Quickstart с Docker
1. Соберите и запустите сервисы (Postgres + приложение). Контейнер выполнит миграции перед стартом сервера:

```zsh
docker compose build --pull
docker compose up -d
```

2. Просмотрите логи и дождитесь готовности сервера:

```zsh
docker compose logs -f app
```

3. Откройте UI в браузере (если он присутствует) или обращайтесь к API по адресу:

http://localhost:8080

Локальная разработка без Docker
--------------------------------
1. На macOS можно установить зависимости через Homebrew:

```zsh
brew install go postgresql node
# запустите postgres и создайте БД/пользователя в соответствии с конфигом или установите переменные окружения
brew services start postgresql
```

2. Примените миграции:

```zsh
make migrate
```

3. Запустите сервер:

```zsh
go run ./cmd/app
```

4. Опционально: запустите Electron desktop (требуется Node.js):

```zsh
cd desktop/electron
npm ci
npm start
```

Визуальное тестирование
-----------------------
- В репозитории есть статический веб‑интерфейс, который используется оболочкой Electron: `desktop/electron/public/index.html`. Для быстрой проверки можно открыть этот файл в браузере и убедиться, что API доступен по `http://localhost:8080`.
- Для точного пользовательского опыта запустите Electron (`npm start`) или соберите инсталляторы через `npm run dist` (требуется electron‑builder и соответствующая среда для сборки).

Дальнейшие улучшения (по запросу)
- могу добавить скрипт, который автоматически откроет браузер на UI после поднятия контейнеров;
- могу подготовить Docker‑образ для визуального предпросмотра в CI;
- могу реализовать транзакционную логику загрузки (временный файл + запись в БД) для предотвращения «осиротевших» файлов.

```
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

Run locally (recommended using Docker)
------------------------------------

Prerequisites (choose one):

- Docker & Docker Compose (recommended) — macOS: https://docs.docker.com/desktop/install/mac-install/
- OR local toolchain: Go 1.20+, PostgreSQL, Node.js/npm (for the Electron desktop wrapper)

Quickstart with Docker
1. Build and start services (Postgres + app). The container will run migrations before starting the server:

```zsh
docker compose build --pull
docker compose up -d
```

2. Check logs and wait for the server to be ready:

```zsh
docker compose logs -f app
```

3. Open the UI in a browser (server serves a web client at /web if present) or access the API at:

http://localhost:8080

Local dev without Docker
------------------------
1. Install prerequisites on macOS (Homebrew):

```zsh
brew install go postgresql node
# start postgres, create DB/user as configured in internal/config defaults or set env vars
brew services start postgresql
```

2. Apply migrations:

```zsh
make migrate
```

3. Run the server:

```zsh
go run ./cmd/app
```

4. Run the Electron desktop (optional, requires Node):

```zsh
cd desktop/electron
npm ci
npm start
```

Notes about visual testing
--------------------------
- The repository contains a static web UI used by the Electron wrapper at `desktop/electron/public/index.html`. For a quick visual, you can open the web UI in a browser if the server serves it under `/web/` (or open `desktop/electron/public/index.html` directly in a browser while the API is running at http://localhost:8080 — the UI will use `localStorage.api_base` or fallback to `http://localhost:8080`).
- To see the desktop experience exactly as end-users will, run the Electron dev flow above (npm start) or build installers with `npm run dist` (requires electron-builder and appropriate build environment per OS).

If you'd like, I can:
- add a tiny script to start the browser pointing at the UI automatically after the container is up, or
- prepare an Electron desktop dev image (for remote CI preview), or
- implement the transactional upload flow (temp file + pending DB record) to avoid orphan files.
