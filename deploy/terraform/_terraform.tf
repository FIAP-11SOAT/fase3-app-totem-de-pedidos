terraform {
  required_version = ">= 1.11.0"

  backend "s3" {
    bucket = "fase3-terraform-state"
    key    = "fase3-app-totem-de-pedidos/terraform.tfstate"
    region = "us-east-1"
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 6.15"
    }
    github = {
      source  = "integrations/github"
      version = "~> 6.0"
    }
  }

}