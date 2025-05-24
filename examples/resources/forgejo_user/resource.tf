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

# User with default settings
resource "forgejo_user" "defaults" {
  login    = "test_user_defaults"
  email    = "test_user_defaults@localhost.localdomain"
  password = var.test_password
}
output "debug_defaults" {
  value     = forgejo_user.defaults
  sensitive = true
}

# User with custom settings
resource "forgejo_user" "non_defaults" {
  login       = "test_user_non_defaults"
  email       = "test_user_non_defaults@localhost.localdomain"
  password    = var.test_password
  full_name   = "Terraform Test User with non-default attributes"
  description = "Purely for testing..."
  website     = "https://forgejo.org/"
  location    = "Mêlée Island"
  visibility  = "private"
}
output "debug_non_defaults" {
  value     = forgejo_user.non_defaults
  sensitive = true
}
