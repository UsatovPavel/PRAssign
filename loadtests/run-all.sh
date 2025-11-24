#!/bin/sh
# portable, explicit and debuggable runner for k6 scripts

set -u

REPORT_DIR="/tests/reports"
mkdir -p "$REPORT_DIR"

failures=0

run_k6() {
  name="$1"
  script="$2"
  logfile="$REPORT_DIR/$(echo "$name" | tr ' /' '__').log"
  echo "=== $name ==="
  echo "log -> $logfile"
  # запуск с логированием в файл (stdout+stderr)
  k6 run "$script" > "$logfile" 2>&1
  ret=$?
  echo "[$name] exit code: $ret"
  if [ "$ret" -ne 0 ]; then
    failures=$((failures+1))
    echo "[$name] FAILED (total failures now: $failures)"
  else
    echo "[$name] OK"
  fi
  echo
  return $ret
}

run_k6 "wait for app" scripts/wait_app.js
run_k6 "INIT" scripts/init.js
run_k6 "CHECK API" scripts/checkapi.js
run_k6 "BASELINE" scripts/baseline.js

echo "=== SUMMARY ==="
if [ "$failures" -ne 0 ]; then
  echo "Some tests failed: $failures"
  echo "Detailed logs: $REPORT_DIR"
  exit 1
fi

echo "All tests passed"
exit 0
