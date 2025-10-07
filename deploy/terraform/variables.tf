data "aws_ecr_authorization_token" "ecr_auth" {}

data "aws_secretsmanager_secret" "master_secrets" {
  name = "terraform-master-credentials"
}

data "aws_secretsmanager_secret_version" "master_secrets" {
  secret_id = data.aws_secretsmanager_secret.master_secrets.id
}

locals {
  aws_master_secrets = jsondecode(data.aws_secretsmanager_secret_version.master_secrets.secret_string)
  aws_ecr_auth_proxy_endpoint = replace(data.aws_ecr_authorization_token.ecr_auth.proxy_endpoint, "https://", "")
}


locals {
  aws_region   = "us-east-1"
  project_name = "fase3-app-totem-de-pedidos"
  eks_cluster_name = "fase3-infra-totem-de-pedidos-eks-cluster"
}