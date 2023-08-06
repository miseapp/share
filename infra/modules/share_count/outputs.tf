output "table_arn" {
  description = "the arn of the count table"
  value       = aws_dynamodb_table.share_count.arn
}
