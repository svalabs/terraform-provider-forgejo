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

# Personal repository
resource "forgejo_repository" "personal" {
  name = "personal_test_repo"
}

# Repository action secret
resource "forgejo_repository_action_secret" "this" {
  repository_id = forgejo_repository.personal.id
  name          = "my_secret"
  data          = "my_secret_value"
}

# Write-only secret value: never stored in the state or plan
# (requires Terraform >= 1.11 or OpenTofu >= 1.11)
variable "repository_secret_value" {
  type      = string
  sensitive = true
  ephemeral = true
}

resource "forgejo_repository_action_secret" "write_only" {
  repository_id   = forgejo_repository.personal.id
  name            = "my_write_only_secret"
  data_wo         = var.repository_secret_value
  data_wo_version = 1 # increment to update the secret
}
