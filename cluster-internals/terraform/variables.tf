variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "availability_zone" {
  description = "AZ to place both nodes in — best practice is a custom VPC with explicit subnets, but this is sufficient for a default VPC lab"
  type        = string
  default     = "us-east-1a"
}

variable "my_ip" {
  description = "Your public IP for kubectl access, without the /32 (e.g. 203.0.113.5)"
  type        = string
  sensitive   = true
}
