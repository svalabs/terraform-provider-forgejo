package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &personalAccessTokenResource{}
	_ resource.ResourceWithConfigure = &personalAccessTokenResource{}
)

// personalAccessTokenResource is the resource implementation.
type personalAccessTokenResource struct {
	client *forgejo.Client
}

// personalAccessTokenResourceModel maps the resource schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3#AccessToken
type personalAccessTokenResourceModel struct {
	UserID         types.Int64  `tfsdk:"user_id"`
	ID             types.Int64  `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Token          types.String `tfsdk:"token"`
	TokenLastEight types.String `tfsdk:"token_last_eight"`
	Scopes         types.Set    `tfsdk:"scopes"`
}

// from is a helper function to load an API struct into Terraform data model.
func (m *personalAccessTokenResourceModel) from(ctx context.Context, t *forgejo.AccessToken) (diags diag.Diagnostics) {
	if t == nil {
		return diags
	}

	m.ID = types.Int64Value(t.ID)
	m.Name = types.StringValue(t.Name)
	m.Token = types.StringValue(t.Token)
	m.TokenLastEight = types.StringValue(t.TokenLastEight)

	var d diag.Diagnostics
	m.Scopes, d = types.SetValueFrom(ctx, types.StringType, t.Scopes)
	diags.Append(d...)

	return diags
}

// to is a helper function to save Terraform data model into an API struct.
func (m *personalAccessTokenResourceModel) to(ctx context.Context, o *forgejo.CreateAccessTokenOption) (diags diag.Diagnostics) {
	if o == nil {
		return diags
	}

	o.Name = m.Name.ValueString()

	scopes := make([]forgejo.AccessTokenScope, 0, len(m.Scopes.Elements()))
	d := m.Scopes.ElementsAs(ctx, &scopes, false)
	diags.Append(d...)

	o.Scopes = scopes

	return diags
}

// Metadata returns the resource type name.
func (r *personalAccessTokenResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_personal_access_token"
}

// Schema defines the schema for the resource.
func (r *personalAccessTokenResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Forgejo repository personal access token resource.

**Note**: Due to an upstream limitation, one cannot create access tokens when authorised with access tokens. Use basic-auth instead.`,

		Attributes: map[string]schema.Attribute{
			"user_id": schema.Int64Attribute{
				Description: "ID of the user.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the personal access token.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"scopes": schema.SetAttribute{
				Description: "Scopes of the personal access token.",
				Required:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the personal access token.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"token": schema.StringAttribute{
				Description: "The personal access token.",
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"token_last_eight": schema.StringAttribute{
				Description: "Last eight characters of the personal access token.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *personalAccessTokenResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*forgejo.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected *forgejo.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *personalAccessTokenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create personal access token resource"))

	var (
		user userResourceModel
		data personalAccessTokenResourceModel
	)

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get user
	usr, diags := getUserByID(
		ctx,
		r.client,
		data.UserID.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	user.from(usr)

	tflog.Info(ctx, "Create personal access token", map[string]any{
		"user":   user.Name.ValueString(),
		"name":   data.Name.ValueString(),
		"scopes": data.Scopes,
	})

	// Generate API request body from plan
	opts := forgejo.CreateAccessTokenOption{}
	diags = data.to(ctx, &opts)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to create new personal access token
	token, res, err := r.client.CreateAccessToken(
		user.Name.ValueString(),
		opts,
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
			case 400:
				msg = fmt.Sprintf(
					"Bad request: %s",
					err,
				)
			case 401:
				msg = fmt.Sprintf(
					"Authentication method is not allowed, use basic-auth: %s",
					err,
				)
			case 403:
				msg = fmt.Sprintf(
					"User %s forbidden: %s",
					user.Name.String(),
					err,
				)
			case 404:
				msg = fmt.Sprintf(
					"User %s not found: %s",
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
		resp.Diagnostics.AddError("Unable to create personal access token", msg)

		return
	}

	// Map response body to model
	diags = data.from(ctx, token)
	resp.Diagnostics.Append(diags...)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *personalAccessTokenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read personal access token resource"))

	var (
		data personalAccessTokenResourceModel
	)

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get personal access token
	token, diags := getPersonalAccessToken(ctx, r.client, data.UserID.ValueInt64(), data.Name.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	data.from(ctx, token)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *personalAccessTokenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update personal access token resource"))

	/*
	 * Personal access tokens can not be updated in-place. All writable attributes have
	 * 'RequiresReplace' plan modifier set.
	 */
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *personalAccessTokenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete personal access token resource"))

	var (
		user userResourceModel
		data personalAccessTokenResourceModel
	)

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get user
	usr, diags := getUserByID(
		ctx,
		r.client,
		data.UserID.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	user.from(usr)

	tflog.Info(ctx, "Delete personal access token", map[string]any{
		"user": user.Name.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Use Forgejo client to delete existing personal access token
	res, err := r.client.DeleteAccessToken(
		user.Name.ValueString(),
		data.ID.ValueInt64(),
	)
	if err == nil {
		return
	}

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
				"User %s forbidden: %s",
				user.Name.String(),
				err,
			)
		case 404:
			msg = fmt.Sprintf(
				"User %s not found: %s",
				user.Name.String(),
				err,
			)
		case 422:
			msg = fmt.Sprintf("Input validation error: %s", err)
		default:
			msg = fmt.Sprintf(
				"Unknown error (status %d): %s",
				res.StatusCode,
				err,
			)
		}
	}
	resp.Diagnostics.AddError("Unable to delete deploy key", msg)
}

// NewpersonalAccessTokenResource is a helper function to simplify the provider implementation.
func NewPersonalAccessTokenResource() resource.Resource {
	return &personalAccessTokenResource{}
}
