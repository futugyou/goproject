package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type DeploymentService service

type CancelDeploymentRequest struct {
	DeploymentId     string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *DeploymentService) CancelDeployment(ctx context.Context, request CancelDeploymentRequest) (*DeploymentInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v12/deployments/%s/cancel", request.DeploymentId),
	}
	params := url.Values{}
	if request.TeamId != nil {
		params.Add("teamId", *request.TeamId)
	}
	if request.TeamSlug != nil {
		params.Add("slug", *request.TeamSlug)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &DeploymentInfo{}
	if err := v.client.http.Patch(ctx, path, nil, response); err != nil {
		return nil, err
	}
	return response, nil
}

type CreateDeploymentRequest struct {
	Name                          string          `json:"name,omitempty"`
	DeploymentID                  string          `json:"deploymentId,omitempty"`
	Files                         []File          `json:"files,omitempty"`
	GitMetadata                   GitMetadata     `json:"gitMetadata,omitempty"`
	GitSource                     GitSource       `json:"gitSource,omitempty"`
	Meta                          string          `json:"meta,omitempty"`
	MonorepoManager               string          `json:"monorepoManager,omitempty"`
	Project                       string          `json:"project,omitempty"`
	ProjectSettings               ProjectSettings `json:"projectSettings,omitempty"`
	Target                        string          `json:"target,omitempty"`
	WithLatestCommit              bool            `json:"withLatestCommit,omitempty"`
	ForceNew                      *string         `json:"-"`
	SkipAutoDetectionConfirmation *string         `json:"-"`
	BaseUrlParameter              `json:"-"`
}

