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
resource "forgejo_repository" "test_repo" {
  name        = "personal_test_repo"
  description = "Terraform Test Repo owned by user"
  auto_init   = true
}

# Tag protection for release tags
# Note: the API requires at least one whitelisted user or team on creation
resource "forgejo_tag_protection" "releases" {
  repository_id       = forgejo_repository.test_repo.id
  name_pattern        = "v*"
  whitelist_usernames = ["tfadmin"]
}

# Example how to import existing tag protections
# id follows the format: <owner>/<repo>/<protection_id>
import {
  id = "tfadmin/personal_test_repo/1"
  to = forgejo_tag_protection.releases
}
