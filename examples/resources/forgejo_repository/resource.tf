terraform {
  required_providers {
    forgejo = {
      source = "svalabs/forgejo"
    }
  }
}

variable "test_password" { sensitive = true }
variable "test_token" { sensitive = true }

provider "forgejo" {
  host = "http://localhost:3000"
}

# Personal repository with default settings
# (owned by the authenticated user)
resource "forgejo_repository" "personal_defaults" {
  name = "personal_test_repo_defaults"
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

  external_tracker = {
    external_tracker_url    = "https://github.com/svalabs/terraform-provider-forgejo/issues"
    external_tracker_format = "https://github.com/svalabs/terraform-provider-forgejo/issues/{index}"
    external_tracker_style  = "numeric"
  }
}

# Organization
resource "forgejo_organization" "owner" {
  name = "test_org"
}

# Organization repository with default settings
# (owned by an organization)
resource "forgejo_repository" "org_defaults" {
  owner = forgejo_organization.owner.name
  name  = "org_test_repo_defaults"
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

  external_tracker = {
    external_tracker_url    = "https://github.com/svalabs/terraform-provider-forgejo/issues"
    external_tracker_format = "https://github.com/svalabs/terraform-provider-forgejo/issues/{index}"
    external_tracker_style  = "numeric"
  }
}

# User
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

  external_tracker = {
    external_tracker_url    = "https://github.com/svalabs/terraform-provider-forgejo/issues"
    external_tracker_format = "https://github.com/svalabs/terraform-provider-forgejo/issues/{index}"
    external_tracker_style  = "numeric"
  }
}

# Clone repository
resource "forgejo_repository" "clone" {
  name       = "clone_test_repo"
  clone_addr = "https://github.com/svalabs/terraform-provider-forgejo"
  auth_token = var.test_token # optional
  mirror     = false
}

# Pull mirror repository
resource "forgejo_repository" "mirror" {
  name            = "mirror_test_repo"
  clone_addr      = "https://github.com/svalabs/terraform-provider-forgejo"
  auth_token      = var.test_token # optional
  mirror          = true
  mirror_interval = "12h0m0s" # optional
}
