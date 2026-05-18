output "cluster_name" {
    description = "GKE cluster name"
    value = google_container_cluster.gateway.name
}

output "cluster_endpoint" {
    description = "GKE cluster endpoint"
    value = google_container_cluster.gateway.endpoint
    sensitive = true
}

output "registry_url" {
    description = "Artifact Registry URL for Docker images"
    value = "${var.region}-docker.pkg.dev/${var.project_id}/${var.registry_name}"
}

output "kubectl_command" {
    description = "Command to connect kubectl to this cluster"
    value = "gcloud container clusters get-credentials ${var.cluster_name} --zone ${var.zone} --project ${var.project_id}"
}