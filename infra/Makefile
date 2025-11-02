export ENV_FILE ?= ./env/secret.env
export TF_VARS_FILE ?= ./env/prod/terraform.tfvars

ssm:
	ENV_FILE=$(ENV_FILE) ./scripts/aws/put_ssm.sh

ecr:
	TF_DIR=. ENV_FILE=$(ENV_FILE) TF_VARS_FILE=$(TF_VARS_FILE) ./scripts/terraform/create_ecr.sh

push:
	ENV_FILE=$(ENV_FILE) ECR_REPO=$(ECR_REPO) ./scripts/aws/push_ecr.sh

deploy:
	TF_DIR=. ENV_FILE=$(ENV_FILE) TF_VARS_FILE=$(TF_VARS_FILE) ./scripts/terraform/deploy.sh

all: ecr push deploy
