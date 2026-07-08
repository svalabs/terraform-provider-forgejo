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

# Existing user
data "forgejo_user" "test_user" {
  login = "test_user"
}

# Existing personal access token
data "forgejo_personal_access_token" "test_token" {
  user_id = data.forgejo_user.test_user.id
  name    = "test token"
}
