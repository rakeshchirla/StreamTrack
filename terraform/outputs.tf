output "cluster_name" {
  value = google_container_cluster.primary.name
}

output "kubeconfig" {
  sensitive = true
  value = <<EOT
gcloud container clusters get-credentials ${google_container_cluster.primary.name} --region ${var.gcp_region}
EOT
}