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
data "forgejo_organization" "organization" {
  name = "test_organization"
}

# Existing team
data "forgejo_team" "team" {
  organization_id = data.forgejo_organization.organization.id
  name            = "test_team"
}

# Existing user
data "forgejo_user" "user" {
  name = "test_user"
}

# Existing member
data "forgejo_team_member" "team_membership" {
  team_id = data.forgejo_team.team.id
  user    = data.forgejo_user.user.login
}
