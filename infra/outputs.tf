output "whoami" {
  description = "Current AWS identity details."
  value = {
    account_id = data.aws_caller_identity.current.account_id
    arn        = data.aws_caller_identity.current.arn
    user_id    = data.aws_caller_identity.current.user_id
  }
}

output "hub_repository" {
  description = "Information about the Hub ECR repository."
  value = {
    name = module.ecr.name
    url  = module.ecr.repository_url
    arn  = module.ecr.arn
  }
}

output "network" {
  description = "Networking resources."
  value = {
    vpc_id             = module.network.vpc_id
    public_subnet_ids  = module.network.public_subnet_ids
    private_subnet_ids = module.network.private_subnet_ids
  }
}

output "api" {
  description = "ECS service and ALB details."
  value = {
    cluster_id             = module.api.cluster_id
    service_name           = module.api.service_name
    target_group_arn       = module.api.target_group_arn
    alb_dns_name           = module.api.alb_dns_name
    alb_security_group_id  = module.api.alb_security_group_id
    service_security_group = module.api.service_security_group_id
    certificate_arn        = module.api.certificate_arn
  }
}

output "database" {
  description = "Aurora Serverless cluster endpoints."
  value = {
    cluster_id      = module.database.cluster_id
    cluster_arn     = module.database.cluster_arn
    writer_endpoint = module.database.endpoint
    reader_endpoint = module.database.reader_endpoint
    security_group  = module.database.security_group_id
  }
  sensitive = true
}

output "web" {
  description = "Static web distribution details."
  value = {
    bucket_name         = module.web.bucket_name
    distribution_id     = module.web.distribution_id
    distribution_domain = module.web.distribution_domain_name
    certificate_arn     = module.web.certificate_arn
  }
}
