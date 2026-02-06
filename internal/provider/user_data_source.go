package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

// userDataSource is the data source implementation.
type userDataSource struct {
	client *forgejo.Client
}

// userDataSourceModel maps the data source schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2#User
type userDataSourceModel struct {
	ID               types.Int64  `tfsdk:"id"`
	Name             types.String `tfsdk:"login"`
	LoginName        types.String `tfsdk:"login_name"`
	SourceID         types.Int64  `tfsdk:"source_id"`
	FullName         types.String `tfsdk:"full_name"`
	Email            types.String `tfsdk:"email"`
	HTMLURL          types.String `tfsdk:"html_url"`
	AvatarURL        types.String `tfsdk:"avatar_url"`
	Language         types.String `tfsdk:"language"`
	Admin            types.Bool   `tfsdk:"admin"`
	LastLogin        types.String `tfsdk:"last_login"`
	Created          types.String `tfsdk:"created_at"`
	Restricted       types.Bool   `tfsdk:"restricted"`
	Active           types.Bool   `tfsdk:"active"`
	ProhibitLogin    types.Bool   `tfsdk:"prohibit_login"`
	Location         types.String `tfsdk:"location"`
	Website          types.String `tfsdk:"website"`
	Description      types.String `tfsdk:"description"`
	Visibility       types.String `tfsdk:"visibility"`
	FollowerCount    types.Int64  `tfsdk:"followers_count"`
	FollowingCount   types.Int64  `tfsdk:"following_count"`
	StarredRepoCount types.Int64  `tfsdk:"starred_repos_count"`
}

// Metadata returns the data source type name.
func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the data source.
func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo user data source.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the user.",
				Computed:    true,
			},
			"login": schema.StringAttribute{
				Description: "Name of the user.",
				Required:    true,
			},
			"login_name": schema.StringAttribute{
				Description: "Login name of the user.",
				Computed:    true,
			},
			"source_id": schema.Int64Attribute{
				Description: "Numeric identifier of the user's authentication source.",
				Computed:    true,
			},
			"full_name": schema.StringAttribute{
				Description: "Full name of the user.",
				Computed:    true,
			},
			"email": schema.StringAttribute{
				Description: "Email address of the user.",
				Computed:    true,
			},
			"avatar_url": schema.StringAttribute{
				Description: "Avatar URL of the user.",
				Computed:    true,
			},
			"html_url": schema.StringAttribute{
				Description: "URL to the user's profile page.",
				Computed:    true,
			},
			"language": schema.StringAttribute{
				Description: "Locale of the user.",
				Computed:    true,
			},
			"admin": schema.BoolAttribute{
				Description: "Is the user an administrator?",
				Computed:    true,
			},
			"last_login": schema.StringAttribute{
				Description: "Time at which the user last logged in.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Time at which the user was created.",
				Computed:    true,
			},
			"restricted": schema.BoolAttribute{
				Description: "Is the user restricted?",
				Computed:    true,
			},
			"active": schema.BoolAttribute{
				Description: "Is the user active?",
				Computed:    true,
			},
			"prohibit_login": schema.BoolAttribute{
				Description: "Are user logins prohibited?",
				Computed:    true,
			},
			"location": schema.StringAttribute{
				Description: "Location of the user.",
				Computed:    true,
			},
			"website": schema.StringAttribute{
				Description: "Website of the user.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the user.",
				Computed:    true,
			},
			"visibility": schema.StringAttribute{
				Description: "Visibility of the user.",
				Computed:    true,
			},
			"followers_count": schema.Int64Attribute{
				Description: "Number of following users.",
				Computed:    true,
			},
			"following_count": schema.Int64Attribute{
				Description: "Number of users followed.",
				Computed:    true,
			},
			"starred_repos_count": schema.Int64Attribute{
				Description: "Number of starred repositories.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer un(trace(ctx, "Read user data source"))

	var data userDataSourceModel

	// Read Terraform configuration data into model
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Read user", map[string]any{
		"name": data.Name.ValueString(),
	})

	// Use Forgejo client to get user by name
	usr, res, err := d.client.GetUserInfo(data.Name.ValueString())
	if err != nil {
		var msg string
		if res == nil {
			msg = fmt.Sprintf("Unknown error with nil response: %s", err)
		} else {
			tflog.Error(ctx, "Error", map[string]any{
				"status": res.Status,
			})

			switch res.StatusCode {
			case 404:
				msg = fmt.Sprintf(
					"User with name %s not found: %s",
					data.Name.String(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to read user", msg)

		return
	}

	// Map response body to model
	data.ID = types.Int64Value(usr.ID)
	data.Name = types.StringValue(usr.UserName)
	data.LoginName = types.StringValue(usr.LoginName)
	data.SourceID = types.Int64Value(usr.SourceID)
	data.FullName = types.StringValue(usr.FullName)
	data.Email = types.StringValue(usr.Email)
	data.HTMLURL = types.StringValue(usr.HTMLURL)
	data.AvatarURL = types.StringValue(usr.AvatarURL)
	data.Language = types.StringValue(usr.Language)
	data.Admin = types.BoolValue(usr.IsAdmin)
	data.LastLogin = types.StringValue(usr.LastLogin.Format(time.RFC3339))
	data.Created = types.StringValue(usr.Created.Format(time.RFC3339))
	data.Restricted = types.BoolValue(usr.Restricted)
	data.Active = types.BoolValue(usr.IsActive)
	data.ProhibitLogin = types.BoolValue(usr.ProhibitLogin)
	data.Location = types.StringValue(usr.Location)
	data.Website = types.StringValue(usr.Website)
	data.Description = types.StringValue(usr.Description)
	data.Visibility = types.StringValue(string(usr.Visibility))
	data.FollowerCount = types.Int64Value(int64(usr.FollowerCount))
	data.FollowingCount = types.Int64Value(int64(usr.FollowingCount))
	data.StarredRepoCount = types.Int64Value(int64(usr.StarredRepoCount))

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// NewUserDataSource is a helper function to simplify the provider implementation.
func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

// getUserByID gets a user by ID.
func getUserByID(ctx context.Context, client *forgejo.Client, userID int64) (*forgejo.User, diag.Diagnostics) {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Getting user by ID", map[string]any{
		"user_id": userID,
	})

	// Use the Forgejo client to get a user by ID
	user, resp, err := client.GetUserByID(userID)
	if err == nil {
		return user, diags
	}

	var msg string
	if resp == nil {
		msg = fmt.Sprintf("Unknown error with nil response: %s", err)
	} else {
		tflog.Error(ctx, "Error", map[string]any{
			"status": resp.Status,
		})

		switch resp.StatusCode {
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
	}
	diags.AddError("Unable to get user by ID", msg)
	return nil, diags
}
