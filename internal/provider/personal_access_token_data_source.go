package provider

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &personalAccessTokenDataSource{}
	_ datasource.DataSourceWithConfigure = &personalAccessTokenDataSource{}
)

// personalAccessTokenDataSource is the data source implementation.
type personalAccessTokenDataSource struct {
	client *forgejo.Client
}

// personalAccessTokenDataSourceModel maps the data source schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3#CreateAccessTokenOption
type personalAccessTokenDataSourceModel struct {
	UserID         types.Int64  `tfsdk:"user_id"`
	ID             types.Int64  `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	TokenLastEight types.String `tfsdk:"token_last_eight"`
	Scopes         types.Set    `tfsdk:"scopes"`
}

// Metadata returns the data source type name.
func (d *personalAccessTokenDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_personal_access_token"
}

// Schema defines the schema for the data source.
func (d *personalAccessTokenDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo personal access token data source.",

		Attributes: map[string]schema.Attribute{
			"user_id": schema.Int64Attribute{
				Description: "ID of the user.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the personal access token.",
				Required:    true,
			},
			"id": schema.Int64Attribute{
				Description: "ID of the personal access token.",
				Computed:    true,
			},
			"token_last_eight": schema.StringAttribute{
				Description: "Last eight characters of the personal access token.",
				Computed:    true,
			},
			"scopes": schema.SetAttribute{
				Description: "Scopes of the personal access token.",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *personalAccessTokenDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*forgejo.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf(
				"Expected *forgejo.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *personalAccessTokenDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer un(trace(ctx, "Read personal access token data source"))

	var (
		data personalAccessTokenDataSourceModel
	)

	// Read Terraform configuration data into model
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get personal access token
	token, diags := getPersonalAccessToken(ctx, d.client, data.UserID.ValueInt64(), data.Name.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	data.ID = types.Int64Value(token.ID)
	data.Name = types.StringValue(token.Name)
	data.TokenLastEight = types.StringValue(token.TokenLastEight)

	data.Scopes, diags = types.SetValueFrom(ctx, types.StringType, token.Scopes)
	resp.Diagnostics.Append(diags...)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func getPersonalAccessToken(
	ctx context.Context,
	client *forgejo.Client,
	userID int64,
	tokenName string) (*forgejo.AccessToken, diag.Diagnostics) {

	var (
		diags diag.Diagnostics
		user  userResourceModel
	)

	// Use Forgejo client to get user
	usr, diags := getUserByID(
		ctx,
		client,
		userID,
	)
	if diags.HasError() {
		return nil, diags
	}

	// Map response body to model
	user.from(usr)

	tflog.Info(ctx, "List personal access tokens", map[string]any{
		"user": user.Name.ValueString(),
	})

	// Use Forgejo client to list personal access tokens
	tokens, res, err := client.ListAccessTokens(
		user.Name.ValueString(),
		forgejo.ListAccessTokensOptions{
			ListOptions: forgejo.ListOptions{
				Page: -1,
			},
		},
	)
	if err != nil {
		var msg string
		if res == nil {
			msg = fmt.Sprintf("Unknown error with nil response: %s", err)
		} else {
			tflog.Error(ctx, "Error", map[string]any{
				"status": res.Status,
			})

			switch res.StatusCode {
			case 403:
				msg = fmt.Sprintf(
					"Personal access tokens from user %s forbidden: %s",
					user.Name.String(),
					err,
				)
			case 404:
				msg = fmt.Sprintf(
					"Personal access tokens from user %s not found: %s",
					user.Name.String(),
					err,
				)
			default:
				msg = fmt.Sprintf(
					"Unknown error (status %d): %s",
					res.StatusCode,
					err,
				)
			}
		}
		diags.AddError("Unable to list personal access tokens", msg)
		return nil, diags
	}

	// Search for personal access token with given name
	idx := slices.IndexFunc(tokens, func(t *forgejo.AccessToken) bool {
		return t.Name == tokenName
	})
	if idx == -1 {
		diags.AddError(
			"Unable to find personal access token by name",
			fmt.Sprintf(
				"Personal access token from user %s and name %s not found",
				user.Name.String(),
				tokenName,
			),
		)

		return nil, diags
	}
	return tokens[idx], diags
}

// NewPersonalAccessTokenDataSource is a helper function to simplify the provider implementation.
func NewPersonalAccessTokenDataSource() datasource.DataSource {
	return &personalAccessTokenDataSource{}
}
