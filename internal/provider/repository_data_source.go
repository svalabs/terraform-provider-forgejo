package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &repositoryDataSource{}
	_ datasource.DataSourceWithConfigure = &repositoryDataSource{}
)

// repositoryDataSource is the data source implementation.
type repositoryDataSource struct {
	client *forgejo.Client
}

// repositoryDataSourceModel maps the data source schema data.
// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo#Repository
type repositoryDataSourceModel struct {
	ID                        types.Int64  `tfsdk:"id"`
	Owner                     types.Object `tfsdk:"owner"`
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
	OriginalURL               types.String `tfsdk:"original_url"`
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
}

// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo#User
type repositoryDataSourceUser struct {
	ID        types.Int64  `tfsdk:"id"`
	UserName  types.String `tfsdk:"login"`
	LoginName types.String `tfsdk:"login_name"`
	FullName  types.String `tfsdk:"full_name"`
	Email     types.String `tfsdk:"email"`
}

func (m repositoryDataSourceUser) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":         types.Int64Type,
		"login":      types.StringType,
		"login_name": types.StringType,
		"full_name":  types.StringType,
		"email":      types.StringType,
	}
}

// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo#Permission
type repositoryDataSourcePermissions struct {
	Admin types.Bool `tfsdk:"admin"`
	Push  types.Bool `tfsdk:"push"`
	Pull  types.Bool `tfsdk:"pull"`
}

func (m repositoryDataSourcePermissions) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"admin": types.BoolType,
		"push":  types.BoolType,
		"pull":  types.BoolType,
	}
}

// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo#InternalTracker
type repositoryDataSourceInternalTracker struct {
	EnableTimeTracker                types.Bool `tfsdk:"enable_time_tracker"`
	AllowOnlyContributorsToTrackTime types.Bool `tfsdk:"allow_only_contributors_to_track_time"`
	EnableIssueDependencies          types.Bool `tfsdk:"enable_issue_dependencies"`
}

func (m repositoryDataSourceInternalTracker) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enable_time_tracker":                   types.BoolType,
		"allow_only_contributors_to_track_time": types.BoolType,
		"enable_issue_dependencies":             types.BoolType,
	}
}

// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo#ExternalTracker
type repositoryDataSourceExternalTracker struct {
	ExternalTrackerURL    types.String `tfsdk:"external_tracker_url"`
	ExternalTrackerFormat types.String `tfsdk:"external_tracker_format"`
	ExternalTrackerStyle  types.String `tfsdk:"external_tracker_style"`
}

func (m repositoryDataSourceExternalTracker) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"external_tracker_url":    types.StringType,
		"external_tracker_format": types.StringType,
		"external_tracker_style":  types.StringType,
	}
}

// https://pkg.go.dev/codeberg.org/mvdkleijn/forgejo-sdk/forgejo#ExternalWiki
type repositoryDataSourceExternalWiki struct {
	ExternalWikiURL types.String `tfsdk:"external_wiki_url"`
}

func (m repositoryDataSourceExternalWiki) attributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"external_wiki_url": types.StringType,
	}
}

// Metadata returns the data source type name.
func (d *repositoryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository"
}

