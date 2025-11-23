#!/bin/sh
set -e

echo "=== wait for app ==="
k6 run scripts/wait_app.js

echo "=== INIT ==="
k6 run scripts/init.js

echo "=== CHECK API ==="
k6 run scripts/checkapi.js

echo "=== BASELINE ==="
k6 run scripts/baseline.js

echo "=== STRESS ==="
k6 run scripts/stress.js

echo "=== PEAK ==="
k6 run scripts/peak.js

echo "=== ENDURANCE ==="
k6 run scripts/endurance.js

echo "=== ALL TESTS COMPLETED SUCCESSFULLY ==="
