#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/lib.sh"

load_env_file "${ENV_FILE:-}"
require_env_vars AWS_REGION AWS_PROFILE APP_NAME STAGE
ensure_command terraform

TF_DIR="${TF_DIR:-${SCRIPT_DIR}/..}"
TF_VARS_FILE="${TF_VARS_FILE:-${SCRIPT_DIR}/../env/prod.tfvars}"

terraform -chdir="$TF_DIR" init -input=false

apply_args=(
  -input=false
  -auto-approve
  -target=aws_ecr_repository.hub
)

if [[ -f "$TF_VARS_FILE" ]]; then
  apply_args+=(-var-file="$TF_VARS_FILE")
fi

apply_args+=(
  -var "region=${AWS_REGION}"
  -var "aws_profile=${AWS_PROFILE}"
  -var "app_name=${APP_NAME}"
  -var "stage=${STAGE}"
)

terraform -chdir="$TF_DIR" apply "${apply_args[@]}"
