terraform {
  required_providers {
    forgejo = {
      source = "svalabs/forgejo"
    }
    tls = {
      source = "hashicorp/tls"
    }
  }
}

provider "forgejo" {
  host = "http://localhost:3000"
}

# Personal repository
resource "forgejo_repository" "personal" {
  name = "personal_test_repo"
}

# Repository webhook
resource "forgejo_repository_webhook" "example" {
  repository_id        = forgejo_repository.personal.id
  authorization_header = "Bearer token123456"
  branch_filter        = "*"
  type                 = "forgejo"
  events               = ["push"]
  config = {
    "content_type" = "json"
    "url"          = "http://example.com/invoke"
  }
}
