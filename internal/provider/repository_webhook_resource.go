package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &repositoryWebhookResource{}
	_ resource.ResourceWithConfigure   = &repositoryWebhookResource{}
	_ resource.ResourceWithImportState = &repositoryWebhookResource{}
)

// repositoryWebhookResource is the resource implementation.
type repositoryWebhookResource struct {
	client *forgejo.Client
}

// repositoryWebhookResourceModel maps the resource schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2#Hook
type repositoryWebhookResourceModel struct {
	WebhookID           types.Int64  `tfsdk:"webhook_id"`
	RepositoryID        types.Int64  `tfsdk:"repository_id"`
	Active              types.Bool   `tfsdk:"active"`
	AuthorizationHeader types.String `tfsdk:"authorization_header"`
	BranchFilter        types.String `tfsdk:"branch_filter"`
	Config              types.Map    `tfsdk:"config"`
	CreatedAt           types.String `tfsdk:"created_at"`
	Events              types.Set    `tfsdk:"events"`
	Type                types.String `tfsdk:"type"`
	UpdatedAt           types.String `tfsdk:"updated_at"`
}

// from is a helper function to load an API struct into Terraform data model.
func (m *repositoryWebhookResourceModel) from(h *forgejo.Hook, ctx context.Context) (diags diag.Diagnostics) {
	if h == nil {
		return diags
	}

	var d diag.Diagnostics

	m.WebhookID = types.Int64Value(h.ID)
	m.Active = types.BoolValue(h.Active)
	m.Config, d = types.MapValueFrom(ctx, types.StringType, h.Config)
	diags.Append(d...)
	m.CreatedAt = types.StringValue(h.Created.Format(time.RFC3339))
	m.Events, d = types.SetValueFrom(ctx, types.StringType, h.Events)
	diags.Append(d...)
	m.Type = types.StringValue(h.Type)
	m.UpdatedAt = types.StringValue(h.Updated.Format(time.RFC3339))

	return diags
}

// to is a helper function to save Terraform data model into an API struct.
func (m *repositoryWebhookResourceModel) to(o *forgejo.EditHookOption, ctx context.Context) (diags diag.Diagnostics) {
	if o == nil {
		o = new(forgejo.EditHookOption)
	}

	var d diag.Diagnostics

	o.Active = m.Active.ValueBoolPointer()
	o.AuthorizationHeader = m.AuthorizationHeader.ValueString()
	o.BranchFilter = m.BranchFilter.ValueString()
	d = m.Config.ElementsAs(ctx, &o.Config, false)
	diags.Append(d...)
	d = m.Events.ElementsAs(ctx, &o.Events, false)
	diags.Append(d...)

	return diags
}

// Metadata returns the resource type name.
func (r *repositoryWebhookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_webhook"
}

