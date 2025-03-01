## 0.2.2 (March 1, 2025)

DEPENDENCIES:

- Update golangci-lint config from template repository

## 0.2.1 (March 1, 2025)

ENHANCEMENTS:

- Improve documentation

DEPENDENCIES:

- Update local test environment to forgejo:9.0 & mariadb:10
- Update GitHub workflows from template repository
- Bump codeberg.org/mvdkleijn/forgejo-sdk/forgejo from v1.2.0 to v2.0.0
- Bump github.com/hashicorp/terraform-plugin-framework from 1.13.0 to 1.14.1
- Bump github.com/hashicorp/terraform-plugin-framework-validators from 0.16.0 to 0.17.0
- Bump golangci/golangci-lint-action from 6.3.1 to 6.5.0
- Bump goreleaser/goreleaser-action from 6.2.0 to 6.2.1

## 0.2.0 (February 8, 2025)

FEATURES:

- **New Resource**: `forgejo_deploy_key` ([documentation](docs/resources/deploy_key.md))
- **New Data Source**: `forgejo_deploy_key` ([documentation](docs/data-sources/deploy_key.md))

DEPENDENCIES:

- Bump golangci/golangci-lint-action from 6.1.1 to 6.3.0
- Bump github.com/hashicorp/terraform-plugin-go from 0.25.0 to 0.26.0
- Bump actions/setup-go from 5.2.0 to 5.3.0
- Bump golang.org/x/net from 0.23.0 to 0.33.0 in /tools

## 0.1.2 (January 17, 2025)

ENHANCEMENTS:

- Improve documentation and add troubleshooting section

DEPENDENCIES:

- Bump golang.org/x/crypto from 0.21.0 to 0.31.0 in /tools
- Bump github.com/golang-jwt/jwt/v4 in /tools
- Bump golang.org/x/crypto from 0.29.0 to 0.31.0
- Bump goreleaser/goreleaser-action from 6.0.0 to 6.1.0
- Bump golangci/golangci-lint-action from 6.1.0 to 6.1.1
- Bump crazy-max/ghaction-import-gpg from 6.1.0 to 6.2.0
- Bump actions/setup-go from 5.0.2 to 5.2.0
- Bump actions/checkout from 4.1.7 to 4.2.2

## 0.1.1 (December 24, 2024) 🎄

ENHANCEMENTS:

- Improve documentation and examples

## 0.1.0 (December 23, 2024)

MINIMUM VIABLE PRODUCT (MVP)

FEATURES:

- Authentication with API token, or with username and password
- **New Resource**: `forgejo_organization` ([documentation](docs/resources/organization.md))
- **New Resource**: `forgejo_repository` ([documentation](docs/resources/repository.md))
- **New Resource**: `forgejo_user` ([documentation](docs/resources/user.md))
- **New Data Source**: `forgejo_organization` ([documentation](docs/data-sources/organization.md))
- **New Data Source**: `forgejo_repository` ([documentation](docs/data-sources/repository.md))
- **New Data Source**: `forgejo_user` ([documentation](docs/data-sources/user.md))
