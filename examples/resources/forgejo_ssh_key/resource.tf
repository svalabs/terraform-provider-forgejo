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

variable "test_password" { sensitive = true }

provider "forgejo" {
  host = "http://localhost:3000"
}

# User
resource "forgejo_user" "test" {
  login    = "test_user"
  email    = "test_user@localhost.localdomain"
  password = var.test_password
}

# Private key
resource "tls_private_key" "ed25519" {
  algorithm = "ED25519"
}

# SSH key
resource "forgejo_ssh_key" "this" {
  user  = forgejo_user.test.login
  key   = trimspace(tls_private_key.ed25519.public_key_openssh)
  title = "test_key"
}
