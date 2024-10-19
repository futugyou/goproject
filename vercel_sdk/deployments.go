package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type DeploymentService service

func (v *DeploymentService) CancelDeployment(ctx context.Context, id string, slug string, teamId string) (*DeploymentInfo, error) {
	path := fmt.Sprintf("/v12/deployments/%s/cancel", id)
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
	result := &DeploymentInfo{}
	err := v.client.http.Patch(ctx, path, nil, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *DeploymentService) CreateDeployment(ctx context.Context, forceNew string, skipAutoDetectionConfirmation string,
	slug string, teamId string, req CreateDeploymentRequest) (*DeploymentInfo, error) {
	path := "/v13/deployments"
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(forceNew) > 0 {
		queryParams.Add("forceNew", forceNew)
	}
	if len(skipAutoDetectionConfirmation) > 0 {
		queryParams.Add("skipAutoDetectionConfirmation", skipAutoDetectionConfirmation)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &DeploymentInfo{}
	err := v.client.http.Post(ctx, path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *DeploymentService) DeleteDeployment(ctx context.Context, id string, deploymentUrl string, slug string, teamId string) (*DeleteDeploymentResponse, error) {
	path := fmt.Sprintf("/v13/deployments/%s", id)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(deploymentUrl) > 0 {
		queryParams.Add("url", deploymentUrl)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &DeleteDeploymentResponse{}
	err := v.client.http.Delete(ctx, path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *DeploymentService) GetDeployment(ctx context.Context, idOrUrl string, slug string, teamId string, withGitRepoInfo string) (*DeploymentInfo, error) {
	path := fmt.Sprintf("/v13/deployments/%s", idOrUrl)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(withGitRepoInfo) > 0 {
		queryParams.Add("withGitRepoInfo", withGitRepoInfo)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &DeploymentInfo{}
	err := v.client.http.Get(ctx, path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *DeploymentService) GetDeploymentEvent(ctx context.Context, idOrUrl string, slug string, teamId string,
	builds string, delimiter string, direction string, follow string, limit string, name string,
	since string, statusCode string, until string) ([]DeploymentEvent, error) {
	path := fmt.Sprintf("/v3/deployments/%s/events", idOrUrl)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(builds) > 0 {
		queryParams.Add("builds", builds)
	}
	if len(delimiter) > 0 {
		queryParams.Add("delimiter", delimiter)
	}
	if len(direction) > 0 {
		queryParams.Add("direction", direction)
	}
	if len(follow) > 0 {
		queryParams.Add("follow", follow)
	}
	if len(until) > 0 {
		queryParams.Add("until", until)
	}
	if len(limit) > 0 {
		queryParams.Add("limit", limit)
	}
	if len(name) > 0 {
		queryParams.Add("name", name)
	}
	if len(since) > 0 {
		queryParams.Add("since", since)
	}
	if len(statusCode) > 0 {
		queryParams.Add("statusCode", statusCode)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := []DeploymentEvent{}
	err := v.client.http.Get(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *DeploymentService) GetDeploymentFile(ctx context.Context, id string, fileId string, slug string, teamId string, filePath string) (*FileTree, error) {
	path := fmt.Sprintf("/v7/deployments/%s/files/%s", id, fileId)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(filePath) > 0 {
		queryParams.Add("path", filePath)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &FileTree{}
	err := v.client.http.Get(ctx, path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *DeploymentService) ListDeployment(ctx context.Context, app string, until string, slug string, teamId string, limit string,
	projectId string, rollbackCandidate string, since string, state string, target string, users string) (*ListDeploymentResponse, error) {
	path := "/v6/deployments"
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(app) > 0 {
		queryParams.Add("app", app)
	}
	if len(until) > 0 {
		queryParams.Add("until", until)
	}
	if len(limit) > 0 {
		queryParams.Add("limit", limit)
	}
	if len(projectId) > 0 {
		queryParams.Add("projectId", projectId)
	}
	if len(rollbackCandidate) > 0 {
		queryParams.Add("rollbackCandidate", rollbackCandidate)
	}
	if len(since) > 0 {
		queryParams.Add("since", since)
	}
	if len(state) > 0 {
		queryParams.Add("state", state)
	}
	if len(target) > 0 {
		queryParams.Add("target", target)
	}
	if len(users) > 0 {
		queryParams.Add("users", users)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &ListDeploymentResponse{}
	err := v.client.http.Get(ctx, path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *DeploymentService) ListDeploymentFile(ctx context.Context, id string, slug string, teamId string) ([]FileTree, error) {
	path := fmt.Sprintf("/v6/deployments/%s/files", id)
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
	result := []FileTree{}
	err := v.client.http.Get(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type ListDeploymentResponse struct {
	Deployments []DeploymentInfo `json:"deployments,omitempty"`
	Pagination  Pagination       `json:"pagination,omitempty"`
	Error       *VercelError     `json:"error,omitempty"`
}

type DeleteDeploymentResponse struct {
	State string       `json:"state,omitempty"`
	Uid   string       `json:"uid,omitempty"`
	Error *VercelError `json:"error,omitempty"`
}

type DeploymentInfo struct {
	Id                            string          `json:"id,omitempty"`
	Meta                          interface{}     `json:"meta,omitempty"`
	Url                           string          `json:"url,omitempty"`
	Alias                         []string        `json:"alias,omitempty"`
	AliasAssigned                 bool            `json:"aliasAssigned,omitempty"`
	AliasAssignedAt               string          `json:"aliasAssignedAt,omitempty"`
	AliasError                    *VercelError    `json:"aliasError,omitempty"`
	AliasFinal                    string          `json:"aliasFinal,omitempty"`
	AliasWarning                  AliasWarning    `json:"aliasWarning,omitempty"`
	AlwaysRefuseToBuild           bool            `json:"alwaysRefuseToBuild,omitempty"`
	AutoAssignCustomDomains       bool            `json:"autoAssignCustomDomains,omitempty"`
	AutomaticAliases              []string        `json:"automaticAliases,omitempty"`
	BootedAt                      int             `json:"bootedAt,omitempty"`
	Build                         Build           `json:"build,omitempty"`
	BuildErrorAt                  int             `json:"buildErrorAt,omitempty"`
	BuildSkipped                  bool            `json:"buildSkipped,omitempty"`
	BuildingAt                    int             `json:"buildingAt,omitempty"`
	Builds                        []string        `json:"builds,omitempty"`
	CanceledAt                    int             `json:"canceledAt,omitempty"`
	ChecksConclusion              string          `json:"checksConclusion,omitempty"`
	ChecksState                   string          `json:"checksState,omitempty"`
	ConnectBuildsEnabled          bool            `json:"connectBuildsEnabled,omitempty"`
	ConnectConfigurationId        string          `json:"connectConfigurationId,omitempty"`
	CreatedAt                     int             `json:"createdAt,omitempty"`
	CreatedIn                     string          `json:"createdIn,omitempty"`
	Creator                       CreatorInfo     `json:"creator,omitempty"`
	Crons                         []Cron          `json:"crons,omitempty"`
	CustomEnvironment             interface{}     `json:"customEnvironment,omitempty"`
	DeletedAt                     int             `json:"deletedAt,omitempty"`
	Env                           []string        `json:"env,omitempty"`
	ErrorCode                     string          `json:"errorCode,omitempty"`
	ErrorLink                     string          `json:"errorLink,omitempty"`
	ErrorMessage                  string          `json:"errorMessage,omitempty"`
	ErrorStep                     string          `json:"errorStep,omitempty"`
	Flags                         interface{}     `json:"flags,omitempty"`
	Functions                     string          `json:"functions,omitempty"`
	GitRepo                       interface{}     `json:"gitRepo,omitempty"`
	GitSource                     GitSource       `json:"gitSource,omitempty"`
	InitReadyAt                   int             `json:"initReadyAt,omitempty"`
	InspectorUrl                  string          `json:"inspectorUrl,omitempty"`
	IsFirstBranchDeployment       bool            `json:"isFirstBranchDeployment,omitempty"`
	IsInConcurrentBuildsQueue     bool            `json:"isInConcurrentBuildsQueue,omitempty"`
	Lambdas                       []interface{}   `json:"lambdas,omitempty"`
	MonorepoManager               string          `json:"monorepoManager,omitempty"`
	Name                          string          `json:"name,omitempty"`
	OidcTokenClaims               interface{}     `json:"oidcTokenClaims,omitempty"`
	OwnerId                       string          `json:"ownerId,omitempty"`
	PassiveConnectConfigurationId string          `json:"passiveConnectConfigurationId,omitempty"`
	PassiveRegions                []string        `json:"passiveRegions,omitempty"`
	Plan                          string          `json:"plan,omitempty"`
	PreviewCommentsEnabled        string          `json:"previewCommentsEnabled,omitempty"`
	Project                       ProjectInfo     `json:"project,omitempty"`
	ProjectId                     string          `json:"projectId,omitempty"`
	ProjectSettings               ProjectSettings `json:"projectSettings,omitempty"`
	Public                        bool            `json:"public,omitempty"`
	Ready                         int             `json:"ready,omitempty"`
	ReadyState                    string          `json:"readyState,omitempty"`
	ReadyStateReason              string          `json:"readyStateReason,omitempty"`
	ReadySubstate                 string          `json:"readySubstate,omitempty"`
	Regions                       []string        `json:"regions,omitempty"`
	Routes                        interface{}     `json:"routes,omitempty"`
	Source                        string          `json:"source,omitempty"`
	Target                        string          `json:"target,omitempty"`
	Team                          Team            `json:"team,omitempty"`
	TtyBuildLogs                  bool            `json:"ttyBuildLogs,omitempty"`
	Type                          string          `json:"type,omitempty"`
	UndeletedAt                   int             `json:"undeletedAt,omitempty"`
	UserAliases                   []string        `json:"userAliases,omitempty"`
	Version                       int             `json:"version,omitempty"`
	Error                         *VercelError    `json:"error,omitempty"`
}

type Build struct {
	Env []string `json:"env,omitempty"`
}

type Cron struct {
	Path     string `json:"path,omitempty"`
	Schedule string `json:"schedule,omitempty"`
}

type Team struct {
	Avatar string `json:"avatar,omitempty"`
	Id     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Slug   string `json:"slug,omitempty"`
}

type CreateDeploymentRequest struct {
	Name                      string          `json:"name,omitempty"`
	CustomEnvironmentSlugOrID string          `json:"customEnvironmentSlugOrId,omitempty"`
	DeploymentID              string          `json:"deploymentId,omitempty"`
	Files                     []File          `json:"files,omitempty"`
	GitMetadata               GitMetadata     `json:"gitMetadata,omitempty"`
	GitSource                 GitSource       `json:"gitSource,omitempty"`
	Meta                      string          `json:"meta,omitempty"`
	MonorepoManager           string          `json:"monorepoManager,omitempty"`
	Project                   string          `json:"project,omitempty"`
	ProjectSettings           ProjectSettings `json:"projectSettings,omitempty"`
	Target                    string          `json:"target,omitempty"`
	WithLatestCommit          bool            `json:"withLatestCommit,omitempty"`
}

type File struct {
	InlinedFile InlinedFile `json:"InlinedFile,omitempty"`
}

type InlinedFile struct {
	Data     string `json:"data,omitempty"`
	Encoding string `json:"encoding,omitempty"`
	File     string `json:"file,omitempty"`
}

type GitMetadata struct {
	RemoteURL        string `json:"remoteUrl,omitempty"`
	CommitAuthorName string `json:"commitAuthorName,omitempty"`
	CommitMessage    string `json:"commitMessage,omitempty"`
	CommitRef        string `json:"commitRef,omitempty"`
	CommitSHA        string `json:"commitSha,omitempty"`
	Dirty            bool   `json:"dirty,omitempty"`
}

type GitSource struct {
	Ref    string `json:"ref,omitempty"`
	RepoID string `json:"repoId,omitempty"`
	SHA    string `json:"sha,omitempty"`
	Type   string `json:"type,omitempty"`
}

type ProjectSettings struct {
	BuildCommand                    string `json:"buildCommand,omitempty"`
	CommandForIgnoringBuildStep     string `json:"commandForIgnoringBuildStep,omitempty"`
	DevCommand                      string `json:"devCommand,omitempty"`
	Framework                       string `json:"framework,omitempty"`
	InstallCommand                  string `json:"installCommand,omitempty"`
	NodeVersion                     string `json:"nodeVersion,omitempty"`
	OutputDirectory                 string `json:"outputDirectory,omitempty"`
	RootDirectory                   string `json:"rootDirectory,omitempty"`
	ServerlessFunctionRegion        string `json:"serverlessFunctionRegion,omitempty"`
	SkipGitConnectDuringLink        bool   `json:"skipGitConnectDuringLink,omitempty"`
	SourceFilesOutsideRootDirectory bool   `json:"sourceFilesOutsideRootDirectory,omitempty"`
}

type DeploymentEvent struct {
	Created      int         `json:"created,omitempty"`
	Payload      interface{} `json:"payload,omitempty"`
	Type         string      `json:"type,omitempty"`
	Date         int         `json:"date,omitempty"`
	DeploymentId string      `json:"deploymentId,omitempty"`
	Id           string      `json:"id,omitempty"`
	Info         interface{} `json:"info,omitempty"`
	Proxy        interface{} `json:"proxy,omitempty"`
	RequestId    string      `json:"requestId,omitempty"`
	Serial       string      `json:"serial,omitempty"`
	Text         string      `json:"text,omitempty"`
	StatusCode   int         `json:"statusCode,omitempty"`
}

type FileTree struct {
	Name        string       `json:"name,omitempty"`
	Type        string       `json:"type,omitempty"`
	Mode        string       `json:"mode,omitempty"`
	Uid         string       `json:"uid,omitempty"`
	ContentType string       `json:"contentType,omitempty"`
	Children    []FileTree   `json:"children,omitempty"`
	Symlink     string       `json:"symlink,omitempty"`
	Error       *VercelError `json:"error,omitempty"`
}
