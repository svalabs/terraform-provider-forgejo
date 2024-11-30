terraform {
  required_providers {
    forgejo = {
      source = "registry.terraform.io/svalabs/forgejo"
    }
  }
}

provider "forgejo" {
  host = "http://localhost:3000"
}

data "forgejo_repository" "this" {
  owner = {
    login = "achim"
  }
  name = "test1"
}

output "debug" {
  value = data.forgejo_repository.this
}
