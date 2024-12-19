package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &repositoryResource{}
	_ resource.ResourceWithConfigure = &repositoryResource{}
)

// repositoryResource is the resource implementation.
type repositoryResource struct {
	client *forgejo.Client
}

// repositoryResourceModel maps the resource schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo#Repository
type repositoryResourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	Owner         types.Object `tfsdk:"owner"`
	Name          types.String `tfsdk:"name"`
	FullName      types.String `tfsdk:"full_name"`
	Description   types.String `tfsdk:"description"`
	Empty         types.Bool   `tfsdk:"empty"`
	Private       types.Bool   `tfsdk:"private"`
	Fork          types.Bool   `tfsdk:"fork"`
	Template      types.Bool   `tfsdk:"template"`
	ParentID      types.Int64  `tfsdk:"parent_id"`
	Mirror        types.Bool   `tfsdk:"mirror"`
	Size          types.Int64  `tfsdk:"size"`
	HTMLURL       types.String `tfsdk:"html_url"`
	SSHURL        types.String `tfsdk:"ssh_url"`
	CloneURL      types.String `tfsdk:"clone_url"`
	OriginalURL   types.String `tfsdk:"original_url"`
	Website       types.String `tfsdk:"website"`
	Stars         types.Int64  `tfsdk:"stars_count"`
	Forks         types.Int64  `tfsdk:"forks_count"`
	Watchers      types.Int64  `tfsdk:"watchers_count"`
	OpenIssues    types.Int64  `tfsdk:"open_issues_count"`
	OpenPulls     types.Int64  `tfsdk:"open_pr_counter"`
	Releases      types.Int64  `tfsdk:"release_counter"`
	DefaultBranch types.String `tfsdk:"default_branch"`
	Archived      types.Bool   `tfsdk:"archived"`
	Created       types.String `tfsdk:"created_at"`
	// updated_at changes with every update, so we're not including it in the model
	// Updated                   types.String `tfsdk:"updated_at"`
	Permissions               types.Object `tfsdk:"permissions"`
	HasIssues                 types.Bool   `tfsdk:"has_issues"`
	InternalTracker           types.Object `tfsdk:"internal_tracker"`
	ExternalTracker           types.Object `tfsdk:"external_tracker"`
	HasWiki                   types.Bool   `tfsdk:"has_wiki"`
	ExternalWiki              types.Object `tfsdk:"external_wiki"`
	HasPullRequests           types.Bool   `tfsdk:"has_pull_requests"`
	HasProjects               types.Bool   `tfsdk:"has_projects"`
	HasReleases               types.Bool   `tfsdk:"has_releases"`
	HasPackages               types.Bool   `tfsdk:"has_packages"`
	HasActions                types.Bool   `tfsdk:"has_actions"`
	IgnoreWhitespaceConflicts types.Bool   `tfsdk:"ignore_whitespace_conflicts"`
	AllowMerge                types.Bool   `tfsdk:"allow_merge_commits"`
	AllowRebase               types.Bool   `tfsdk:"allow_rebase"`
	AllowRebaseMerge          types.Bool   `tfsdk:"allow_rebase_explicit"`
	AllowSquash               types.Bool   `tfsdk:"allow_squash_merge"`
	AvatarURL                 types.String `tfsdk:"avatar_url"`
	Internal                  types.Bool   `tfsdk:"internal"`
	MirrorInterval            types.String `tfsdk:"mirror_interval"`
	MirrorUpdated             types.String `tfsdk:"mirror_updated"`
	DefaultMergeStyle         types.String `tfsdk:"default_merge_style"`
	IssueLabels               types.String `tfsdk:"issue_labels"`
	AutoInit                  types.Bool   `tfsdk:"auto_init"`
	Gitignores                types.String `tfsdk:"gitignores"`
	License                   types.String `tfsdk:"license"`
	Readme                    types.String `tfsdk:"readme"`
	TrustModel                types.String `tfsdk:"trust_model"`
}

