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
resource "gpg_key_pair" "test" {
  identities = [{
    name  = "Test User"
    email = "test_user@localhost.localdomain"
  }]
  passphrase = var.gpg_key_password
}

# Forgejo GPG key
resource "forgejo_gpg_key" "this" {
  # There is no user for this resource, GPG keys can only be managed for the authenticated user.
  armored_public_key = gpg_key_pair.test.public_key
}
