// Package provider exposes the Forgejo Terraform provider constructor for
// embedding the provider in other Go programs.
//
// The import path and the signature and behavior of [New] follow Semantic
// Versioning and won't change in a backwards-incompatible way within a major
// version. The concrete type behind the returned [frameworkprovider.Provider],
// the required Terraform Plugin Framework version, and the minimum Go version
// are not covered by this guarantee.
//
// The provider falls back to the FORGEJO_HOST, FORGEJO_API_TOKEN,
// FORGEJO_USERNAME, and FORGEJO_PASSWORD environment variables of the host
// process for any value not set in the provider configuration.
package provider

import (
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"

	internalprovider "github.com/svalabs/terraform-provider-forgejo/internal/provider"
)

// New returns a factory function that creates Forgejo Terraform provider
// instances, reporting the given version.
//
// The returned function can be passed to providerserver.Serve, or called to
// obtain a provider instance for providerserver.NewProtocol6WithError.
//
// Use [frameworkprovider.Provider] interface methods only; the concrete type
// is an implementation detail.
func New(version string) func() frameworkprovider.Provider {
	return internalprovider.New(version)
}
