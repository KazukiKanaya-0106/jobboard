data "aws_caller_identity" "current" {}

output "whoami" {
  value = {
    account_id = data.aws_caller_identity.current.account_id
    arn        = data.aws_caller_identity.current.arn
    user_id    = data.aws_caller_identity.current.user_id
  }
}

output "hub_repository" {
  value = {
    name = aws_ecr_repository.hub.name
    url  = aws_ecr_repository.hub.repository_url
  }
}
