# Developing the Terraform Provider for Forgejo

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23

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

To add a new dependency `github.com/author/dependency` to your Terraform provider:

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
