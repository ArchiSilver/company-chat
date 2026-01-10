#!/usr/bin/env zsh
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
COMPOSE="docker compose"
TIMEOUT=${TIMEOUT:-60}

echo "Starting docker-compose services..."
cd "$ROOT_DIR"
$COMPOSE up -d --build

echo "Waiting for backend to become healthy on http://localhost:8080/health (timeout ${TIMEOUT}s)"
START=$(date +%s)
while true; do
  if curl -sSf http://localhost:8080/health >/dev/null 2>&1; then
    echo "backend is up"
    break
  fi
  NOW=$(date +%s)
  if [ $((NOW-START)) -ge $TIMEOUT ]; then
    echo "Timeout waiting for backend"
    exit 1
  fi
  sleep 1
done

echo "Opening UI preview..."
zsh scripts/preview.sh
