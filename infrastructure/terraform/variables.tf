variable "compute_flavor" {
  type        = string
  default     = "STD2-2-4"
  description = "Flavor name for the instance"
}

variable "availability_zone_name" {
  type        = string
  default     = "MS1"
  description = "Availability zone name"
}

variable "instance_name" {
  type        = string
  default     = "my-instance"
  description = "Name of the virtual machine"
}
