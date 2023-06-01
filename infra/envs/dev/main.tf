// -- add provider --
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.9"
    }
  }

  required_version = ">= 0.14.9"
}

provider "aws" {
  region = var.aws_region

  // use local credentials
  access_key = var.local ? "test" : null
  secret_key = var.local ? "test" : null

  // use localstack-compatible apis
  s3_use_path_style           = var.local
  skip_credentials_validation = var.local
  skip_metadata_api_check     = var.local
  skip_requesting_account_id  = var.local

  // use localstack url for all services
  endpoints {
    dynamodb = var.local ? var.local_url : null
    iam      = var.local ? var.local_url : null
    lambda   = var.local ? var.local_url : null
    s3       = var.local ? var.local_url : null
  }
}

// -- add modules --
module "share_add" {
  source = "../../modules/share_add"

  // dev only
  local = var.local

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
  source        = "../../modules/share_files"
  name          = var.share_files_name
  force_destroy = true
}
