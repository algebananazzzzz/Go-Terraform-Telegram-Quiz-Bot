locals {
  security_groups = ["sg-allownat"]
  subnet_names    = ["mysubnet"]
}

data "aws_subnets" "private" {
  filter {
    name   = "tag:Name"
    values = local.subnet_names
  }
}

data "aws_security_groups" "allow_nat" {
  filter {
    name   = "tag:Name"
    values = local.security_groups
  }
}
