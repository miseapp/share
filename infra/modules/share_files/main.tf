// -- add bucket --
resource "aws_s3_bucket" "share_files" {
  bucket        = var.name
  force_destroy = var.force_destroy
}
