export ENV_FILE ?= ./env/secret.env
export TF_VARS_FILE ?= ./env/prod.tfvars

ssm:
	ENV_FILE=$(ENV_FILE) ./scripts/put_ssm.sh

ecr:
	TF_DIR=. ENV_FILE=$(ENV_FILE) TF_VARS_FILE=$(TF_VARS_FILE) ./scripts/create_ecr.sh

push:
	ENV_FILE=$(ENV_FILE) ECR_REPO=$(ECR_REPO) ./scripts/push_ecr.sh

deploy:
	TF_DIR=. ENV_FILE=$(ENV_FILE) TF_VARS_FILE=$(TF_VARS_FILE) ./scripts/deploy.sh

all: ecr push deploy
