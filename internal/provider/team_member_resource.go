package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v4"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &teamMemberResource{}
	_ resource.ResourceWithConfigure   = &teamMemberResource{}
	_ resource.ResourceWithImportState = &teamMemberResource{}
)

// teamMemberResource is the resource implementation.
type teamMemberResource struct {
	client *forgejo.Client
}

// teamMemberResourceModel maps the resource schema data.
type teamMemberResourceModel struct {
	ID     types.String `tfsdk:"id"`
	TeamID types.Int64  `tfsdk:"team_id"`
	User   types.String `tfsdk:"user"`
}

// Metadata returns the resource type name.
func (r *teamMemberResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team_member"
}

// Schema defines the schema for the resource.
func (r *teamMemberResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo team member resource.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the membership in team_id/username format.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"team_id": schema.Int64Attribute{
				Description: "Numeric identifier of the team. Changing this forces a new resource to be created.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"user": schema.StringAttribute{
				Description: "Username of the team member. Changing this forces a new resource to be created.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *teamMemberResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *teamMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create team member resource"))

	var data teamMemberResourceModel

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to set a team member
	diags = setTeamMember(
		ctx,
		r.client,
		data.TeamID.ValueInt64(),
		data.User.ValueString(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	data.ID = types.StringValue(teamMemberID(data.TeamID.ValueInt64(), data.User.ValueString()))
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *teamMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read team member resource"))

	var data teamMemberResourceModel

	// Read Terraform prior state into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to check a team member
	res, diags := checkTeamMember(
		ctx,
		r.client,
		data.TeamID.ValueInt64(),
		data.User.ValueString(),
	)
	if res != nil && res.StatusCode == 404 {
		resp.State.RemoveResource(ctx)

		return
	}
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	data.ID = types.StringValue(teamMemberID(data.TeamID.ValueInt64(), data.User.ValueString()))
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *teamMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update team member resource"))

	/*
	 * Team members can not be updated in-place. All writable attributes have
	 * 'RequiresReplace' plan modifier set.
	 */
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *teamMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete team member resource"))

	var data teamMemberResourceModel

	// Read Terraform prior state into the model.
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to delete a team member
	diags = deleteTeamMember(
		ctx,
		r.client,
		data.TeamID.ValueInt64(),
		data.User.ValueString(),
	)
	resp.Diagnostics.Append(diags...)
}

// ImportState verifies an existing membership without adding the user to the team.
func (r *teamMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer un(trace(ctx, "Import team member resource"))

	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
		resp.Diagnostics.AddError("Unable to parse import identifier", "Expected import identifier with format: 'team_id/username'.")

		return
	}
	teamID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || teamID <= 0 {
		resp.Diagnostics.AddError("Unable to parse import identifier", "team_id must be a positive integer in 'team_id/username'.")

		return
	}

	_, diags := checkTeamMember(ctx, r.client, teamID, parts[1])
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &teamMemberResourceModel{
		ID:     types.StringValue(teamMemberID(teamID, parts[1])),
		TeamID: types.Int64Value(teamID),
		User:   types.StringValue(parts[1]),
	})...)
}

// NewTeamMemberResource is a helper function to simplify the provider implementation.
func NewTeamMemberResource() resource.Resource {
	return &teamMemberResource{}
}

func teamMemberID(teamID int64, user string) string {
	return fmt.Sprintf("%d/%s", teamID, user)
}

// setTeamMember is a helper function to add a user to a team.
func setTeamMember(ctx context.Context, client *forgejo.Client, teamID int64, userName string) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Create team member", map[string]any{
		"team_id": teamID,
		"user":    userName,
	})

	// Use Forgejo client to set a team member
	res, err := client.AddTeamMember(teamID, userName)
	if err == nil {
		return diags
	}

	// Handle errors
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
				"Either user '%s' or team with ID %d not found: %s",
				userName,
				teamID,
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
	diags.AddError("Unable to create team member", msg)

	return diags
}

// deleteTeamMember is a helper function to delete a user from a team.
func deleteTeamMember(ctx context.Context, client *forgejo.Client, teamID int64, userName string) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Delete team member", map[string]any{
		"team_id":   teamID,
		"user_name": userName,
	})

	// Use Forgejo client to delete a team member
	res, err := client.RemoveTeamMember(teamID, userName)
	if err == nil || (res != nil && res.StatusCode == 404) {
		return diags
	}

	// Handle errors
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
				"Either user '%s' or team with ID %d not found: %s",
				userName,
				teamID,
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
	diags.AddError("Unable to delete team member", msg)

	return diags
}
