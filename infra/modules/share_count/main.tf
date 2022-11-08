// -- add kv store --
resource "aws_dynamodb_table" "share_count" {
  name           = var.name
  billing_mode   = "PROVISIONED"
  read_capacity  = 25
  write_capacity = 25
  hash_key       = "Id"

  attribute {
    name = "Id"
    type = "S"
  }
}
