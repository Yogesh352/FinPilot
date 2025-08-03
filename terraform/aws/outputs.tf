output "ec2_public_ip" {
  value = aws_instance.api_service.public_ip
}

output "rds_endpoint" {
  value = aws_db_instance.postgres.endpoint
}