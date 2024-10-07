# Terraform Forgejo Provider

This repository contains a Terraform provider for [Forgejo](https://forgejo.org/) — self-hosted lightweight software forge.

## Contents

It is in a **very early** stage and currently contains the following...

Resources:

- `forgejo_organization` ([documentation](docs/resources/organization.md))

Data Sources:

- `forgejo_organization` ([documentation](docs/data-sources/organization.md))

### Directory Layout

```shell
terraform-provider-forgejo/
├── docker/    # Example Forgejo installation for local development
├── docs/      # Generated documentation
├── examples/  # Provider usage examples
├── internal/  # Provider source code and tests
└── tools/     # Scripts for generating documentation
```

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23

## Building the Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the Provider

```terraform
terraform {
  required_providers {
    forgejo = {
      source = "registry.terraform.io/svalabs/forgejo"
    }
  }
}

provider "forgejo" {
  host      = "http://localhost:3000"
  api_token = "1234567890abcdefghijklmnopqrstuvwxyz1234"
}

data "forgejo_organization" "example" {
  name = "existing_org"
}

resource "forgejo_organization" "example" {
  name        = "new_org"
  full_name   = "Terraform Test Org"
  description = "Purely for testing..."
  website     = "https://forgejo.org/"
  location    = "Mêlée Island"
  visibility  = "private"
}
```

Refer to the `examples/` directory for more usage examples.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

_Note:_ Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

## Copyright and License

Copyright (c) 2024 SVA System Vertrieb Alexander GmbH.

Released under the terms of the [Mozilla Public License (MPL-2.0)](LICENSE).
