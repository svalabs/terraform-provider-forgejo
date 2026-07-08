terraform {
  required_providers {
    forgejo = {
      source = "svalabs/forgejo"
    }
  }
}

variable "test_password" { sensitive = true }

provider "forgejo" {
  host = "http://localhost:3000"

  username = "admin"
  password = var.forgejo_password
  # ...or use the FORGEJO_USERNAME / FORGEJO_PASSWORD environment variables
}

resource "forgejo_user" "test_user" {
  login    = "test_user"
  email    = "test_user@localhost.localdomain"
  password = var.test_password
}

resource "forgejo_personal_access_token" "test_token" {
  user_id = forgejo_user.test_user.id
  name    = "test token"
  scopes = [
    "all",
    "read:repository"
  ]
}
