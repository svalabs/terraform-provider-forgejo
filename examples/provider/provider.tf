terraform {
  required_providers {
    forgejo = {
      source = "svalabs/forgejo"
    }
  }
}

# Authenticate with API token
provider "forgejo" {
  alias     = "apiToken"
  host      = "http://localhost:3000"
  api_token = "1234567890abcdefghijklmnopqrstuvwxyz1234"
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
  password = "passw0rd"
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
