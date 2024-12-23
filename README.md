# Terraform Provider for Forgejo

This repository contains a [Terraform](https://www.terraform.io/) provider for [Forgejo](https://forgejo.org/) — self-hosted lightweight software forge.

## Contents

The Forgejo Terraform Provider allows managing resources within Forgejo. It is in an **early** stage and currently provides the following...

Resources:

- `forgejo_organization` ([documentation](docs/resources/organization.md))
- `forgejo_repository` ([documentation](docs/resources/repository.md))
- `forgejo_user` ([documentation](docs/resources/user.md))

Data Sources:

- `forgejo_organization` ([documentation](docs/data-sources/organization.md))
- `forgejo_repository` ([documentation](docs/data-sources/repository.md))
- `forgejo_user` ([documentation](docs/data-sources/user.md))

## Using the Provider

Import the provider into your Terraform configuration:

```terraform
terraform {
  required_providers {
    forgejo = {
      source = "svalabs/forgejo"
      version = "~> 0.1.0"
    }
  }
}
```

There are two methods for authenticating to the Forgejo API: using an API token, or with username and password.

It is recommended to supply an API token to authenticate with a given Forgejo host:

```terraform
provider "forgejo" {
  host      = "http://localhost:3000"
  api_token = "1234567890abcdefghijklmnopqrstuvwxyz1234"
  # or use the FORGEJO_API_TOKEN environment variable
}
```

Alternatively, supply username and password to authenticate:

```terraform
provider "forgejo" {
  host     = "http://localhost:3000"
  username = "admin"
  password = "passw0rd"
  # or use the FORGEJO_USERNAME / FORGEJO_PASSWORD environment variables
}
```

A **personal repository** can be created like so:

```terraform
resource "forgejo_repository" "example" {
  name        = "new_personal_repo"
  description = "Purely for testing..."
}
```

A **user repository** can be created like so (requires administrative privileges):

```terraform
resource "forgejo_user" "example" {
  name = "new_user"
}

resource "forgejo_repository" "example" {
  owner       = forgejo_user.example.name
  name        = "new_user_repo"
  description = "Purely for testing..."
}
```

An **organization repository** can be created like so:

```terraform
resource "forgejo_organization" "example" {
  name = "new_org"
}

resource "forgejo_repository" "example" {
  owner       = forgejo_organization.example.name
  name        = "new_org_repo"
  description = "Purely for testing..."
}
```

These examples create repositories with most attributes set to their default values. However, many settings can be customized:

```terraform
resource "forgejo_repository" "example" {
  owner          = forgejo_organization.example.name
  name           = "new_org_repo"
  description    = "Purely for testing..."
  private        = true
  default_branch = "dev"
  auto_init      = true
  trust_model    = "collaborator"

  internal_tracker = {
    enable_time_tracker                   = false
    allow_only_contributors_to_track_time = false
    enable_issue_dependencies             = false
  }
}
```

Refer to the `examples/` directory for more usage examples.

## Copyright and License

Copyright (c) 2024 SVA System Vertrieb Alexander GmbH.

Released under the terms of the [Mozilla Public License (MPL-2.0)](LICENSE).
