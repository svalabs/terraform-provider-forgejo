terraform {
  required_providers {
    forgejo = {
      source = "svalabs/forgejo"
    }
  }
}

provider "forgejo" {
  host = "http://localhost:3000"
}

# Existing GPG key
data "forgejo_gpg_key" "this" {
  user   = "test_user" # Optional, uses authenticated user if not provided
  key_id = "test_key"
}
