terraform {
  required_providers {
    forgejo = {
      source = "svalabs/forgejo"
    }
  }
}

variable "forgejo_api_token" { sensitive = true }
variable "forgejo_password" { sensitive = true }

# Authenticate with API token
provider "forgejo" {
  alias     = "apiToken"
  host      = "http://localhost:3000"
  api_token = var.forgejo_api_token
  # ...or use the FORGEJO_API_TOKEN environment variable
}
resource "forgejo_repository" "example_personal" {
  provider    = forgejo.apiToken
  name        = "new_personal_repo"
  description = "Purely for testing..."
}

# Authenticate with username and password
provider "forgejo" {
  alias    = "username"
  host     = "http://localhost:3000"
  username = "admin"
  password = var.forgejo_password
  # ...or use the FORGEJO_USERNAME / FORGEJO_PASSWORD environment variables
}
resource "forgejo_organization" "example" {
  provider = forgejo.username
  name     = "new_org"
}
resource "forgejo_repository" "example_org" {
  provider    = forgejo.username
  owner       = forgejo_organization.example.name
  name        = "new_org_repo"
  description = "Purely for testing..."
}
