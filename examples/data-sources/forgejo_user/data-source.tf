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

data "forgejo_user" "this" {
  login = "achim"
}

output "debug" {
  value = data.forgejo_user.this
}
