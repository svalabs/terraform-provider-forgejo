package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"terraform-provider-forgejo/internal/provider"
)

const (
	// forgejoTest* defines the test environment used to configure the Forgejo
	// provider during acceptance testing.
	forgejoTestHost  = "http://localhost:3000"
	forgejoTestUser  = "tfadmin"
	forgejoTestEmail = forgejoTestUser + "@localhost"

	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Forgejo client is properly configured.
	// It is also possible to use the FORGEJO_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `provider "forgejo" {
		host     = "` + forgejoTestHost + `"
		username = ""
		password = ""
	}
`
	providerBasicAuthConfig = `provider "forgejo" {
		alias     = "basicAuth"
		host      = "` + forgejoTestHost + `"
		api_token = ""
	}
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
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
	if v := os.Getenv("FORGEJO_API_TOKEN"); v == "" {
		t.Fatal("FORGEJO_API_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("FORGEJO_USERNAME"); v == "" {
		t.Fatal("FORGEJO_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("FORGEJO_PASSWORD"); v == "" {
		t.Fatal("FORGEJO_PASSWORD must be set for acceptance tests")
	}
}
