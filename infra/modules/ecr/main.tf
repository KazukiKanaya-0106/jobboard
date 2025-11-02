locals {
  repository_name = (
    var.repository_name_override != null && trimspace(var.repository_name_override) != ""
  ) ? trimspace(var.repository_name_override) : "${var.app_name}-hub"
}

resource "aws_ecr_repository" "this" {
  name                 = local.repository_name
  image_tag_mutability = var.image_tag_mutability

  image_scanning_configuration {
    scan_on_push = var.enable_image_scan
  }

  dynamic "encryption_configuration" {
    for_each = var.encryption_type == "KMS" && var.kms_key_arn != null ? [1] : []

    content {
      encryption_type = var.encryption_type
      kms_key         = var.kms_key_arn
    }
  }

  dynamic "encryption_configuration" {
    for_each = var.encryption_type == "KMS" && var.kms_key_arn == null ? [1] : []

    content {
      encryption_type = var.encryption_type
    }
  }

  dynamic "encryption_configuration" {
    for_each = var.encryption_type == "AES256" ? [1] : []

    content {
      encryption_type = var.encryption_type
    }
  }

  tags = var.additional_tags
}
