// any variables with no default are sourced from .env
variable "aws_region" {
  description = "the aws region"
  type = string
}

variable "share_add_name" {
  description = "the name of the add service"
  type = string
}

variable "share_add_binary" {
  description = "the path in the archive to the handler"
  type = string
}

variable "share_add_archive" {
  description = "the path in the filesystem to the archive"
  type = string
}

variable "share_count_name" {
  description = "the name of the count service"
  type = string
}

variable "share_files_name" {
  description = "the name of the files service"
  type = string
}

variable "share_files_host" {
  description = "the host url for the files"
  type = string
}