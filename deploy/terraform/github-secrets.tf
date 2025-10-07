resource "github_actions_variable" "aws_ecr_auth_endpoint" {
  repository    = local.project_name
  value         = local.aws_ecr_auth_proxy_endpoint
  variable_name = "AWS_ECR_AUTH_ENDPOINT"
}

resource "github_actions_variable" "aws_ecr_repository_url" {
  repository    = local.project_name
  value         = aws_ecr_repository.app.repository_url
  variable_name = "AWS_ECR_REPOSITORY_URL"
}

resource "github_actions_variable" "aws_lambda_function_name" {
  repository    = local.project_name
  value         = local.eks_cluster_name
  variable_name = "AWS_EKS_CLUSTER_NAME"
}