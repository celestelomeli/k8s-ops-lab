terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket         = "k8s-ops-lab-tfstate"
    key            = "cluster-internals/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "k8s-ops-lab-state-lock"
    encrypt        = true
  }
}

provider "aws" {
  region = var.region
}
