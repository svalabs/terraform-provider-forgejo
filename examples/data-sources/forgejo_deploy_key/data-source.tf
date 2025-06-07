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

# Existing deploy key
data "forgejo_deploy_key" "this" {
  repository_id = data.forgejo_repository.user.id
  title         = "test_key"
}
