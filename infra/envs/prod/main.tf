// -- add provider --
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.9"
    }
  }

  backend "s3" {
    bucket = "mise--share-infra"
    key    = "prod/state"
    region = "us-east-1"
  }

  required_version = ">= 0.14.9"
}

provider "aws" {
  region = var.aws_region
}

// -- add modules --
module "share_add" {
  source = "../../modules/share_add"

  // local
  name    = var.share_add_name
  binary  = var.share_add_binary
  archive = var.share_add_archive

  // external
  share_count_name = var.share_count_name
  share_files_name = var.share_files_name
  share_files_host = var.share_files_host
}

module "share_count" {
  source = "../../modules/share_count"
  name   = var.share_count_name
}

module "share_files" {
  source = "../../modules/share_files_cdn"
  name   = var.share_files_name
}
