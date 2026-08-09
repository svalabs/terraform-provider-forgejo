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
  // Due to an upstream limitation, one cannot create access tokens when authorized with an access token.
  // Use basic-auth instead (FORGEJO_USERNAME / FORGEJO_PASSWORD environment variables).
}

resource "forgejo_user" "test_user" {
  login    = "test_user"
  email    = "test_user@localhost.localdomain"
  password = var.test_password
}

resource "forgejo_personal_access_token" "test_token" {
  user = forgejo_user.test_user.login
  name = "test token"
  scopes = [
    "read:repository"
  ]
}
