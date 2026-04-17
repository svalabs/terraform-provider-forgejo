package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &repositoryActionVariableResource{}
	_ resource.ResourceWithConfigure = &repositoryActionVariableResource{}
)

// repositoryActionVariableResource is the resource implementation.
type repositoryActionVariableResource struct {
	client *forgejo.Client
}

// repositoryActionVariableResourceModel maps the resource schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3#CreateVariableOption
type repositoryActionVariableResourceModel struct {
	RepositoryID types.Int64  `tfsdk:"repository_id"`
	Name         types.String `tfsdk:"name"`
	Data         types.String `tfsdk:"data"`
}

// from is a helper function to load an API struct into Terraform data model.
func (m *repositoryActionVariableResourceModel) from(v *forgejo.ActionVariable) {
	if v == nil {
		return
	}

	// name is omitted here, to maintain the user's configuration casing
	m.Data = types.StringValue(v.Data)
}

// to is a helper function to save Terraform data model into an API struct.
func (m *repositoryActionVariableResourceModel) to(o *forgejo.CreateVariableOption) {
	if o == nil {
		return
	}

	o.Name = m.Name.ValueString()
	o.Data = m.Data.ValueString()
}

// Metadata returns the resource type name.
func (r *repositoryActionVariableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_action_variable"
}

// Schema defines the schema for the resource.
func (r *repositoryActionVariableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo repository action variable resource.",

		Attributes: map[string]schema.Attribute{
			"repository_id": schema.Int64Attribute{
				Description: "Numeric identifier of the repository. Changing this forces a new resource to be created.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
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
func (r *repositoryActionVariableResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *repositoryActionVariableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create repository action variable resource"))

	var (
		repo repositoryResourceModel
		data repositoryActionVariableResourceModel
	)

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get repository
	rep, diags := getRepositoryByID(
		ctx,
		r.client,
		data.RepositoryID.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	repo.from(rep)

	tflog.Info(ctx, "Create repository action variable", map[string]any{
		"repository_id": data.RepositoryID.ValueInt64(),
		"user":          repo.Owner.ValueString(),
		"repo":          repo.Name.ValueString(),
		"name":          data.Name.ValueString(),
		"data":          data.Data.ValueString(),
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

	// Use Forgejo client to create new repository action variable
	res, err := r.client.CreateRepoActionVariable(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
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
					"Repository with owner %s and name %s not found: %s",
					repo.Owner.String(),
					repo.Name.String(),
					err,
				)
			case 409:
				msg = fmt.Sprintf(
					"Action variable with owner %s, repo %s and name %s conflict: %s",
					repo.Owner.String(),
					repo.Name.String(),
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
		resp.Diagnostics.AddError("Unable to create repository action variable", msg)

		return
	}

	// Use Forgejo client to get repository action variable
	variable, res, err := r.client.GetRepoActionVariable(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
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
					"Action variable with owner %s, repo %s and name %s not found: %s",
					repo.Owner.String(),
					repo.Name.String(),
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
		resp.Diagnostics.AddError("Unable to read repository action variable", msg)

		return
	}

	// Map response body to model
	data.from(variable)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *repositoryActionVariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read repository action variable resource"))

	var (
		repo repositoryResourceModel
		data repositoryActionVariableResourceModel
	)

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get repository
	rep, diags := getRepositoryByID(
		ctx,
		r.client,
		data.RepositoryID.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	repo.from(rep)

	tflog.Info(ctx, "Read repository action variable", map[string]any{
		"repository_id": data.RepositoryID.ValueInt64(),
		"user":          repo.Owner.ValueString(),
		"repo":          repo.Name.ValueString(),
		"name":          data.Name.ValueString(),
	})

	// Use Forgejo client to get repository action variable
	variable, res, err := r.client.GetRepoActionVariable(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
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
					"Action variable with owner %s, repo %s and name %s not found: %s",
					repo.Owner.String(),
					repo.Name.String(),
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
		resp.Diagnostics.AddError("Unable to read repository action variable", msg)

		return
	}

	// Map response body to model
	data.from(variable)

	/*
	 * The variable exists, so we re-save the state from the prior state data.
	 * This is to signal to Terraform that the resource still exists without
	 * overriding the user's configuration casing for the name.
	 */
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *repositoryActionVariableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update repository action variable resource"))

	var (
		state repositoryActionVariableResourceModel
		plan  repositoryActionVariableResourceModel
		repo  repositoryResourceModel
	)

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform plan data into model
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get repository
	rep, diags := getRepositoryByID(
		ctx,
		r.client,
		plan.RepositoryID.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	repo.from(rep)

	tflog.Info(ctx, "Update repository action variable", map[string]any{
		"repository_id": plan.RepositoryID.ValueInt64(),
		"user":          repo.Owner.ValueString(),
		"repo":          repo.Name.ValueString(),
		"old_name":      state.Name.ValueString(),
		"new_name":      plan.Name.ValueString(),
		"data":          plan.Data.ValueString(),
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

	// Use Forgejo client to update repository action variable
	res, err := r.client.UpdateRepoActionVariable(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
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
					"Action variable with owner %s, repo %s and name %s not found: %s",
					repo.Owner.String(),
					repo.Name.String(),
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
		resp.Diagnostics.AddError("Unable to update repository action variable", msg)

		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *repositoryActionVariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete repository action variable resource"))

	var (
		repo repositoryResourceModel
		data repositoryActionVariableResourceModel
	)

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get repository
	rep, diags := getRepositoryByID(
		ctx,
		r.client,
		data.RepositoryID.ValueInt64(),
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response body to model
	repo.from(rep)

	tflog.Info(ctx, "Delete repository action variable", map[string]any{
		"repository_id": data.RepositoryID.ValueInt64(),
		"user":          repo.Owner.ValueString(),
		"repo":          repo.Name.ValueString(),
		"name":          data.Name.ValueString(),
	})

	// Use Forgejo client to delete existing repository action variable
	res, err := r.client.DeleteRepoActionVariable(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
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
					"Action variable with owner %s, repo %s and name %s not found: %s",
					repo.Owner.String(),
					repo.Name.String(),
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
		resp.Diagnostics.AddError("Unable to delete repository action variable", msg)

		return
	}
}

// NewRepositoryActionVariableResource is a helper function to simplify the provider implementation.
func NewRepositoryActionVariableResource() resource.Resource {
	return &repositoryActionVariableResource{}
}