func (m *repositoryResourceModel) from(r *forgejo.Repository) {
	m.ID = types.Int64Value(r.ID)
	m.Name = types.StringValue(r.Name)
	m.FullName = types.StringValue(r.FullName)
	m.Description = types.StringValue(r.Description)
	m.Empty = types.BoolValue(r.Empty)
	m.Private = types.BoolValue(r.Private)
	m.Fork = types.BoolValue(r.Fork)
	m.Template = types.BoolValue(r.Template)
	if r.Parent != nil {
		m.ParentID = types.Int64Value(r.Parent.ID)
	} else {
		m.ParentID = types.Int64Null()
	}
	m.Mirror = types.BoolValue(r.Mirror)
	m.Size = types.Int64Value(int64(r.Size))
	m.HTMLURL = types.StringValue(r.HTMLURL)
	m.SSHURL = types.StringValue(r.SSHURL)
	m.CloneURL = types.StringValue(r.CloneURL)
	m.OriginalURL = types.StringValue(r.OriginalURL)
	m.Website = types.StringValue(r.Website)
	m.Stars = types.Int64Value(int64(r.Stars))
	m.Forks = types.Int64Value(int64(r.Forks))
	m.Watchers = types.Int64Value(int64(r.Watchers))
	m.OpenIssues = types.Int64Value(int64(r.OpenIssues))
	m.OpenPulls = types.Int64Value(int64(r.OpenPulls))
	m.Releases = types.Int64Value(int64(r.Releases))
	m.DefaultBranch = types.StringValue(r.DefaultBranch)
	m.Archived = types.BoolValue(r.Archived)
	m.Created = types.StringValue(r.Created.String())
	// m.Updated = types.StringValue(r.Updated.String())
	m.HasIssues = types.BoolValue(r.HasIssues)
	m.HasWiki = types.BoolValue(r.HasWiki)
	m.HasPullRequests = types.BoolValue(r.HasPullRequests)
	m.HasProjects = types.BoolValue(r.HasProjects)
	m.HasReleases = types.BoolValue(r.HasReleases)
	m.HasPackages = types.BoolValue(r.HasPackages)
	m.HasActions = types.BoolValue(r.HasActions)
	m.IgnoreWhitespaceConflicts = types.BoolValue(r.IgnoreWhitespaceConflicts)
	m.AllowMerge = types.BoolValue(r.AllowMerge)
	m.AllowRebase = types.BoolValue(r.AllowRebase)
	m.AllowRebaseMerge = types.BoolValue(r.AllowRebaseMerge)
	m.AllowSquash = types.BoolValue(r.AllowSquash)
	m.AvatarURL = types.StringValue(r.AvatarURL)
	m.Internal = types.BoolValue(r.Internal)
	m.MirrorInterval = types.StringValue(r.MirrorInterval)
	m.MirrorUpdated = types.StringValue(r.MirrorUpdated.String())
	m.DefaultMergeStyle = types.StringValue(string(r.DefaultMergeStyle))
}

// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo#User
type repositoryResourceUser struct {
	ID        types.Int64  `tfsdk:"id"`
	UserName  types.String `tfsdk:"login"`
	LoginName types.String `tfsdk:"login_name"`
	FullName  types.String `tfsdk:"full_name"`
	Email     types.String `tfsdk:"email"`
}

func (m repositoryResourceUser) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":         types.Int64Type,
		"login":      types.StringType,
		"login_name": types.StringType,
		"full_name":  types.StringType,
		"email":      types.StringType,
	}
}
func (m *repositoryResourceModel) ownerFrom(ctx context.Context, r *forgejo.Repository) diag.Diagnostics {
	if r.Owner == nil {
		m.Owner = types.ObjectNull(
			repositoryResourceUser{}.attributeTypes(),
		)
		return nil
	}

	ownerElement := repositoryResourceUser{
		ID:        types.Int64Value(r.Owner.ID),
		UserName:  types.StringValue(r.Owner.UserName),
		LoginName: types.StringValue(r.Owner.LoginName),
		FullName:  types.StringValue(r.Owner.FullName),
		Email:     types.StringValue(r.Owner.Email),
	}

	ownerValue, diags := types.ObjectValueFrom(
		ctx,
		ownerElement.attributeTypes(),
		ownerElement,
	)

	if !diags.HasError() {
		m.Owner = ownerValue
	}

	return diags
}

// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo#Permission
type repositoryResourcePermissions struct {
	Admin types.Bool `tfsdk:"admin"`
	Push  types.Bool `tfsdk:"push"`
	Pull  types.Bool `tfsdk:"pull"`
}

