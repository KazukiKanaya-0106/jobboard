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
