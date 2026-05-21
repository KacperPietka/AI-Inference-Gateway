variable "project_id" {
    description = "GCP project ID"
    type = string
}

variable "region" {
    description = "GCP region"
    type = string
}

variable "zone" {
    description = "GCP zone"
    type = string
}

variable "cluster_name" {
    description = "GKE cluster name"
    type = string
}

variable "registry_name" {
    description = "Artifact Registry respository name"
    type = string
}

variable "node_count" {
    description = "Number of nodes in the cluster"
    type = number
}

variable "machine_type" {
    description = "GKE node machine type"
    type = string
}

variable "project_number" {
  description = "GCP project number"
  type = string
}