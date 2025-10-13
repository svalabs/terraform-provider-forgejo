# Terraform Provider for Forgejo

![Tests](https://github.com/svalabs/terraform-provider-forgejo/actions/workflows/test.yml/badge.svg)
![Release](https://github.com/svalabs/terraform-provider-forgejo/actions/workflows/release.yml/badge.svg)

This repository contains a [Terraform](https://www.terraform.io/) and [OpenTofu](https://opentofu.org/) provider for [Forgejo](https://forgejo.org/) — self-hosted lightweight software forge.
The project is based on the awsome [Forgejo SDK for Go](https://codeberg.org/mvdkleijn/forgejo-sdk) by [Martijn van der Kleijn](https://vanderkleijn.net/).
Thanks, Martijn, for the heavy lifting!

## Contents

The Forgejo Terraform/OpenTofu Provider allows managing resources and data source within Forgejo instances.
It currently provides the following...

Resources:

- `forgejo_collaborator` ([documentation](docs/resources/collaborator.md))
- `forgejo_deploy_key` ([documentation](docs/resources/deploy_key.md))
- `forgejo_organization` ([documentation](docs/resources/organization.md))
- `forgejo_repository` ([documentation](docs/resources/repository.md))
- `forgejo_repository_action_secret` ([documentation](docs/resources/repository_action_secret.md))
- `forgejo_user` ([documentation](docs/resources/user.md))

Data Sources:

- `forgejo_collaborator` ([documentation](docs/data-sources/collaborator.md))
- `forgejo_deploy_key` ([documentation](docs/data-sources/deploy_key.md))
- `forgejo_organization` ([documentation](docs/data-sources/organization.md))
- `forgejo_repository` ([documentation](docs/data-sources/repository.md))
- `forgejo_user` ([documentation](docs/data-sources/user.md))

## Using the Provider

Import the provider into your Terraform/OpenTofu configuration:

```terraform
terraform {
  required_providers {
    forgejo = {
      source  = "svalabs/forgejo"
      version = "~> 0.4.0"
    }
  }
}
```

There are two methods for authenticating to the Forgejo API: using an API token, or with username and password.

It is recommended to supply an **API token** to authenticate with a given Forgejo host:

```terraform
provider "forgejo" {
  host      = "http://localhost:3000"
  api_token = "<<<your_api_key>>>"
  # ...or use the FORGEJO_API_TOKEN environment variable
}
```

API tokens can be generated through the Forgejo web interface, by navigating to Settings → Applications → Access tokens → Generate new token.

The following API token permissions are required:

- `write:organization`
- `write:repository`
- `write:user`

Optionally, for administrative privileges (required to manage users and user repositories):

- `write:admin`

Alternatively, supply **username** and **password** to authenticate:

```terraform
provider "forgejo" {
  host     = "http://localhost:3000"
  username = "<<<your_username>>>"
  password = "<<<your_password>>>"
  # ...or use the FORGEJO_USERNAME / FORGEJO_PASSWORD environment variables
}
```

> **Important**: The Forgejo API client does not (currently) allow ignoring certificate errors.
> When connecting through `https://`, the Forgejo host must supply certificates trusted by the Terraform/OpenTofu host.
> Hence, self-signed certificates must be imported locally.
> This can be achieved by running the following command:
>
> ```shell
> echo quit | openssl s_client -showcerts -servername <<<forgejo_host>>> -connect <<<forgejo_host>>> > /etc/ssl/certs/cacert.pem
> ```

A **personal repository** can be created like so:

```terraform
resource "forgejo_repository" "example" {
  name        = "new_personal_repo"
  description = "Purely for testing..."
}
```

A **user repository** can be created like so (requires administrative privileges):

```terraform
resource "forgejo_user" "owner" {
  login = "new_user"
}

resource "forgejo_repository" "example" {
  owner       = forgejo_user.owner.login
  name        = "new_user_repo"
  description = "Purely for testing..."
}
```

An **organization repository** can be created like so:

```terraform
resource "forgejo_organization" "owner" {
  name = "new_org"
}

resource "forgejo_repository" "example" {
  owner       = forgejo_organization.owner.name
  name        = "new_org_repo"
  description = "Purely for testing..."
}
```

A **clone repository** can be created like so:

```terraform
resource "forgejo_repository" "clone" {
  name       = "clone_test_repo"
  clone_addr = "https://github.com/svalabs/terraform-provider-forgejo"
  mirror     = false
}
```

A **pull mirror repository** can be created like so:

```terraform
resource "forgejo_repository" "mirror" {
  name            = "mirror_test_repo"
  clone_addr      = "https://github.com/svalabs/terraform-provider-forgejo"
  mirror          = true
  mirror_interval = "12h0m0s"
}
```

These examples create repositories with most attributes set to their default values.
However, many settings can be customized:

```terraform
resource "forgejo_repository" "example" {
  owner          = forgejo_organization.owner.name
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

## Troubleshooting

### Error: failed to verify certificate: certificate signed by unknown authority

In case of the following error message:

```
Error: Unable to Create Forgejo API Client

    An unexpected error occurred when creating the Forgejo API client. If the
    error is not clear, please contact the provider developers.

    Forgejo Client Error: Get "https://.../api/v1/version":
    tls: failed to verify certificate: x509: certificate signed by unknown
    authority
```

Extract the self-signed certificate from the Forgejo host and import it locally:

```shell
echo quit | openssl s_client -showcerts -servername <<<forgejo_host>>> -connect <<<forgejo_host>>> > /etc/ssl/certs/cacert.pem
```

### Error: token does not have at least one of required scope(s)

In case of the following error message:

```
Error: Unable to get repository by id

    Unknown error: token does not have at least one of required scope(s):
    [read:repository]
```

Re-generate the API token used for authentication, and make sure to select the following permissions:

- `write:organization`
- `write:repository`
- `write:user`
- Optional, for managing users and user repositories: `write:admin`

## Developing & Contributing to the Provider

The [CONTRIBUTING.md](CONTRIBUTING.md) file is a basic outline on how to build and develop the provider.

## Copyright and License

Copyright (c) 2024 SVA System Vertrieb Alexander GmbH.

Released under the terms of the [Mozilla Public License (MPL-2.0)](LICENSE).
