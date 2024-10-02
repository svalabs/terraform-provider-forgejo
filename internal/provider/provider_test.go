package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"terraform-provider-forgejo/internal/provider"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Forgejo client is properly configured.
	// It is also possible to use the FORGEJO_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `provider "forgejo" { host = "http://localhost:3000" }
`
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"forgejo": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
	if v := os.Getenv("FORGEJO_API_TOKEN"); v == "" {
		t.Fatal("FORGEJO_API_TOKEN must be set for acceptance tests")
	}
}
