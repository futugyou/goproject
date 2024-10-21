package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type ProjectService service

type AddDomainToProjectRequest struct {
	Name               string `json:"name,omitempty"`
	GitBranch          string `json:"gitBranch,omitempty"`
	Redirect           string `json:"redirect,omitempty"`
	RedirectStatusCode int    `json:"redirectStatusCode,omitempty"`
	IdOrName           string `json:"-"`
	BaseUrlParameter   `json:"-"`
}

func (v *ProjectService) AddDomainToProject(ctx context.Context, request AddDomainToProjectRequest) (*ProjectDomainInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v10/projects/%s/domains", request.IdOrName),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ProjectDomainInfo{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type UpsertProjectRequest struct {
	Name                                 string                `json:"name"`
	BuildCommand                         string                `json:"buildCommand,omitempty"`
	CommandForIgnoringBuildStep          string                `json:"commandForIgnoringBuildStep,omitempty"`
	DevCommand                           string                `json:"devCommand,omitempty"`
	EnvironmentVariables                 []EnvironmentVariable `json:"environmentVariables,omitempty"`
	Framework                            string                `json:"framework,omitempty"`
	GitRepository                        GitRepository         `json:"gitRepository,omitempty"`
	InstallCommand                       string                `json:"installCommand,omitempty"`
	OutputDirectory                      string                `json:"outputDirectory,omitempty"`
	PublicSource                         bool                  `json:"publicSource,omitempty"`
	RootDirectory                        string                `json:"rootDirectory,omitempty"`
	ServerlessFunctionRegion             string                `json:"serverlessFunctionRegion,omitempty"`
	ServerlessFunctionZeroConfigFailover string                `json:"serverlessFunctionZeroConfigFailover,omitempty"`
	SkipGitConnectDuringLink             bool                  `json:"skipGitConnectDuringLink,omitempty"`
	IdOrName                             string                `json:"-"`
	BaseUrlParameter                     `json:"-"`
}

