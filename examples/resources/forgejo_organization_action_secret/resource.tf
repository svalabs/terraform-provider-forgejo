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

# Organization
resource "forgejo_organization" "test" {
  name = "test_org"
}

# Organization action secret
resource "forgejo_organization_action_secret" "this" {
  organization = forgejo_organization.test.name
  name         = "my_secret"
  data         = "my_secret_value"
}
