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

