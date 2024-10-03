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

resource "forgejo_organization" "this" {
  name        = "tftest"
  full_name   = "Terraform Test Org"
  description = "Purely for testing..."
}

output "debug" {
  value = forgejo_organization.this
}
