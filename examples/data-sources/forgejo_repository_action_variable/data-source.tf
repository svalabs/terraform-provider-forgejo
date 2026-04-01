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
data "forgejo_repository" "test" {
  name = "test_repo"
}

# Existing repository action variable
data "forgejo_repository_action_variable" "this" {
  repository_id = data.forgejo_repository.test.id
  name          = "my_variable"
}
