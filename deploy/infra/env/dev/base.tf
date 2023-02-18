terraform {
  required_providers {
    kubernetes = {
      source = "hashicorp/kubernetes"
    }
	 helm = {
      source = "hashicorp/helm"
      version = "2.9.0"
    }
  }
}

provider "kubernetes" {
  config_path    = "~/.kube/config"
  config_context = local.k8_context
}

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config"
  }
}
