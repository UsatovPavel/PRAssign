#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "${ROOT_DIR}"

# load env
set -a
[ -f .env ] && . ./.env
set +a

COMPOSE_BASE="-f docker-compose.base.yaml -f docker-compose.yaml"
APP_CID="$(docker compose ${COMPOSE_BASE} ps -q app | head -n1)"
if [ -z "${APP_CID}" ]; then
  echo "[resilence] no app container found" >&2
  exit 1
fi

echo "[resilence] stopping app ${APP_CID}..."
docker stop "${APP_CID}"
trap 'echo "[resilence] starting app back..."; docker start "${APP_CID}"' EXIT

echo "[resilence] running factorial e2e (single)..."
API_BASE_URL="http://localhost:${SERVER_PORT:-8080}" \
	AUTH_KEY="${AUTH_KEY:-}" \
	go test ./tests/end-to-end -run "FactorialEndToEnd" -count=1

echo "[resilence] ok"

