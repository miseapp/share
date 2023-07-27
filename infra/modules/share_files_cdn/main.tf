// -- add modules --
module "share_files" {
  source = "../share_files"
  name   = var.name
}

// -- locals --
locals {
  origin_id = "${var.name}--cdn-origin"
}

// -- add cdn --
resource "aws_cloudfront_distribution" "share_files" {
  # aliases = ["share.miseapp.co"]

  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"
  price_class         = "PriceClass_100"

  origin {
    domain_name              = module.share_files.bucket_regional_domain_name
    origin_id                = local.origin_id
    origin_access_control_id = aws_cloudfront_origin_access_control.share_files.id
  }

  default_cache_behavior {
    target_origin_id = local.origin_id
    allowed_methods  = ["HEAD", "GET"]
    cached_methods   = ["HEAD", "GET"]
    cache_policy_id  = aws_cloudfront_cache_policy.share_files.id

    viewer_protocol_policy = "allow-all"
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }
}

resource "aws_cloudfront_origin_access_control" "share_files" {
  name                              = "${var.name}--cdn-origin-acl"
  description                       = "cdn access control to the shared files bucket"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

// https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/Expiration.html#ExpirationDownloadDist
resource "aws_cloudfront_cache_policy" "share_files" {
  name    = "${var.name}--cdn-caching"
  comment = "caches for 1 year (the maximum)"

  min_ttl     = 31536000
  max_ttl     = 31536000
  default_ttl = 31536000

  parameters_in_cache_key_and_forwarded_to_origin {
    cookies_config {
      cookie_behavior = "all"
    }

    headers_config {
      header_behavior = "none"
    }

    query_strings_config {
      query_string_behavior = "all"
    }
  }
}

// -- add cdn -> bucket access control
resource "aws_s3_bucket_public_access_block" "share_files" {
  bucket                  = module.share_files.bucket_id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_policy" "share_files" {
  bucket = module.share_files.bucket_id
  policy = jsonencode({
    Version : "2012-10-17",
    Statement : {
      Sid : "AllowCloudFrontServicePrincipalRead",
      Effect : "Allow",
      Principal : {
        Service : "cloudfront.amazonaws.com"
      },
      Action : [
        "s3:GetObject",
      ],
      Resource : [
        "${module.share_files.bucket_arn}/*"
      ],
      Condition : {
        StringEquals : {
          "AWS:SourceArn" : aws_cloudfront_distribution.share_files.arn
        }
      }
    }
  })
}
