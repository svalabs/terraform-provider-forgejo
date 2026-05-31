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
  name           = "personal_test_repo"
  description    = "Terraform Test Repo owned by user"
  website        = "http://localhost:3000"
  private        = true
  template       = true
  default_branch = "custom"

  issue_labels = "Default"
  auto_init    = true
  readme       = "Default"
  trust_model  = "collaborator"
}


# Branch protection with default settings
resource "forgejo_branch_protection" "defaults" {
  branch_name   = "main"
  repository_id = forgejo_repository.test_repo.id
}

# Branch protection with custom settings
resource "forgejo_branch_protection" "non_defaults" {
  branch_name   = "main"
  repository_id = forgejo_repository.test_repo.id

  block_on_outdated_branch  = true
  block_on_rejected_reviews = true
  dismiss_stale_approvals   = true
  required_approvals        = 1

  enable_approvals_whitelist = true
  approvals_whitelist_usernames = [
    "alice",
    "bob",
  ]

  enable_merge_whitelist = true
  merge_whitelist_usernames = [
    "alice",
    "bob",
  ]

  enable_push           = true
  enable_push_whitelist = true
  push_whitelist_usernames = [
    "alice",
    "bob",
  ]
}

# Example how to import existing branch protections
# id follows the format: <owner>/<repo>/<branch>
import {
  id = "tfadmin/personal_test_repo/main"
  to = forgejo_branch_protection.defaults
}
