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

# Organization action variable
resource "forgejo_organization_action_variable" "this" {
  organization_id = forgejo_organization.test.id
  name            = "my_variable"
  data            = "my_variable_value"
}
