data "aws_caller_identity" "current" {}

module "ecr" {
  source = "./modules/ecr"

  app_name                 = var.app_name
  stage                    = var.stage
  repository_name_override = var.ecr_repository_name
  additional_tags          = local.additional_tags
}
