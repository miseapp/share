output "bucket_id" {
  description = "the id of the website bucket"
  value = aws_s3_bucket.share_files.id
}

output "bucket_arn" {
  description = "the arn of the website bucket"
  value = aws_s3_bucket.share_files.arn
}