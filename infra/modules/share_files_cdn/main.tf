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
    origin_id                = local.origin_id
    origin_access_control_id = aws_cloudfront_origin_access_control.share_files.id
    domain_name              = module.share_files.endpoint
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
  name                              = "${var.name}--cdn-oac"
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

// -- add bucket access control
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
          "AWS:SourceArn" : "arn:aws:cloudfront::<AWS account ID>:distribution/<CloudFront distribution ID>"
        }
      }
    }
  })
}
