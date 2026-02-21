package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"terraform-provider-forgejo/internal/provider"
	"terraform-provider-forgejo/internal/testing/fixtures"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Forgejo client is properly configured.
	// It is also possible to use the FORGEJO_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `provider "forgejo" { }
`
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"forgejo": providerserver.NewProtocol6WithError(
		provider.New("test")(),
	),
}

func testAccPreCheck(t *testing.T) {
	containers, err := fixtures.GetTestContainers(t.Context())
	if err != nil {
		t.Fatalf("Error getting test containers: %s", err)
	}

	forgejoHost, err := containers.ForgejoContainer.GetHost(t.Context())
	if err != nil {
		t.Fatalf("Error getting Forgejo host: %s", err)
	}
	t.Setenv("FORGEJO_HOST", forgejoHost)

	forgejoAPIToken, err := containers.ForgejoContainer.GetAPIToken(t.Context())
	if err != nil {
		t.Fatalf("Error getting Forgejo API token: %s", err)
	}
	t.Setenv("FORGEJO_API_TOKEN", forgejoAPIToken)
}
