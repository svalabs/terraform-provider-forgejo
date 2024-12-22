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

data "forgejo_repository" "user" {
  owner = {
    login = "achim"
  }
  name = "user_test_repo_1"
}
output "user_debug" {
  value = data.forgejo_repository.user
}

data "forgejo_repository" "org" {
  owner = {
    login = "test_org_1"
  }
  name = "org_test_repo_1"
}
output "org_debug" {
  value = data.forgejo_repository.org
}
