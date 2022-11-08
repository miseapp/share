// any variables with no default are sourced from .env
variable "local" {
  type = bool
  description = "if the service is running locally"
}

variable "local_url" {
  type = string
  description = "the url for the local aws container"
}

variable "aws_region" {
  type = string
  description = "the aws region"
}

variable "share_add_name" {
  type = string
  description = "the name of the add service"
}

variable "share_add_binary" {
  type = string
  description = "the path in the archive to the handler"
}

variable "share_add_archive" {
  type = string
  description = "the path in the filesystem to the archive"
}

variable "share_count_name" {
  type = string
  description = "the name of the count service"
}

variable "share_files_name" {
  type = string
  description = "the name of the files service"
}

variable "share_files_host" {
  type = string
  description = "the host url for the files"
}