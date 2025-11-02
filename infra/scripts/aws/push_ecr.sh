#!/usr/bin/env bash
set -euo pipefail

THIS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${THIS_DIR}/../lib.sh"

load_env_file "${ENV_FILE:-}"
require_env_vars AWS_REGION AWS_PROFILE APP_NAME STAGE
ensure_command aws
ensure_command docker

ECR_REPO="${ECR_REPO:-${APP_NAME}-hub}"
DOCKERFILE="${DOCKERFILE:-Dockerfile}"
IMAGE_CONTEXT="${IMAGE_CONTEXT:-.}"
IMAGE_TAG="${IMAGE_TAG:-latest}"

ACCOUNT_ID="$(aws sts get-caller-identity --query Account --output text --profile "$AWS_PROFILE")"
ECR_REGISTRY="${ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com"
ECR_URI="${ECR_REGISTRY}/${ECR_REPO}"

aws ecr get-login-password --region "$AWS_REGION" --profile "$AWS_PROFILE" \
  | docker login --username AWS --password-stdin "$ECR_REGISTRY"

docker build -f "$DOCKERFILE" -t "${ECR_REPO}:${IMAGE_TAG}" "$IMAGE_CONTEXT"
docker tag "${ECR_REPO}:${IMAGE_TAG}" "${ECR_URI}:${IMAGE_TAG}"
docker push "${ECR_URI}:${IMAGE_TAG}"

echo "pushed: ${ECR_URI}:${IMAGE_TAG}"
