terraform {
  required_providers {
    forgejo = {
      source = "registry.terraform.io/svalabs/forgejo"
    }
  }
}

provider "forgejo" {
  host     = "http://localhost:3000"
  username = "achim"
  password = "password"
}