func (m repositoryResourcePermissions) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"admin": types.BoolType,
		"push":  types.BoolType,
		"pull":  types.BoolType,
	}
}
func (m repositoryResourcePermissions) defaultObject() map[string]attr.Value {
	return map[string]attr.Value{
		"admin": types.BoolValue(true),
		"push":  types.BoolValue(true),
		"pull":  types.BoolValue(true),
	}
}
func (m *repositoryResourceModel) permissionsFrom(ctx context.Context, r *forgejo.Repository) diag.Diagnostics {
	if r.Permissions == nil {
		m.Permissions = types.ObjectNull(
			repositoryResourcePermissions{}.attributeTypes(),
		)
		return nil
	}

	permsElement := repositoryResourcePermissions{
		Admin: types.BoolValue(r.Permissions.Admin),
		Push:  types.BoolValue(r.Permissions.Push),
		Pull:  types.BoolValue(r.Permissions.Pull),
	}

	permsValue, diags := types.ObjectValueFrom(
		ctx,
		permsElement.attributeTypes(),
		permsElement,
	)

	if !diags.HasError() {
		m.Permissions = permsValue
	}

	return diags
}

// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo#InternalTracker
type repositoryResourceInternalTracker struct {
	EnableTimeTracker                types.Bool `tfsdk:"enable_time_tracker"`
	AllowOnlyContributorsToTrackTime types.Bool `tfsdk:"allow_only_contributors_to_track_time"`
	EnableIssueDependencies          types.Bool `tfsdk:"enable_issue_dependencies"`
}

func (m repositoryResourceInternalTracker) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enable_time_tracker":                   types.BoolType,
		"allow_only_contributors_to_track_time": types.BoolType,
		"enable_issue_dependencies":             types.BoolType,
	}
}
func (m repositoryResourceInternalTracker) defaultObject() map[string]attr.Value {
	return map[string]attr.Value{
		"enable_time_tracker":                   types.BoolValue(true),
		"allow_only_contributors_to_track_time": types.BoolValue(true),
		"enable_issue_dependencies":             types.BoolValue(true),
	}
}
func (m *repositoryResourceModel) internalTrackerFrom(ctx context.Context, r *forgejo.Repository) diag.Diagnostics {
	if r.InternalTracker == nil {
		m.InternalTracker = types.ObjectNull(
			repositoryResourceInternalTracker{}.attributeTypes(),
		)
		return nil
	}

	intTrackerElement := repositoryResourceInternalTracker{
		EnableTimeTracker:                types.BoolValue(r.InternalTracker.EnableTimeTracker),
		AllowOnlyContributorsToTrackTime: types.BoolValue(r.InternalTracker.AllowOnlyContributorsToTrackTime),
		EnableIssueDependencies:          types.BoolValue(r.InternalTracker.EnableIssueDependencies),
	}

	intTrackerValue, diags := types.ObjectValueFrom(
		ctx,
		intTrackerElement.attributeTypes(),
		intTrackerElement,
	)

	if !diags.HasError() {
		m.InternalTracker = intTrackerValue
	}

	return diags
}
func (m *repositoryResourceModel) internalTrackerTo(ctx context.Context, it *forgejo.InternalTracker) diag.Diagnostics {
	if m.InternalTracker.IsNull() {
		return nil
	}

	var intTracker repositoryResourceInternalTracker
	diags := m.InternalTracker.As(ctx, &intTracker, basetypes.ObjectAsOptions{})

	if !diags.HasError() {
		it = &forgejo.InternalTracker{
			EnableTimeTracker:                intTracker.EnableTimeTracker.ValueBool(),
			AllowOnlyContributorsToTrackTime: intTracker.AllowOnlyContributorsToTrackTime.ValueBool(),
			EnableIssueDependencies:          intTracker.EnableIssueDependencies.ValueBool(),
		}
	}

	return diags
}

// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo#ExternalTracker
type repositoryResourceExternalTracker struct {
	ExternalTrackerURL    types.String `tfsdk:"external_tracker_url"`
	ExternalTrackerFormat types.String `tfsdk:"external_tracker_format"`
	ExternalTrackerStyle  types.String `tfsdk:"external_tracker_style"`
}

