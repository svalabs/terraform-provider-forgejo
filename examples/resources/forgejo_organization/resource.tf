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

# Organization with default settings
resource "forgejo_organization" "defaults" {
  name = "test_org_defaults"
}
output "debug_defaults" {
  value = forgejo_organization.defaults
}

# Organization with custom settings
resource "forgejo_organization" "non_defaults" {
  name        = "test_org_non_defaults"
  full_name   = "Terraform Test Org with non-default attributes"
  description = "Purely for testing..."
  website     = "https://forgejo.org/"
  location    = "Mêlée Island"
  visibility  = "private"
}
output "debug_non_defaults" {
  value = forgejo_organization.non_defaults
}
