variable "region" {
  type = string
}
variable "aws_profile" {
  type = string
}
variable "app_name" {
  type = string
}
variable "stage" {
  type = string
}
variable "ecr_repository_name" {
  type        = string
  default     = null
  description = "Optional override for the ECR repository name."
}
