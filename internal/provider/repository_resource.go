package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
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
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2#Repository
type repositoryResourceModel struct {
	ID                        types.Int64  `tfsdk:"id"`
	Owner                     types.String `tfsdk:"owner"`
	Name                      types.String `tfsdk:"name"`
	FullName                  types.String `tfsdk:"full_name"`
	Description               types.String `tfsdk:"description"`
	Empty                     types.Bool   `tfsdk:"empty"`
	Private                   types.Bool   `tfsdk:"private"`
	Fork                      types.Bool   `tfsdk:"fork"`
	Template                  types.Bool   `tfsdk:"template"`
	ParentID                  types.Int64  `tfsdk:"parent_id"`
	Mirror                    types.Bool   `tfsdk:"mirror"`
	Size                      types.Int64  `tfsdk:"size"`
	HTMLURL                   types.String `tfsdk:"html_url"`
	SSHURL                    types.String `tfsdk:"ssh_url"`
	CloneURL                  types.String `tfsdk:"clone_url"`
	Website                   types.String `tfsdk:"website"`
	Stars                     types.Int64  `tfsdk:"stars_count"`
	Forks                     types.Int64  `tfsdk:"forks_count"`
	Watchers                  types.Int64  `tfsdk:"watchers_count"`
	OpenIssues                types.Int64  `tfsdk:"open_issues_count"`
	OpenPulls                 types.Int64  `tfsdk:"open_pr_counter"`
	Releases                  types.Int64  `tfsdk:"release_counter"`
	DefaultBranch             types.String `tfsdk:"default_branch"`
	Archived                  types.Bool   `tfsdk:"archived"`
	Created                   types.String `tfsdk:"created_at"`
	Updated                   types.String `tfsdk:"updated_at"`
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
	CloneAddr                 types.String `tfsdk:"clone_addr"`
	AuthToken                 types.String `tfsdk:"auth_token"`
}

