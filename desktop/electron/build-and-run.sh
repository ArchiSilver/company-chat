#!/usr/bin/env bash
set -euo pipefail

# Helper: install deps, start app (dev) or build production artifacts
# Usage:
#   ./build-and-run.sh dev   # npm ci && npm start
#   ./build-and-run.sh build # npm ci && npm run dist

MODE=${1:-dev}
cd "$(dirname "$0")"

echo "Node/npm must be installed on this machine."

if [ "$MODE" = "dev" ]; then
  npm ci
  npm start
elif [ "$MODE" = "build" ]; then
  npm ci
  npm run dist
else
  echo "Unknown mode: $MODE" >&2
  exit 2
fi
