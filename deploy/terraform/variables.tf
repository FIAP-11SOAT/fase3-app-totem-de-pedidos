data "aws_ecr_authorization_token" "ecr_auth" {}

locals {
  aws_region   = "us-east-1"
  project_name = "fase3-totem-de-pedidos"
  eks_cluster_name = "fase3-totem-de-pedidos-eks-cluster"
}