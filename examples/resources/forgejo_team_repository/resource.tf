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
resource "forgejo_organization" "owner" {
  name = "test_org"
}

# Repository
resource "forgejo_repository" "repo" {
  owner = forgejo_organization.owner.name
  name  = "some_repo"
}

# Team
resource "forgejo_team" "team" {
  organization_id           = forgejo_organization.owner.id
  name                      = "org_test_team"
  includes_all_repositories = false

  units_map = {
    "repo.code" = "read"
  }
}

# Team repository
resource "forgejo_team_repository" "example" {
  team_id       = forgejo_team.team.id
  repository_id = forgejo_repository.repo.id
}

# Import team repository
# id follows the format: <org_name>/<team_name>/<repo_name>
import {
  id = "test_org/org_test_team/some_repo"
  to = forgejo_team_repository.example
}
