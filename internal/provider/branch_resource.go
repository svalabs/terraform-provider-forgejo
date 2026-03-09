package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &branchResource{}
	_ resource.ResourceWithConfigure   = &branchResource{}
	_ resource.ResourceWithImportState = &branchResource{}
)

// branchResource is the resource implementation.
type branchResource struct {
	client *forgejo.Client
}

// branchResourceModel maps the resource schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2#Branch
type branchResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	OldBranchName types.String `tfsdk:"old_branch_name"`
	Repository    types.String `tfsdk:"repository"`
	Owner         types.String `tfsdk:"owner"`
}

// from is a helper function to load an API struct into Terraform data model.
func (m *branchResourceModel) from(r *forgejo.Branch) {
	if r == nil {
		return
	}

	m.ID = types.StringValue(fmt.Sprintf("%s/%s/%s", m.Owner, m.Repository, r.Name))
	m.Name = types.StringValue(r.Name)
}

// to is a helper function to save Terraform data model into an API struct.
func (m *branchResourceModel) to(o *forgejo.CreateBranchOption) {
	if o == nil {
		return
	}

	o.BranchName = m.Name.ValueString()
	o.OldBranchName = m.OldBranchName.ValueString()
}

// Metadata returns the resource type name.
func (r *branchResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_branch"
}

// Schema defines the schema for the resource.
func (r *branchResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Forgejo branch resource.`,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier of the branch.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"owner": schema.StringAttribute{
				Description: "Owner of the repository (user or organization). Changing this forces a new resource to be created.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"repository": schema.StringAttribute{
				Description: "Name of the repository. Changing this forces a new resource to be created.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the repository.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"old_branch_name": schema.StringAttribute{
				Description: "Name of the old branch to create from (optional).",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *branchResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *branchResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create branch resource"))

	var data branchResourceModel

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		rep *forgejo.Branch
		res *forgejo.Response
		err error
	)

	tflog.Info(ctx, "Creating branch resource", map[string]any{
		"name":  data.Name.ValueString(),
		"owner": data.Owner.ValueString(),
	})

	copts := forgejo.CreateBranchOption{
		BranchName: data.Name.ValueString(),
	}

	err = copts.Validate()
	if err != nil {
		resp.Diagnostics.AddError("Input validation error", err.Error())
		return
	}

	rep, res, err = r.client.CreateBranch(data.Owner.ValueString(), data.Repository.ValueString(), copts)

	if err != nil {
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
					"Repository with owner %s and name %s forbidden: %s",
					data.Owner.String(),
					data.Name.String(),
					err,
				)
			case 404:
				msg = fmt.Sprintf(
					"Repository owner with name %s not found: %s",
					data.Owner.String(),
					err,
				)
			case 409:
				msg = fmt.Sprintf(
					"Repository with name %s already exists: %s",
					data.Name.String(),
					err,
				)
			case 413:
				msg = fmt.Sprintf("Quota exceeded: %s", err)
			case 422:
				msg = fmt.Sprintf("Input validation error: %s", err)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to create branch", msg)

		return
	}

	// Map response body to model
	data.from(rep)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *branchResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read branch resource"))

	var data branchResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		rep *forgejo.Branch
		res *forgejo.Response
		err error
	)

	rep, res, err = r.client.GetRepoBranch(data.Owner.ValueString(), data.Repository.ValueString(), data.Name.ValueString())
	if err != nil {
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
					"Repository with owner %s and name %s forbidden: %s",
					data.Owner.String(),
					data.Name.String(),
					err,
				)
			case 404:
				msg = fmt.Sprintf(
					"Repository owner with name %s not found: %s",
					data.Owner.String(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to read branch", msg)

		return
	}

	// Map response body to model
	data.from(rep)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *branchResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// The code shouldn't go here. Ever.
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *branchResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete branch resource"))

	var data branchResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		v   bool
		res *forgejo.Response
		err error
	)

	tflog.Info(ctx, "Delete branch", map[string]any{
		"owner": data.Owner.ValueString(),
		"name":  data.Name.ValueString(),
	})

	v, res, err = r.client.DeleteRepoBranch(
		data.Owner.ValueString(),
		data.Repository.ValueString(),
		data.Name.ValueString(),
	)

	if err != nil || v == false {
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
					"Repository with owner %s and name %s forbidden: %s",
					data.Owner.String(),
					data.Name.String(),
					err,
				)
			case 404:
				msg = fmt.Sprintf(
					"Repository with owner %s and name %s not found: %s",
					data.Owner.String(),
					data.Name.String(),
					err,
				)
			case 422:
				msg = fmt.Sprintf("Input validation error: %s", err)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to delete branch", msg)

		return
	}
}

// ImportState reads an existing resource and adds it to Terraform state on success.
func (r *branchResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer un(trace(ctx, "Import branch resource"))

	var state branchResourceModel

	// Parse import identifier
	cmp := strings.Split(req.ID, "/")
	if len(cmp) != 3 {
		resp.Diagnostics.AddError(req.ID, "Import ID must be in format 'owner/repository/name'")

		return
	}
	owner, repositoryName, branchName := cmp[0], cmp[1], cmp[2]

	tflog.Info(ctx, "Read repository", map[string]any{
		"owner":      owner,
		"repository": repositoryName,
		"name":       branchName,
	})

	rep, res, err := r.client.GetRepoBranch(owner, repositoryName, branchName)
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
					"Repository with owner '%s' and name '%s' not found: %s",
					owner,
					repositoryName,
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to read repository", msg)

		return
	}

	// Map response body to model
	if resp.Diagnostics.HasError() {
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%s/%s/%s", owner, repositoryName, branchName))
	state.Name = types.StringValue(rep.Name)
	state.Repository = types.StringValue(repositoryName)
	state.Owner = types.StringValue(owner)
	resp.State.Set(ctx, &state)
}

// NewBranchResource is a helper function to simplify the provider implementation.
func NewBranchResource() resource.Resource {
	return &branchResource{}
}