func (v *ProjectService) CreateProject(ctx context.Context, request UpsertProjectRequest) (*ProjectInfo, error) {
	u := &url.URL{
		Path: "/v10/projects",
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ProjectInfo{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type CreateEnvRequest struct {
	IdOrName         string  `json:"-"`
	Upsert           *string `json:"-"`
	BaseUrlParameter `json:"-"`
	Envs             []ProjectEnv
}

func (v *ProjectService) CreateEnvironmentVariables(ctx context.Context, request CreateEnvRequest) (*CreateProjectEnvResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v10/projects/%s/env", request.IdOrName),
	}
	params := request.GetUrlValues()
	if request.Upsert != nil {
		params.Add("upsert", *request.Upsert)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &CreateProjectEnvResponse{}
	if err := v.client.http.Post(ctx, path, request.Envs, response); err != nil {
		return nil, err
	}
	return response, nil
}

type DeleteProjectRequest struct {
	IdOrName         string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *ProjectService) DeleteProject(ctx context.Context, request DeleteProjectRequest) (*string, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v9/projects/%s", request.IdOrName),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := ""
	if err := v.client.http.Delete(ctx, path, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

type EditEnvRequest struct {
	Id                   string `json:"-"`
	IdOrName             string `json:"-"`
	BaseUrlParameter     `json:"-"`
	Key                  string   `json:"key,omitempty"`
	Target               string   `json:"target,omitempty"`
	Type                 string   `json:"type,omitempty"`
	Value                string   `json:"value,omitempty"`
	GitBranch            string   `json:"gitBranch,omitempty"`
	Comment              string   `json:"comment,omitempty"`
	CustomEnvironmentIds []string `json:"customEnvironmentIds,omitempty"`
}

func (v *ProjectService) EditEnvironmentVariable(ctx context.Context, request EditEnvRequest) (*ProjectEnv, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v9/projects/%s/env/%s", request.IdOrName, request.Id),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ProjectEnv{}
	if err := v.client.http.Patch(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type ListEnvParameter struct {
	IdOrName         string `json:"-"`
	BaseUrlParameter `json:"-"`
	GitBranch        *string
	Source           *string
}

func (v *ProjectService) ListEnvironmentVariable(ctx context.Context, request ListEnvParameter) ([]ProjectEnv, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v9/projects/%s/env", request.IdOrName),
	}
	params := request.GetUrlValues()
	if request.GitBranch != nil {
		params.Add("gitBranch", *request.GitBranch)
	}
	if request.Source != nil {
		params.Add("source", *request.Source)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := []ProjectEnv{}
	if err := v.client.http.Get(ctx, path, &response); err != nil {
		return nil, err
	}
	return response, nil
}

type GetProjectParameter struct {
	IdOrName         string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *ProjectService) GetProject(ctx context.Context, request GetProjectParameter) (*ProjectInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v9/projects/%s", request.IdOrName),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ProjectInfo{}
	if err := v.client.http.Get(ctx, path, &response); err != nil {
		return nil, err
	}
	return response, nil
}

type GetProjectDomainParameter struct {
	IdOrName         string `json:"-"`
	Domain           string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *ProjectService) GetProjectDomain(ctx context.Context, request GetProjectDomainParameter) (*ProjectDomainInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v9/projects/%s/domains/%s", request.IdOrName, request.Domain),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ProjectDomainInfo{}
	if err := v.client.http.Get(ctx, path, &response); err != nil {
		return nil, err
	}
	return response, nil
}

type ListProjectDomainParameter struct {
	IdOrName         string `json:"-"`
	BaseUrlParameter `json:"-"`
	GitBranch        *string
	Limit            *string
	Order            *string // ASC DESC
	Production       *string // true false
	Redirect         *string
	Redirects        *string // true false
	Since            *string
	Target           *string // production preview
	Until            *string
	Verified         *string // true false
}

func (v *ProjectService) ListProjectDomain(ctx context.Context, request ListProjectDomainParameter) (*ListProjectDomainResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v9/projects/%s/domains", request.IdOrName),
	}
	params := request.GetUrlValues()
	if request.GitBranch != nil {
		params.Add("gitBranch", *request.GitBranch)
	}
	if request.Limit != nil {
		params.Add("limit", *request.Limit)
	}
	if request.Order != nil {
		params.Add("order", *request.Order)
	}
	if request.Production != nil {
		params.Add("production", *request.Production)
	}
	if request.Redirect != nil {
		params.Add("redirect", *request.Redirect)
	}
	if request.Redirects != nil {
		params.Add("redirects", *request.Redirects)
	}
	if request.Since != nil {
		params.Add("since", *request.Since)
	}
	if request.Target != nil {
		params.Add("target", *request.Target)
	}
	if request.Until != nil {
		params.Add("until", *request.Until)
	}
	if request.Verified != nil {
		params.Add("verified", *request.Verified)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ListProjectDomainResponse{}
	if err := v.client.http.Get(ctx, path, &response); err != nil {
		return nil, err
	}
	return response, nil
}

type GetEnvParameter struct {
	Id               string `json:"-"`
	IdOrName         string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *ProjectService) GetEnvironmentVariable(ctx context.Context, request GetEnvParameter) (*ProjectEnv, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/projects/%s/env/%s", request.IdOrName, request.Id),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ProjectEnv{}
	if err := v.client.http.Get(ctx, path, &response); err != nil {
		return nil, err
	}
	return response, nil
}

type ListProjectParameter struct {
	IdOrName          string `json:"-"`
	BaseUrlParameter  `json:"-"`
	Deprecated        *string
	EdgeConfigId      *string
	EdgeConfigTokenId *string
	ExcludeRepos      *string
	From              *string
	GitForkProtection *string
	Limit             *string
	Repo              *string
	RepoId            *string
	RepoUrl           *string
	Search            *string
}

func (v *ProjectService) ListProject(ctx context.Context, request ListProjectParameter) (*ListProjectResponse, error) {
	u := &url.URL{
		Path: "/v9/projects",
	}
	params := request.GetUrlValues()
	if request.Deprecated != nil {
		params.Add("deprecated", *request.Deprecated)
	}
	if request.EdgeConfigId != nil {
		params.Add("edgeConfigId", *request.EdgeConfigId)
	}
	if request.EdgeConfigTokenId != nil {
		params.Add("edgeConfigTokenId", *request.EdgeConfigTokenId)
	}
	if request.ExcludeRepos != nil {
		params.Add("excludeRepos", *request.ExcludeRepos)
	}
	if request.From != nil {
		params.Add("from", *request.From)
	}
	if request.GitForkProtection != nil {
		params.Add("gitForkProtection", *request.GitForkProtection)
	}
	if request.Limit != nil {
		params.Add("limit", *request.Limit)
	}
	if request.Repo != nil {
		params.Add("repo", *request.Repo)
	}
	if request.RepoId != nil {
		params.Add("repoId", *request.RepoId)
	}
	if request.RepoUrl != nil {
		params.Add("repoUrl", *request.RepoUrl)
	}
	if request.Search != nil {
		params.Add("search", *request.Search)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ListProjectResponse{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type ListProjectAliasParameter struct {
	projectId        string `json:"-"`
	BaseUrlParameter `json:"-"`
	FailedOnly       *string
	Limit            *string
	Since            *string
	Until            *string
}

func (v *ProjectService) ListProjectAlias(ctx context.Context, request ListProjectAliasParameter) (*ListProjectAliasResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/projects/%s/promote/aliases", request.projectId),
	}
	params := request.GetUrlValues()
	if request.FailedOnly != nil {
		params.Add("failedOnly", *request.FailedOnly)
	}
	if request.Limit != nil {
		params.Add("limit", *request.Limit)
	}
	if request.Since != nil {
		params.Add("since", *request.Since)
	}
	if request.Until != nil {
		params.Add("until", *request.Until)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ListProjectAliasResponse{}
	if err := v.client.http.Get(ctx, path, &response); err != nil {
		return nil, err
	}
	return response, nil
}

type PauseProjectRequest struct {
	ProjectId        string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *ProjectService) PauseProject(ctx context.Context, request PauseProjectRequest) (*string, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/projects/%s/pause", request.ProjectId),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := ""
	if err := v.client.http.Post(ctx, path, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

type RemoveDomainFromProjectRequest struct {
	ProjectId        string `json:"-"`
	Domain           string
	BaseUrlParameter `json:"-"`
}

func (v *ProjectService) RemoveDomainFromProject(ctx context.Context, request RemoveDomainFromProjectRequest) (*string, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v9/projects/%s/domains/%s", request.ProjectId, request.Domain),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := ""
	if err := v.client.http.Delete(ctx, path, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

type RemoveEnvRequest struct {
	IdOrName         string
	Id               string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *ProjectService) RemoveEnvironmentVariable(ctx context.Context, request RemoveEnvRequest) (*string, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v9/projects/%s/env/%s", request.IdOrName, request.Id),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := ""
	if err := v.client.http.Delete(ctx, path, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

type PointsProductionDomainsRequest struct {
	ProjectId        string `json:"projectId"`
	DeploymentId     string `json:"deploymentId"`
	BaseUrlParameter `json:"-"`
}

func (v *ProjectService) PointsProductionDomains(ctx context.Context, request PointsProductionDomainsRequest) (*string, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v10/projects/%s/promote/%s", request.ProjectId, request.DeploymentId),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := ""
	if err := v.client.http.Post(ctx, path, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

type UnpauseProjectRequest struct {
	ProjectId        string `json:"projectId"`
	BaseUrlParameter `json:"-"`
}

func (v *ProjectService) UnpauseProject(ctx context.Context, request UnpauseProjectRequest) (*string, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v10/projects/%s/unpause", request.ProjectId),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := ""
	if err := v.client.http.Post(ctx, path, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (v *ProjectService) UpdateProject(ctx context.Context, request UpsertProjectRequest) (*ProjectInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v9/projects/%s", request.IdOrName),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ProjectInfo{}
	if err := v.client.http.Patch(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type UpdateProjectDataCacheRequest struct {
	ProjectId        string `json:""`
	Disabled         string `json:"disabled"`
	BaseUrlParameter `json:"-"`
}

func (v *ProjectService) UpdateProjectDataCache(ctx context.Context, request UpdateProjectDataCacheRequest) (*ProjectInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/data-cache/projects/%s", request.ProjectId),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ProjectInfo{}
	if err := v.client.http.Patch(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type UpdateProjectDomainRequest struct {
	GitBranch          string `json:"gitBranch,omitempty"`
	Redirect           string `json:"redirect,omitempty"`
	RedirectStatusCode int    `json:"redirectStatusCode,omitempty"`
	Domain             string `json:"-"`
	IdOrName           string `json:"-"`
	BaseUrlParameter   `json:"-"`
}

func (v *ProjectService) UpdateProjectDomain(ctx context.Context, request UpdateProjectDomainRequest) (*ProjectDomainInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v9/projects/%s/domains/%s", request.IdOrName, request.Domain),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ProjectDomainInfo{}
	if err := v.client.http.Patch(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type UpdateProtectionBypassRequest struct {
	IdOrName         string `json:"-"`
	BaseUrlParameter `json:"-"`
	Revoke           ProtectionBypass   `json:"revoke,omitempty"`
	Generate         ProtectionGenerate `json:"generate,omitempty"`
}

func (v *ProjectService) UpdateProtectionBypass(ctx context.Context, request UpdateProtectionBypassRequest) (*string, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/projects/%s/protection-bypass", request.IdOrName),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := ""
	if err := v.client.http.Patch(ctx, path, request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

type VerifyProjectDomainRequest struct {
	IdOrName         string `json:"-"`
	Domain           string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *ProjectService) VerifyProjectDomain(ctx context.Context, request VerifyProjectDomainRequest) (*ProjectDomainInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v9/projects/%s/domains/%s/verify", request.IdOrName, request.Domain),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ProjectDomainInfo{}
	if err := v.client.http.Post(ctx, path, nil, response); err != nil {
		return nil, err
	}
	return response, nil
}

type EnvironmentVariable struct {
	Key       string `json:"key,omitempty"`
	Target    string `json:"target,omitempty"`
	GitBranch string `json:"gitBranch,omitempty"`
	Type      string `json:"type,omitempty"`
	Value     string `json:"value,omitempty"`
}

type GitRepository struct {
	Repo string `json:"repo,omitempty"`
	Type string `json:"type,omitempty"`
}

type ProjectDomainInfo struct {
	Name                string         `json:"name,omitempty"`
	GitBranch           string         `json:"gitBranch,omitempty"`
	Redirect            string         `json:"redirect,omitempty"`
	RedirectStatusCode  int            `json:"redirectStatusCode,omitempty"`
	ApexName            string         `json:"apexName,omitempty"`
	CreatedAt           int            `json:"createdAt,omitempty"`
	CustomEnvironmentId string         `json:"customEnvironmentId,omitempty"`
	ProjectId           string         `json:"projectId,omitempty"`
	UpdatedAt           int            `json:"updatedAt,omitempty"`
	Verified            bool           `json:"verified,omitempty"`
	Verification        []Verification `json:"verification,omitempty"`
	Error               *VercelError   `json:"error,omitempty"`
}

type Verification struct {
	Domain string `json:"domain,omitempty"`
	Reason string `json:"reason,omitempty"`
	Type   string `json:"type,omitempty"`
	Value  string `json:"value,omitempty"`
}

type ProjectInfo struct {
	AccountId                         string       `json:"accountId,omitempty"`
	AutoAssignCustomDomains           bool         `json:"autoAssignCustomDomains,omitempty"`
	AutoAssignCustomDomainsUpdatedBy  string       `json:"autoAssignCustomDomainsUpdatedBy,omitempty"`
	AutoExposeSystemEnvs              bool         `json:"autoExposeSystemEnvs,omitempty"`
	BuildCommand                      string       `json:"buildCommand,omitempty"`
	CommandForIgnoringBuildStep       string       `json:"commandForIgnoringBuildStep,omitempty"`
	ConcurrencyBucketName             string       `json:"concurrencyBucketName,omitempty"`
	ConnectBuildsEnabled              bool         `json:"connectBuildsEnabled,omitempty"`
	ConnectConfigurationId            string       `json:"connectConfigurationId,omitempty"`
	CreatedAt                         int          `json:"createdAt,omitempty"`
	CustomerSupportCodeVisibility     bool         `json:"customerSupportCodeVisibility,omitempty"`
	DevCommand                        string       `json:"devCommand,omitempty"`
	DirectoryListing                  bool         `json:"directoryListing,omitempty"`
	EnableAffectedProjectsDeployments bool         `json:"enableAffectedProjectsDeployments,omitempty"`
	EnablePreviewFeedback             bool         `json:"enablePreviewFeedback,omitempty"`
	Crons                             interface{}  `json:"crons,omitempty"`
	DataCache                         interface{}  `json:"dataCache,omitempty"`
	DeploymentExpiration              interface{}  `json:"deploymentExpiration,omitempty"`
	Env                               interface{}  `json:"env,omitempty"`
	Framework                         string       `json:"framework,omitempty"`
	GitComments                       interface{}  `json:"gitComments,omitempty"`
	GitForkProtection                 bool         `json:"gitForkProtection,omitempty"`
	GitLFS                            bool         `json:"gitLFS,omitempty"`
	HasActiveBranches                 bool         `json:"hasActiveBranches,omitempty"`
	HasFloatingAliases                bool         `json:"hasFloatingAliases,omitempty"`
	Id                                string       `json:"id,omitempty"`
	Name                              string       `json:"name,omitempty"`
	NodeVersion                       string       `json:"nodeVersion,omitempty"`
	OutputDirectory                   string       `json:"outputDirectory,omitempty"`
	Error                             *VercelError `json:"error,omitempty"`
}

type ProjectEnv struct {
	Key                  string       `json:"key,omitempty"`
	Target               string       `json:"target,omitempty"`
	Type                 string       `json:"type,omitempty"`
	Value                string       `json:"value,omitempty"`
	GitBranch            string       `json:"gitBranch,omitempty"`
	Comment              string       `json:"comment,omitempty"`
	CustomEnvironmentIds []string     `json:"customEnvironmentIds,omitempty"`
	Error                *VercelError `json:"error,omitempty"`
}

type CreateProjectEnvResponse struct {
	Created interface{}   `json:"created,omitempty"`
	Failed  []interface{} `json:"failed,omitempty"`
	Error   *VercelError  `json:"error,omitempty"`
}

type ListProjectDomainResponse struct {
	Domains    []ProjectDomainInfo `json:"domains,omitempty"`
	Pagination Pagination          `json:"pagination,omitempty"`
	Error      *VercelError        `json:"error,omitempty"`
}

type ListProjectResponse struct {
	Projects   []ProjectInfo `json:"projects,omitempty"`
	Pagination Pagination    `json:"pagination,omitempty"`
	Error      *VercelError  `json:"error,omitempty"`
}

type ListProjectAliasResponse struct {
	Aliases    []AliasInfo  `json:"aliases,omitempty"`
	Pagination Pagination   `json:"pagination,omitempty"`
	Error      *VercelError `json:"error,omitempty"`
}

type ProtectionBypass struct {
	Regenerate bool   `json:"regenerate,omitempty"`
	Ssecret    string `json:"secret,omitempty"`
}

type ProtectionGenerate struct {
	Ssecret string `json:"secret,omitempty"`
}
