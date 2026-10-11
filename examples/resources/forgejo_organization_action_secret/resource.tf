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
  organization_id = forgejo_organization.test.id
  name            = "my_secret"
  data            = "my_secret_value"
}

# Write-only secret value: never stored in the state or plan
# (requires Terraform >= 1.11 or OpenTofu >= 1.11)
variable "organization_secret_value" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "forgejo_organization_action_secret" "write_only" {
  organization_id = forgejo_organization.test.id
  name            = "my_write_only_secret"
  data_wo         = var.organization_secret_value
  data_wo_version = 1 # increment to update the secret
}
