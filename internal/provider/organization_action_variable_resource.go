package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &organizationActionVariableResource{}
	_ resource.ResourceWithConfigure = &organizationActionVariableResource{}
)

// organizationActionVariableResource is the resource implementation.
type organizationActionVariableResource struct {
	client *forgejo.Client
}

// organizationActionVariableResourceModel maps the resource schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3#CreateVariableOption
type organizationActionVariableResourceModel struct {
	Organization   types.String `tfsdk:"organization"`
	OrganizationID types.Int64  `tfsdk:"organization_id"`
	Name           types.String `tfsdk:"name"`
	Data           types.String `tfsdk:"data"`
}

// from is a helper function to load an API struct into Terraform data model.
func (m *organizationActionVariableResourceModel) from(v *forgejo.ActionVariable) {
	if v == nil {
		return
	}

	// name is omitted here, to maintain the user's configuration casing
	m.Data = types.StringValue(v.Data)
}

// to is a helper function to save Terraform data model into an API struct.
func (m *organizationActionVariableResourceModel) to(o *forgejo.CreateVariableOption) {
	if o == nil {
		return
	}

	o.Name = m.Name.ValueString()
	o.Data = m.Data.ValueString()
}

// Metadata returns the resource type name.
func (r *organizationActionVariableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_action_variable"
}

// Schema defines the schema for the resource.
func (r *organizationActionVariableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Forgejo organization action variable resource.

**Note**: The authenticated user must be a member of the managed organization(s) or have administrative privileges!`,

		Attributes: map[string]schema.Attribute{
			"organization": schema.StringAttribute{
				MarkdownDescription: "Name of the owning organization. Changing this forces a new resource to be created. **Note**: One of `organization` or `organization_id` must be specified.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("organization_id"),
					}...),
				},
			},
			"organization_id": schema.Int64Attribute{
				MarkdownDescription: "Numeric identifier of the owning organization. Changing this forces a new resource to be created. **Note**: One of `organization` or `organization_id` must be specified.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplaceIfConfigured(),
				},
				Validators: []validator.Int64{
					int64validator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("organization"),
					}...),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the variable.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"data": schema.StringAttribute{
				Description: "Data of the variable.",
				Required:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *organizationActionVariableResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *organizationActionVariableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create organization action variable resource"))

	var data organizationActionVariableResourceModel

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get organization name from ID if not provided
	if data.Organization.IsNull() || data.Organization.IsUnknown() {
		org, diags := getOrganizationByID(
			ctx,
			r.client,
			data.OrganizationID.ValueInt64(),
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Map response body to model
		data.Organization = types.StringValue(org.UserName)
	} else {
		// Clear organization ID if name is provided
		data.OrganizationID = types.Int64Null()
	}

	tflog.Info(ctx, "Create organization action variable", map[string]any{
		"organization": data.Organization.ValueString(),
		"name":         data.Name.ValueString(),
		"data":         data.Data.ValueString(),
	})

	// Generate API request body from plan
	opts := forgejo.CreateVariableOption{}
	data.to(&opts)

	// Validate API request body
	err := opts.Validate()
	if err != nil {
		resp.Diagnostics.AddError("Input validation error", err.Error())

		return
	}

	// Use Forgejo client to create new organization action variable
	res, err := r.client.CreateOrgActionVariable(
		data.Organization.ValueString(),
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
				msg = fmt.Sprintf("Bad request: %s", err)
			case 404:
				msg = fmt.Sprintf(
					"Organization with name %s not found: %s",
					data.Organization.String(),
					err,
				)
			case 409:
				msg = fmt.Sprintf(
					"Action variable with org %s and name %s conflict: %s",
					data.Organization.String(),
					data.Name.String(),
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
		resp.Diagnostics.AddError("Unable to create organization action variable", msg)

		return
	}

	// Use Forgejo client to get organization action variable
	variable, diags := r.getVariable(
		ctx,
		data.Organization.ValueString(),
		data.Name.ValueString(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	data.from(variable)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *organizationActionVariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read organization action variable resource"))

	var data organizationActionVariableResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get organization action variable
	variable, diags := r.getVariable(
		ctx,
		data.Organization.ValueString(),
		data.Name.ValueString(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	data.from(variable)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *organizationActionVariableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update organization action variable resource"))

	var (
		state organizationActionVariableResourceModel
		plan  organizationActionVariableResourceModel
	)

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform plan data into the model
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get organization name from ID if not provided
	if plan.Organization.IsNull() || plan.Organization.IsUnknown() {
		org, diags := getOrganizationByID(
			ctx,
			r.client,
			plan.OrganizationID.ValueInt64(),
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Map response body to model
		plan.Organization = types.StringValue(org.UserName)
	} else {
		// Clear organization ID if name is provided
		plan.OrganizationID = types.Int64Null()
	}

	tflog.Info(ctx, "Update organization action variable", map[string]any{
		"organization": plan.Organization.ValueString(),
		"old_name":     state.Name.ValueString(),
		"new_name":     plan.Name.ValueString(),
		"data":         plan.Data.ValueString(),
	})

	// Generate API request body from plan
	opts := forgejo.CreateVariableOption{}
	plan.to(&opts)

	// Validate API request body
	err := opts.Validate()
	if err != nil {
		resp.Diagnostics.AddError("Input validation error", err.Error())

		return
	}

	// Use Forgejo client to update organization action variable
	res, err := r.client.UpdateOrgActionVariable(
		state.Organization.ValueString(),
		state.Name.ValueString(),
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
				msg = fmt.Sprintf("Bad request: %s", err)
			case 404:
				msg = fmt.Sprintf(
					"Action variable with org %s and name %s not found: %s",
					state.Organization.String(),
					state.Name.String(),
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
		resp.Diagnostics.AddError("Unable to update organization action variable", msg)

		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *organizationActionVariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete organization action variable resource"))

	var data organizationActionVariableResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Delete organization action variable", map[string]any{
		"organization": data.Organization.ValueString(),
		"name":         data.Name.ValueString(),
	})

	// Use Forgejo client to delete existing organization action variable
	res, err := r.client.DeleteOrgActionVariable(
		data.Organization.ValueString(),
		data.Name.ValueString(),
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
				msg = fmt.Sprintf("Bad request: %s", err)
			case 404:
				msg = fmt.Sprintf(
					"Action variable with org %s and name %s not found: %s",
					data.Organization.String(),
					data.Name.String(),
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
		resp.Diagnostics.AddError("Unable to delete organization action variable", msg)

		return
	}
}

// NewOrganizationActionVariableResource is a helper function to simplify the provider implementation.
func NewOrganizationActionVariableResource() resource.Resource {
	return &organizationActionVariableResource{}
}

// getVariable returns the variable with the given name from the organization.
func (r *organizationActionVariableResource) getVariable(ctx context.Context, org, name string) (*forgejo.ActionVariable, diag.Diagnostics) {
	var diags diag.Diagnostics

	tflog.Info(ctx, "Read organization action variable", map[string]any{
		"org":  org,
		"name": name,
	})

	// Use Forgejo client to get organization action variable
	variable, res, err := r.client.GetOrgActionVariable(
		org,
		name,
	)
	if err == nil {
		return variable, diags
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
		case 400:
			msg = fmt.Sprintf("Bad request: %s", err)
		case 404:
			msg = fmt.Sprintf(
				"Action variable with organization '%s' and name '%s' not found: %s",
				org,
				name,
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
	diags.AddError("Unable to read organization action variable", msg)

	return nil, diags
}
