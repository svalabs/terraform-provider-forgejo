terraform {
  required_providers {
    forgejo = {
      source = "registry.terraform.io/svalabs/forgejo"
    }
  }
}

provider "forgejo" {
  host     = "http://localhost:3000"
  username = "achim"
  password = "password"
}

data "forgejo_organization" "this" {
  name = "test1"
}

output "debug" {
  value = data.forgejo_organization.this
}
