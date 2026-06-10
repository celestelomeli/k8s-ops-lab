output "cluster_name" {
  value = aws_eks_cluster.main.name
}

output "cluster_endpoint" {
  value = aws_eks_cluster.main.endpoint
}

output "kubeconfig_command" {
  description = "Run this to configure kubectl for the cluster"
  value       = "aws eks update-kubeconfig --region ${var.region} --name ${var.cluster_name}"
}

output "inventory_api_repository_url" {
  value = aws_ecr_repository.inventory_api.repository_url
}

output "notifier_repository_url" {
  value = aws_ecr_repository.notifier.repository_url
}

output "eso_role_arn" {
  description = "IAM role ARN for External Secrets Operator — use to annotate the ESO service account"
  value       = aws_iam_role.eso.arn
}
