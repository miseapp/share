// TODO: figure out how to share env between terraform, makefile,
// go test code
variable "share_region" {
  type    = string
  default = "us-east-1"
}

// TODO: figure out if "share.files" works
variable "share_files-name" {
  type    = string
  default = "share-files"
}

variable "share_count-name" {
  type    = string
  default = "share.count"
}

variable "share_add-name" {
  type    = string
  default = "share.add"
}

variable "share_add-iam" {
  type    = string
  default = "share.add-role"
}

variable "share_add-archive" {
  type    = string
  default = "../build/share.add.zip"
}
