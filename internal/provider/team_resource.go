package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
	Organization            types.String `tfsdk:"organization"`
	CanCreateOrgRepo        types.Bool   `tfsdk:"can_create_org_repo"`
	Description             types.String `tfsdk:"description"`
	ImportIfExists          types.Bool   `tfsdk:"import_if_exists"`
	IncludesAllRepositories types.Bool   `tfsdk:"includes_all_repositories"`
	Permission              types.String `tfsdk:"permission"`
	Units                   types.Set    `tfsdk:"units"`
}

// from is a helper function to populate Terraform data model from an API struct.
func (m *teamResourceModel) from(t *forgejo.Team, ctx context.Context) (diags diag.Diagnostics) {
	m.ID = types.Int64Value(t.ID)
	m.Name = types.StringValue(t.Name)
	m.Description = types.StringValue(t.Description)
	m.Organization = types.StringValue(t.Organization.UserName)
	m.Permission = types.StringValue(string(t.Permission))
	m.CanCreateOrgRepo = types.BoolValue(t.CanCreateOrgRepo)
	m.IncludesAllRepositories = types.BoolValue(t.IncludesAllRepositories)
	m.Units, diags = types.SetValueFrom(ctx, types.StringType, t.Units)

	return
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
	diags = m.Units.ElementsAs(ctx, o.Units, false)

	return
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
			"organization": schema.StringAttribute{
				Description: "Name of the owning organization.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"can_create_org_repo": schema.BoolAttribute{
				Description: "Can create repositories?",
				Computed:    true,
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the team.",
				Computed:    true,
				Optional:    true,
			},
			"import_if_exists": schema.BoolAttribute{
				Description: "Import the team if it already exists?",
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"includes_all_repositories": schema.BoolAttribute{
				Description: "Has access to all repositories?",
				Computed:    true,
				Optional:    true,
			},
			"permission": schema.StringAttribute{
				Description: "Organization permission.",
				Computed:    true,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"none",
						"read",
						"write",
						"admin",
						"owner",
					),
				},
			},
			"units": schema.SetAttribute{
				Description: "Set of units.",
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
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

	team, err := getOrgTeamByName(ctx, r.client, data.Organization, data.Name)
	if err != nil {
		resp.Diagnostics.AddError("Unable to get team by name", err.Error())
		return
	}

	if team != nil {
		if !data.ImportIfExists.ValueBool() {
			resp.Diagnostics.AddError(fmt.Sprintf("Team '%s' already exists:", data.Name.ValueString()), err.Error())
			return
		}
	} else {
		team, err = createTeam(ctx, r.client, data.Organization.ValueString(), data.Name.ValueString())
		resp.Diagnostics.AddError("Team creation error", err.Error())
		return
	}

	tflog.Info(ctx, "Add team to org", map[string]any{
		"team":         data.Name.ValueString(),
		"organization": data.Organization.ValueString(),
	})

	opts := forgejo.EditTeamOption{}
	diags = data.to(&opts, ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = editTeam(ctx, r.client, team.ID, opts)
	if err != nil {
		resp.Diagnostics.AddError("Unable to edit team", err.Error())
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

	team, err := getOrgTeamByName(ctx, r.client, data.Organization, data.Name)
	if err != nil {
		resp.Diagnostics.AddError("Unable to get team by name", err.Error())
		return
	}

	tflog.Info(ctx, "Get team from org", map[string]any{
		"team":         team.Name,
		"organization": team.Organization.UserName,
	})

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

	// Read Terraform plan data into the model.
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Update team from org", map[string]any{
		"team":         data.Name.ValueString(),
		"organization": data.Organization.ValueString(),
	})

	opts := forgejo.EditTeamOption{}
	diags = data.to(&opts, ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := editTeam(ctx, r.client, data.ID.ValueInt64(), opts)
	if err != nil {
		resp.Diagnostics.AddError("Unable to edit team", err.Error())
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

	tflog.Info(ctx, "Delete team from repository", map[string]any{
		"team":         data.Name.ValueString(),
		"organization": data.Organization.ValueString(),
	})

	// Use Forgejo client to delete existing team.
	res, err := r.client.DeleteTeam(data.ID.ValueInt64())
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		switch res.StatusCode {
		case 404:
			err = fmt.Errorf(
				"Team with name '%s' not found: %s",
				data.Name.String(),
				err,
			)
		default:
			err = fmt.Errorf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to delete team", err.Error())
		return
	}
}

// NewTeamResource is a helper function to simplify the provider implementation.
func NewTeamResource() resource.Resource {
	return &teamResource{}
}

func createTeam(ctx context.Context, client *forgejo.Client, organization string, teamName string) (team *forgejo.Team, err error) {
	opts := forgejo.CreateTeamOption{
		Name: teamName,
	}

	err = opts.Validate()
	if err != nil {
		err = fmt.Errorf("Input validation error: %s", err.Error())
		return
	}

	team, resp, err := client.CreateTeam(organization, opts)
	if err == nil {
		return
	}

	tflog.Error(ctx, "Error", map[string]any{
		"status": resp.Status,
	})

	switch resp.StatusCode {
	case 403:
		err = fmt.Errorf(
			"Team with owner '%s' and name '%s' forbidden: %s",
			organization,
			teamName,
			err,
		)
	case 404:
		err = fmt.Errorf(
			"Organization with name '%s' not found: %s",
			organization,
			err,
		)
	case 422:
		err = fmt.Errorf("Input validation error: %s", err)
	default:
		err = fmt.Errorf("Unknown error: %s", err)
	}
	return
}

func editTeam(ctx context.Context, client *forgejo.Client, teamID int64, opts forgejo.EditTeamOption) (err error) {
	err = opts.Validate()
	if err != nil {
		err = fmt.Errorf("Input validation error: %s", err)
		return
	}

	resp, err := client.EditTeam(teamID, opts)
	if err == nil {
		return
	}

	tflog.Error(ctx, "Error", map[string]any{
		"status": resp.Status,
	})

	switch resp.StatusCode {
	case 404:
		err = fmt.Errorf(
			"Team with ID '%d' not found: %s",
			teamID,
			err,
		)
	default:
		err = fmt.Errorf("Unknown error: %s", err)
	}
	return
}
