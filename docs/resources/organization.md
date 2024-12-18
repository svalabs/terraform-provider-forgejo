---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "forgejo_organization Resource - forgejo"
subcategory: ""
description: |-
  Forgejo organization resource
---

# forgejo_organization (Resource)

Forgejo organization resource

## Example Usage

```terraform
terraform {
  required_providers {
    forgejo = {
      source = "registry.terraform.io/svalabs/forgejo"
    }
  }
}

provider "forgejo" {
  host = "http://localhost:3000"
}

resource "forgejo_organization" "defaults" {
  name = "tftest_defaults"
}
resource "forgejo_organization" "non_defaults" {
  name        = "tftest_non_defaults"
  full_name   = "Terraform Test Org with non-default attributes"
  description = "Purely for testing..."
  website     = "https://forgejo.org/"
  location    = "Mêlée Island"
  visibility  = "private"
}

output "debug_defaults" {
  value = forgejo_organization.defaults
}
output "debug_non_defaults" {
  value = forgejo_organization.non_defaults
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the organization.

### Optional

- `avatar_url` (String) Avatar URL of the organization.
- `description` (String) Description of the organization.
- `full_name` (String) Full name of the organization.
- `location` (String) Location of the organization.
- `visibility` (String) Visibility of the organization. Possible values are 'public' (default), 'limited', or 'private'.
- `website` (String) Website of the organization.

### Read-Only

- `id` (Number) Numeric identifier of the organization.
