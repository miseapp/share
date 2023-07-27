// -- add modules --
module "share_files" {
  source        = "../share_files"
  name          = var.name
  force_destroy = true
}

// -- add bucket web hosting
resource "aws_s3_bucket_website_configuration" "share_files" {
  bucket = module.share_files.bucket_id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "index.html"
  }
}

// -- add bucket public read
resource "aws_s3_bucket_policy" "share_files" {
  bucket = module.share_files.bucket_id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid       = "PublicReadGetObject"
        Effect    = "Allow"
        Principal = "*"
        Action    = "s3:GetObject"
        Resource = [
          "${module.share_files.bucket_arn}/*",
        ]
      },
    ]
  })
}
