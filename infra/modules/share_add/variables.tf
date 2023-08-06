variable "local" {
  description = "if the service is running locally"
  type        = bool
  default     = false
}

variable "name" {
  description = "the name of the add service"
  type        = string
}

variable "binary" {
  description = "the path in the archive to the handler"
  type        = string
}

variable "archive" {
  description = "the path in the filesystem to the archive"
  type        = string
}

variable "share_count_name" {
  description = "the name of the count service"
  type        = string
}

variable "share_count_table_arn" {
  description = "the arn of the count table"
  type        = string
}

variable "share_files_name" {
  description = "the name of the files service"
  type        = string
}

variable "share_files_host" {
  description = "the host url for the files"
  type        = string
}

variable "share_files_bucket_arn" {
  description = "the arn of the files bucket"
  type        = string
}
