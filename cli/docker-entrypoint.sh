#!/bin/sh
set -euo pipefail

echo "[cli] building jobboard binary..."
mkdir -p /app/bin
CGO_ENABLED=0 go build -o /app/bin/jobboard ./cmd/jobboard
echo "[cli] build complete -> bin/jobboard"

exec tail -f /dev/null
