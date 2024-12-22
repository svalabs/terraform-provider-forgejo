terraform {
  required_providers {
    forgejo = {
      source = "svalabs/forgejo"
    }
  }
}

provider "forgejo" {
  alias    = "username"
  host     = "http://localhost:3000"
  username = "achim"
  password = "password"
}
data "forgejo_organization" "username" {
  provider = forgejo.username
  name     = "test1"
}

provider "forgejo" {
  alias     = "apiToken"
  host      = "http://localhost:3000"
  api_token = "c754bf42c0728e3031e8245a70dbdda0419aff44"
}
data "forgejo_organization" "apiToken" {
  provider = forgejo.apiToken
  name     = "test1"
}
