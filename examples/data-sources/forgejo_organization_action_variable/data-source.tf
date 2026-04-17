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

# Existing organization
data "forgejo_organization" "test" {
  name = "test_org"
}

# Existing organization action variable
data "forgejo_organization_action_variable" "this" {
  organization_id = data.forgejo_organization.test.id
  name            = "my_variable"
}