func (v *DeploymentService) CreateDeployment(ctx context.Context, request CreateDeploymentRequest) (*DeploymentInfo, error) {
	u := &url.URL{
		Path: "/v13/deployments",
	}
	params := url.Values{}
	if request.TeamId != nil {
		params.Add("teamId", *request.TeamId)
	}
	if request.TeamSlug != nil {
		params.Add("slug", *request.TeamSlug)
	}
	if request.SkipAutoDetectionConfirmation != nil {
		params.Add("skipAutoDetectionConfirmation", *request.SkipAutoDetectionConfirmation)
	}
	if request.ForceNew != nil {
		params.Add("forceNew", *request.ForceNew)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &DeploymentInfo{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type DeleteDeploymentRequest struct {
	DeploymentID     string  `json:"deploymentId"`
	Url              *string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *DeploymentService) DeleteDeployment(ctx context.Context, request DeleteDeploymentRequest) (*DeleteDeploymentResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v13/deployments/%s", request.DeploymentID),
	}
	params := url.Values{}
	if request.TeamId != nil {
		params.Add("teamId", *request.TeamId)
	}
	if request.TeamSlug != nil {
		params.Add("slug", *request.TeamSlug)
	}
	if request.Url != nil {
		params.Add("url", *request.Url)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &DeleteDeploymentResponse{}
	if err := v.client.http.Delete(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type GetDeploymentParameter struct {
	IdOrUrl          string
	WithGitRepoInfo  *string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *DeploymentService) GetDeployment(ctx context.Context, request GetDeploymentParameter) (*DeploymentInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v13/deployments/%s", request.IdOrUrl),
	}
	params := url.Values{}
	if request.TeamId != nil {
		params.Add("teamId", *request.TeamId)
	}
	if request.TeamSlug != nil {
		params.Add("slug", *request.TeamSlug)
	}
	if request.WithGitRepoInfo != nil {
		params.Add("withGitRepoInfo", *request.WithGitRepoInfo)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &DeploymentInfo{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type GetDeploymentEventParameter struct {
	IdOrUrl          string
	Builds           *string // Allowed value: 1
	Delimiter        *string // Allowed value: 1
	Direction        *string // Allowed value: backward forward
	Follow           *string // Allowed value: 1
	Limit            *string
	Name             *string
	Since            *string
	StatusCode       *string
	Until            *string
	BaseUrlParameter `json:"-"`
}

func (v *DeploymentService) GetDeploymentEvent(ctx context.Context, request GetDeploymentEventParameter) ([]DeploymentEvent, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v13/deployments/%s/events", request.IdOrUrl),
	}
	params := url.Values{}
	if request.TeamId != nil {
		params.Add("teamId", *request.TeamId)
	}
	if request.TeamSlug != nil {
		params.Add("slug", *request.TeamSlug)
	}
	if request.Builds != nil {
		params.Add("builds", *request.Builds)
	}
	if request.Delimiter != nil {
		params.Add("delimiter", *request.Delimiter)
	}
	if request.Direction != nil {
		params.Add("direction", *request.Direction)
	}
	if request.Follow != nil {
		params.Add("follow", *request.Follow)
	}
	if request.Until != nil {
		params.Add("until", *request.Until)
	}
	if request.Limit != nil {
		params.Add("limit", *request.Limit)
	}
	if request.Name != nil {
		params.Add("name", *request.Name)
	}
	if request.Since != nil {
		params.Add("since", *request.Since)
	}
	if request.StatusCode != nil {
		params.Add("statusCode", *request.StatusCode)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := []DeploymentEvent{}
	if err := v.client.http.Get(ctx, path, &response); err != nil {
		return nil, err
	}
	return response, nil
}

type GetDeploymentFileParameter struct {
	DeploymentID     string
	FileId           string
	Path             *string
	BaseUrlParameter `json:"-"`
}

func (v *DeploymentService) GetDeploymentFile(ctx context.Context, request GetDeploymentFileParameter) (*FileTree, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v7/deployments/%s/files/%s", request.DeploymentID, request.FileId),
	}
	params := url.Values{}
	if request.TeamId != nil {
		params.Add("teamId", *request.TeamId)
	}
	if request.TeamSlug != nil {
		params.Add("slug", *request.TeamSlug)
	}
	if request.Path != nil {
		params.Add("path", *request.Path)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &FileTree{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type ListDeploymentParameter struct {
	App               *string
	Limit             *string
	ProjectId         *string
	RollbackCandidate *string // Allowed value: true false
	Since             *string
	State             *string
	Target            *string
	Until             *string
	Users             *string
	BaseUrlParameter  `json:"-"`
}

func (v *DeploymentService) ListDeployment(ctx context.Context, request ListDeploymentParameter) (*ListDeploymentResponse, error) {
	u := &url.URL{
		Path: "/v6/deployments",
	}
	params := url.Values{}
	if request.TeamId != nil {
		params.Add("teamId", *request.TeamId)
	}
	if request.TeamSlug != nil {
		params.Add("slug", *request.TeamSlug)
	}
	if request.App != nil {
		params.Add("app", *request.App)
	}
	if request.Until != nil {
		params.Add("until", *request.Until)
	}
	if request.Limit != nil {
		params.Add("limit", *request.Limit)
	}
	if request.ProjectId != nil {
		params.Add("projectId", *request.ProjectId)
	}
	if request.RollbackCandidate != nil {
		params.Add("rollbackCandidate", *request.RollbackCandidate)
	}
	if request.Since != nil {
		params.Add("since", *request.Since)
	}
	if request.State != nil {
		params.Add("state", *request.State)
	}
	if request.Target != nil {
		params.Add("target", *request.Target)
	}
	if request.Users != nil {
		params.Add("users", *request.Users)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ListDeploymentResponse{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type ListDeploymentFileParameter struct {
	DeploymentID     string
	BaseUrlParameter `json:"-"`
}

func (v *DeploymentService) ListDeploymentFile(ctx context.Context, request ListDeploymentFileParameter) ([]FileTree, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v6/deployments/%s/files", request.DeploymentID),
	}
	params := url.Values{}
	if request.TeamId != nil {
		params.Add("teamId", *request.TeamId)
	}
	if request.TeamSlug != nil {
		params.Add("slug", *request.TeamSlug)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := []FileTree{}
	if err := v.client.http.Get(ctx, path, &response); err != nil {
		return nil, err
	}
	return response, nil
}

// TODO: need common http file upload method first.
func (v *DeploymentService) UploadDeploymentFiles(ctx context.Context, request interface{}) ([]interface{}, error) {
	return nil, nil
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

type File struct {
	InlinedFile InlinedFile `json:"InlinedFile,omitempty"`
}

type InlinedFile struct {
	Data     string `json:"data,omitempty"`
	Encoding string `json:"encoding,omitempty"`
	File     string `json:"file,omitempty"`
	Sha      string `json:"sha,omitempty"`
	Size     int    `json:"size,omitempty"`
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
	Ref           string      `json:"ref,omitempty"`
	RepoID        string      `json:"repoId,omitempty"`
	Repo          string      `json:"repo,omitempty"`
	SHA           string      `json:"sha,omitempty"`
	Type          string      `json:"type,omitempty"`
	Org           string      `json:"org,omitempty"`
	ProjectId     interface{} `json:"projectId,omitempty"`
	RepoUuid      string      `json:"repoUuid,omitempty"`
	WorkspaceUuid string      `json:"workspaceUuid,omitempty"`
	Owner         string      `json:"owner,omitempty"`
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