func (m repositoryResourceExternalTracker) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"external_tracker_url":    types.StringType,
		"external_tracker_format": types.StringType,
		"external_tracker_style":  types.StringType,
	}
}
func (m *repositoryResourceModel) externalTrackerFrom(ctx context.Context, r *forgejo.Repository) diag.Diagnostics {
	if r.ExternalTracker == nil {
		m.ExternalTracker = types.ObjectNull(
			repositoryResourceExternalTracker{}.attributeTypes(),
		)
		return nil
	}

	extTrackerElement := repositoryResourceExternalTracker{
		ExternalTrackerURL:    types.StringValue(r.ExternalTracker.ExternalTrackerURL),
		ExternalTrackerFormat: types.StringValue(r.ExternalTracker.ExternalTrackerFormat),
		ExternalTrackerStyle:  types.StringValue(r.ExternalTracker.ExternalTrackerStyle),
	}

	extTrackerValue, diags := types.ObjectValueFrom(
		ctx,
		extTrackerElement.attributeTypes(),
		extTrackerElement,
	)

	if !diags.HasError() {
		m.ExternalTracker = extTrackerValue
	}

	return diags
}
func (m *repositoryResourceModel) externalTrackerTo(ctx context.Context, et *forgejo.ExternalTracker) diag.Diagnostics {
	if m.ExternalTracker.IsNull() {
		return nil
	}

	var extTracker repositoryResourceExternalTracker
	diags := m.ExternalTracker.As(ctx, &extTracker, basetypes.ObjectAsOptions{})

	if !diags.HasError() {
		et = &forgejo.ExternalTracker{
			ExternalTrackerURL:    extTracker.ExternalTrackerURL.ValueString(),
			ExternalTrackerFormat: extTracker.ExternalTrackerFormat.ValueString(),
			ExternalTrackerStyle:  extTracker.ExternalTrackerStyle.ValueString(),
		}
	}

	return diags
}

// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo#ExternalWiki
type repositoryResourceExternalWiki struct {
	ExternalWikiURL types.String `tfsdk:"external_wiki_url"`
}

func (m repositoryResourceExternalWiki) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"external_wiki_url": types.StringType,
	}
}
func (m *repositoryResourceModel) externalWikiFrom(ctx context.Context, r *forgejo.Repository) diag.Diagnostics {
	if r.ExternalWiki == nil {
		m.ExternalWiki = types.ObjectNull(
			repositoryResourceExternalWiki{}.attributeTypes(),
		)
		return nil
	}

	wikiElement := repositoryResourceExternalWiki{
		ExternalWikiURL: types.StringValue(r.ExternalWiki.ExternalWikiURL),
	}

	wikiValue, diags := types.ObjectValueFrom(
		ctx,
		wikiElement.attributeTypes(),
		wikiElement,
	)

	if !diags.HasError() {
		m.ExternalWiki = wikiValue
	}

	return diags
}
func (m *repositoryResourceModel) externalWikiTo(ctx context.Context, ew *forgejo.ExternalWiki) diag.Diagnostics {
	if m.ExternalWiki.IsNull() {
		return nil
	}

	var extWiki repositoryResourceExternalWiki
	diags := m.ExternalWiki.As(ctx, &extWiki, basetypes.ObjectAsOptions{})

	if !diags.HasError() {
		ew = &forgejo.ExternalWiki{
			ExternalWikiURL: extWiki.ExternalWikiURL.ValueString(),
		}
	}

	return diags
}