// Schema defines the schema for the data source.
func (d *repositoryDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Forgejo repository data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the repository.",
				Computed:    true,
			},
			"owner": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						Description: "Numeric identifier of the user.",
						Computed:    true,
					},
					"login": schema.StringAttribute{
						Description: "Name of the user.",
						Required:    true,
					},
					"login_name": schema.StringAttribute{
						Description: "Login name of the user.",
						Computed:    true,
					},
					"full_name": schema.StringAttribute{
						Description: "Full name of the user.",
						Computed:    true,
					},
					"email": schema.StringAttribute{
						Description: "Email address of the user.",
						Computed:    true,
					},
				},
				Description: "Owner of the repository.",
				Required:    true,
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
				Computed:    true,
			},
			"empty": schema.BoolAttribute{
				Description: "Is the repository empty?",
				Computed:    true,
			},
			"private": schema.BoolAttribute{
				Description: "Is the repository private?",
				Computed:    true,
			},
			"fork": schema.BoolAttribute{
				Description: "Is the repository a fork?",
				Computed:    true,
			},
			"template": schema.BoolAttribute{
				Description: "Is the repository a template?",
				Computed:    true,
			},
			"parent_id": schema.Int64Attribute{
				Description: "Numeric identifier of the parent repository.",
				Computed:    true,
			},
			"mirror": schema.BoolAttribute{
				Description: "Is the repository a mirror?",
				Computed:    true,
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
			"original_url": schema.StringAttribute{
				Description: "Original URL of the repository.",
				Computed:    true,
			},
			"website": schema.StringAttribute{
				Description: "Website of the repository.",
				Computed:    true,
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
				Computed:    true,
			},
			"archived": schema.BoolAttribute{
				Description: "Is the repository archived?",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Time at which the repository was created.",
				Computed:    true,
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
				Computed:    true,
			},
			"internal_tracker": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"enable_time_tracker": schema.BoolAttribute{
						Description: "Enable time tracking.",
						Computed:    true,
					},
					"allow_only_contributors_to_track_time": schema.BoolAttribute{
						Description: "Let only contributors track time.",
						Computed:    true,
					},
					"enable_issue_dependencies": schema.BoolAttribute{
						Description: "Enable dependencies for issues and pull requests.",
						Computed:    true,
					},
				},
				Description: "Settings for built-in issue tracker.",
				Computed:    true,
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
			},
			"has_wiki": schema.BoolAttribute{
				Description: "Is the repository wiki enabled?",
				Computed:    true,
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
			},
			"has_pull_requests": schema.BoolAttribute{
				Description: "Are repository pull requests enabled?",
				Computed:    true,
			},
			"has_projects": schema.BoolAttribute{
				Description: "Are repository projects enabled?",
				Computed:    true,
			},
			"has_releases": schema.BoolAttribute{
				Description: "Are repository releases enabled?",
				Computed:    true,
			},
			"has_packages": schema.BoolAttribute{
				Description: "Is the repository package registry enabled?",
				Computed:    true,
			},
			"has_actions": schema.BoolAttribute{
				Description: "Are integrated CI/CD pipelines enabled?",
				Computed:    true,
			},
			"ignore_whitespace_conflicts": schema.BoolAttribute{
				Description: "Are whitespace conflicts ignored?",
				Computed:    true,
			},
			"allow_merge_commits": schema.BoolAttribute{
				Description: "Allowed to create merge commit?",
				Computed:    true,
			},
			"allow_rebase": schema.BoolAttribute{
				Description: "Allowed to rebase then fast-forward?",
				Computed:    true,
			},
			"allow_rebase_explicit": schema.BoolAttribute{
				Description: "Allowed to rebase then create merge commit?",
				Computed:    true,
			},
			"allow_squash_merge": schema.BoolAttribute{
				Description: "Allowed to create squash commit?",
				Computed:    true,
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
				Computed:    true,
			},
			"mirror_updated": schema.StringAttribute{
				Description: "Time at which the repository mirror was updated.",
				Computed:    true,
			},
			"default_merge_style": schema.StringAttribute{
				Description: "Default merge style of the repository.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *repositoryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*forgejo.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf(
				"Expected *forgejo.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	d.client = client
}

