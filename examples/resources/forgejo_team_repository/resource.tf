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

resource "forgejo_organization" "owner" {
  name = "test_org"
}

resource "forgejo_team" "developers" {
  name                      = "developers"
  organization              = forgejo_organization.owner.name
  includes_all_repositories = false
  units_map = {
    "repo.code" = "read"
  }
}

resource "forgejo_repository" "example" {
  owner   = forgejo_organization.owner.name
  name    = "example"
  private = true
}

resource "forgejo_team_repository" "assignment" {
  team_id       = forgejo_team.developers.id
  repository_id = forgejo_repository.example.id
}
