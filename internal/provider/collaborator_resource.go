package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	_ resource.Resource              = &collaboratorResource{}
	_ resource.ResourceWithConfigure = &collaboratorResource{}
)

// collaboratorResource is the resource implementation.
type collaboratorResource struct {
	client *forgejo.Client
}

// collaboratorResourceModel maps the resource schema data.
type collaboratorResourceModel struct {
	RepositoryID types.Int64  `tfsdk:"repository_id"`
	User         types.String `tfsdk:"user"`
	Permission   types.String `tfsdk:"permission"`
}

// Metadata returns the resource type name.
func (r *collaboratorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collaborator"
}

// Schema defines the schema for the resource.
func (r *collaboratorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo collaborator resource.",

		Attributes: map[string]schema.Attribute{
			"repository_id": schema.Int64Attribute{
				Description: "Numeric identifier of the repository.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"user": schema.StringAttribute{
				Description: "Username of the collaborator.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"permission": schema.StringAttribute{
				Description: "Repository permissions of the collaborator. Must be one of 'read', 'write', 'admin'.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"read",
						"write",
						"admin",
					),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *collaboratorResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *collaboratorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create collaborator resource"))

	var (
		repo repositoryResourceModel
		data collaboratorResourceModel
	)

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Get repository by id", map[string]any{
		"id": data.RepositoryID.ValueInt64(),
	})

	// Use Forgejo client to get repository by id
	rep, res, err := r.client.GetRepoByID(data.RepositoryID.ValueInt64())
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf(
				"Repository with id %d not found: %s",
				data.RepositoryID.ValueInt64(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to get repository by id", msg)

		return
	}

	// Map response body to model
	repo.from(rep)

	tflog.Info(ctx, "Add collaborator to repository", map[string]any{
		"owner":        repo.Owner.ValueString(),
		"repo":         repo.Name.ValueString(),
		"collaborator": data.User.ValueString(),
		"permission":   data.Permission.ValueString(),
	})

	// Generate API request body from plan
	am := forgejo.AccessMode(data.Permission.ValueString())
	opts := forgejo.AddCollaboratorOption{Permission: &am}

	// Validate API request body
	err = opts.Validate()
	if err != nil {
		resp.Diagnostics.AddError("Input validation error", err.Error())

		return
	}

	// Use Forgejo client to add new collaborator
	res, err = r.client.AddCollaborator(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		data.User.ValueString(),
		opts,
	)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 403:
			msg = fmt.Sprintf(
				"Collaborator with user %s repo %s and name %s forbidden: %s",
				repo.Owner.String(),
				repo.Name.String(),
				data.User.String(),
				err,
			)
		case 404:
			msg = fmt.Sprintf(
				"Collaborator with user %s repo %s and name %s not found: %s",
				repo.Owner.String(),
				repo.Name.String(),
				data.User.String(),
				err,
			)
		case 422:
			msg = fmt.Sprintf("Input validation error: %s", err)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to add collaborator", msg)

		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *collaboratorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read collaborator resource"))

	var (
		repo repositoryResourceModel
		data collaboratorResourceModel
	)

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Get repository by id", map[string]any{
		"id": data.RepositoryID.ValueInt64(),
	})

	// Use Forgejo client to get repository by id
	rep, res, err := r.client.GetRepoByID(data.RepositoryID.ValueInt64())
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf(
				"Repository with id %d not found: %s",
				data.RepositoryID.ValueInt64(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to get repository by id", msg)

		return
	}

	// Map response body to model
	repo.from(rep)

	tflog.Info(ctx, "Get collaborator permission", map[string]any{
		"owner":        repo.Owner.ValueString(),
		"repo":         repo.Name.ValueString(),
		"collaborator": data.User.ValueString(),
	})

	// Use Forgejo client to get collaborator permission
	perms, res, err := r.client.CollaboratorPermission(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		data.User.ValueString(),
	)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 403:
			msg = fmt.Sprintf(
				"Collaborator with user %s repo %s and name %s forbidden: %s",
				repo.Owner.String(),
				repo.Name.String(),
				data.User.String(),
				err,
			)
		case 404:
			msg = fmt.Sprintf(
				"Collaborator with user %s repo %s and name %s not found: %s",
				repo.Owner.String(),
				repo.Name.String(),
				data.User.String(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to get collaborator permission", msg)

		return
	}

	// Map response body to model
	data.Permission = types.StringValue(string(perms.Permission))

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *collaboratorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update collaborator resource"))

	var (
		repo repositoryResourceModel
		data collaboratorResourceModel
	)

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Get repository by id", map[string]any{
		"id": data.RepositoryID.ValueInt64(),
	})

	// Use Forgejo client to get repository by id
	rep, res, err := r.client.GetRepoByID(data.RepositoryID.ValueInt64())
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf(
				"Repository with id %d not found: %s",
				data.RepositoryID.ValueInt64(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to get repository by id", msg)

		return
	}

	// Map response body to model
	repo.from(rep)

	tflog.Info(ctx, "Update collaborator", map[string]any{
		"owner":        repo.Owner.ValueString(),
		"repo":         repo.Name.ValueString(),
		"collaborator": data.User.ValueString(),
		"permission":   data.Permission.ValueString(),
	})

	// Generate API request body from plan
	am := forgejo.AccessMode(data.Permission.ValueString())
	opts := forgejo.AddCollaboratorOption{Permission: &am}

	// Validate API request body
	err = opts.Validate()
	if err != nil {
		resp.Diagnostics.AddError("Input validation error", err.Error())

		return
	}

	res, err = r.client.AddCollaborator(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		data.User.ValueString(),
		opts,
	)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 403:
			msg = fmt.Sprintf(
				"Collaborator with user %s repo %s and name %s forbidden: %s",
				repo.Owner.String(),
				repo.Name.String(),
				data.User.String(),
				err,
			)
		case 404:
			msg = fmt.Sprintf(
				"Collaborator with user %s repo %s and name %s not found: %s",
				repo.Owner.String(),
				repo.Name.String(),
				data.User.String(),
				err,
			)
		case 422:
			msg = fmt.Sprintf("Input validation error: %s", err)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to update collaborator", msg)

		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *collaboratorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete collaborator resource"))

	var (
		repo repositoryResourceModel
		data collaboratorResourceModel
	)

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Get repository by id", map[string]any{
		"id": data.RepositoryID.ValueInt64(),
	})

	// Use Forgejo client to get repository by id
	rep, res, err := r.client.GetRepoByID(data.RepositoryID.ValueInt64())
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf(
				"Repository with id %d not found: %s",
				data.RepositoryID.ValueInt64(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to get repository by id", msg)

		return
	}

	// Map response body to model
	repo.from(rep)

	tflog.Info(ctx, "Delete collaborator from repository", map[string]any{
		"owner":        repo.Owner.ValueString(),
		"repo":         repo.Name.ValueString(),
		"collaborator": data.User.ValueString(),
	})

	// Use Forgejo client to delete existing collaborator
	res, err = r.client.DeleteCollaborator(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		data.User.ValueString(),
	)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf(
				"Collaborator with user %s repo %s and name %s not found: %s",
				repo.Owner.String(),
				repo.Name.String(),
				data.User.String(),
				err,
			)
		case 422:
			msg = fmt.Sprintf("Input validation error: %s", err)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to delete collaborator", msg)

		return
	}
}

// NewCollaboratorResource is a helper function to simplify the provider implementation.
func NewCollaboratorResource() resource.Resource {
	return &collaboratorResource{}
}
