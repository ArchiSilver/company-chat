# Running tests locally

This document explains how to run unit and integration tests locally for the Company Chat project.

Prerequisites
- Docker (for local Postgres)
- Go 1.24

1. Start Postgres via docker-compose

```bash
docker-compose up -d
```

2. Apply migrations

```bash
go run ./cmd/migrate
```

3. Run tests

```bash
# run all tests (integrations will run if DB/migrations present)
go test ./... -v
```

Notes
- Integration tests are written to skip when migrations / expected schema are not present. When running locally with docker-compose and migrations applied they will execute.
- If you want to run only unit tests:

```bash
go test ./internal/... -v
```

- CI is configured to run migrations and tests automatically.

Tips:
- Use `make test` as a shortcut for `go test ./... -v`.
- If integration tests are skipped, ensure:
	- `docker-compose up -d` is running
	- `go run ./cmd/migrate` applied migrations successfully

