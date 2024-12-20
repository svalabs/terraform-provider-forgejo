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

resource "forgejo_user" "defaults" {
  login    = "tftest_user_defaults"
  email    = "tftest_user_defaults@localhost.localdomain"
  password = "passw0rd"
}
resource "forgejo_user" "non_defaults" {
  login       = "tftest_user_non_defaults"
  email       = "tftest_user_non_defaults@localhost.localdomain"
  password    = "passw0rd"
  full_name   = "Terraform Test User with non-default attributes"
  description = "Purely for testing..."
  website     = "https://forgejo.org/"
  location    = "Mêlée Island"
  visibility  = "private"
}

output "debug_defaults" {
  value = forgejo_user.defaults
}
output "debug_non_defaults" {
  value = forgejo_user.non_defaults
}
