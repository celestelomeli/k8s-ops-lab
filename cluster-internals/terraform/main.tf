data "aws_vpc" "default" {
  default = true
}

data "aws_subnet" "target" {
  filter {
    name   = "availabilityZone"
    values = [var.availability_zone]
  }

  filter {
    name   = "default-for-az"
    values = ["true"]
  }
}

data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"] # Canonical

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }
}