// from is a helper function to load an API struct into Terraform data model.
func (m *repositoryResourceModel) from(r *forgejo.Repository) {
	m.ID = types.Int64Value(r.ID)
	m.Owner = types.StringValue(r.Owner.UserName)
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
	m.CloneAddr = types.StringValue(r.OriginalURL)
	m.Website = types.StringValue(r.Website)
	m.Stars = types.Int64Value(int64(r.Stars))
	m.Forks = types.Int64Value(int64(r.Forks))
	m.Watchers = types.Int64Value(int64(r.Watchers))
	m.OpenIssues = types.Int64Value(int64(r.OpenIssues))
	m.OpenPulls = types.Int64Value(int64(r.OpenPulls))
	m.Releases = types.Int64Value(int64(r.Releases))
	m.DefaultBranch = types.StringValue(r.DefaultBranch)

	if !m.Mirror.ValueBool() {
		// cannot archive/un-archive repository mirrors
		m.Archived = types.BoolValue(r.Archived)
	}

	m.Created = types.StringValue(r.Created.String())
	m.Updated = types.StringValue(r.Updated.String())
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

// to is a helper function to save Terraform data model into an API struct.
func (m *repositoryResourceModel) to(o *forgejo.EditRepoOption) {
	if o == nil {
		o = new(forgejo.EditRepoOption)
	}

	o.Name = m.Name.ValueStringPointer()
	o.Description = m.Description.ValueStringPointer()
	o.Website = m.Website.ValueStringPointer()
	o.Private = m.Private.ValueBoolPointer()
	o.Template = m.Template.ValueBoolPointer()
	o.HasIssues = m.HasIssues.ValueBoolPointer()
	o.HasWiki = m.HasWiki.ValueBoolPointer()
	o.DefaultBranch = m.DefaultBranch.ValueStringPointer()
	o.HasPullRequests = m.HasPullRequests.ValueBoolPointer()
	o.HasProjects = m.HasProjects.ValueBoolPointer()
	o.HasReleases = m.HasReleases.ValueBoolPointer()
	o.HasPackages = m.HasPackages.ValueBoolPointer()
	o.HasActions = m.HasActions.ValueBoolPointer()
	o.IgnoreWhitespaceConflicts = m.IgnoreWhitespaceConflicts.ValueBoolPointer()
	o.AllowMerge = m.AllowMerge.ValueBoolPointer()
	o.AllowRebase = m.AllowRebase.ValueBoolPointer()
	o.AllowRebaseMerge = m.AllowRebaseMerge.ValueBoolPointer()
	o.AllowSquash = m.AllowSquash.ValueBoolPointer()

	if !m.Mirror.ValueBool() {
		// cannot archive/un-archive repository mirrors
		o.Archived = m.Archived.ValueBoolPointer()
	}

	if m.MirrorInterval.ValueString() != "" {
		o.MirrorInterval = m.MirrorInterval.ValueStringPointer()
	}

	// o.AllowManualMerge = m.AllowManualMerge.ValueBoolPointer()
	// o.AutodetectManualMerge = m.AutodetectManualMerge.ValueBoolPointer()
	// o.DefaultMergeStyle =
}

// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2#Permission
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

// permissionsFrom is a helper function to load an API struct into Terraform data model.
func (m *repositoryResourceModel) permissionsFrom(ctx context.Context, p *forgejo.Permission) diag.Diagnostics {
	if p == nil {
		m.Permissions = types.ObjectNull(
			repositoryResourcePermissions{}.attributeTypes(),
		)
		return nil
	}

	permsElement := repositoryResourcePermissions{
		Admin: types.BoolValue(p.Admin),
		Push:  types.BoolValue(p.Push),
		Pull:  types.BoolValue(p.Pull),
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

// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2#InternalTracker
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

// internalTrackerFrom is a helper function to load an API struct into Terraform data model.
func (m *repositoryResourceModel) internalTrackerFrom(ctx context.Context, it *forgejo.InternalTracker) diag.Diagnostics {
	if it == nil {
		m.InternalTracker = types.ObjectNull(
			repositoryResourceInternalTracker{}.attributeTypes(),
		)
		return nil
	}

	intTrackerElement := repositoryResourceInternalTracker{
		EnableTimeTracker:                types.BoolValue(it.EnableTimeTracker),
		AllowOnlyContributorsToTrackTime: types.BoolValue(it.AllowOnlyContributorsToTrackTime),
		EnableIssueDependencies:          types.BoolValue(it.EnableIssueDependencies),
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

// internalTrackerTo is a helper function to save Terraform data model into an API struct.
func (m *repositoryResourceModel) internalTrackerTo(ctx context.Context, o *forgejo.EditRepoOption) diag.Diagnostics {
	if m.InternalTracker.IsNull() || m.InternalTracker.IsUnknown() {
		return nil
	}

	var intTracker repositoryResourceInternalTracker
	diags := m.InternalTracker.As(ctx, &intTracker, basetypes.ObjectAsOptions{})

	if !diags.HasError() {
		if o.InternalTracker == nil {
			o.InternalTracker = new(forgejo.InternalTracker)
		}
		o.InternalTracker.EnableTimeTracker = intTracker.EnableTimeTracker.ValueBool()
		o.InternalTracker.AllowOnlyContributorsToTrackTime = intTracker.AllowOnlyContributorsToTrackTime.ValueBool()
		o.InternalTracker.EnableIssueDependencies = intTracker.EnableIssueDependencies.ValueBool()
	}

	return diags
}

// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2#ExternalTracker
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

// externalTrackerFrom is a helper function to load an API struct into Terraform data model.
func (m *repositoryResourceModel) externalTrackerFrom(ctx context.Context, et *forgejo.ExternalTracker) diag.Diagnostics {
	if et == nil {
		m.ExternalTracker = types.ObjectNull(
			repositoryResourceExternalTracker{}.attributeTypes(),
		)
		return nil
	}

	extTrackerElement := repositoryResourceExternalTracker{
		ExternalTrackerURL:    types.StringValue(et.ExternalTrackerURL),
		ExternalTrackerFormat: types.StringValue(et.ExternalTrackerFormat),
		ExternalTrackerStyle:  types.StringValue(et.ExternalTrackerStyle),
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

// externalTrackerTo is a helper function to save Terraform data model into an API struct.
func (m *repositoryResourceModel) externalTrackerTo(ctx context.Context, o *forgejo.EditRepoOption) diag.Diagnostics {
	if m.ExternalTracker.IsNull() || m.ExternalTracker.IsUnknown() {
		return nil
	}

	var extTracker repositoryResourceExternalTracker
	diags := m.ExternalTracker.As(ctx, &extTracker, basetypes.ObjectAsOptions{})

	if !diags.HasError() {
		if o.ExternalTracker == nil {
			o.ExternalTracker = new(forgejo.ExternalTracker)
		}
		o.ExternalTracker.ExternalTrackerURL = extTracker.ExternalTrackerURL.ValueString()
		o.ExternalTracker.ExternalTrackerFormat = extTracker.ExternalTrackerFormat.ValueString()
		o.ExternalTracker.ExternalTrackerStyle = extTracker.ExternalTrackerStyle.ValueString()
	}

	return diags
}

// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2#ExternalWiki
type repositoryResourceExternalWiki struct {
	ExternalWikiURL types.String `tfsdk:"external_wiki_url"`
}

func (m repositoryResourceExternalWiki) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"external_wiki_url": types.StringType,
	}
}

// externalWikiFrom is a helper function to load an API struct into Terraform data model.
func (m *repositoryResourceModel) externalWikiFrom(ctx context.Context, ew *forgejo.ExternalWiki) diag.Diagnostics {
	if ew == nil {
		m.ExternalWiki = types.ObjectNull(
			repositoryResourceExternalWiki{}.attributeTypes(),
		)
		return nil
	}

	wikiElement := repositoryResourceExternalWiki{
		ExternalWikiURL: types.StringValue(ew.ExternalWikiURL),
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

// externalWikiTo is a helper function to save Terraform data model into an API struct.
func (m *repositoryResourceModel) externalWikiTo(ctx context.Context, o *forgejo.EditRepoOption) diag.Diagnostics {
	if m.ExternalWiki.IsNull() || m.ExternalWiki.IsUnknown() {
		return nil
	}

	var extWiki repositoryResourceExternalWiki
	diags := m.ExternalWiki.As(ctx, &extWiki, basetypes.ObjectAsOptions{})

	if !diags.HasError() {
		if o.ExternalWiki == nil {
			o.ExternalWiki = new(forgejo.ExternalWiki)
		}
		o.ExternalWiki.ExternalWikiURL = extWiki.ExternalWikiURL.ValueString()
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
		Description: "Forgejo repository resource.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the repository.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"owner": schema.StringAttribute{
				Description: "Owner of the repository.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
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
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"template": schema.BoolAttribute{
				Description: "Is the repository a template?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"parent_id": schema.Int64Attribute{
				Description: "Numeric identifier of the parent repository.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"mirror": schema.BoolAttribute{
				Description: "Is the repository a mirror?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
				Validators: []validator.Bool{
					boolvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("clone_addr"),
					}...),
				},
			},
			"size": schema.Int64Attribute{
				Description: "Size of the repository in KiB.",
				Computed:    true,
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
			"website": schema.StringAttribute{
				Description: "Website of the repository.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"stars_count": schema.Int64Attribute{
				Description: "Number of stars of the repository.",
				Computed:    true,
			},
			"forks_count": schema.Int64Attribute{
				Description: "Number of forks of the repository.",
				Computed:    true,
			},
			"watchers_count": schema.Int64Attribute{
				Description: "Number of watchers of the repository.",
				Computed:    true,
			},
			"open_issues_count": schema.Int64Attribute{
				Description: "Number of open issues of the repository.",
				Computed:    true,
			},
			"open_pr_counter": schema.Int64Attribute{
				Description: "Number of open pull requests of the repository.",
				Computed:    true,
			},
			"release_counter": schema.Int64Attribute{
				Description: "Number of releases of the repository.",
				Computed:    true,
			},
			"default_branch": schema.StringAttribute{
				Description: "Default branch of the repository.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"archived": schema.BoolAttribute{
				Description: "Is the repository archived?",
				Optional:    true,
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
			"updated_at": schema.StringAttribute{
				Description: "Time at which the repository was updated.",
				Computed:    true,
			},
			"permissions": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"admin": schema.BoolAttribute{
						Description: "Allowed to administer?",
						Computed:    true,
					},
					"push": schema.BoolAttribute{
						Description: "Allowed to push?",
						Computed:    true,
					},
					"pull": schema.BoolAttribute{
						Description: "Allowed to pull?",
						Computed:    true,
					},
				},
				Description: "Permissions of the repository.",
				Computed:    true,
			},
			"has_issues": schema.BoolAttribute{
				Description: "Is the repository issue tracker enabled?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"internal_tracker": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"enable_time_tracker": schema.BoolAttribute{
						Description: "Enable time tracking?",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"allow_only_contributors_to_track_time": schema.BoolAttribute{
						Description: "Let only contributors track time?",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"enable_issue_dependencies": schema.BoolAttribute{
						Description: "Enable dependencies for issues and pull requests?",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
				},
				Description: "Settings for built-in issue tracker.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.Object{
					objectvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("external_tracker"),
					}...),
				},
			},
			"external_tracker": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"external_tracker_url": schema.StringAttribute{
						Description: "URL of external issue tracker.",
						Required:    true,
					},
					"external_tracker_format": schema.StringAttribute{
						Description: "External Issue Tracker URL Format. Use the placeholders {user}, {repo} and {index} for the username, repository name and issue index.",
						Required:    true,
					},
					"external_tracker_style": schema.StringAttribute{
						Description: "External Issue Tracker Number Format, either `numeric` or `alphanumeric`.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("numeric"),
						Validators: []validator.String{
							stringvalidator.OneOf([]string{"numeric", "alphanumeric"}...),
						},
					},
				},
				Description: "Settings for external issue tracker.",
				Optional:    true,
				Validators: []validator.Object{
					objectvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("internal_tracker"),
					}...),
				},
			},
			"has_wiki": schema.BoolAttribute{
				Description: "Is the repository wiki enabled?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"external_wiki": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"external_wiki_url": schema.StringAttribute{
						Description: "URL of external wiki.",
						Required:    true,
					},
				},
				Description: "Settings for external wiki.",
				Optional:    true,
			},
			"has_pull_requests": schema.BoolAttribute{
				Description: "Are repository pull requests enabled?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"has_projects": schema.BoolAttribute{
				Description: "Are repository projects enabled?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"has_releases": schema.BoolAttribute{
				Description: "Are repository releases enabled?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"has_packages": schema.BoolAttribute{
				Description: "Is the repository package registry enabled?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"has_actions": schema.BoolAttribute{
				Description: "Are integrated CI/CD pipelines enabled?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"ignore_whitespace_conflicts": schema.BoolAttribute{
				Description: "Are whitespace conflicts ignored?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"allow_merge_commits": schema.BoolAttribute{
				Description: "Allowed to create merge commit?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"allow_rebase": schema.BoolAttribute{
				Description: "Allowed to rebase then fast-forward?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"allow_rebase_explicit": schema.BoolAttribute{
				Description: "Allowed to rebase then create merge commit?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"allow_squash_merge": schema.BoolAttribute{
				Description: "Allowed to create squash commit?",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"avatar_url": schema.StringAttribute{
				Description: "Avatar URL of the repository.",
				Computed:    true,
			},
			"internal": schema.BoolAttribute{
				Description: "Is the repository internal?",
				Computed:    true,
			},
			"mirror_interval": schema.StringAttribute{
				Description: "Mirror interval of the repository.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.All(
						stringvalidator.AlsoRequires(path.Expressions{
							path.MatchRoot("mirror"),
						}...),
						stringvalidator.RegexMatches(
							regexp.MustCompile("^(0|[1-9][0-9]*)h[1-5]?[0-9]m[1-5]?[0-9]s$"),
							"must follow '23h59m59s' format",
						),
					),
				},
			},
			"mirror_updated": schema.StringAttribute{
				Description: "Time at which the repository mirror was updated.",
				Computed:    true,
			},
			"default_merge_style": schema.StringAttribute{
				Description: "Default merge style of the repository.",
				Computed:    true,
			},
			"issue_labels": schema.StringAttribute{
				Description: "Issue Label set to use.",
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
				Description: "Gitignores to use.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"license": schema.StringAttribute{
				Description: "License to use.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"readme": schema.StringAttribute{
				Description: "Readme of the repository to create.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"trust_model": schema.StringAttribute{
				Description: "TrustModel of the repository.",
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
			"clone_addr": schema.StringAttribute{
				Description: "Migrate / clone from URL.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"auth_token": schema.StringAttribute{
				Description: "API token for authenticating with migrate / clone URL.",
				Optional:    true,
				Sensitive:   true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("clone_addr"),
					}...),
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

	var data repositoryResourceModel

	// Read Terraform plan data into model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var (
		rep *forgejo.Repository
		res *forgejo.Response
		err error
	)

	// Determine if repository is to be migrated
	if data.CloneAddr.ValueString() != "" {
		tflog.Info(ctx, "Migrate repository", map[string]any{
			"name":              data.Name.ValueString(),
			"owner":             data.Owner.ValueString(),
			"clone_addr":        data.CloneAddr.ValueString(),
			"auth_token":        data.AuthToken.ValueString(),
			"mirror":            data.Mirror.ValueBool(),
			"private":           data.Private.ValueBool(),
			"description":       data.Description.ValueString(),
			"has_wiki":          data.HasWiki.ValueBool(),
			"has_issues":        data.HasIssues.ValueBool(),
			"has_pull_requests": data.HasPullRequests.ValueBool(),
			"has_releases":      data.HasReleases.ValueBool(),
			"mirror_interval":   data.MirrorInterval.ValueString(),
		})

		// Generate API request body from plan
		copts := forgejo.MigrateRepoOption{
			RepoName:  data.Name.ValueString(),
			RepoOwner: data.Owner.ValueString(),
			CloneAddr: data.CloneAddr.ValueString(),
			// Service:      forgejo.GitServiceType(""),
			// AuthUsername: "",
			// AuthPassword: "",
			AuthToken:   data.AuthToken.ValueString(),
			Mirror:      data.Mirror.ValueBool(),
			Private:     data.Private.ValueBool(),
			Description: data.Description.ValueString(),
			Wiki:        data.HasWiki.ValueBool(),
			// Milestones:     false,
			// Labels:         false,
			Issues:         data.HasIssues.ValueBool(),
			PullRequests:   data.HasPullRequests.ValueBool(),
			Releases:       data.HasReleases.ValueBool(),
			MirrorInterval: data.MirrorInterval.ValueString(),
			// LFS:            false,
			// LFSEndpoint:    "",
		}

		// Validate API request body
		err = copts.Validate(r.client)
		if err != nil {
			resp.Diagnostics.AddError("Input validation error", err.Error())

			return
		}

		// Use Forgejo client to create new repository migration
		rep, res, err = r.client.MigrateRepo(copts)
	} else {
		tflog.Info(ctx, "Create repository", map[string]any{
			"owner":          data.Owner.ValueString(),
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
		copts := forgejo.CreateRepoOption{
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
		err = copts.Validate(r.client)
		if err != nil {
			resp.Diagnostics.AddError("Input validation error", err.Error())

			return
		}

		// Determine type of repository
		var ownerType string
		if data.Owner.ValueString() == "" {
			// No owner -> personal repository
			ownerType = "personal"
		} else {
			// Use Forgejo client to check if owner is org
			_, res, _ := r.client.GetOrg(data.Owner.ValueString())
			if res.StatusCode == 404 {
				// Use Forgejo client to check if owner is user
				_, res, _ = r.client.GetUserInfo(data.Owner.ValueString())
				if res.StatusCode == 404 {
					resp.Diagnostics.AddError(
						"Owner not found",
						fmt.Sprintf(
							"Neither organization nor user with name %s exists.",
							data.Owner.String(),
						),
					)

					return
				}
				// User exists -> user repository
				ownerType = "user"
			} else {
				// Org exists -> org repository
				ownerType = "org"
			}
		}

		switch ownerType {
		case "org":
			// Use Forgejo client to create new org repository
			rep, res, err = r.client.CreateOrgRepo(
				data.Owner.ValueString(),
				copts,
			)
		case "personal":
			// Use Forgejo client to create new personal repository
			rep, res, err = r.client.CreateRepo(copts)
		case "user":
			// Use Forgejo client to create new user repository
			rep, res, err = r.client.AdminCreateRepo(
				data.Owner.ValueString(),
				copts,
			)
		}
	}

	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
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
		resp.Diagnostics.AddError("Unable to create repository", msg)

		return
	}

	tflog.Info(ctx, "Update repository", map[string]any{
		"owner":                       rep.Owner.UserName,
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
	eopts := forgejo.EditRepoOption{}
	data.to(&eopts)
	diags = data.internalTrackerTo(ctx, &eopts)
	diags.Append(data.externalTrackerTo(ctx, &eopts)...)
	diags.Append(data.externalWikiTo(ctx, &eopts)...)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate API request body
	// err := eopts.Validate()
	// if err != nil {
	// 	resp.Diagnostics.AddError("Input validation error", err.Error())

	// 	return
	// }

	// Use Forgejo client to update existing repository
	rep, res, err = r.client.EditRepo(
		rep.Owner.UserName,
		data.Name.ValueString(),
		eopts,
	)
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 403:
			msg = fmt.Sprintf(
				"Repository with owner '%s' and name %s forbidden: %s",
				rep.Owner.UserName,
				data.Name.String(),
				err,
			)
		case 404:
			msg = fmt.Sprintf(
				"Repository with owner '%s' and name %s not found: %s",
				rep.Owner.UserName,
				data.Name.String(),
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
	diags = data.permissionsFrom(ctx, rep.Permissions)
	diags.Append(data.internalTrackerFrom(ctx, rep.InternalTracker)...)
	diags.Append(data.externalTrackerFrom(ctx, rep.ExternalTracker)...)
	diags.Append(data.externalWikiFrom(ctx, rep.ExternalWiki)...)
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

	var data repositoryResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Get repository by id", map[string]any{
		"id": data.ID.ValueInt64(),
	})

	// Use Forgejo client to get repository by id
	rep, res, err := r.client.GetRepoByID(data.ID.ValueInt64())
	if err != nil {
		tflog.Error(ctx, "Error", map[string]any{
			"status": res.Status,
		})

		var msg string
		switch res.StatusCode {
		case 404:
			msg = fmt.Sprintf(
				"Repository with id %d not found: %s",
				data.ID.ValueInt64(),
				err,
			)
		default:
			msg = fmt.Sprintf("Unknown error: %s", err)
		}
		resp.Diagnostics.AddError("Unable to get repository by id", msg)

		return
	}

	// Map response body to model
	data.from(rep)
	diags = data.permissionsFrom(ctx, rep.Permissions)
	diags.Append(data.internalTrackerFrom(ctx, rep.InternalTracker)...)
	diags.Append(data.externalTrackerFrom(ctx, rep.ExternalTracker)...)
	diags.Append(data.externalWikiFrom(ctx, rep.ExternalWiki)...)
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
	owner := data.Owner.ValueString()
	if owner == "" {
		owner = state.Owner.ValueString()
	}

	tflog.Info(ctx, "Update repository", map[string]any{
		"owner":                       owner,
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
	opts := forgejo.EditRepoOption{}
	data.to(&opts)
	diags = data.internalTrackerTo(ctx, &opts)
	diags.Append(data.externalTrackerTo(ctx, &opts)...)
	diags.Append(data.externalWikiTo(ctx, &opts)...)
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
		owner,
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
				"Repository with owner '%s' and name %s forbidden: %s",
				owner,
				state.Name.String(),
				err,
			)
		case 404:
			msg = fmt.Sprintf(
				"Repository with owner '%s' and name %s not found: %s",
				owner,
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
	diags = data.permissionsFrom(ctx, rep.Permissions)
	diags.Append(data.internalTrackerFrom(ctx, rep.InternalTracker)...)
	diags.Append(data.externalTrackerFrom(ctx, rep.ExternalTracker)...)
	diags.Append(data.externalWikiFrom(ctx, rep.ExternalWiki)...)
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

	var data repositoryResourceModel

	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Delete repository", map[string]any{
		"owner": data.Owner.ValueString(),
		"name":  data.Name.ValueString(),
	})

	// Use Forgejo client to delete existing repository
	res, err := r.client.DeleteRepo(
		data.Owner.ValueString(),
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
