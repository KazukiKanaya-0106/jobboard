data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

module "network" {
  source = "./modules/network"

  app_name        = var.app_name
  stage           = var.stage
  vpc_cidr        = var.network_vpc_cidr
  az_count        = var.network_az_count
  additional_tags = local.additional_tags
}

module "ecr" {
  source = "./modules/ecr"

  app_name                 = var.app_name
  stage                    = var.stage
  repository_name_override = var.ecr_repository_name
  additional_tags          = local.additional_tags
}

locals {
  api_container_image = coalesce(
    var.api_container_image,
    "${module.ecr.repository_url}:latest"
  )
  api_hosted_zone_id = coalesce(var.api_hosted_zone_id, var.hosted_zone_id)
  web_hosted_zone_id = coalesce(var.web_hosted_zone_id, var.hosted_zone_id)
}

module "api" {
  source = "./modules/api_service"

  app_name                   = var.app_name
  stage                      = var.stage
  vpc_id                     = module.network.vpc_id
  public_subnet_ids          = module.network.public_subnet_ids
  container_image            = local.api_container_image
  container_port             = var.api_container_port
  desired_count              = var.api_desired_count
  task_cpu                   = var.api_task_cpu
  task_memory                = var.api_task_memory
  environment                = var.api_environment
  assign_public_ip           = var.api_assign_public_ip
  health_check_path          = var.api_health_check_path
  load_balancer_idle_timeout = var.api_load_balancer_idle_timeout
  certificate_arn            = var.api_certificate_arn
  domain_name                = var.api_domain_name
  hosted_zone_id             = local.api_hosted_zone_id
  additional_tags            = local.additional_tags
}

module "database" {
  source = "./modules/aurora_serverless"

  app_name                     = var.app_name
  stage                        = var.stage
  vpc_id                       = module.network.vpc_id
  private_subnet_ids           = module.network.private_subnet_ids
  allowed_security_group_ids   = concat([module.api.service_security_group_id], var.aurora_additional_allowed_security_group_ids)
  database_name                = var.aurora_database_name
  master_username              = var.aurora_master_username
  master_password              = var.aurora_master_password
  serverless_min_capacity      = var.aurora_min_capacity
  serverless_max_capacity      = var.aurora_max_capacity
  backup_retention_period      = var.aurora_backup_retention_period
  preferred_backup_window      = var.aurora_preferred_backup_window
  preferred_maintenance_window = var.aurora_preferred_maintenance_window
  instance_count               = var.aurora_instance_count
  enable_data_api              = var.aurora_enable_data_api
  deletion_protection          = var.aurora_deletion_protection
  skip_final_snapshot          = var.aurora_skip_final_snapshot
  additional_tags              = local.additional_tags
}

module "web" {
  source = "./modules/web_spa"

  providers = {
    aws           = aws
    aws.us_east_1 = aws.us_east_1
  }

  app_name            = var.app_name
  stage               = var.stage
  domain_name         = var.web_domain_name
  hosted_zone_id      = local.web_hosted_zone_id
  bucket_name         = var.web_bucket_name
  acm_certificate_arn = var.web_certificate_arn
  default_root_object = var.web_default_root_object
  price_class         = var.web_price_class
  compress            = var.web_compress
  additional_tags     = local.additional_tags
}
