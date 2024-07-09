package vercel

import (
	"fmt"
	"log"
	"net/url"
)

func (v *VercelClient) AddDomainToProject(idOrName string, slug string, teamId string, req AddDomainToProjectRequest) (*AddDomainToProjectResponse, error) {
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
	result := &AddDomainToProjectResponse{}
	err := v.http.Post(path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) CreateProject(name string, slug string, teamId string, req CreateProjectRequest) (*ProjectInfo, error) {
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
	err := v.http.Post(path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) CreateEnvironmentVariables(idOrName string, slug string, teamId string, req []ProjectEnv) (*CreateProjectEnvResponse, error) {
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
	err := v.http.Post(path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) DeleteProject(idOrName string, slug string, teamId string) (*string, error) {
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
	err := v.http.Delete(path, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *VercelClient) EditEnvironmentVariable(idOrName string, id string, slug string, teamId string, req ProjectEnv) (*ProjectEnv, error) {
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
	err := v.http.Patch(path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) RetrieveEnvironmentVariable(idOrName string, slug string, teamId string, decrypt string,
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
	err := v.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) GetProject(idOrName string, slug string, teamId string) (*ProjectInfo, error) {
	path := fmt.Sprintf("/v9/projects/%s", idOrName)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}

	result := &ProjectInfo{}
	err := v.http.Get(path, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) GetProjects() string {
	path := "/v9/projects"
	result := "[]"
	err := v.http.Get(path, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
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

type AddDomainToProjectResponse struct {
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
