package provider

import (
	"context"
	"fmt"
	"slices"
	"strings"

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
	_ resource.Resource              = &repositoryActionSecretResource{}
	_ resource.ResourceWithConfigure = &repositoryActionSecretResource{}
)

// repositoryActionSecretResource is the resource implementation.
type repositoryActionSecretResource struct {
	client *forgejo.Client
}

// repositoryActionSecretResourceModel maps the resource schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2#CreateSecretOption
type repositoryActionSecretResourceModel struct {
	RepositoryID types.Int64  `tfsdk:"repository_id"`
	Name         types.String `tfsdk:"name"`
	Data         types.String `tfsdk:"data"`
}

// Metadata returns the resource type name.
func (r *repositoryActionSecretResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_action_secret"
}

// Schema defines the schema for the resource.
func (r *repositoryActionSecretResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo repository action secret resource.",

		Attributes: map[string]schema.Attribute{
			"repository_id": schema.Int64Attribute{
				Description: "Numeric identifier of the repository.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the secret.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"data": schema.StringAttribute{
				Description: "The data of the secret.",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *repositoryActionSecretResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *repositoryActionSecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create repository action secret resource"))

	var (
		repo repositoryResourceModel
		data repositoryActionSecretResourceModel
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

	tflog.Info(ctx, "Create repository action secret", map[string]any{
		"user": repo.Owner.ValueString(),
		"repo": repo.Name.ValueString(),
		"name": data.Name.ValueString(),
		"data": data.Data.ValueString(),
	})

	// Generate API request body from plan
	opts := forgejo.CreateSecretOption{
		Name: data.Name.ValueString(),
		Data: data.Data.ValueString(),
	}

	// Validate API request body
	err = opts.Validate()
	if err != nil {
		resp.Diagnostics.AddError("Input validation error", err.Error())

		return
	}

	// Use Forgejo client to create new repository action secret
	res, err = r.client.CreateRepoActionSecret(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		opts,
	)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 400:
			msg = fmt.Sprintf("Generic error: %s", err)
		case 404:
			msg = fmt.Sprintf(
				"Repository with owner %s and name %s not found: %s",
				repo.Owner.String(),
				repo.Name.String(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to create repository action secret", msg)

		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *repositoryActionSecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read repository action secret resource"))

	var (
		repo repositoryResourceModel
		data repositoryActionSecretResourceModel
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

	tflog.Info(ctx, "List repository action secrets", map[string]any{
		"user": repo.Owner.ValueString(),
		"repo": repo.Name.ValueString(),
		"name": data.Name.ValueString(),
	})

	// Use Forgejo client to list repository action secrets
	secrets, res, err := r.client.ListRepoActionSecret(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		forgejo.ListRepoActionSecretOption{},
	)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf(
				"Repository action secrets with user %s and repo %s not found: %s",
				repo.Owner.String(),
				repo.Name.String(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to list repository action secrets", msg)

		return
	}

	// Search for repository action secrets with given name
	idx := slices.IndexFunc(secrets, func(s *forgejo.Secret) bool {
		return strings.EqualFold(s.Name, data.Name.ValueString())
	})
	if idx == -1 {
		resp.Diagnostics.AddError(
			"Unable to get repository action secret by name",
			fmt.Sprintf(
				"Repository action secret with user %s repo %s and name %s not found.",
				repo.Owner.String(),
				repo.Name.String(),
				data.Name.String(),
			),
		)

		return
	}

	/*
	 * The secret exists, so we re-save the state from the prior state data.
	 * This is to signal to Terraform that the resource still exists without
	 * overriding the user's configuration casing for the name.
	 */
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *repositoryActionSecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update repository action secret resource"))

	var (
		repo repositoryResourceModel
		data repositoryActionSecretResourceModel
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

	tflog.Info(ctx, "Update repository action secret", map[string]any{
		"user": repo.Owner.ValueString(),
		"repo": repo.Name.ValueString(),
		"name": data.Name.ValueString(),
		"data": data.Data.ValueString(),
	})

	// Generate API request body from plan
	opts := forgejo.CreateSecretOption{
		Name: data.Name.ValueString(),
		Data: data.Data.ValueString(),
	}

	// Validate API request body
	err = opts.Validate()
	if err != nil {
		resp.Diagnostics.AddError("Input validation error", err.Error())

		return
	}

	// Use Forgejo client to update repository action secret
	res, err = r.client.CreateRepoActionSecret(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		opts,
	)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 400:
			msg = fmt.Sprintf("Generic error: %s", err)
		case 404:
			msg = fmt.Sprintf(
				"Repository with owner %s and name %s not found: %s",
				repo.Owner.String(),
				repo.Name.String(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to update repository action secret", msg)

		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *repositoryActionSecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete repository action secret resource"))

	var (
		repo repositoryResourceModel
		data repositoryActionSecretResourceModel
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

	tflog.Info(ctx, "Delete repository action secret", map[string]any{
		"user": repo.Owner.ValueString(),
		"repo": repo.Name.ValueString(),
		"name": data.Name.ValueString(),
	})

	resp.Diagnostics.AddWarning(
		"Resource cannot be deleted from Forgejo",
		fmt.Sprintf(
			"The Forgejo SDK does not currently support deleting repository action secrets. "+
				"Secret with owner %s repo %s and name %s will be removed from Terraform state, but will remain in Forgejo.",
			repo.Owner.String(),
			repo.Name.String(),
			data.Name.String(),
		),
	)
}

// NewRepositoryActionSecretResource is a helper function to simplify the provider implementation.
func NewRepositoryActionSecretResource() resource.Resource {
	return &repositoryActionSecretResource{}
}
