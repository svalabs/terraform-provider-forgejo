terraform {
  required_providers {
    forgejo = {
      source = "registry.terraform.io/svalabs/forgejo"
    }
  }
}

provider "forgejo" {
  host = "http://localhost:3000"
}

resource "forgejo_repository" "personal_defaults" {
  owner = {}
  name  = "personal_tftest_defaults1"
}
resource "forgejo_repository" "personal_non_defaults" {
  owner       = {}
  name        = "personal_tftest_non_defaults"
  description = "Terraform Test Repo owned by user with non-default attributes"
  # website        = "http://localhost:3000"
  private        = true
  template       = true
  default_branch = "custom"
  issue_labels   = "Default"
  auto_init      = false
  readme         = "Default"
  trust_model    = "collaborator"
}

resource "forgejo_repository" "org_defaults" {
  owner = {
    login = "test_org_1"
  }
  name = "org_tftest_defaults"
}
resource "forgejo_repository" "org_non_defaults" {
  owner = {
    login = "test_org_1"
  }
  name        = "org_tftest_non_defaults"
  description = "Terraform Test Repo owned by org with non-default attributes"
  # website        = "http://localhost:3000"
  private        = true
  template       = true
  default_branch = "custom"
  issue_labels   = "Default"
  auto_init      = false
  readme         = "Default"
  trust_model    = "collaborator"
}

resource "forgejo_repository" "user_defaults" {
  owner = {
    login = "test_user_1"
  }
  name = "user_tftest_defaults"
}
resource "forgejo_repository" "user_non_defaults" {
  owner = {
    login = "test_user_1"
  }
  name        = "user_tftest_non_defaults"
  description = "Terraform Test Repo owned by user with non-default attributes"
  # website        = "http://localhost:3000"
  private        = true
  template       = true
  default_branch = "custom"
  issue_labels   = "Default"
  auto_init      = false
  readme         = "Default"
  trust_model    = "collaborator"
}

output "personal_debug_defaults" {
  value = forgejo_repository.personal_defaults
}
output "personal_debug_non_defaults" {
  value = forgejo_repository.personal_non_defaults
}

output "org_debug_defaults" {
  value = forgejo_repository.org_defaults
}
output "org_debug_non_defaults" {
  value = forgejo_repository.org_non_defaults
}

output "user_debug_defaults" {
  value = forgejo_repository.user_defaults
}
output "user_debug_non_defaults" {
  value = forgejo_repository.user_non_defaults
}
