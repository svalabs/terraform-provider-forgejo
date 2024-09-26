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
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &forgejoProvider{
			version: version,
		}
	}
}

// Metadata returns the provider type name.
func (p *forgejoProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "forgejo"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *forgejoProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *forgejoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Forgejo client")

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
			"Unknown Forgejo API Host",
			"The provider cannot create the Forgejo API client as there is an unknown configuration value for the Forgejo API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the FORGEJO_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Forgejo API Username",
			"The provider cannot create the Forgejo API client as there is an unknown configuration value for the Forgejo API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the FORGEJO_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Forgejo API Password",
			"The provider cannot create the Forgejo API client as there is an unknown configuration value for the Forgejo API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the FORGEJO_PASSWORD environment variable.",
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

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Forgejo API Host",
			"The provider cannot create the Forgejo API client as there is a missing or empty value for the Forgejo API host. "+
				"Set the host value in the configuration or use the FORGEJO_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Forgejo API Username",
			"The provider cannot create the Forgejo API client as there is a missing or empty value for the Forgejo API username. "+
				"Set the username value in the configuration or use the FORGEJO_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Forgejo API Password",
			"The provider cannot create the Forgejo API client as there is a missing or empty value for the Forgejo API password. "+
				"Set the password value in the configuration or use the FORGEJO_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "forgejo_host", host)
	ctx = tflog.SetField(ctx, "forgejo_username", username)
	ctx = tflog.SetField(ctx, "forgejo_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "forgejo_password")

	tflog.Debug(ctx, "Creating Forgejo client")

	// Create a new Forgejo client using the configuration values
	client, err := forgejo.NewClient(host, forgejo.SetBasicAuth(username, password))
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

	tflog.Info(ctx, "Configured Forgejo client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *forgejoProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewOrganizationDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *forgejoProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
