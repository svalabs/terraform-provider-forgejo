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
    # The "secret" key is write-only: Forgejo accepts it on create/update but
    # never returns it. The provider preserves it from configuration so that
    # managing it here does not cause "inconsistent result after apply" errors.
    "secret" = "supersecret"
  }
}

# Import repository webhook
# id follows the format: <owner>/<repo>/<webhook_id>
import {
  id = "tfadmin/personal_test_repo/123"
  to = forgejo_repository_webhook.example
}
