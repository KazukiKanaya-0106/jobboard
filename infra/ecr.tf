locals {
  ecr_repository_name = (
    var.ecr_repository_name != null && trimspace(var.ecr_repository_name) != ""
  ) ? trimspace(var.ecr_repository_name) : "${var.app_name}-hub"
}

resource "aws_ecr_repository" "hub" {
  name                 = local.ecr_repository_name
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  encryption_configuration {
    encryption_type = "KMS"
  }
}
