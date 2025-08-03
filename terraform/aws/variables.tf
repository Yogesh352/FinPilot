variable "key_name" {
  description = "Name of your AWS SSH key pair"
  type        = string
}

variable "db_password" {
  description = "Password for RDS Postgres"
  type        = string
  sensitive   = true
}

variable "ami_id" {
    description = "AMI ID for EC2 Instance"
    type = string
    default = "ami-08a6efd148b1f7504"
}

variable "ec2_instance_tpye" {
    description = "EC2 Instance"
    type = string
    default = "t3.micro"
}