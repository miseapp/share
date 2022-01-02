terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.27"
    }
  }

  required_version = ">= 0.14.9"
}

// configure localstack provider
provider "aws" {
  region                      = var.share_region
  s3_force_path_style         = true
  access_key                  = "test"
  secret_key                  = "test"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  // configure all services to use localstack url
  endpoints {
    dynamodb = "http://localhost:4566"
    lambda   = "http://localhost:4566"
    s3       = "http://localhost:4566"
    iam      = "http://localhost:4566"
  }
}

// s3: share url bucket
resource "aws_s3_bucket" "share_files" {
  bucket = var.share_files_name
  acl    = "public-read"

  website {
    index_document = "index.html"
    error_document = "index.html"
  }
}

resource "aws_s3_bucket_policy" "share_files" {
  bucket = aws_s3_bucket.share_files.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "PublicReadGetObject"
        Effect    = "Allow"
        Principal = "*"
        Action    = "s3:GetObject"
        Resource = [
          "${aws_s3_bucket.share_files.arn}",
          "${aws_s3_bucket.share_files.arn}/*",
        ]
      },
    ]
  })
}

// dynamo: share kv store
resource "aws_dynamodb_table" "share_count" {
  name           = var.share_count_name
  billing_mode   = "PROVISIONED"
  read_capacity  = 25
  write_capacity = 25
  hash_key       = "Id"

  attribute {
    name = "Id"
    type = "S"
  }
}

// lambda: share endpoint
resource "aws_lambda_function" "share_add" {
  runtime          = "go1.x"
  function_name    = var.share_add_name
  handler          = var.share_add_binary
  filename         = var.share_add_archive
  source_code_hash = filebase64sha256(var.share_add_archive)
  role             = aws_iam_role.share_add.arn
}

resource "aws_iam_role" "share_add" {
  name = var.share_add_iam

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}
