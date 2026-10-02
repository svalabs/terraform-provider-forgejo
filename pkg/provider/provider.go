// Package provider exposes the Forgejo Terraform provider constructor.
package provider

import (
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"

	internalprovider "github.com/svalabs/terraform-provider-forgejo/internal/provider"
)

// New returns a function that creates a Forgejo Terraform provider.
func New(version string) func() frameworkprovider.Provider {
	return internalprovider.New(version)
}