// Read refreshes the Terraform state with the latest data.
func (d *repositoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	defer un(trace(ctx, "Read repository data source"))

	var (
		data  repositoryDataSourceModel
		owner repositoryDataSourceUser
	)

	// Read Terraform configuration data into model
	diags := req.Config.Get(ctx, &data)
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
	rep, res, err := d.client.GetRepo(
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
	data.ID = types.Int64Value(rep.ID)
	data.FullName = types.StringValue(rep.FullName)
	data.Description = types.StringValue(rep.Description)
	data.Empty = types.BoolValue(rep.Empty)
	data.Private = types.BoolValue(rep.Private)
	data.Fork = types.BoolValue(rep.Fork)
	data.Template = types.BoolValue(rep.Template)
	if rep.Parent != nil {
		data.ParentID = types.Int64Value(rep.Parent.ID)
	}
	data.Mirror = types.BoolValue(rep.Mirror)
	data.Size = types.Int64Value(int64(rep.Size))
	data.HTMLURL = types.StringValue(rep.HTMLURL)
	data.SSHURL = types.StringValue(rep.SSHURL)
	data.CloneURL = types.StringValue(rep.CloneURL)
	data.OriginalURL = types.StringValue(rep.OriginalURL)
	data.Website = types.StringValue(rep.Website)
	data.Stars = types.Int64Value(int64(rep.Stars))
	data.Forks = types.Int64Value(int64(rep.Forks))
	data.Watchers = types.Int64Value(int64(rep.Watchers))
	data.OpenIssues = types.Int64Value(int64(rep.OpenIssues))
	data.OpenPulls = types.Int64Value(int64(rep.OpenPulls))
	data.Releases = types.Int64Value(int64(rep.Releases))
	data.DefaultBranch = types.StringValue(rep.DefaultBranch)
	data.Archived = types.BoolValue(rep.Archived)
	data.Created = types.StringValue(rep.Created.String())
	data.Updated = types.StringValue(rep.Updated.String())
	data.HasIssues = types.BoolValue(rep.HasIssues)
	data.HasWiki = types.BoolValue(rep.HasWiki)
	data.HasPullRequests = types.BoolValue(rep.HasPullRequests)
	data.HasProjects = types.BoolValue(rep.HasProjects)
	data.HasReleases = types.BoolValue(rep.HasReleases)
	data.HasPackages = types.BoolValue(rep.HasPackages)
	data.HasActions = types.BoolValue(rep.HasActions)
	data.IgnoreWhitespaceConflicts = types.BoolValue(rep.IgnoreWhitespaceConflicts)
	data.AllowMerge = types.BoolValue(rep.AllowMerge)
	data.AllowRebase = types.BoolValue(rep.AllowRebase)
	data.AllowRebaseMerge = types.BoolValue(rep.AllowRebaseMerge)
	data.AllowSquash = types.BoolValue(rep.AllowSquash)
	data.AvatarURL = types.StringValue(rep.AvatarURL)
	data.Internal = types.BoolValue(rep.Internal)
	data.MirrorInterval = types.StringValue(rep.MirrorInterval)
	data.MirrorUpdated = types.StringValue(rep.MirrorUpdated.String())
	data.DefaultMergeStyle = types.StringValue(string(rep.DefaultMergeStyle))

	// Repository owner
	if rep.Owner != nil {
		ownerElement := repositoryDataSourceUser{
			ID:        types.Int64Value(rep.Owner.ID),
			UserName:  types.StringValue(rep.Owner.UserName),
			LoginName: types.StringValue(rep.Owner.LoginName),
			FullName:  types.StringValue(rep.Owner.FullName),
			Email:     types.StringValue(rep.Owner.Email),
		}
		ownerValue, diags := types.ObjectValueFrom(
			ctx,
			ownerElement.attributeTypes(),
			ownerElement,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Owner = ownerValue
	}

	// Repository permissions
	if rep.Permissions != nil {
		perms := repositoryDataSourcePermissions{
			Admin: types.BoolValue(rep.Permissions.Admin),
			Push:  types.BoolValue(rep.Permissions.Push),
			Pull:  types.BoolValue(rep.Permissions.Pull),
		}
		permsValue, diags := types.ObjectValueFrom(
			ctx,
			perms.attributeTypes(),
			perms,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Permissions = permsValue
	}

	// Internal issue tracker
	if rep.InternalTracker != nil {
		intTracker := repositoryDataSourceInternalTracker{
			EnableTimeTracker:                types.BoolValue(rep.InternalTracker.EnableTimeTracker),
			AllowOnlyContributorsToTrackTime: types.BoolValue(rep.InternalTracker.AllowOnlyContributorsToTrackTime),
			EnableIssueDependencies:          types.BoolValue(rep.InternalTracker.EnableIssueDependencies),
		}
		intTrackerValue, diags := types.ObjectValueFrom(
			ctx,
			intTracker.attributeTypes(),
			intTracker,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.InternalTracker = intTrackerValue
	}

	// External issue tracker
	if rep.ExternalTracker != nil {
		extTracker := repositoryDataSourceExternalTracker{
			ExternalTrackerURL:    types.StringValue(rep.ExternalTracker.ExternalTrackerURL),
			ExternalTrackerFormat: types.StringValue(rep.ExternalTracker.ExternalTrackerFormat),
			ExternalTrackerStyle:  types.StringValue(rep.ExternalTracker.ExternalTrackerStyle),
		}
		extTrackerValue, diags := types.ObjectValueFrom(
			ctx,
			extTracker.attributeTypes(),
			extTracker,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.ExternalTracker = extTrackerValue
	}

	// External wiki
	if rep.ExternalWiki != nil {
		wiki := repositoryDataSourceExternalWiki{
			ExternalWikiURL: types.StringValue(rep.ExternalWiki.ExternalWikiURL),
		}
		wikiValue, diags := types.ObjectValueFrom(
			ctx,
			wiki.attributeTypes(),
			wiki,
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.ExternalWiki = wikiValue
	}

	// Save data into Terraform state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// NewRepositoryDataSource is a helper function to simplify the provider implementation.
func NewRepositoryDataSource() datasource.DataSource {
	return &repositoryDataSource{}
}
