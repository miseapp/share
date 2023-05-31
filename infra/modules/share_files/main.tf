// -- add bucket --
resource "aws_s3_bucket" "share_files" {
  bucket = var.name
  force_destroy = var.force_destroy
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
