variable "project_id" {
    description = "GCP project ID"
    type = string
    default = "project-192493f2-913e-4971-86c"
}

variable "region" {
    description = "GCP region"
    type = string
    default = "europe-west1"
}

variable "zone" {
    description = "GCP zone"
    type = string
    default = "europe-west1-b"
}

variable "cluster_name" {
    description = "GKE cluster name"
    type = string
    default = "kuar-cluster"
}

variable "registry_name" {
    description = "Artifact Registry respository name"
    type = string
    default = "inference-gateway"
}

variable "node_count" {
    description = "Number of nodes in the cluster"
    type = number
    default = 3
}

variable "machine_type" {
    description = "GKE node machine type"
    type = string
    default = "e2-medium"
}

variable "project_number" {
  description = "GCP project number"
  type = string
  default = "626881549611"
}