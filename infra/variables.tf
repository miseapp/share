// any variables with no default are sourced from .env
variable "share_region" {
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

variable "share_add_iam" {
  type    = string
  default = "share.add-role"
}
