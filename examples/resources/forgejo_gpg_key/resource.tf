terraform {
  required_providers {
    forgejo = {
      source = "svalabs/forgejo"
    }
    gpg = {
      source = "terraform-provider-gpg/gpg"
    }
  }
}

variable "gpg_key_password" { sensitive = true }

provider "forgejo" {
  host = "http://localhost:3000"
}

# GPG key
resource "gpg_key" "test" {
  identities = [{
    name  = "Test User"
    email = "test_user@localhost.localdomain"
  }]
  passphrase = var.gpg_key_password
}

# Forgejo GPG key
resource "forgejo_gpg_key" "this" {
  armored_public_key = gpg_key.test.public_key
}
