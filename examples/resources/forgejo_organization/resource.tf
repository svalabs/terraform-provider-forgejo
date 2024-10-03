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

resource "forgejo_organization" "defaults" {
  name = "tftest_defaults"
}
resource "forgejo_organization" "non_defaults" {
  name        = "tftest_non_defaults"
  full_name   = "Terraform Test Org with non-default attributes"
  description = "Purely for testing..."
  website     = "https://forgejo.org/"
  location    = "Mêlée Island"
  visibility  = "private"
}

output "debug" {
  value = forgejo_organization.defaults
}
