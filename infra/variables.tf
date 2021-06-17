variable "share_region" {
  type    = string
  default = "us-east-1"
}

variable "share_bucket-name" {
  type    = string
  default = "share.files"
}

variable "share_table-name" {
  type    = string
  default = "share.count"
}

variable "share_lambda-name" {
  type    = string
  default = "share.add"
}

variable "share_lambda-iam" {
  type    = string
  default = "share.add-role"
}

variable "share_lambda-filename" {
  type    = string
  default = "../build/share.zip"
}
