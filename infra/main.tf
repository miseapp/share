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
  access_key                  = "mock_access_key"
  region                      = var.share_region
  s3_force_path_style         = true
  secret_key                  = "mock_secret_key"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  // configure all services to use localstack url
  endpoints {
    dynamodb = "http://localhost:4566"
    lambda   = "http://localhost:4566"
    s3       = "http://localhost:4566"
  }
}

// s3: share url bucket
resource "aws_s3_bucket" "share" {
  bucket = var.share_bucket-name
  acl    = "public-read"

  website {
    index_document = "index.html"
    error_document = "index.html"
  }
}

resource "aws_s3_bucket_policy" "share" {
  bucket = aws_s3_bucket.share.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "PublicReadGetObject"
        Effect    = "Allow"
        Principal = "*"
        Action    = "s3:GetObject"
        Resource = [
          "${aws_s3_bucket.share.arn}",
          "${aws_s3_bucket.share.arn}/*",
        ]
      },
    ]
  })
}
