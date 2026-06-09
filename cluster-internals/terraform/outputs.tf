output "control_plane_instance_id" {
  description = "Instance ID for SSM: aws ssm start-session --target <value>"
  value       = aws_instance.control_plane.id
}

output "control_plane_private_ip" {
  description = "Private IP for Ansible inventory"
  value       = aws_instance.control_plane.private_ip
}

output "worker_instance_id" {
  description = "Instance ID for SSM: aws ssm start-session --target <value>"
  value       = aws_instance.worker.id
}

output "worker_private_ip" {
  description = "Private IP for Ansible inventory"
  value       = aws_instance.worker.private_ip
}
