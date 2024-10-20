package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type ProjectService service

func (v *ProjectService) AddDomainToProject(ctx context.Context, idOrName string, slug string, teamId string, req AddDomainToProjectRequest) (*ProjectDomainInfo, error) {
	path := fmt.Sprintf("/v10/projects/%s/domains", idOrName)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &ProjectDomainInfo{}
	err := v.client.http.Post(ctx, path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *ProjectService) CreateProject(ctx context.Context, slug string, teamId string, req CreateProjectRequest) (*ProjectInfo, error) {
	path := "/v10/projects"
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}

	result := &ProjectInfo{}
	err := v.client.http.Post(ctx, path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *ProjectService) CreateEnvironmentVariables(ctx context.Context, idOrName string, slug string, teamId string, req []ProjectEnv) (*CreateProjectEnvResponse, error) {
	path := fmt.Sprintf("/v10/projects/%s/env", idOrName)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}

	result := &CreateProjectEnvResponse{}
	err := v.client.http.Post(ctx, path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *ProjectService) DeleteProject(ctx context.Context, idOrName string, slug string, teamId string) (*string, error) {
	path := fmt.Sprintf("/v9/projects/%s", idOrName)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}

	result := ""
	err := v.client.http.Delete(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *ProjectService) EditEnvironmentVariable(ctx context.Context, idOrName string, id string, slug string, teamId string, req ProjectEnv) (*ProjectEnv, error) {
	path := fmt.Sprintf("/v9/projects/%s/env/%s", idOrName, id)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}

	result := &ProjectEnv{}
	err := v.client.http.Patch(ctx, path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *ProjectService) ListEnvironmentVariable(ctx context.Context, idOrName string, slug string, teamId string, decrypt string,
	gitBranch string, source string) ([]ProjectEnv, error) {
	path := fmt.Sprintf("/v9/projects/%s/env", idOrName)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(decrypt) > 0 {
		queryParams.Add("decrypt", decrypt)
	}
	if len(gitBranch) > 0 {
		queryParams.Add("gitBranch", gitBranch)
	}
	if len(source) > 0 {
		queryParams.Add("source", source)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}

	result := []ProjectEnv{}
	err := v.client.http.Get(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *ProjectService) GetProject(ctx context.Context, idOrName string, slug string, teamId string) (*ProjectInfo, error) {
	path := fmt.Sprintf("/v9/projects/%s", idOrName)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}

	result := &ProjectInfo{}
	err := v.client.http.Get(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *ProjectService) GetProjectDomain(ctx context.Context, idOrName string, domain string, slug string, teamId string) (*ProjectDomainInfo, error) {
	path := fmt.Sprintf("/v9/projects/%s/domains/%s", idOrName, domain)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}

	result := &ProjectDomainInfo{}
	err := v.client.http.Get(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *ProjectService) ListProjectDomain(ctx context.Context, idOrName string, slug string, teamId string) (*ListProjectDomainResponse, error) {
	path := fmt.Sprintf("/v9/projects/%s/domains", idOrName)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}

	result := &ListProjectDomainResponse{}
	err := v.client.http.Get(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *ProjectService) GetEnvironmentVariable(ctx context.Context, idOrName string, id string, slug string, teamId string) (*ProjectEnv, error) {
	path := fmt.Sprintf("/v1/projects/%s/env/%s", idOrName, id)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}

	result := &ProjectEnv{}
	err := v.client.http.Get(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *ProjectService) ListProject(ctx context.Context, slug string, teamId string) (*ListProjectResponse, error) {
	path := "/v9/projects"
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}

	result := &ListProjectResponse{}
	err := v.client.http.Get(ctx, path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *ProjectService) ListProjectAlias(ctx context.Context, projectId string, slug string, teamId string) (*ListProjectAliasResponse, error) {
	path := fmt.Sprintf("/v1/projects/%s/promote/aliases", projectId)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}

	result := &ListProjectAliasResponse{}
	err := v.client.http.Get(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *ProjectService) PauseProject(ctx context.Context, projectId string, slug string, teamId string) (*string, error) {
	path := fmt.Sprintf("/v1/projects/%s/pause", projectId)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}

	result := ""
	err := v.client.http.Post(ctx, path, nil, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *ProjectService) RemoveDomainFromProject(ctx context.Context, idOrName string, domain string, slug string, teamId string) (*string, error) {
	path := fmt.Sprintf("/v9/projects/%s/domains/%s", idOrName, domain)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}

	result := ""
	err := v.client.http.Delete(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *ProjectService) RemoveEnvironmentVariable(ctx context.Context, idOrName string, id string, slug string, teamId string) (*string, error) {
	path := fmt.Sprintf("/v9/projects/%s/env/%s", idOrName, id)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}

	result := ""
	err := v.client.http.Delete(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *ProjectService) PointsProductionDomains(ctx context.Context, projectId string, deploymentId string, slug string, teamId string) (*string, error) {
	path := fmt.Sprintf("/v10/projects/%s/promote/%s", projectId, deploymentId)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}

	result := ""
	err := v.client.http.Post(ctx, path, nil, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *ProjectService) UnpauseProject(ctx context.Context, projectId string, slug string, teamId string) (*string, error) {
	path := fmt.Sprintf("/v1/projects/%s/unpause", projectId)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}

	result := ""
	err := v.client.http.Post(ctx, path, nil, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *ProjectService) UpdateProject(ctx context.Context, idOrName string, slug string, teamId string, req CreateProjectRequest) (*ProjectInfo, error) {
	path := fmt.Sprintf("/v9/projects/%s", idOrName)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}

	result := &ProjectInfo{}
	err := v.client.http.Patch(ctx, path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *ProjectService) UpdateProjectDataCache(ctx context.Context, idOrName string, slug string, teamId string, disabled bool) (*ProjectInfo, error) {
	path := fmt.Sprintf("/v1/data-cache/projects/%s", idOrName)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	req := struct {
		Disabled bool `json:"disabled"`
	}{
		Disabled: disabled,
	}
	result := &ProjectInfo{}
	err := v.client.http.Patch(ctx, path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *ProjectService) UpdateProjectDomain(ctx context.Context, idOrName string, domain string,
	slug string, teamId string, req UpdateProjectDomainRequest) (*ProjectDomainInfo, error) {
	path := fmt.Sprintf("/v9/projects/%s/domains/%s", idOrName, domain)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	result := &ProjectDomainInfo{}
	err := v.client.http.Patch(ctx, path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *ProjectService) UpdateProtectionBypass(ctx context.Context, idOrName string, slug string, teamId string, req UpdateProtectionBypassRequest,
) (*string, error) {
	path := fmt.Sprintf("/v1/projects/%s/protection-bypass", idOrName)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	result := ""
	err := v.client.http.Patch(ctx, path, req, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *ProjectService) VerifyProjectDomain(ctx context.Context, idOrName string, domain string,
	slug string, teamId string) (*ProjectDomainInfo, error) {
	path := fmt.Sprintf("/v9/projects/%s/domains/%s/verify", idOrName, domain)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	result := &ProjectDomainInfo{}
	err := v.client.http.Post(ctx, path, nil, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type CreateProjectRequest struct {
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

type AddDomainToProjectRequest struct {
	Name               string `json:"name,omitempty"`
	GitBranch          string `json:"gitBranch,omitempty"`
	Redirect           string `json:"redirect,omitempty"`
	RedirectStatusCode int    `json:"redirectStatusCode,omitempty"`
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

type UpdateProjectDomainRequest struct {
	GitBranch           string `json:"gitBranch,omitempty"`
	Redirect            string `json:"redirect,omitempty"`
	RedirectStatusCode  int    `json:"redirectStatusCode,omitempty"`
	CustomEnvironmentId string `json:"customEnvironmentId,omitempty"`
}

type UpdateProtectionBypassRequest struct {
	Revoke ProtectionBypass `json:"revoke,omitempty"`
}

type ProtectionBypass struct {
	Regenerate bool   `json:"regenerate,omitempty"`
	Ssecret    string `json:"secret,omitempty"`
}
