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
		Description: "Forgejo team resource.",
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

	// Read Terraform plan data into model.
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	team, err := createTeam(ctx, r.client, data.OrganizationID, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Team creation error", err.Error())
		return
	}

	opts := forgejo.EditTeamOption{}
	diags = data.to(&opts, ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	team, err = editTeam(ctx, r.client, team.ID, opts)
	if err != nil {
		resp.Diagnostics.AddError("Unable to edit team", err.Error())
		return
	}

	diags = data.from(team, ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

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

	tflog.Info(ctx, "Get team from org", map[string]any{
		"team":            data.Name,
		"organization_id": data.OrganizationID,
	})

	team, err := getOrgTeamByName(ctx, r.client, data.OrganizationID, data.Name)
	if err != nil {
		resp.Diagnostics.AddError("Unable to get team by name", err.Error())
		return
	}

	if team == nil {
		resp.Diagnostics.AddError("Unable to get team by name", fmt.Sprintf("No team found called %s within organisation ID %s.", data.Name, data.OrganizationID))
		return
	}

	diags = data.from(team, ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *teamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update team resource"))

	var data teamResourceModel

	// Read Terraform plan data into the model.
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Update team from org", map[string]any{
		"team":            data.Name.ValueString(),
		"organization_id": data.OrganizationID.ValueInt64(),
	})

	opts := forgejo.EditTeamOption{}
	diags = data.to(&opts, ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	team, err := editTeam(ctx, r.client, data.ID.ValueInt64(), opts)
	if err != nil {
		resp.Diagnostics.AddError("Unable to edit team", err.Error())
		return
	}

	if team == nil {
		resp.Diagnostics.AddError("Unable edit team", fmt.Sprintf("No team found called %s within organisation ID %d.", data.Name.ValueString(), data.OrganizationID.ValueInt64()))
		return
	}

	diags = data.from(team, ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state.
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *teamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete team resource"))

	var data teamResourceModel

	// Read Terraform prior state data into the model.
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Delete team from org", map[string]any{
		"team":            data.Name.ValueString(),
		"organization_id": data.OrganizationID.ValueInt64(),
	})

	// Use Forgejo client to delete existing team.
	res, err := r.client.DeleteTeam(data.ID.ValueInt64())
	if err == nil {
		return
	}
	tflog.Error(ctx, "Error", map[string]any{
		"status": res.Status,
	})

	switch res.StatusCode {
	case 404:
		err = fmt.Errorf(
			"the Team with name '%s' was not found: %s",
			data.Name.String(),
			err,
		)
	default:
		err = fmt.Errorf("unknown error: %s", err)
	}
	resp.Diagnostics.AddError("Unable to delete team", err.Error())
}

// NewTeamResource is a helper function to simplify the provider implementation.
func NewTeamResource() resource.Resource {
	return &teamResource{}
}

func createTeam(ctx context.Context, client *forgejo.Client, organizationID types.Int64, teamName string) (team *forgejo.Team, err error) {
	tflog.Info(ctx, "Add team to org", map[string]any{
		"team":            teamName,
		"organization_id": organizationID,
	})

	opts := forgejo.CreateTeamOption{
		Name:       teamName,
		Permission: forgejo.AccessMode("read"),
		Units:      []forgejo.RepoUnitType{forgejo.RepoUnitType("repo.code")},
	}

	err = opts.Validate()
	if err != nil {
		err = fmt.Errorf("input validation error: %s", err.Error())
		return
	}

	organization, err := getOrganizationByID(ctx, client, organizationID)
	if err != nil {
		return nil, err
	}

	if organization == nil {
		err = fmt.Errorf(
			"no Organization with id '%d' was found",
			organizationID.ValueInt64(),
		)
		return nil, err
	}

	team, resp, err := client.CreateTeam(organization.UserName, opts)
	if err == nil {
		return team, nil
	}

	tflog.Error(ctx, "Error", map[string]any{
		"status": resp.Status,
	})

	switch resp.StatusCode {
	case 403:
		err = fmt.Errorf(
			"the Team with owner '%s' and name '%s' is forbidden: %s",
			organization.UserName,
			teamName,
			err,
		)
	case 404:
		err = fmt.Errorf(
			"the Organization with name '%s' was not found: %s",
			organization.UserName,
			err,
		)
	case 422:
		err = fmt.Errorf("input validation error: %s", err)
	default:
		err = fmt.Errorf("unknown error: %s", err)
	}
	return team, err
}

func editTeam(ctx context.Context, client *forgejo.Client, teamID int64, opts forgejo.EditTeamOption) (team *forgejo.Team, err error) {
	tflog.Info(ctx, "Edit team", map[string]any{
		"team_id": teamID,
	})

	err = opts.Validate()
	if err != nil {
		err = fmt.Errorf("input validation error: %s", err)
		return nil, err
	}

	resp, err := client.EditTeam(teamID, opts)
	if err == nil {
		team, resp, err = client.GetTeam(teamID)
		if err == nil {
			return team, nil
		}
	}

	tflog.Error(ctx, "Error", map[string]any{
		"status": resp.Status,
	})

	switch resp.StatusCode {
	case 404:
		err = fmt.Errorf(
			"the Team with ID '%d' was not found: %s",
			teamID,
			err,
		)
	default:
		err = fmt.Errorf("unknown error: %s", err)
	}
	return nil, err
}
