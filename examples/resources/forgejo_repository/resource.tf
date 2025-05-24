terraform {
  required_providers {
    forgejo = {
      source = "svalabs/forgejo"
    }
  }
}

variable "test_password" { sensitive = true }

provider "forgejo" {
  host = "http://localhost:3000"
}

# Personal repository with default settings
# (owned by the authenticated user)
resource "forgejo_repository" "personal_defaults" {
  name = "personal_test_repo_defaults"
}
output "personal_debug_defaults" {
  value = forgejo_repository.personal_defaults
}

# Personal repository with custom settings
resource "forgejo_repository" "personal_non_defaults" {
  name           = "personal_test_repo_non_defaults"
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

resource "forgejo_organization" "owner" {
  name = "test_org"
}

# Organization repository with default settings
# (owned by an organization)
resource "forgejo_repository" "org_defaults" {
  owner = forgejo_organization.owner.name
  name  = "org_test_repo_defaults"
}
output "org_debug_defaults" {
  value = forgejo_repository.org_defaults
}

# Organization repository with custom settings
resource "forgejo_repository" "org_non_defaults" {
  owner          = forgejo_organization.owner.name
  name           = "org_test_repo_non_defaults"
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

resource "forgejo_user" "owner" {
  login    = "test_user"
  email    = "test_user@localhost.localdomain"
  password = var.test_password
}

# User repository with default settings
# (owned by a different user)
resource "forgejo_repository" "user_defaults" {
  owner = forgejo_user.owner.login
  name  = "user_test_repo_defaults"
}
output "user_debug_defaults" {
  value = forgejo_repository.user_defaults
}

# User repository with custom settings
resource "forgejo_repository" "user_non_defaults" {
  owner          = forgejo_user.owner.login
  name           = "user_test_repo_non_defaults"
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
