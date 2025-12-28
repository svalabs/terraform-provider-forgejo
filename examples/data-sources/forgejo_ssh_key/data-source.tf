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

# Existing SSH key
data "forgejo_ssh_key" "this" {
  user  = "test_user"
  title = "test_key"
}
