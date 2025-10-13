## 0.5.0 (October 13, 2025)

FEATURES:

- **New Resource**: `forgejo_repository_action_secret` ([documentation](docs/resources/repository_action_secret.md))
- **New Resource**: `forgejo_organization_action_secret` ([documentation](docs/resources/organization_action_secret.md))

BUG FIXES:

- `forgejo_repository`: improve schema validation for tracker and wiki attributes ([documentation](docs/resources/repository.md))

ENHANCEMENTS:

- Add more test cases
- Improve documentation

DEPENDENCIES:

- Bump actions/setup-go from 5.5.0 to 6.0.0
- Bump github.com/hashicorp/terraform-plugin-framework from 1.15.1 to 1.16.1
- Bump github.com/hashicorp/terraform-plugin-go from 0.28.0 to 0.29.0

## 0.4.0 (August 24, 2025)

FEATURES:

- **New Resource**: `forgejo_collaborator` ([documentation](docs/resources/collaborator.md))
- **New Data Source**: `forgejo_collaborator` ([documentation](docs/data-sources/collaborator.md))

DEPENDENCIES:

- Update to Go version 1.24.6
- Bump actions/checkout from 4.2.2 to 5.0.0
- Bump github.com/hashicorp/terraform-plugin-framework 1.15.0 to 1.15.1
- Bump github.com/hashicorp/terraform-plugin-testing from 1.13.2 to 1.13.3
- Bump goreleaser/goreleaser-action from 6.3.0 to 6.4.0

## 0.3.1 (June 29, 2025)

FEATURES:

- `forgejo_repository`: add token authentication for repository migrations (clone & pull mirror repos) ([documentation](docs/resources/repository.md))

DEPENDENCIES:

- Update to go 1.23.10
- Bump codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2 from 2.1.0 to 2.2.0
- Bump github.com/cloudflare/circl from 1.6.0 to 1.6.1
- Bump github.com/cloudflare/circl from 1.3.7 to 1.6.1 in /tools
- Bump github.com/hashicorp/terraform-plugin-testing from 1.13.1 to 1.13.2

## 0.3.0 (June 8, 2025)

FEATURES:

- `forgejo_repository`: add support for repository migration (clone & pull mirror repos) ([documentation](docs/resources/repository.md))

ENHANCEMENTS:

- Automatically run acceptance tests in CI

DOCUMENTATION:

- Add open in dev container badge

DEPENDENCIES:

- Use Go version 1.23.9 consistently
- Update test environment to mariadb:lts
- Remove unneeded dependencies
- Bump github.com/hashicorp/terraform-plugin-go from 0.27.0 to 0.28.0
- Bump github.com/hashicorp/terraform-plugin-testing from 1.13.0 to 1.13.1

## 0.2.4 (May 24, 2025)

DOCUMENTATION:

- Add workflow status badges
- Update examples to use variables instead of hardcoded secrets

DEPENDENCIES:

- Update local test environment to forgejo:10.0
- Migrate golangci-lint configuration from v1 to v2
- Bump actions/setup-go from 5.3.0 to 5.5.0
- Bump codeberg.org/mvdkleijn/forgejo-sdk/forgejo from v2.0.0 to v2.1.0
- Bump crazy-max/ghaction-import-gpg from 6.2.0 to 6.3.0
- Bump github.com/hashicorp/terraform-plugin-framework
- Bump github.com/hashicorp/terraform-plugin-framework-validators
- Bump github.com/hashicorp/terraform-plugin-testing
- Bump golang.org/x/net from 0.36.0 to 0.38.0
- Bump golang.org/x/net from 0.36.0 to 0.38.0 in /tools
- Bump golangci/golangci-lint-action from 6.5.1 to 8.0.0
- Bump goreleaser/goreleaser-action from 6.2.1 to 6.3.0

## 0.2.3 (March 22, 2025)

ENHANCEMENTS:

- Improve documentation

DEPENDENCIES:

- Bump golang.org/x/net from 0.33.0 to 0.36.0 in /tools
- Bump golang.org/x/net from 0.34.0 to 0.36.0
- Bump golangci/golangci-lint-action from 6.5.0 to 6.5.1
- Bump github.com/golang-jwt/jwt/v4 in /tools

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
