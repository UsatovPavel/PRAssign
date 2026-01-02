#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "${ROOT_DIR}"

# load env
set -a
[ -f .env ] && . ./.env
set +a

BASE_URL="http://localhost:${SERVER_PORT:-8080}"
COMPOSE_BASE="-f docker-compose.base.yaml -f docker-compose.yaml"

# pick any running app container
APP_CID="$(docker compose ${COMPOSE_BASE} ps -q app | head -n1)"
if [ -z "${APP_CID}" ]; then
  echo "[resilence] no app container found" >&2
  exit 1
fi

echo "[resilence] stopping app ${APP_CID}..."
docker stop "${APP_CID}"
trap 'echo "[resilence] starting app back..."; docker start "${APP_CID}"' EXIT

echo "[resilence] hitting /health x3..."
for i in 1 2 3; do
  curl -sf "${BASE_URL}/health" >/dev/null
done

echo "[resilence] hitting /auth/token x3..."
for i in 1 2 3; do
  token_resp=$(curl -s -X POST "${BASE_URL}/auth/token" \
    -H "Content-Type: application/json" \
    -d '{"username":"resilence-user"}')
  echo "${token_resp}" | grep -q "token" || { echo "auth/token failed: ${token_resp}"; exit 1; }
done

echo "[resilence] ok"

