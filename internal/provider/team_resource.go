package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &teamResource{}
	_ resource.ResourceWithConfigure = &teamResource{}
)

// teamResource is the resource implementation.
type teamResource struct {
	client *forgejo.Client
}

// teamResourceModel maps the resource schema data.
type teamResourceModel struct {
	ID                      types.Int64  `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	OrganizationID          types.Int64  `tfsdk:"organization_id"`
	CanCreateOrgRepo        types.Bool   `tfsdk:"can_create_org_repo"`
	Description             types.String `tfsdk:"description"`
	IncludesAllRepositories types.Bool   `tfsdk:"includes_all_repositories"`
	Permission              types.String `tfsdk:"permission"`
	Units                   types.Set    `tfsdk:"units"`
}

// from is a helper function to populate Terraform data model from an API struct.
func (m *teamResourceModel) from(t *forgejo.Team, ctx context.Context) (diags diag.Diagnostics) {
	if t == nil {
		return diags
	}

	m.ID = types.Int64Value(t.ID)
	m.Name = types.StringValue(t.Name)
	m.Description = types.StringValue(t.Description)
	if t.Organization != nil {
		m.OrganizationID = types.Int64Value(t.Organization.ID)
	}
	m.Permission = types.StringValue(string(t.Permission))
	m.CanCreateOrgRepo = types.BoolValue(t.CanCreateOrgRepo)
	m.IncludesAllRepositories = types.BoolValue(t.IncludesAllRepositories)
	m.Units, diags = types.SetValueFrom(ctx, types.StringType, t.Units)

	return diags
}

// to is a helper function to save Terraform data model into an API struct.
func (m *teamResourceModel) to(o *forgejo.EditTeamOption, ctx context.Context) (diags diag.Diagnostics) {
	if o == nil {
		o = new(forgejo.EditTeamOption)
	}

	o.Name = m.Name.ValueString()
	o.Description = m.Description.ValueStringPointer()
	o.Permission = forgejo.AccessMode(m.Permission.ValueString())
	o.CanCreateOrgRepo = m.CanCreateOrgRepo.ValueBoolPointer()
	o.IncludesAllRepositories = m.IncludesAllRepositories.ValueBoolPointer()
	diags = m.Units.ElementsAs(ctx, &o.Units, false)

	return diags
}

// Metadata returns the resource type name.
func (r *teamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

// Schema defines the schema for the resource.
func (r *teamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `Forgejo team resource.
Note: Managing teams requires administrative privileges!`,

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the team.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the team.",
				Required:    true,
			},
			"organization_id": schema.Int64Attribute{
				Description: "ID of the owning organization.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"can_create_org_repo": schema.BoolAttribute{
				Description: "Can create repositories?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the team.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString(""),
			},
			"includes_all_repositories": schema.BoolAttribute{
				Description: "Has access to all repositories?",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"permission": schema.StringAttribute{
				Description: "Permissions within the owning organization. **Note**: If you set `admin` or `owner` here, make sure to set all units. This is due to an SDK limitation.",
				Computed:    true,
				Optional:    true,
				Default:     stringdefault.StaticString("read"),
				Validators: []validator.String{
					stringvalidator.OneOf(
						"read",
						"write",
						"admin",
						"owner",
					),
				},
			},
			"units": schema.SetAttribute{
				Description: "Set of units. **Note**: If the permission is `admin` or `owner` this should include all units due to an SDK limitation.",
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("repo.code"),
						},
					),
				),
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(
							"repo.code",
							"repo.issues",
							"repo.pulls",
							"repo.ext_issues",
							"repo.wiki",
							"repo.ext_wiki",
							"repo.releases",
							"repo.projects",
							"repo.packages",
							"repo.actions",
						),
					),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *teamResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *teamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create team resource"))

	var data teamResourceModel

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to create new team
	team, diags := createTeam(
		ctx,
		r.client,
		data.OrganizationID,
		data.Name.ValueString(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	data.ID = types.Int64Value(team.ID)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	opts := forgejo.EditTeamOption{}
	diags = data.to(&opts, ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to update existing team
	team, diags = editTeam(
		ctx,
		r.client,
		team.ID,
		opts,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	diags = data.from(team, ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *teamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read team resource"))

	var data teamResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Read team", map[string]any{
		"name":            data.Name.ValueString(),
		"organization_id": data.OrganizationID.ValueInt64(),
	})

	// Use Forgejo client to read existing team
	team, diags := getOrgTeamByName(
		ctx,
		r.client,
		data.OrganizationID,
		data.Name,
	)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Map response body to model
	diags = data.from(team, ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *teamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update team resource"))

	var data teamResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	opts := forgejo.EditTeamOption{}
	diags = data.to(&opts, ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to update existing team
	team, diags := editTeam(
		ctx,
		r.client,
		data.ID.ValueInt64(),
		opts,
	)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Map response body to model
	diags = data.from(team, ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *teamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete team resource"))

	var data teamResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Delete team", map[string]any{
		"name":            data.Name.ValueString(),
		"organization_id": data.OrganizationID.ValueInt64(),
	})

	// Use Forgejo client to delete existing team
	res, err := r.client.DeleteTeam(data.ID.ValueInt64())
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
					"Team with name %s not found: %s",
					data.Name.String(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to delete team", msg)

		return
	}
}

// NewTeamResource is a helper function to simplify the provider implementation.
func NewTeamResource() resource.Resource {
	return &teamResource{}
}

// createTeam is a helper function to create a team.
func createTeam(ctx context.Context, client *forgejo.Client, organizationID types.Int64, teamName string) (*forgejo.Team, diag.Diagnostics) {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Create team", map[string]any{
		"name":            teamName,
		"organization_id": organizationID,
	})

	// Generate API request body
	opts := forgejo.CreateTeamOption{
		Name:       teamName,
		Permission: forgejo.AccessMode("read"),
		Units:      []forgejo.RepoUnitType{forgejo.RepoUnitType("repo.code")},
	}

	// Validate API request body
	err := opts.Validate()
	if err != nil {
		diags.AddError("Input validation error", err.Error())

		return nil, diags
	}

	// Use Forgejo client to get organization
	organization, diags := getOrganizationByID(
		ctx,
		client,
		organizationID,
	)
	if diags.HasError() {
		return nil, diags
	}

	// Use Forgejo client to create new team
	team, res, err := client.CreateTeam(organization.UserName, opts)
	if err == nil {
		return team, diags
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
		case 403:
			msg = fmt.Sprintf(
				"Team with owner '%s' and name '%s' forbidden: %s",
				organization.UserName,
				teamName,
				err,
			)
		case 404:
			msg = fmt.Sprintf(
				"Organization with name '%s' not found: %s",
				organization.UserName,
				err,
			)
		case 422:
			msg = fmt.Sprintf("Input validation error: %s", err)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
	}
	diags.AddError("Unable to create team", msg)

	return nil, diags
}

// editTeam is a helper function to update an existing team.
func editTeam(ctx context.Context, client *forgejo.Client, teamID int64, opts forgejo.EditTeamOption) (*forgejo.Team, diag.Diagnostics) {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Update team", map[string]any{
		"id": teamID,
	})

	// Validate API request body
	err := opts.Validate()
	if err != nil {
		diags.AddError("Input validation error", err.Error())

		return nil, diags
	}

	// Use Forgejo client to update existing team
	res, err := client.EditTeam(teamID, opts)
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
					"Team with ID %d not found: %s",
					teamID,
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		diags.AddError("Unable to update team", msg)

		return nil, diags
	}

	// Use Forgejo client to fetch updated team
	team, res, err := client.GetTeam(teamID)
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
					"Team with ID %d not found: %s",
					teamID,
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		diags.AddError("Unable to read team", msg)

		return nil, diags
	}

	return team, diags
}
