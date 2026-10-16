package provider

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	_ resource.Resource                = &tagProtectionResource{}
	_ resource.ResourceWithConfigure   = &tagProtectionResource{}
	_ resource.ResourceWithImportState = &tagProtectionResource{}
)

// tagProtectionResource is the resource implementation.
type tagProtectionResource struct {
	client *forgejo.Client
}

// tagProtectionResourceModel maps the resource schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v3#TagProtection
type tagProtectionResourceModel struct {
	ID                 types.Int64  `tfsdk:"id"`
	RepositoryID       types.Int64  `tfsdk:"repository_id"`
	NamePattern        types.String `tfsdk:"name_pattern"`
	WhitelistUsernames types.Set    `tfsdk:"whitelist_usernames"`
	WhitelistTeams     types.Set    `tfsdk:"whitelist_teams"`
	CreatedAt          types.String `tfsdk:"created_at"`
	UpdatedAt          types.String `tfsdk:"updated_at"`
}

// Metadata returns the resource type name.
func (r *tagProtectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tag_protection"
}

// Schema defines the schema for the resource.
func (r *tagProtectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo tag protection resource.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the tag protection.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"repository_id": schema.Int64Attribute{
				Description: "Numeric identifier of the repository. Changing this forces a new resource to be created.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"name_pattern": schema.StringAttribute{
				Description: "Glob pattern of the tag names to protect, e.g. `v*`.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\S(?s:.*\S)?$`),
						"must not be empty or have leading or trailing whitespace",
					),
				},
			},
			"whitelist_usernames": schema.SetAttribute{
				Description: "Whitelisted users allowed to create or delete protected tags. **Note**: The Forgejo API requires at least one whitelisted user or team when creating a tag protection, and silently drops whitelisted users without access to the repository.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.StringType,
						[]attr.Value{},
					),
				),
			},
			"whitelist_teams": schema.SetAttribute{
				Description: "Whitelisted teams allowed to create or delete protected tags. **Note**: The Forgejo API requires at least one whitelisted user or team when creating a tag protection, and silently drops whitelisted teams without access to the repository.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default: setdefault.StaticValue(
					types.SetValueMust(
						types.StringType,
						[]attr.Value{},
					),
				),
			},
			"created_at": schema.StringAttribute{
				Description: "Time at which the tag protection was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated_at": schema.StringAttribute{
				Description: "Time at which the tag protection was updated.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *tagProtectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *tagProtectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create tag protection resource"))

	var (
		repo repositoryResourceModel
		data tagProtectionResourceModel
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

	// Generate API request body from plan
	opts := r.toCreateOption(ctx, &data)

	// The API rejects a tag protection without any whitelist entries on creation
	if len(opts.WhitelistUsernames) == 0 && len(opts.WhitelistTeams) == 0 {
		resp.Diagnostics.AddError(
			"Unable to create tag protection",
			"The Forgejo API requires at least one whitelisted user or team when creating a tag protection. Whitelists can be emptied with a subsequent update.",
		)

		return
	}

	tflog.Info(ctx, "Create tag protection", map[string]any{
		"repository_id":       data.RepositoryID.ValueInt64(),
		"repository_name":     repo.Name.ValueString(),
		"repository_owner":    repo.Owner.ValueString(),
		"name_pattern":        opts.NamePattern,
		"whitelist_usernames": opts.WhitelistUsernames,
		"whitelist_teams":     opts.WhitelistTeams,
	})

	// Use Forgejo client to create tag protection
	protection, res, err := r.client.CreateTagProtection(
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
			case 403:
				msg = fmt.Sprintf(
					"Repository with owner %s and name %s forbidden: %s",
					repo.Owner.String(),
					repo.Name.String(),
					err,
				)
			case 404:
				msg = fmt.Sprintf(
					"Repository with owner %s and name %s not found: %s",
					repo.Owner.String(),
					repo.Name.String(),
					err,
				)
			case 422:
				msg = fmt.Sprintf("Input validation error: %s", err)
			case 423:
				msg = fmt.Sprintf(
					"Repository with owner %s and name %s is archived: %s",
					repo.Owner.String(),
					repo.Name.String(),
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
		resp.Diagnostics.AddError("Unable to create tag protection", msg)

		return
	}

	// Report whitelist entries silently dropped by the API
	diags = r.checkWhitelists(opts.WhitelistUsernames, opts.WhitelistTeams, protection, "Unable to create tag protection")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update model with response data to ensure computed fields are correctly populated
	diags = r.from(ctx, protection, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *tagProtectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read tag protection resource"))

	var (
		repo repositoryResourceModel
		data tagProtectionResourceModel
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

	tflog.Info(ctx, "Read tag protection", map[string]any{
		"owner": repo.Owner.ValueString(),
		"repo":  repo.Name.ValueString(),
		"id":    data.ID.ValueInt64(),
	})

	// Use Forgejo client to get tag protection
	protection, res, err := r.client.GetTagProtection(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		data.ID.ValueInt64(),
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
				// Remove the tag protection from state, so that Terraform
				// recreates it if it has been deleted outside of Terraform
				resp.State.RemoveResource(ctx)

				return
			default:
				msg = fmt.Sprintf(
					"Unknown error (status %d): %s",
					res.StatusCode,
					err,
				)
			}
		}
		resp.Diagnostics.AddError("Unable to read tag protection", msg)

		return
	}

	// Update model with response data
	diags = r.from(ctx, protection, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *tagProtectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update tag protection resource"))

	var (
		repo repositoryResourceModel
		data tagProtectionResourceModel
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

	// Convert model to API request
	opts := r.toEditOption(ctx, &data)

	tflog.Info(ctx, "Update tag protection", map[string]any{
		"repository_id":       data.RepositoryID.ValueInt64(),
		"repository_name":     repo.Name.ValueString(),
		"repository_owner":    repo.Owner.ValueString(),
		"id":                  data.ID.ValueInt64(),
		"name_pattern":        data.NamePattern.ValueString(),
		"whitelist_usernames": opts.WhitelistUsernames,
		"whitelist_teams":     opts.WhitelistTeams,
	})

	// Use Forgejo client to update tag protection
	protection, res, err := r.client.EditTagProtection(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		data.ID.ValueInt64(),
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
			case 403:
				msg = fmt.Sprintf(
					"Repository with owner %s and name %s forbidden: %s",
					repo.Owner.String(),
					repo.Name.String(),
					err,
				)
			case 404:
				msg = fmt.Sprintf(
					"Tag protection with owner %s, repo %s and ID %d not found: %s",
					repo.Owner.String(),
					repo.Name.String(),
					data.ID.ValueInt64(),
					err,
				)
			case 422:
				msg = fmt.Sprintf("Input validation error: %s", err)
			case 423:
				msg = fmt.Sprintf(
					"Repository with owner %s and name %s is archived: %s",
					repo.Owner.String(),
					repo.Name.String(),
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
		resp.Diagnostics.AddError("Unable to update tag protection", msg)

		return
	}

	// Report whitelist entries silently dropped by the API
	diags = r.checkWhitelists(opts.WhitelistUsernames, opts.WhitelistTeams, protection, "Unable to update tag protection")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update model with response data
	diags = r.from(ctx, protection, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *tagProtectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete tag protection resource"))

	var (
		repo repositoryResourceModel
		data tagProtectionResourceModel
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

	tflog.Info(ctx, "Delete tag protection", map[string]any{
		"repository_id":    data.RepositoryID.ValueInt64(),
		"repository_name":  repo.Name.ValueString(),
		"repository_owner": repo.Owner.ValueString(),
		"id":               data.ID.ValueInt64(),
	})

	// Use Forgejo client to delete tag protection
	res, err := r.client.DeleteTagProtection(
		repo.Owner.ValueString(),
		repo.Name.ValueString(),
		data.ID.ValueInt64(),
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
					"Repository with owner %s and name %s forbidden: %s",
					repo.Owner.String(),
					repo.Name.String(),
					err,
				)
			case 404:
				msg = fmt.Sprintf(
					"Tag protection with owner %s, repo %s and ID %d not found: %s",
					repo.Owner.String(),
					repo.Name.String(),
					data.ID.ValueInt64(),
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
		resp.Diagnostics.AddError("Unable to delete tag protection", msg)

		return
	}
}

// ImportState reads an existing resource and adds it to Terraform state on success.
func (r *tagProtectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	defer un(trace(ctx, "Import tag protection resource"))

	var data tagProtectionResourceModel

	// Parse import identifier
	parts := strings.Split(req.ID, "/")
	if len(parts) != 3 {
		resp.Diagnostics.AddError(
			"Unable to parse import identifier",
			fmt.Sprintf(
				"Expected import identifier with format: 'owner/repo/protectionID', got: '%s'",
				req.ID,
			),
		)

		return
	}
	owner, repo, protectionIDstr := parts[0], parts[1], parts[2]

	// Parse protection ID
	protectionID, err := strconv.ParseInt(protectionIDstr, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to parse import identifier",
			fmt.Sprintf(
				"Failed to parse protection ID: %s",
				err,
			),
		)

		return
	}

	tflog.Info(ctx, "Read tag protection", map[string]any{
		"owner": owner,
		"repo":  repo,
		"id":    protectionID,
	})

	// Use Forgejo client to get tag protection
	protection, res, err := r.client.GetTagProtection(
		owner,
		repo,
		protectionID,
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
					"Tag protection with owner '%s', repo '%s' and ID %d not found: %s",
					owner,
					repo,
					protectionID,
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
		resp.Diagnostics.AddError("Unable to read tag protection", msg)

		return
	}

	// Use Forgejo client to get repository
	rep, diags := getRepositoryByName(
		ctx,
		r.client,
		owner,
		repo,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response to model
	data.RepositoryID = types.Int64Value(rep.ID)
	diags = r.from(ctx, protection, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// NewTagProtectionResource is a helper function to simplify the provider implementation.
func NewTagProtectionResource() resource.Resource {
	return &tagProtectionResource{}
}

// checkWhitelists reports whitelist entries that were requested but are
// missing from the API response. The Forgejo API silently drops whitelisted
// users and teams without access to the repository, which would otherwise
// surface as a generic "inconsistent result after apply" error.
func (r *tagProtectionResource) checkWhitelists(usernames, teams []string, protection *forgejo.TagProtection, summary string) (diags diag.Diagnostics) {
	if protection == nil {
		return diags
	}

	if missing := missingEntries(usernames, protection.WhitelistUsernames); len(missing) > 0 {
		diags.AddAttributeError(
			path.Root("whitelist_usernames"),
			summary,
			fmt.Sprintf(
				"The following whitelisted users have no access to the repository: %s",
				strings.Join(missing, ", "),
			),
		)
	}

	if missing := missingEntries(teams, protection.WhitelistTeams); len(missing) > 0 {
		diags.AddAttributeError(
			path.Root("whitelist_teams"),
			summary,
			fmt.Sprintf(
				"The following whitelisted teams have no access to the repository: %s",
				strings.Join(missing, ", "),
			),
		)
	}

	return diags
}

// missingEntries returns the entries of requested that are not contained in actual.
func missingEntries(requested, actual []string) (missing []string) {
	for _, entry := range requested {
		if !slices.Contains(actual, entry) {
			missing = append(missing, entry)
		}
	}

	return missing
}

// Helper function to convert model to CreateTagProtectionOption.
func (r *tagProtectionResource) toCreateOption(ctx context.Context, data *tagProtectionResourceModel) forgejo.CreateTagProtectionOption {
	opts := forgejo.CreateTagProtectionOption{
		NamePattern: data.NamePattern.ValueString(),
	}

	var usernames []string
	data.WhitelistUsernames.ElementsAs(ctx, &usernames, false)
	opts.WhitelistUsernames = usernames

	var teams []string
	data.WhitelistTeams.ElementsAs(ctx, &teams, false)
	opts.WhitelistTeams = teams

	return opts
}

// Helper function to convert model to EditTagProtectionOption.
func (r *tagProtectionResource) toEditOption(ctx context.Context, data *tagProtectionResourceModel) forgejo.EditTagProtectionOption {
	opts := forgejo.EditTagProtectionOption{
		NamePattern: data.NamePattern.ValueStringPointer(),
	}

	var usernames []string
	data.WhitelistUsernames.ElementsAs(ctx, &usernames, false)
	opts.WhitelistUsernames = usernames

	var teams []string
	data.WhitelistTeams.ElementsAs(ctx, &teams, false)
	opts.WhitelistTeams = teams

	return opts
}

// from is a helper function to load an API struct into Terraform data model.
func (r *tagProtectionResource) from(ctx context.Context, protection *forgejo.TagProtection, data *tagProtectionResourceModel) (diags diag.Diagnostics) {
	if protection == nil || data == nil {
		return diags
	}

	data.ID = types.Int64Value(protection.ID)
	data.NamePattern = types.StringValue(protection.NamePattern)
	data.CreatedAt = types.StringValue(protection.Created.Format(time.RFC3339))
	data.UpdatedAt = types.StringValue(protection.Updated.Format(time.RFC3339))

	// Handle Sets; empty API lists are mapped to known empty sets, not null
	var d diag.Diagnostics
	usernames := protection.WhitelistUsernames
	if usernames == nil {
		usernames = []string{}
	}
	data.WhitelistUsernames, d = types.SetValueFrom(ctx, types.StringType, usernames)
	diags.Append(d...)

	teams := protection.WhitelistTeams
	if teams == nil {
		teams = []string{}
	}
	data.WhitelistTeams, d = types.SetValueFrom(ctx, types.StringType, teams)
	diags.Append(d...)

	// All fields of the API struct are mapped; nothing is intentionally omitted

	return diags
}
