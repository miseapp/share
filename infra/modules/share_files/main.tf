// -- add bucket --
resource "aws_s3_bucket" "share_files" {
  bucket = var.name
}

// -- add bucket web hosting
resource "aws_s3_bucket_website_configuration" "share_files" {
  bucket = aws_s3_bucket.share_files.bucket

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "index.html"
  }
}

// -- add bucket policy
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