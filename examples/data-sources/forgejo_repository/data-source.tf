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

# Existing user repository
data "forgejo_repository" "user" {
  owner = {
    login = "test_user"
  }
  name = "user_test_repo"
}
output "user_debug" {
  value = data.forgejo_repository.user
}

# Existing organization repository
data "forgejo_repository" "org" {
  owner = {
    login = "test_org"
  }
  name = "org_test_repo"
}
output "org_debug" {
  value = data.forgejo_repository.org
}
