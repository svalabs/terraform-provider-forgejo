# Developing the Terraform Provider for Forgejo

[![Open in Dev Container](https://img.shields.io/static/v1?label=Dev%20Container&message=Open&color=blue)](https://vscode.dev/redirect?url=vscode://ms-vscode-remote.remote-containers/cloneInVolume?url=https://github.com/svalabs/terraform-provider-forgejo)

## Requirements

- [Forgejo](https://forgejo.org/docs/latest/admin/installation/) >= 15.0
- [Terraform](https://developer.hashicorp.com/terraform/install) >= 1.13
- [OpenTofu](https://opentofu.org/docs/intro/install/) >= 1.10
- [Go](https://golang.org/doc/install) >= 1.25

## Building the Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the `make install` command:

```shell
make install
```

### Directory Layout

```shell
terraform-provider-forgejo/
├── docker/    # Example Forgejo installation for local development
├── docs/      # Generated documentation
├── examples/  # Provider usage examples
├── internal/  # Provider source code and tests
└── tools/     # Scripts for generating documentation
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform/OpenTofu provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `make install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```shell
make install
```

To generate or update documentation, run `make generate`.

```shell
make generate
```

In order to run the full suite of Acceptance tests, run `make testacc`.

_Note:_ Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

## AI Guidance

Using AI tools to help write your PR is acceptable, but as the author, you are responsible for understanding every change. If you used AI tools in preparing your PR, you must disclose this in the description of your PR. For example, including “This PR was written in part with the assistance of generative AI,” in the PR description is sufficient.

Listing AI tooling as a co-author, co-signing commits using an AI tool, or using the `assisted-by`, `co-developed` or similar commit trailer is not allowed.
