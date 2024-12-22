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

data "forgejo_organization" "this" {
  name = "test1"
}
output "debug" {
  value = data.forgejo_organization.this
}
