output "bucket_id" {
  description = "the id of the bucket"
  value       = aws_s3_bucket.share_files.id
}

output "bucket_arn" {
  description = "the arn of the bucket"
  value       = aws_s3_bucket.share_files.arn
}

output "bucket_regional_domain_name" {
  description = "the regional domain name of the bucket"
  value       = aws_s3_bucket.share_files.bucket_regional_domain_name
}
