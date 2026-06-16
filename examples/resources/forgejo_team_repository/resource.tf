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

# Team
resource "forgejo_team" "team" {
  organization_id = forgejo_organization.owner.id
  name            = "org_test_team"

  units_map = {
    "repo.code" = "read"
  }
}

# Repository
resource "forgejo_repository" "repo" {
  owner = forgejo_organization.owner.name
  name  = "test_repo"
}

# Team repository
resource "forgejo_team_repository" "membership" {
  team_id    = forgejo_team.team.id
  owner      = forgejo_organization.owner.name
  repository = forgejo_repository.repo.name
}
