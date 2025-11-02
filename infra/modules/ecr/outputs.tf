output "id" {
  description = "Repository identifier."
  value       = aws_ecr_repository.this.id
}

output "arn" {
  description = "Repository ARN."
  value       = aws_ecr_repository.this.arn
}

output "name" {
  description = "Repository name."
  value       = aws_ecr_repository.this.name
}

output "repository_url" {
  description = "Repository URL for Docker pushes."
  value       = aws_ecr_repository.this.repository_url
}
