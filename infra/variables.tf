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

variable "hosted_zone_id" {
  description = "Primary Route53 hosted zone ID for application domains."
  type        = string
}

variable "api_hosted_zone_id" {
  description = "Optional override hosted zone ID for the API domain."
  type        = string
  nullable    = true
  default     = null
}

variable "web_hosted_zone_id" {
  description = "Optional override hosted zone ID for the web domain."
  type        = string
  nullable    = true
  default     = null
}

variable "web_domain_name" {
  description = "Fully-qualified domain name for the SPA."
  type        = string
}

variable "api_domain_name" {
  description = "Fully-qualified domain name for the API."
  type        = string
}

variable "web_certificate_arn" {
  description = "Optional pre-existing ACM certificate ARN (us-east-1) for the SPA."
  type        = string
  nullable    = true
  default     = null
}

variable "api_certificate_arn" {
  description = "Optional pre-existing ACM certificate ARN for the API."
  type        = string
  nullable    = true
  default     = null
}

variable "web_bucket_name" {
  description = "Optional explicit S3 bucket name for the SPA."
  type        = string
  nullable    = true
  default     = null
}

variable "web_default_root_object" {
  description = "Default root object for the SPA distribution."
  type        = string
  default     = "index.html"
}

variable "web_price_class" {
  description = "CloudFront price class for the SPA distribution."
  type        = string
  default     = "PriceClass_100"
}

variable "web_compress" {
  description = "Enable compression for the SPA distribution."
  type        = bool
  default     = true
}

variable "api_container_image" {
  description = "ECR image (including tag) to deploy for the API. Defaults to the hub repository's latest tag."
  type        = string
  nullable    = true
  default     = null
}

variable "api_container_port" {
  description = "Container port exposed by the API service."
  type        = number
  default     = 8080
}

variable "api_desired_count" {
  description = "Desired count for the ECS service."
  type        = number
  default     = 1
}

variable "api_task_cpu" {
  description = "Fargate CPU units for the API task."
  type        = number
  default     = 256
}

variable "api_task_memory" {
  description = "Fargate memory (MiB) for the API task."
  type        = number
  default     = 512
}

variable "api_environment" {
  description = "Environment variables passed to the API container."
  type        = map(string)
  default     = {}
}

variable "api_assign_public_ip" {
  description = "Assign public IPs to ECS tasks."
  type        = bool
  default     = true
}

variable "api_health_check_path" {
  description = "ALB health check path."
  type        = string
  default     = "/health"
}

variable "api_load_balancer_idle_timeout" {
  description = "ALB idle timeout in seconds."
  type        = number
  default     = 60
}

variable "aurora_database_name" {
  description = "Initial database name."
  type        = string
  default     = "jobboard"
}

variable "aurora_master_username" {
  description = "Master username for the Aurora cluster."
  type        = string
  default     = "jobboard"
}

variable "aurora_master_password" {
  description = "Master password for the Aurora cluster."
  type        = string
  sensitive   = true
}

variable "aurora_min_capacity" {
  description = "Minimum ACUs for Aurora Serverless v2 (AWS currently enforces >= 0.5 ACU)."
  type        = number
  default     = 0.5

  validation {
    condition     = var.aurora_min_capacity >= 0.5
    error_message = "Aurora Serverless v2 minimum capacity must be at least 0.5 ACU."
  }
}

variable "aurora_max_capacity" {
  description = "Maximum ACUs for Aurora Serverless v2."
  type        = number
  default     = 2

  validation {
    condition     = var.aurora_max_capacity >= var.aurora_min_capacity
    error_message = "Aurora Serverless v2 max capacity must be greater than or equal to min capacity."
  }
}

variable "aurora_backup_retention_period" {
  description = "Automated backup retention in days."
  type        = number
  default     = 1
}

variable "aurora_preferred_backup_window" {
  description = "Preferred backup window."
  type        = string
  default     = "04:00-06:00"
}

variable "aurora_preferred_maintenance_window" {
  description = "Preferred maintenance window."
  type        = string
  default     = "sun:06:00-sun:07:00"
}

variable "aurora_instance_count" {
  description = "Number of Aurora instances (serverless v2 requires at least one)."
  type        = number
  default     = 1
}

variable "aurora_enable_data_api" {
  description = "Enable the Aurora Data API endpoint."
  type        = bool
  default     = true
}

variable "aurora_additional_allowed_security_group_ids" {
  description = "Additional security groups allowed to access the Aurora cluster."
  type        = list(string)
  default     = []
}

variable "aurora_deletion_protection" {
  description = "Enable deletion protection for the Aurora cluster."
  type        = bool
  default     = false
}

variable "aurora_skip_final_snapshot" {
  description = "Skip final snapshot on destroy."
  type        = bool
  default     = true
}

variable "network_vpc_cidr" {
  description = "CIDR block for the application VPC."
  type        = string
  default     = "10.0.0.0/16"
}

variable "network_az_count" {
  description = "Number of AZs to span."
  type        = number
  default     = 2
}
