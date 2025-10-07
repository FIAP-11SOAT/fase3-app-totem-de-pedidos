output "ecr_repository_url" {
  description = "ECR repository URL"
  value       = aws_ecr_repository.app.repository_url
}

output "ecr_repository_arn" {
  description = "ECR repository ARN"  
  value       = aws_ecr_repository.app.arn
}

output "eks_cluster_name" {
  description = "EKS cluster name"
  value       = local.eks_cluster_name
}

output "project_name" {
  description = "Project name"
  value       = local.project_name
}

output "docker_image_url" {
  description = "Full Docker image URL"
  value       = "${aws_ecr_repository.app.repository_url}:latest"
}