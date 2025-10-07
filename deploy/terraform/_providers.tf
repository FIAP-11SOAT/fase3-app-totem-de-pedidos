provider "aws" {
  region = local.aws_region

  default_tags {
    tags = {
      Project   = local.project_name
      Terraform = "true"
    }
  }

}

provider "github" {
  token = local.aws_master_secrets["GITHUB_ACCESS_TOKEN"]
  owner = local.aws_master_secrets["GITHUB_ORG"]
}

provider "kubernetes" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority[0].data)
  token                  = data.aws_eks_cluster_auth.cluster.token
}

data "aws_eks_cluster" "cluster" {
  name = local.eks_cluster_name
}

data "aws_eks_cluster_auth" "cluster" {
  name = data.aws_eks_cluster.cluster.name
}

