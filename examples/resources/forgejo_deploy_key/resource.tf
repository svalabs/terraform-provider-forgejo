terraform {
  required_providers {
    forgejo = {
      source = "svalabs/forgejo"
    }
    tls = {
      source = "hashicorp/tls"
    }
  }
}

provider "forgejo" {
  host = "http://localhost:3000"
}

# Private key
resource "tls_private_key" "ed25519" {
  algorithm = "ED25519"
}

# Personal repository
resource "forgejo_repository" "personal_defaults" {
  name = "personal_test_repo_defaults"
}

# Deploy key
resource "forgejo_deploy_key" "this" {
  repository_id = forgejo_repository.personal_defaults.id
  key           = trimspace(tls_private_key.ed25519.public_key_openssh)
  title         = "test_key"
  read_only     = false
}
output "debug" {
  value = forgejo_deploy_key.this
}
