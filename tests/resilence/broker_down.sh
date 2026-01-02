#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "${ROOT_DIR}"

# load env for AUTH_KEY / SERVER_PORT etc.
set -a
[ -f .env ] && . ./.env
set +a

COMPOSE_BASE="-f docker-compose.base.yaml -f docker-compose.yaml"

echo "[resilence] stopping kafka1..."
docker compose ${COMPOSE_BASE} stop kafka1

trap 'echo "[resilence] starting kafka1 back..."; docker compose ${COMPOSE_BASE} start kafka1' EXIT

echo "[resilence] running e2e..."
API_BASE_URL="http://localhost:${SERVER_PORT:-8080}" \
	AUTH_KEY="${AUTH_KEY:-}" \
	go test ./tests/end-to-end/... -count=1

echo "[resilence] starting kafka1 back..."
docker compose ${COMPOSE_BASE} start kafka1

echo "[resilence] done"

