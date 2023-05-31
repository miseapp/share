// -- add modules --
module "share_files" {
  source = "../share_files"

  // variables
  name = var.name
}

// -- add bucket access control
resource "aws_s3_bucket_acl" "b_acl" {
  bucket = share_files.bucket_id
  acl    = "private"
}