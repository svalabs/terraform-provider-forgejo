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

# Personal repository
resource "forgejo_repository" "personal" {
  name = "personal_test_repo"
}

# Repository action secret
resource "forgejo_repository_action_secret" "this" {
  repository_id = forgejo_repository.personal.id
  name          = "my_secret"
  data          = "my_secret_value"
}
