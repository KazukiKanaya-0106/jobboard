#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${SCRIPT_DIR}/lib.sh"

load_env_file "${ENV_FILE:-}"
require_env_vars AWS_REGION AWS_PROFILE APP_NAME STAGE DB_PASSWORD AUTH_JWT_SECRET
ensure_command aws

put_secure() {
  local name="$1"
  local value="$2"

  echo "Putting secure parameter: $name"
  aws ssm put-parameter \
    --name "$name" \
    --value "$value" \
    --type "SecureString" \
    --overwrite \
    --region "$AWS_REGION" \
    --profile "$AWS_PROFILE"
}

declare -A SECURE_PARAMS=(
  ["DB_PASSWORD"]="$DB_PASSWORD"
  ["AUTH_JWT_SECRET"]="$AUTH_JWT_SECRET"
)

echo "Registering SSM parameters... (region=$AWS_REGION)"
echo "App $APP_NAME | Stage $STAGE | Profile $AWS_PROFILE"

for key in "${!SECURE_PARAMS[@]}"; do
  put_secure "/${APP_NAME}/${STAGE}/${key}" "${SECURE_PARAMS[$key]}"
done

echo "All parameters registered successfully."
