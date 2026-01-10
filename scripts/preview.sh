#!/usr/bin/env zsh
# Простой скрипт для предпросмотра статического UI в `desktop/electron/public`
# Открывает локальный http сервер и запускает браузер (macOS `open`).

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
PUBLIC_DIR="$ROOT_DIR/desktop/electron/public"
PORT=${PORT:-9000}

if [ ! -d "$PUBLIC_DIR" ]; then
  echo "Папка с UI не найдена: $PUBLIC_DIR"
  exit 1
fi

echo "Запускаю статический сервер для $PUBLIC_DIR на порту $PORT"

# Попробуем python3, затем python
if command -v python3 >/dev/null 2>&1; then
  (cd "$PUBLIC_DIR" && python3 -m http.server "$PORT") &
  SERVER_PID=$!
elif command -v python >/dev/null 2>&1; then
  (cd "$PUBLIC_DIR" && python -m SimpleHTTPServer "$PORT") &
  SERVER_PID=$!
else
  echo "python3 или python не найдены. Установите Python 3 или используйте Electron (npm start)."
  exit 2
fi

# Открыть в браузере (macOS)
URL="http://localhost:$PORT"
echo "Открываю $URL в браузере..."
open "$URL" || echo "Не удалось автоматически открыть браузер. Перейдите по адресу: $URL"

echo "Нажмите Ctrl-C чтобы остановить сервер"
wait $SERVER_PID
