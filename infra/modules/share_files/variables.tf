variable "name" {
  description = "the name of the files service"
  type = string
}

variable "force_destroy" {
  description = "if nonempty buckets should be destroyed"
  type = bool
  default = false
}