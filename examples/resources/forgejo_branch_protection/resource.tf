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

# Example how to import existing branch protections
# id follows the format: <owner>/<repo>/<branch>
import {
  id = "tfadmin/personal_test_repo/main"
  to = forgejo_branch_protection.main
}

# Branch protection on the main branch
resource "forgejo_branch_protection" "main" {
  branch_name   = "main"
  repository_id = forgejo_repository.test_repo.id
}