// Schema defines the schema for the resource.
func (r *repositoryWebhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `Forgejo repository webhook resource.`,

		Attributes: map[string]schema.Attribute{
			"repository_id": schema.Int64Attribute{
				Description: "Numeric identifier of the repository. Changing this forces a new resource to be created.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"webhook_id": schema.Int64Attribute{
				Description: "Numeric identifier of the webhook.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"active": schema.BoolAttribute{
				Description: "Boolean indicating if the webhook is active.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
				Default: booldefault.StaticBool(false),
			},
			"authorization_header": schema.StringAttribute{
				Description: "Authorization header to send to the target.",
				Optional:    true,
			},
			"branch_filter": schema.StringAttribute{
				Description: "List of allowed branches for push, branch creation and branch deletion events, specified as glob pattern. If empty or *, events for all branches are reported.",
				Optional:    true,
			},
			"config": schema.MapAttribute{
				Description: "Map of configuration settings.",
				Required:    true,
				ElementType: types.StringType,
			},
			"created_at": schema.StringAttribute{
				Description: "Time at which the webhook was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"events": schema.SetAttribute{
				Description: "List of events which trigger the webhook.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.StringType,
						[]attr.Value{
							types.StringValue("push"),
						},
					),
				),
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(
							"action_run_failure",
							"action_run_recover",
							"action_run_success",
							"create",
							"delete",
							"fork",
							"issue_assign",
							"issue_comment",
							"issue_label",
							"issue_milestone",
							"issues",
							"package",
							"pull_request",
							"pull_request_assign",
							"pull_request_comment",
							"pull_request_label",
							"pull_request_milestone",
							"pull_request_review_approved",
							"pull_request_review_comment",
							"pull_request_review_rejected",
							"pull_request_review_request",
							"pull_request_sync",
							"push",
							"release",
							"repository",
							"wiki",
						),
					),
				},
			},
			"type": schema.StringAttribute{
				Description: "Type of webhook. Changing this forces a new resource to be created.",
				Required:    true,
				Computed:    false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Description: "Time at which the webhook was updated.",
				Computed:    true,
				// PlanModifiers: []planmodifier.String{
				// 	stringplanmodifier.UseStateForUnknown(),
				// },
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *repositoryWebhookResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *repositoryWebhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create repository webhook resource"))

	var (
		data repositoryWebhookResourceModel
		repo repositoryResourceModel
	)

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get repository by id
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

	tflog.Info(ctx, "Create repository webhook", map[string]any{
		"repo": repo.Name.ValueString(),
	})

	// Generate API request body from plan
	opts := forgejo.CreateHookOption{
		Active:              data.Active.ValueBool(),
		AuthorizationHeader: data.AuthorizationHeader.ValueString(),
		BranchFilter:        data.BranchFilter.ValueString(),
		Type:                forgejo.HookType(data.Type.ValueString()),
	}
	var events []string
	data.Events.ElementsAs(ctx, &events, false)
	opts.Events = events
	var config map[string]string
	data.Config.ElementsAs(ctx, &config, false)
	opts.Config = config

	// Validate API request body
	err := opts.Validate()
	if err != nil {
		resp.Diagnostics.AddError("Input validation error", err.Error())

		return
	}

	// Use Forgejo client to create new repository webhook
	hook, res, err := r.client.CreateRepoHook(
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
			case 404:
				msg = fmt.Sprintf(
					"Repository with owner %s and name %s not found: %s",
					repo.Owner.String(),
					repo.Name.String(),
					err,
				)
			case 422:
				msg = fmt.Sprintf("Input validation error: %s", err)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to create repository webhook", msg)

		return
	}

	// Map response body to model
	diags = data.from(hook, ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *repositoryWebhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read repository webhook resource"))

	var (
		repo repositoryResourceModel
		data repositoryWebhookResourceModel
	)

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get repository by id
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

	tflog.Info(ctx, "Read repository webhook", map[string]any{
		"repo":    repo.Name.ValueString(),
		"hook_id": data.WebhookID.ValueInt64(),
	})

	// Use Forgejo client to get repository webhook
	hook, res, err := r.client.GetRepoHook(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		data.WebhookID.ValueInt64(),
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
			case 404:
				msg = fmt.Sprintf(
					"Repo webhook with repo owner %s repo name %s and webhook id %d not found: %s",
					repo.Owner.String(),
					repo.Name.String(),
					data.WebhookID.ValueInt64(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to read repository webhook", msg)

		return
	}

	// Map response body to model
	diags = data.from(hook, ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *repositoryWebhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update repository webhook resource"))

	var (
		data repositoryWebhookResourceModel
		repo repositoryResourceModel
	)

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get repository by id
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

	tflog.Info(ctx, "Update repository webhook", map[string]any{
		"id": data.WebhookID.ValueInt64(),
	})

	// Generate API request body from plan
	opts := forgejo.EditHookOption{}
	diags = data.to(&opts, ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to update existing repository webhook
	res, err := r.client.EditRepoHook(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		data.WebhookID.ValueInt64(),
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
			case 404:
				msg = fmt.Sprintf(
					"Repo webhook with repo owner %s repo name %s and webhook id %d not found: %s",
					repo.Owner.String(),
					repo.Name.String(),
					data.WebhookID.ValueInt64(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to update repository webhook", msg)

		return
	}

	tflog.Info(ctx, "Read repository webhook", map[string]any{
		"id": data.WebhookID.ValueInt64(),
	})

	// Use Forgejo client to fetch updated repository webhook
	hook, res, err := r.client.GetRepoHook(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		data.WebhookID.ValueInt64(),
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
			case 404:
				msg = fmt.Sprintf(
					"Repo webhook with repo owner %s repo name %s and webhook id %d not found: %s",
					repo.Owner.String(),
					repo.Name.String(),
					data.WebhookID.ValueInt64(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to read repository webhook", msg)

		return

	}

	// Map response body to model
	diags = data.from(hook, ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *repositoryWebhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete repository webhook resource"))

	var (
		repo repositoryResourceModel
		data repositoryWebhookResourceModel
	)

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get repository by id
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

	tflog.Info(ctx, "Delete repository webhook", map[string]any{
		"owner":  repo.Owner.ValueString(),
		"repo":   repo.Name.ValueString(),
		"key_id": data.WebhookID.ValueInt64(),
	})

	// Use Forgejo client to delete existing deploy key
	res, err := r.client.DeleteRepoHook(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		data.WebhookID.ValueInt64(),
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
			case 403:
				msg = fmt.Sprintf(
					"Repo webhook with repo owner %s repo name %s and webhook id %d forbidden: %s",
					repo.Owner.String(),
					repo.Name.String(),
					data.WebhookID.ValueInt64(),
					err,
				)
			case 404:
				msg = fmt.Sprintf(
					"Repo webhook with repo owner %s repo name %s and webhook id %d not found: %s",
					repo.Owner.String(),
					repo.Name.String(),
					data.WebhookID.ValueInt64(),
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to delete repository webhook", msg)

		return
	}

}

// ImportState reads an existing resource and adds it to Terraform state on success.
func (r *repositoryWebhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer un(trace(ctx, "Import repository webhook resource"))

	var state repositoryWebhookResourceModel

	// Parse import identifier
	cmp := strings.Split(req.ID, "/")
	if len(cmp) != 3 {
		resp.Diagnostics.AddError(
			"Unable to parse import identifier",
			fmt.Sprintf(
				"Expected import identifier with format: 'owner/name/webhookID', got: '%s'",
				req.ID,
			),
		)
		return
	}
	owner, repositoryName, hookID := cmp[0], cmp[1], cmp[2]

	id, err := strconv.ParseInt(hookID, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("Failed to parse webhook ID: %s", err)
		resp.Diagnostics.AddError("Unable to import repository webhook", msg)

		return
	}

	tflog.Info(ctx, "Read repository webhook", map[string]any{
		"repo":    repositoryName,
		"hook_id": id,
	})

	// Use Forgejo client to get repository by name
	rep, diags := getRepositoryByName(
		ctx,
		r.client,
		owner,
		repositoryName,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use Forgejo client to get repository webhook
	hook, res, err := r.client.GetRepoHook(owner, repositoryName, id)
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
					"Repository webhook '%s' not found: %s",
					req.ID,
					err,
				)
			default:
				msg = fmt.Sprintf("Unknown error: %s", err)
			}
		}
		resp.Diagnostics.AddError("Unable to read repository webhook", msg)

		return
	}

	// Map response body to model
	state.from(hook, ctx)

	state.RepositoryID = types.Int64Value(rep.ID)

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

}

// NewRepositoryResource is a helper function to simplify the provider implementation.
func NewRepositoryWebhookResource() resource.Resource {
	return &repositoryWebhookResource{}
}
