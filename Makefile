## Makefile — common development tasks

GO=go

.PHONY: fmt test run migrate docker-build

fmt:
	${GO} fmt ./...

test:
	${GO} test ./... -v

run:
	${GO} run ./cmd/app

migrate:
	${GO} run ./cmd/migrate

docker-build:
	docker build -t company-chat:local .
.PHONY: run build test

run:
	go run ./cmd/app

build:
	go build -o bin/app ./cmd/app

test:
	go test ./... -v

migrate-up:
	go run ./cmd/migrate

