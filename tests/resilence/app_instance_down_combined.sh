#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "${ROOT_DIR}"

set -a
[ -f .env ] && . ./.env
set +a

COMPOSE_BASE="-f docker-compose.base.yaml -f docker-compose.yaml"
BASE_URL="http://localhost:${SERVER_PORT:-8080}"

APP_COUNT="$(docker compose ${COMPOSE_BASE} ps -q app | wc -l | tr -d ' ')"
if [ "${APP_COUNT}" -lt 2 ]; then
  echo "[resilence] need >=2 app replicas to run this test (found ${APP_COUNT})" >&2
  exit 1
fi

APP_CID="$(docker compose ${COMPOSE_BASE} ps -q app | head -n1)"
echo "[resilence] stopping app ${APP_CID}..."
docker stop "${APP_CID}"
trap 'echo "[resilence] starting app back..."; docker start "${APP_CID}"' EXIT

# restart proxy to refresh upstream list
echo "[resilence] restarting proxy to refresh upstream..."
docker compose ${COMPOSE_BASE} restart proxy
sleep 3

echo "[resilence] waiting for proxy health on ${BASE_URL}/health ..."
for i in {1..30}; do
  if curl -sf "${BASE_URL}/health" >/dev/null; then
    break
  fi
  sleep 1
  if [ "$i" -eq 30 ]; then
    echo "[resilence] health check failed after stop" >&2
    exit 1
  fi
done

echo "[resilence] obtaining admin token..."
TOKEN_JSON=$(curl -s -X POST "${BASE_URL}/auth/token" -H "Content-Type: application/json" -d '{"username":"admin"}')
TOKEN=$(python3 - <<'PY'
import json,sys
try:
    print(json.loads(sys.stdin.read())["token"])
except Exception:
    sys.exit(1)
PY
<<<"${TOKEN_JSON}") || { echo "[resilence] failed to get token: ${TOKEN_JSON}" >&2; exit 1; }

hdr=(-H "token: ${TOKEN}" -H "Content-Type: application/json")

TEAM_ID="team-$(date +%s)"
AUTHOR_ID="author-$(date +%s)"
REV_ID="rev-$(date +%s)"
echo "[resilence] creating team..."
resp=$(curl -s -o /dev/null -w "%{http_code}" -X POST "${BASE_URL}/team/add" "${hdr[@]}" \
  -d "{\"team_name\":\"${TEAM_ID}\",\"members\":[{\"user_id\":\"${AUTHOR_ID}\",\"username\":\"Author\",\"is_active\":true},{\"user_id\":\"${REV_ID}\",\"username\":\"Rev\",\"is_active\":true}]}") || true
if [ "${resp}" != "201" ]; then
  echo "[resilence] team/add failed, status=${resp}" >&2
  exit 1
fi

PR_ID="pr-$(date +%s)"
echo "[resilence] creating PR..."
resp=$(curl -s -o /dev/null -w "%{http_code}" -X POST "${BASE_URL}/pullRequest/create" "${hdr[@]}" \
  -d "{\"pull_request_id\":\"${PR_ID}\",\"pull_request_name\":\"Feature\",\"author_id\":\"${AUTHOR_ID}\"}") || true
if [ "${resp}" != "201" ]; then
  echo "[resilence] pr/create failed, status=${resp}" >&2
  exit 1
fi

echo "[resilence] fetching statistics..."
resp=$(curl -s -o /dev/null -w "%{http_code}" -X GET "${BASE_URL}/statistics/assignments/pullrequests" "${hdr[@]}") || true
if [ "${resp}" != "200" ]; then
  echo "[resilence] stats pullrequests failed, status=${resp}" >&2
  exit 1
fi
resp=$(curl -s -o /dev/null -w "%{http_code}" -X GET "${BASE_URL}/statistics/assignments/users" "${hdr[@]}") || true
if [ "${resp}" != "200" ]; then
  echo "[resilence] stats users failed, status=${resp}" >&2
  exit 1
fi

echo "[resilence] running factorial e2e single..."
API_BASE_URL="${BASE_URL}" AUTH_KEY="${AUTH_KEY:-}" go test ./tests/end-to-end -run "FactorialEndToEnd" -count=1

echo "[resilence] success"

