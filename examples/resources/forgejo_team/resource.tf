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

# Team with default settings
resource "forgejo_team" "default_team" {
  organization_id = forgejo_organization.owner.id
  name            = "org_test_team_defaults"

  units_map = {
    "repo.code" = "read"
  }
}

# Team with custom settings
resource "forgejo_team" "custom_team" {
  organization_id           = forgejo_organization.owner.id
  name                      = "org_test_team_non_defaults"
  can_create_org_repo       = true
  description               = "A team with non-default parameters."
  includes_all_repositories = true
  permission                = "read"

  units_map = {
    "repo.code"   = "read"
    "repo.issues" = "write"
    "repo.pulls"  = "read"
  }
}
