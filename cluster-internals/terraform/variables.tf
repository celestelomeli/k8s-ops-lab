variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "availability_zone" {
  description = "AZ to place both nodes in" 
  type        = string
  default     = "us-east-1a"
}

variable "my_ip" {
  description = "Public IP for kubectl access, without the /32"
  type        = string
  sensitive   = true
}
