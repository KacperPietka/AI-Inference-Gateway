terraform {
    required_version = ">= 1.7.0"
    
    required_providers {
        google = {
            source = "hashicorp/google"
            version = "~> 5.0"
        }
    }

    backend "gcs" {
        bucket = "tf-state-inference-gateway"
        prefix = "terraform/state"
    }
}