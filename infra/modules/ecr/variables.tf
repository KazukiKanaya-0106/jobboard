variable "app_name" {
  description = "Application identifier used for naming."
  type        = string
}

variable "stage" {
  description = "Deployment stage (e.g. dev, stg, prod)."
  type        = string
}

variable "repository_name_override" {
  description = "Optional override for the full ECR repository name."
  type        = string
  nullable    = true
  default     = null
}

variable "enable_image_scan" {
  description = "Whether to enable image scanning on push."
  type        = bool
  default     = true
}

variable "image_tag_mutability" {
  description = "Image tag mutability setting for the repository."
  type        = string
  default     = "MUTABLE"

  validation {
    condition     = contains(["MUTABLE", "IMMUTABLE"], var.image_tag_mutability)
    error_message = "image_tag_mutability must be either MUTABLE or IMMUTABLE."
  }
}

variable "encryption_type" {
  description = "Repository encryption type."
  type        = string
  default     = "KMS"

  validation {
    condition     = contains(["KMS", "AES256"], var.encryption_type)
    error_message = "encryption_type must be KMS or AES256."
  }
}

variable "kms_key_arn" {
  description = "Optional custom KMS key ARN for repository encryption."
  type        = string
  nullable    = true
  default     = null
}

variable "additional_tags" {
  description = "Extra tags to apply to the repository in addition to provider default tags."
  type        = map(string)
  default     = {}
}