// Metadata returns the resource type name.
func (r *repositoryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

// Schema defines the schema for the resource.
func (r *repositoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo repository resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the repository.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"owner": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "Numeric identifier of the user.",
						Computed:    true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"login": schema.StringAttribute{
						Description: "Name of the user.",
						Required:    true,
					},
					"login_name": schema.StringAttribute{
						Description: "Login name of the user.",
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"full_name": schema.StringAttribute{
						Description: "Full name of the user.",
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"email": schema.StringAttribute{
						Description: "Email address of the user.",
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
				Description: "Owner of the repository.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the repository.",
				Required:    true,
			},
			"full_name": schema.StringAttribute{
				Description: "Full name of the repository.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description of the repository.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"empty": schema.BoolAttribute{
				Description: "Is the repository empty?",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"private": schema.BoolAttribute{
				Description: "Is the repository private?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"fork": schema.BoolAttribute{
				Description: "Is the repository a fork?",
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"template": schema.BoolAttribute{
				Description: "Is the repository a template?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"parent_id": schema.Int64Attribute{
				Description: "Numeric identifier of the parent repository.",
				Optional:    true,
			},
			"mirror": schema.BoolAttribute{
				Description: "Is the repository a mirror?",
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"size": schema.Int64Attribute{
				Description: "Size of the repository in KiB.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"html_url": schema.StringAttribute{
				Description: "HTML URL of the repository.",
				Computed:    true,
			},
			"ssh_url": schema.StringAttribute{
				Description: "SSH URL of the repository.",
				Computed:    true,
			},
			"clone_url": schema.StringAttribute{
				Description: "Clone URL of the repository.",
				Computed:    true,
			},
			"original_url": schema.StringAttribute{
				Description: "Original URL of the repository.",
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"website": schema.StringAttribute{
				Description: "Website of the repository.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"stars_count": schema.Int64Attribute{
				Description: "Number of stars of the repository.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"forks_count": schema.Int64Attribute{
				Description: "Number of forks of the repository.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"watchers_count": schema.Int64Attribute{
				Description: "Number of watchers of the repository.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"open_issues_count": schema.Int64Attribute{
				Description: "Number of open issues of the repository.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"open_pr_counter": schema.Int64Attribute{
				Description: "Number of open pull requests of the repository.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"release_counter": schema.Int64Attribute{
				Description: "Number of releases of the repository.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"default_branch": schema.StringAttribute{
				Description: "Default branch of the repository.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("main"),
			},
			"archived": schema.BoolAttribute{
				Description: "Is the repository archived?",
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"created_at": schema.StringAttribute{
				Description: "Time at which the repository was created.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// "updated_at": schema.StringAttribute{
			// 	Description: "Time at which the repository was updated.",
			// 	Computed:    true,
			// },
			"permissions": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"admin": schema.BoolAttribute{
						Description: "Allowed to administer?",
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"push": schema.BoolAttribute{
						Description: "Allowed to push?",
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"pull": schema.BoolAttribute{
						Description: "Allowed to pull?",
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
				},
				Description: "Permissions of the repository.",
				Computed:    true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						repositoryResourcePermissions{}.attributeTypes(),
						repositoryResourcePermissions{}.defaultObject(),
					),
				),
			},
			"has_issues": schema.BoolAttribute{
				Description: "Is the repository issue tracker enabled?",
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"internal_tracker": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"enable_time_tracker": schema.BoolAttribute{
						Description: "Enable time tracking.",
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"allow_only_contributors_to_track_time": schema.BoolAttribute{
						Description: "Let only contributors track time.",
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"enable_issue_dependencies": schema.BoolAttribute{
						Description: "Enable dependencies for issues and pull requests.",
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
				},
				Description: "Settings for built-in issue tracker.",
				Computed:    true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						repositoryResourceInternalTracker{}.attributeTypes(),
						repositoryResourceInternalTracker{}.defaultObject(),
					),
				),
			},
			"external_tracker": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"external_tracker_url": schema.StringAttribute{
						Description: "URL of external issue tracker.",
						Computed:    true,
					},
					"external_tracker_format": schema.StringAttribute{
						Description: "External issue tracker URL format.",
						Computed:    true,
					},
					"external_tracker_style": schema.StringAttribute{
						Description: "External issue tracker number format.",
						Computed:    true,
					},
				},
				Description: "Settings for external issue tracker.",
				Computed:    true,
				Default: objectdefault.StaticValue(
					types.ObjectNull(
						repositoryResourceExternalTracker{}.attributeTypes(),
					),
				),
			},
			"has_wiki": schema.BoolAttribute{
				Description: "Is the repository wiki enabled?",
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"external_wiki": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"external_wiki_url": schema.StringAttribute{
						Description: "URL of external wiki.",
						Computed:    true,
					},
				},
				Description: "Settings for external wiki.",
				Computed:    true,
				Default: objectdefault.StaticValue(
					types.ObjectNull(
						repositoryResourceExternalWiki{}.attributeTypes(),
					),
				),
			},
			"has_pull_requests": schema.BoolAttribute{
				Description: "Are repository pull requests enabled?",
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"has_projects": schema.BoolAttribute{
				Description: "Are repository projects enabled?",
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"has_releases": schema.BoolAttribute{
				Description: "Are repository releases enabled?",
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"has_packages": schema.BoolAttribute{
				Description: "Is the repository package registry enabled?",
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"has_actions": schema.BoolAttribute{
				Description: "Are integrated CI/CD pipelines enabled?",
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"ignore_whitespace_conflicts": schema.BoolAttribute{
				Description: "Are whitespace conflicts ignored?",
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"allow_merge_commits": schema.BoolAttribute{
				Description: "Allowed to create merge commit?",
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"allow_rebase": schema.BoolAttribute{
				Description: "Allowed to rebase then fast-forward?",
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"allow_rebase_explicit": schema.BoolAttribute{
				Description: "Allowed to rebase then create merge commit?",
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"allow_squash_merge": schema.BoolAttribute{
				Description: "Allowed to create squash commit?",
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"avatar_url": schema.StringAttribute{
				Description: "Avatar URL of the repository.",
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"internal": schema.BoolAttribute{
				Description: "Is the repository internal?",
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"mirror_interval": schema.StringAttribute{
				Description: "Mirror interval of the repository.",
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"mirror_updated": schema.StringAttribute{
				Description: "Time at which the repository mirror was updated.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"default_merge_style": schema.StringAttribute{
				Description: "Default merge style of the repository.",
				Computed:    true,
				Default:     stringdefault.StaticString("merge"),
			},
			"issue_labels": schema.StringAttribute{
				Description: "Issue Label set to use",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"auto_init": schema.BoolAttribute{
				Description: "Whether the repository should be auto-intialized?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"gitignores": schema.StringAttribute{
				Description: "Gitignores to use",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"license": schema.StringAttribute{
				Description: "License to use",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"readme": schema.StringAttribute{
				Description: "Readme of the repository to create",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"trust_model": schema.StringAttribute{
				Description: "TrustModel of the repository",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOf(
						"default",
						"collaborator",
						"committer",
						"collaboratorcommitter",
					),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *repositoryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *repositoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	defer un(trace(ctx, "Create repository resource"))

	var (
		data  repositoryResourceModel
		owner repositoryResourceUser
	)

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read repository owner into model
	if !data.Owner.IsUnknown() {
		diags = data.Owner.As(ctx, &owner, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Info(ctx, "Create repository", map[string]any{
		"owner":          owner.UserName.ValueString(),
		"name":           data.Name.ValueString(),
		"description":    data.Description.ValueString(),
		"private":        data.Private.ValueBool(),
		"issue_labels":   data.IssueLabels.ValueString(),
		"auto_init":      data.AutoInit.ValueBool(),
		"template":       data.Template.ValueBool(),
		"gitignores":     data.Gitignores.ValueString(),
		"license":        data.License.ValueString(),
		"readme":         data.Readme.ValueString(),
		"default_branch": data.DefaultBranch.ValueString(),
		"trust_model":    forgejo.TrustModel(data.TrustModel.ValueString()),
	})

	// Generate API request body from plan
	opts := forgejo.CreateRepoOption{
		Name:          data.Name.ValueString(),
		Description:   data.Description.ValueString(),
		Private:       data.Private.ValueBool(),
		IssueLabels:   data.IssueLabels.ValueString(),
		AutoInit:      data.AutoInit.ValueBool(),
		Template:      data.Template.ValueBool(),
		Gitignores:    data.Gitignores.ValueString(),
		License:       data.License.ValueString(),
		Readme:        data.Readme.ValueString(),
		DefaultBranch: data.DefaultBranch.ValueString(),
		TrustModel:    forgejo.TrustModel(data.TrustModel.ValueString()),
	}

	// Validate API request body
	err := opts.Validate(r.client)
	if err != nil {
		resp.Diagnostics.AddError("Input validation error", err.Error())

		return
	}

	// Determine type of repository
	var ownerType string
	if owner.UserName.ValueString() == "" {
		// No owner -> personal repository
		ownerType = "personal"
	} else {
		// Use Forgejo client to check if owner is org
		_, res, _ := r.client.GetOrg(owner.UserName.ValueString())
		if res.StatusCode == 404 {
			// Use Forgejo client to check if owner is user
			// TODO: check to see if owner is user or raise an error

			resp.Diagnostics.AddError(
				"Owner not found",
				fmt.Sprintf(
					"Neither organization nor user with name %s exists.",
					owner.UserName.String(),
				),
			)

			return

			// User exists -> user repository
			// ownerType = "user"
		} else {
			// Org exists -> org repository
			ownerType = "org"
		}
	}

	var (
		rep *forgejo.Repository
		res *forgejo.Response
	)

	switch ownerType {
	case "org":
		// Use Forgejo client to create new org repository
		rep, res, err = r.client.CreateOrgRepo(
			owner.UserName.ValueString(),
			opts,
		)
	case "personal":
		// Use Forgejo client to create new personal repository
		rep, res, err = r.client.CreateRepo(opts)
	}
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 409:
			msg = fmt.Sprintf(
				"Repository with name %s already exists: %s",
				data.Name.String(),
				err,
			)
		case 422:
			msg = fmt.Sprintf("Input validation error: %s", err)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to create repository", msg)

		return
	}

	// TODO: Call API again to modify remaining attributes (e.g. website)

	// Map response body to model
	data.from(rep)
	diags = data.ownerFrom(ctx, rep)
	resp.Diagnostics.Append(diags...)
	diags = data.permissionsFrom(ctx, rep)
	resp.Diagnostics.Append(diags...)
	diags = data.internalTrackerFrom(ctx, rep)
	resp.Diagnostics.Append(diags...)
	diags = data.externalTrackerFrom(ctx, rep)
	resp.Diagnostics.Append(diags...)
	diags = data.externalWikiFrom(ctx, rep)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *repositoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	defer un(trace(ctx, "Read repository resource"))

	var (
		data  repositoryResourceModel
		owner repositoryResourceUser
	)

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read repository owner into model
	diags = data.Owner.As(ctx, &owner, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Get repository by name", map[string]any{
		"owner": owner.UserName.ValueString(),
		"name":  data.Name.ValueString(),
	})

	// Use Forgejo client to get repository by owner and name
	rep, res, err := r.client.GetRepo(
		owner.UserName.ValueString(),
		data.Name.ValueString(),
	)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf(
				"Repository with owner %s and name %s not found: %s",
				owner.UserName.String(),
				data.Name.String(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to get repository by name", msg)

		return
	}

	// Map response body to model
	data.from(rep)
	diags = data.ownerFrom(ctx, rep)
	resp.Diagnostics.Append(diags...)
	diags = data.permissionsFrom(ctx, rep)
	resp.Diagnostics.Append(diags...)
	diags = data.internalTrackerFrom(ctx, rep)
	resp.Diagnostics.Append(diags...)
	diags = data.externalTrackerFrom(ctx, rep)
	resp.Diagnostics.Append(diags...)
	diags = data.externalWikiFrom(ctx, rep)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *repositoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	defer un(trace(ctx, "Update repository resource"))

	var (
		data  repositoryResourceModel
		state repositoryResourceModel
		owner repositoryResourceUser
	)

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform prior state data into the model
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read repository owner into model
	if !data.Owner.IsUnknown() {
		diags = data.Owner.As(ctx, &owner, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
	} else {
		diags = state.Owner.As(ctx, &owner, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Update repository", map[string]any{
		"owner":                       owner.UserName.ValueString(),
		"name":                        data.Name.ValueString(),
		"description":                 data.Description.ValueString(),
		"website":                     data.Website.ValueString(),
		"private":                     data.Private.ValueBool(),
		"template":                    data.Template.ValueBool(),
		"has_issues":                  data.HasIssues.ValueBool(),
		"internal_tracker":            data.InternalTracker.String(),
		"external_tracker":            data.ExternalTracker.String(),
		"has_wiki":                    data.HasWiki.ValueBool(),
		"external_wiki":               data.ExternalWiki.String(),
		"default_branch":              data.DefaultBranch.ValueString(),
		"has_pull_requests":           data.HasPullRequests.ValueBool(),
		"has_projects":                data.HasProjects.ValueBool(),
		"has_releases":                data.HasReleases.ValueBool(),
		"has_packages":                data.HasPackages.ValueBool(),
		"has_actions":                 data.HasActions.ValueBool(),
		"ignore_whitespace_conflicts": data.IgnoreWhitespaceConflicts.ValueBool(),
		"allow_merge_commits":         data.AllowMerge.ValueBool(),
		"allow_rebase":                data.AllowRebase.ValueBool(),
		"allow_rebase_explicit":       data.AllowRebaseMerge.ValueBool(),
		"allow_squash_merge":          data.AllowSquash.ValueBool(),
		"archived":                    data.Archived.ValueBool(),
		"mirror_interval":             data.MirrorInterval.ValueString(),
		// "allow_manual_merge": data.AllowManualMerge.ValueBool(),
		// "autodetect_manual_merge": data.AutodetectManualMerge.ValueBool(),
		// "default_merge_style":
	})

	// Generate API request body from plan
	opts := forgejo.EditRepoOption{
		Name:                      data.Name.ValueStringPointer(),
		Description:               data.Description.ValueStringPointer(),
		Website:                   data.Website.ValueStringPointer(),
		Private:                   data.Private.ValueBoolPointer(),
		Template:                  data.Template.ValueBoolPointer(),
		HasIssues:                 data.HasIssues.ValueBoolPointer(),
		HasWiki:                   data.HasWiki.ValueBoolPointer(),
		DefaultBranch:             data.DefaultBranch.ValueStringPointer(),
		HasPullRequests:           data.HasPullRequests.ValueBoolPointer(),
		HasProjects:               data.HasProjects.ValueBoolPointer(),
		HasReleases:               data.HasReleases.ValueBoolPointer(),
		HasPackages:               data.HasPackages.ValueBoolPointer(),
		HasActions:                data.HasActions.ValueBoolPointer(),
		IgnoreWhitespaceConflicts: data.IgnoreWhitespaceConflicts.ValueBoolPointer(),
		AllowMerge:                data.AllowMerge.ValueBoolPointer(),
		AllowRebase:               data.AllowRebase.ValueBoolPointer(),
		AllowRebaseMerge:          data.AllowRebaseMerge.ValueBoolPointer(),
		AllowSquash:               data.AllowSquash.ValueBoolPointer(),
		Archived:                  data.Archived.ValueBoolPointer(),
		MirrorInterval:            data.MirrorInterval.ValueStringPointer(),
		// AllowManualMerge: data.AllowManualMerge.ValueBoolPointer(),
		// AutodetectManualMerge: data.AutodetectManualMerge.ValueBoolPointer(),
		// DefaultMergeStyle:
	}

	// Read objects into request body
	diags = data.internalTrackerTo(ctx, opts.InternalTracker)
	resp.Diagnostics.Append(diags...)
	diags = data.externalTrackerTo(ctx, opts.ExternalTracker)
	resp.Diagnostics.Append(diags...)
	diags = data.externalWikiTo(ctx, opts.ExternalWiki)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate API request body
	// err := opts.Validate()
	// if err != nil {
	// 	resp.Diagnostics.AddError("Input validation error", err.Error())

	// 	return
	// }

	// Use Forgejo client to update existing repository
	rep, res, err := r.client.EditRepo(
		owner.UserName.ValueString(),
		state.Name.ValueString(),
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
				"Repository with owner %s and name %s forbidden: %s",
				owner.UserName.String(),
				state.Name.String(),
				err,
			)
		case 404:
			msg = fmt.Sprintf(
				"Repository with owner %s and name %s not found: %s",
				owner.UserName.String(),
				state.Name.String(),
				err,
			)
		case 422:
			msg = fmt.Sprintf("Input validation error: %s", err)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to update repository", msg)

		return
	}

	// Map response body to model
	data.from(rep)
	diags = data.ownerFrom(ctx, rep)
	resp.Diagnostics.Append(diags...)
	diags = data.permissionsFrom(ctx, rep)
	resp.Diagnostics.Append(diags...)
	diags = data.internalTrackerFrom(ctx, rep)
	resp.Diagnostics.Append(diags...)
	diags = data.externalTrackerFrom(ctx, rep)
	resp.Diagnostics.Append(diags...)
	diags = data.externalWikiFrom(ctx, rep)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *repositoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	defer un(trace(ctx, "Delete repository resource"))

	var (
		data  repositoryResourceModel
		owner repositoryResourceUser
	)

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read repository owner into model
	diags = data.Owner.As(ctx, &owner, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Delete repository", map[string]any{
		"owner": owner.UserName.ValueString(),
		"name":  data.Name.ValueString(),
	})

	// Use Forgejo client to delete existing repository
	res, err := r.client.DeleteRepo(
		owner.UserName.ValueString(),
		data.Name.ValueString(),
	)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 403:
			msg = fmt.Sprintf(
				"Repository with owner %s and name %s forbidden: %s",
				owner.UserName.String(),
				data.Name.String(),
				err,
			)
		case 404:
			msg = fmt.Sprintf(
				"Repository with owner %s and name %s not found: %s",
				owner.UserName.String(),
				data.Name.String(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to delete repository", msg)

		return
	}
}

// NewRepositoryResource is a helper function to simplify the provider implementation.
func NewRepositoryResource() resource.Resource {
	return &repositoryResource{}
}
