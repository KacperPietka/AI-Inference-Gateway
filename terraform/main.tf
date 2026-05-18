provider "google" {
    project = var.project_id
    region = var.region
    zone = var.zone
}

resource "google_artifact_registry_repository" "gateway" {
    repository_id = var.registry_name
    format = "DOCKER"
    location = var.region
    description = "Inference Gateway Docker images"

    cleanup_policies {
        id = "keep-minimum-versions"
        action = "KEEP"
        most_recent_versions {
          keep_count = 10
        }
    }
}

# IAM -> allow GKE to pull from Artifact Registry
resource "google_project_iam_member" "gke_artifact_registry" {
  project = var.project_id
  role    = "roles/artifactregistry.reader"
  member  = "serviceAccount:${var.project_number}-compute@developer.gserviceaccount.com"
}

resource "google_project_iam_member" "gke_node_service_account" {
  project = var.project_id
  role = "roles/container.defaultNodeServiceAccount"
  member = "serviceAccount:${var.project_number}-compute@developer.gserviceaccount.com"
}

# GKE Cluster
resource "google_container_cluster" "gateway" {
  name = var.cluster_name
  location = var.zone

  remove_default_node_pool = true
  initial_node_count = 1
  lifecycle {
    ignore_changes = [
      node_config,
      node_pool,
      initial_node_count,
    ]
  }

  workload_identity_config {
    workload_pool = "${var.project_id}.svc.id.goog"
  }
}

# GKE Node Pool
resource "google_container_node_pool" "gateway_nodes" {
  name = "default-pool"
  cluster = google_container_cluster.gateway.name
  location = var.zone
  node_count = var.node_count

  node_config {
    machine_type = var.machine_type

    oauth_scopes = [
      "https://www.googleapis.com/auth/devstorage.read_only",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/service.management.readonly",
      "https://www.googleapis.com/auth/servicecontrol",
      "https://www.googleapis.com/auth/trace.append",
    ]

    labels = {
      app = "inference-gateway"
    }
  }
  lifecycle {
    ignore_changes = [
      node_config[0].metadata,
      node_config[0].resource_labels,
      node_config[0].kubelet_config,
      node_config[0].shielded_instance_config,
    ]
  }

  management {
    auto_repair = true
    auto_upgrade = true
  }
}