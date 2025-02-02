package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo"
)

// Ensure the implementation satisfies the expected interfaces.
var _ provider.Provider = &forgejoProvider{}

// forgejoProvider defines the provider implementation.
type forgejoProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// forgejoProviderModel describes the provider data model.
type forgejoProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	ApiToken types.String `tfsdk:"api_token"`
}

// Metadata returns the provider type name.
func (p *forgejoProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "forgejo"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *forgejoProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for Forgejo — self-hosted lightweight software forge",

		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URI for Forgejo API. May also be provided via FORGEJO_HOST environment variable.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username for Forgejo API. May also be provided via FORGEJO_USERNAME environment variable.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for Forgejo API. May also be provided via FORGEJO_PASSWORD environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"api_token": schema.StringAttribute{
				Description: "Token for Forgejo API. May also be provided via FORGEJO_API_TOKEN environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *forgejoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	defer un(trace(ctx, "Configure provider"))

	// Retrieve provider data from configuration
	var config forgejoProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Forgejo Host",
			"The provider cannot create the Forgejo API client as there is an unknown configuration value for the Forgejo host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the FORGEJO_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Forgejo Username",
			"The provider cannot create the Forgejo API client as there is an unknown configuration value for the Forgejo username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the FORGEJO_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Forgejo Password",
			"The provider cannot create the Forgejo API client as there is an unknown configuration value for the Forgejo password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the FORGEJO_PASSWORD environment variable.",
		)
	}

	if config.ApiToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Unknown Forgejo API Token",
			"The provider cannot create the Forgejo API client as there is an unknown configuration value for the Forgejo API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the FORGEJO_API_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("FORGEJO_HOST")
	username := os.Getenv("FORGEJO_USERNAME")
	password := os.Getenv("FORGEJO_PASSWORD")
	token := os.Getenv("FORGEJO_API_TOKEN")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if !config.ApiToken.IsNull() {
		token = config.ApiToken.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Forgejo Host",
			"The provider cannot create the Forgejo API client as there is a missing or empty value for the Forgejo host. "+
				"Set the host value in the configuration or use the FORGEJO_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" && token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Forgejo Username or API Token",
			"The provider cannot create the Forgejo API client as there is a missing or empty value for the Forgejo username and API token. "+
				"Set the username value in the configuration or use the FORGEJO_USERNAME environment variable. "+
				"Alternatively, set the api_token value in the configuration or use the FORGEJO_API_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username != "" && password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Forgejo Password",
			"The provider cannot create the Forgejo API client as there is a missing or empty value for the Forgejo password. "+
				"Set the password value in the configuration or use the FORGEJO_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username != "" && token != "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Both Forgejo Username and API Token",
			"The provider cannot create the Forgejo API client as both, Forgejo username and API token are set. "+
				"Set *either* the username value / FORGEJO_USERNAME environment variable, *or* the api_token value / FORGEJO_API_TOKEN environment variable, but not both.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	var (
		client *forgejo.Client
		err    error
	)

	// Create a new Forgejo client using the configuration values
	if username != "" {
		tflog.Info(ctx, "Create Forgejo client", map[string]any{
			"forgejo_host":     host,
			"forgejo_username": username,
			"forgejo_password": "***",
		})

		client, err = forgejo.NewClient(host, forgejo.SetBasicAuth(username, password))
	}
	if token != "" {
		tflog.Info(ctx, "Create Forgejo client", map[string]any{
			"forgejo_host":      host,
			"forgejo_api_token": "***",
		})

		client, err = forgejo.NewClient(host, forgejo.SetToken(token))
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Forgejo API Client",
			"An unexpected error occurred when creating the Forgejo API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Forgejo Client Error: "+err.Error(),
		)

		return
	}

	// Make the Forgejo client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *forgejoProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDeployKeyDataSource,
		NewOrganizationDataSource,
		NewRepositoryDataSource,
		NewUserDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *forgejoProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDeployKeyResource,
		NewOrganizationResource,
		NewRepositoryResource,
		NewUserResource,
	}
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &forgejoProvider{
			version: version,
		}
	}
}

// Internal helper functions for tracing function execution.
func trace(ctx context.Context, s string) (context.Context, string) {
	tflog.Trace(ctx, s+" - begin")
	return ctx, s
}
func un(ctx context.Context, s string) {
	tflog.Trace(ctx, s+" - end")
}
