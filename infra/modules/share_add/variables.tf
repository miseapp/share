variable "local" {
  type = bool
  default = false
  description = "if the service is running locally"
}

variable "name" {
  type = string
  description = "the name of the add service"
}

variable "binary" {
  type = string
  description = "the path in the archive to the handler"
}

variable "archive" {
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