#!/bin/sh
failures=0

echo "=== wait for app ==="
k6 run scripts/wait_app.js

echo "=== INIT ==="
k6 run scripts/init.js

echo "=== CHECK API ==="
k6 run scripts/checkapi.js

echo "=== BASELINE ==="
k6 run scripts/baseline.js


echo "=== BASIC TESTS COMPLETED SUCCESSFULLY ==="
