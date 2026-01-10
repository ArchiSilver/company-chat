#!/usr/bin/env bash
set -euo pipefail

BASE="http://localhost:8080"

echo "Starting smoke test against $BASE"

# register (ignore errors if exists)
curl -s -X POST "$BASE/api/auth/register" -H "Content-Type: application/json" \
  -d '{"email":"smoke+test@example.com","username":"smoke","password":"pass123"}' || true

echo "Login..."
TOKEN=$(curl -s -X POST "$BASE/api/auth/login" -H "Content-Type: application/json" \
  -d '{"email":"smoke+test@example.com","password":"pass123"}' | grep -o '"access_token"[^,}]*' | sed -E 's/.*:"([^"]+)"/\1/')

if [ -z "$TOKEN" ]; then
  echo "Failed to obtain token; full response:"
  curl -s -X POST "$BASE/api/auth/login" -H "Content-Type: application/json" \
    -d '{"email":"smoke+test@example.com","password":"pass123"}' || true
  exit 1
fi

echo "Token obtained (len=${#TOKEN})"

echo "Create chat..."
CHAT=$(curl -s -X POST "$BASE/api/chats" -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"smoke-room","type":"room"}' | grep -o '"id"[^,}]*' | sed -E 's/.*:"?([^"]+)"?/\1/')

if [ -z "$CHAT" ]; then
  echo "Failed to create chat; response:"
  curl -s -X POST "$BASE/api/chats" -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" \
    -d '{"name":"smoke-room","type":"room"}'
  exit 1
fi

echo "Chat id: $CHAT"

echo "Post REST message..."
curl -s -X POST "$BASE/api/chats/$CHAT/messages" -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"content":"hello from REST smoke"}' | jq || true

echo "Metrics snapshot (filtered):"
curl -s $BASE/metrics | grep company_chat_messages_total -n || true

echo "Listing recent chat_participants (for debugging):"
psql "postgres://app:password@localhost:5432/app?sslmode=disable" -c "SELECT chat_id, user_id, joined_at FROM chat_participants ORDER BY joined_at DESC LIMIT 10;" || true

echo "Done REST part. Now attempt WebSocket smoke (non-interactive)..."

# small delay to avoid race where WS is attempted before DB writes/participant insertion propagate
sleep 1

if command -v websocat >/dev/null 2>&1; then
  echo "websocat found; performing simple ws send"
  printf '{"type":"message","content":"hello from websocat"}\n' | websocat -1 "ws://localhost:8080/ws?token=$TOKEN&chat_id=$CHAT"
else
  echo "websocat not installed; skipping interactive ws step. You can run:"
  echo "  websocat 'ws://localhost:8080/ws?token=$TOKEN&chat_id=$CHAT'"
fi

# Try Go WS client if available (non-interactive). This will attempt to send a message
if [ -f "scripts/ws_smoke.go" ]; then
  echo "Running Go WS client..."
  if go version >/dev/null 2>&1; then
    go run scripts/ws_smoke.go -token "$TOKEN" -chat "$CHAT" || echo "Go WS client finished with non-fatal error"
  else
    echo "go not available; skipping Go WS client"
  fi
fi

echo "Check metrics after WS (filtered):"
curl -s $BASE/metrics | grep company_chat_messages_total -n || true

echo "Smoke finished"
