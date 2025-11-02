#!/usr/bin/env bash
set -euo pipefail

THIS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${THIS_DIR}/../.." && pwd)"

source "${THIS_DIR}/../lib.sh"

load_env_file "${ENV_FILE:-}"
require_env_vars AWS_REGION AWS_PROFILE APP_NAME STAGE
ensure_command terraform

TF_DIR="${TF_DIR:-${ROOT_DIR}}"
TF_VARS_FILE="${TF_VARS_FILE:-${ROOT_DIR}/env/prod/terraform.tfvars}"

terraform -chdir="$TF_DIR" init -input=false

apply_args=(
  -input=false
  -auto-approve
  -target=module.ecr.aws_ecr_repository.this
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
