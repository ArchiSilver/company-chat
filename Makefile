## Makefile — common development tasks

GO=go

.PHONY: fmt test run build migrate docker-build migrate-up

fmt:
	${GO} fmt ./...

test:
	${GO} test ./... -v

run:
	${GO} run ./cmd/app

build:
	${GO} build -o bin/app ./cmd/app

migrate:
	${GO} run ./cmd/migrate

migrate-up: migrate

docker-build:
	docker build -t company-chat:local .

.PHONY: preview
preview:
	@echo "Preview UI in browser (serves desktop/electron/public)"
	@zsh scripts/preview.sh

.PHONY: dev-up
dev-up:
	@echo "Bring up docker-compose, wait for backend and open UI preview"
	@zsh scripts/dev-up.sh

# Desktop build (Wails-based multi-platform desktop app)
.PHONY: desktop-dev desktop-build

desktop-dev:
	@which wails >/dev/null 2>&1 || (echo "wails CLI not found. Install it: https://wails.io/" && exit 1)
	wails dev

desktop-build:
	@which wails >/dev/null 2>&1 || (echo "wails CLI not found. Install it: https://wails.io/" && exit 1)
	wails build

.PHONY: desktop-electron

desktop-electron:
	@echo "Starting Electron app (install deps first: cd desktop/electron && npm install)"
	(cd desktop/electron && npm install --no-audit --no-fund && npm start)

