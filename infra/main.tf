terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.9"
    }
  }

  required_version = ">= 0.14.9"
}

// localstack
provider "aws" {
  region                      = var.share_region

  access_key                  = var.local ? "test" : null
  secret_key                  = var.local ? "test" : null
  s3_use_path_style           = var.local
  skip_credentials_validation = var.local
  skip_metadata_api_check     = var.local
  skip_requesting_account_id  = var.local

  // configure all services to use localstack url
  endpoints {
    dynamodb     = var.local ? "http://localhost:4566" : null
    iam          = var.local ? "http://localhost:4566" : null
    lambda       = var.local ? "http://localhost:4566" : null
    s3           = var.local ? "http://s3.localhost.localstack.cloud:4566" : null
  }
}

// s3: share url bucket
resource "aws_s3_bucket" "share_files" {
  bucket = var.share_files_name
}

resource "aws_s3_bucket_acl" "share_files" {
  bucket = aws_s3_bucket.share_files.id
  acl    = "public-read"
}

resource "aws_s3_bucket_website_configuration" "share_files" {
  bucket = aws_s3_bucket.share_files.bucket

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "index.html"
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

// lambda: share service
resource "aws_lambda_function" "share_add" {
  runtime          = "go1.x"
  function_name    = var.share_add_name
  handler          = var.share_add_binary
  filename         = var.share_add_archive
  source_code_hash = filebase64sha256(var.share_add_archive)
  role             = aws_iam_role.share_add.arn
}

resource "aws_lambda_function_url" "share_add" {
  function_name      = aws_lambda_function.share_add.arn
  authorization_type = "NONE"

  cors {
    allow_credentials = true
    allow_origins     = ["*"]
    allow_methods     = ["*"]
    allow_headers     = ["date", "keep-alive"]
    expose_headers    = ["keep-alive", "date"]
    max_age           = 86400
  }
}

output "share_add_url" {
  value = "${aws_lambda_function_url.share_add.function_url}"
}

resource "aws_iam_role" "share_add" {
  name = "${var.share_add_name}-iam"

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