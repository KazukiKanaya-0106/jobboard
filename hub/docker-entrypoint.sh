#!/bin/sh
set -euo pipefail

if [ "${RUN_MIGRATIONS:-true}" != "false" ]; then
  echo "[hub] running database migrations..."
  go run ./cmd/migrate --cmd up
fi

exec "$@"
