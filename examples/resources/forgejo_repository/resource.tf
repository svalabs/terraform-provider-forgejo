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

#
# Personal repository
#
resource "forgejo_repository" "personal_defaults" {
  name = "personal_tftest_defaults"
}
output "personal_debug_defaults" {
  value = forgejo_repository.personal_defaults
}

resource "forgejo_repository" "personal_non_defaults" {
  name           = "personal_tftest_non_defaults"
  description    = "Terraform Test Repo owned by user with non-default attributes"
  website        = "http://localhost:3000"
  private        = true
  template       = true
  default_branch = "custom"
  issue_labels   = "Default"
  auto_init      = false
  readme         = "Default"
  trust_model    = "collaborator"
  archived       = true

  internal_tracker = {
    enable_time_tracker                   = false
    allow_only_contributors_to_track_time = false
    enable_issue_dependencies             = false
  }
}
output "personal_debug_non_defaults" {
  value = forgejo_repository.personal_non_defaults
}

#
# Organization repository
#
resource "forgejo_organization" "owner" {
  name = "tftest_org"
}

resource "forgejo_repository" "org_defaults" {
  owner = forgejo_organization.owner.name
  name  = "org_tftest_defaults"
}
output "org_debug_defaults" {
  value = forgejo_repository.org_defaults
}

resource "forgejo_repository" "org_non_defaults" {
  owner          = forgejo_organization.owner.name
  name           = "org_tftest_non_defaults"
  description    = "Terraform Test Repo owned by org with non-default attributes"
  website        = "http://localhost:3000"
  private        = true
  template       = true
  default_branch = "custom"
  issue_labels   = "Default"
  auto_init      = false
  readme         = "Default"
  trust_model    = "collaborator"
  archived       = true

  internal_tracker = {
    enable_time_tracker                   = false
    allow_only_contributors_to_track_time = false
    enable_issue_dependencies             = false
  }
}
output "org_debug_non_defaults" {
  value = forgejo_repository.org_non_defaults
}

#
# User repository
#
resource "forgejo_user" "owner" {
  login    = "tftest_user"
  email    = "tftest_user@localhost.localdomain"
  password = "passw0rd"
}

resource "forgejo_repository" "user_defaults" {
  owner = forgejo_user.owner.login
  name  = "user_tftest_defaults"
}
output "user_debug_defaults" {
  value = forgejo_repository.user_defaults
}

resource "forgejo_repository" "user_non_defaults" {
  owner          = forgejo_user.owner.login
  name           = "user_tftest_non_defaults"
  description    = "Terraform Test Repo owned by user with non-default attributes"
  website        = "http://localhost:3000"
  private        = true
  template       = true
  default_branch = "custom"
  issue_labels   = "Default"
  auto_init      = false
  readme         = "Default"
  trust_model    = "collaborator"
  archived       = true

  internal_tracker = {
    enable_time_tracker                   = false
    allow_only_contributors_to_track_time = false
    enable_issue_dependencies             = false
  }
}
output "user_debug_non_defaults" {
  value = forgejo_repository.user_non_defaults
}
