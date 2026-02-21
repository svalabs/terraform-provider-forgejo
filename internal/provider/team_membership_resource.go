package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &teamMembershipResource{}
	_ resource.ResourceWithConfigure = &teamMembershipResource{}
)

// teamMembershipResource is the resource implementation.
type teamMembershipResource struct {
	client *forgejo.Client
}

// teamMembershipResourceModel maps the resource schema data.
type teamMembershipResourceModel struct {
	TeamID   types.Int64  `tfsdk:"team_id"`
	UserID   types.Int64  `tfsdk:"user_id"`
	UserName types.String `tfsdk:"user_name"`
}

// Metadata returns the resource type name.
func (r *teamMembershipResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_membership"
}

// Schema defines the schema for the resource.
func (r *teamMembershipResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Forgejo team membership resource.
Note: Managing teams requires administrative privileges!`,

		Attributes: map[string]schema.Attribute{
			"team_id": schema.Int64Attribute{
				Description: "Numeric identifier of the team.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"user_id": schema.Int64Attribute{
				Description: "Numeric identifier of the user.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"user_name": schema.StringAttribute{
				Description: "Name of the user.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *teamMembershipResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *teamMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create team membership resource"))

	var data teamMembershipResourceModel

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get user by ID
	user, diags := getUserByID(ctx, r.client, data.UserID.ValueInt64())
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	tflog.Info(ctx, "Create user's team membership", map[string]any{
		"team_id":   data.TeamID,
		"user_id":   data.UserID,
		"user_name": user.UserName,
	})

	// Use Forgejo client to set a team membership
	diags = setTeamMembership(ctx, r.client, data.TeamID.ValueInt64(), user.UserName)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	data.UserName = types.StringValue(user.UserName)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *teamMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read team membership resource"))

	var data teamMembershipResourceModel

	// Read Terraform prior state into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, diags := getUserByID(ctx, r.client, data.UserID.ValueInt64())
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	tflog.Info(ctx, "Read user's team membership", map[string]any{
		"team_id":   data.TeamID,
		"user_id":   data.UserID,
		"user_name": user.UserName,
	})

	_, diags = getTeamMembership(ctx, r.client,
		data.TeamID.ValueInt64(),
		user.UserName,
	)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	data.UserName = types.StringValue(user.UserName)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *teamMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update team membership resource"))

	/*
	 * Team memberships can not be updated in-place. All writable attributes have
	 * 'RequiresReplace' plan modifier set.
	 */
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *teamMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete team membership resource"))

	var data teamMembershipResourceModel

	// Read Terraform prior state into the model.
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, diags := getUserByID(ctx, r.client, data.UserID.ValueInt64())
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	diags = deleteTeamMembership(ctx, r.client,
		data.TeamID.ValueInt64(),
		user.UserName,
	)
	resp.Diagnostics.Append(diags...)
}

// NewTeamMembershipResource is a helper function to simplify the provider implementation.
func NewTeamMembershipResource() resource.Resource {
	return &teamMembershipResource{}
}

// setTeamMembership adds a user by name to a team by ID.
func setTeamMembership(ctx context.Context, client *forgejo.Client, teamID int64, userName string) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Setting team membership", map[string]any{
		"team_id":   teamID,
		"user_name": userName,
	})

	// Use Forgejo client to set a team membership
	res, err := client.AddTeamMember(teamID, userName)
	if err == nil {
		return diags
	}

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
				"Either User %s or Team with ID %d not found: %s",
				userName,
				teamID,
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
	}
	diags.AddError("Unable to add team membership", msg)
	return diags
}

// deleteTeamMembership deletes a user by name from a team by ID.
func deleteTeamMembership(ctx context.Context, client *forgejo.Client, teamID int64, userName string) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Deleting team membership", map[string]any{
		"team_id":   teamID,
		"user_name": userName,
	})

	res, err := client.RemoveTeamMember(teamID, userName)
	if err == nil {
		return diags
	}

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
				"Either User %s or Team with ID %d not found: %s",
				userName,
				teamID,
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
	}
	diags.AddError("Unable to delete team membership", msg)
	return diags
}
