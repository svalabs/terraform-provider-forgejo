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

# Existing repository
data "forgejo_repository" "user" {
  owner = {
    login = "test_user"
  }
  name = "user_test_repo"
}

# Existing collaborator
data "forgejo_collaborator" "this" {
  repository_id = data.forgejo_repository.user.id
  user          = "test_collaborator"
}
