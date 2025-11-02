#!/usr/bin/env bash
# shellcheck disable=SC1090
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEFAULT_ENV_FILE="${SCRIPT_DIR}/../env/secret.env"

load_env_file() {
  local env_file="${1:-$DEFAULT_ENV_FILE}"
  if [[ ! -f "$env_file" ]]; then
    echo "Missing env file: $env_file" >&2
    exit 1
  fi

  set -a
  source "$env_file"
  set +a
}

require_env_vars() {
  local missing=0
  for name in "$@"; do
    if [[ -z "${!name:-}" ]]; then
      echo "Required environment variable '${name}' is not set." >&2
      missing=1
    fi
  done
  if [[ $missing -ne 0 ]]; then
    exit 1
  fi
}

ensure_command() {
  local cmd="$1"
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "Required command not found: $cmd" >&2
    exit 1
  fi
}
