// any variables with no default are sourced from .env
variable "local" {
  type = bool
}

variable "local_url" {
  type = string
}

variable "aws_region" {
  type = string
}

variable "share_files_name" {
  type = string
}

variable "share_count_name" {
  type = string
}

variable "share_add_name" {
  type = string
}

variable "share_add_binary" {
  type = string
}

variable "share_add_archive" {
  type = string
}