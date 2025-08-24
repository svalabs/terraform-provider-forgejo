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

# Organization
resource "forgejo_organization" "owner" {
  name = "test_org"
}

# Organization repository
resource "forgejo_repository" "org" {
  owner = forgejo_organization.owner.name
  name  = "org_test_repo"
}

# User
resource "forgejo_user" "test" {
  login    = "test_user"
  email    = "test_user@localhost.localdomain"
  password = var.test_password
}

# Collaborator
resource "forgejo_collaborator" "admin" {
  repository_id = forgejo_repository.org.id
  user          = forgejo_user.test.login
  permission    = "admin"
}
