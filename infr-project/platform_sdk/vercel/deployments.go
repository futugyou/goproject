package vercel

import (
	"fmt"
	"net/url"
)

func (v *VercelClient) CancelDeployment(id string, slug string, teamId string) (*DeploymentInfo, error) {
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
	err := v.http.Patch(path, nil, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type DeploymentInfo struct {
	Id                            string        `json:"id,omitempty"`
	Meta                          interface{}   `json:"meta,omitempty"`
	Url                           string        `json:"url,omitempty"`
	Alias                         []string      `json:"alias,omitempty"`
	AliasAssigned                 bool          `json:"aliasAssigned,omitempty"`
	AliasAssignedAt               string        `json:"aliasAssignedAt,omitempty"`
	AliasError                    *VercelError  `json:"aliasError,omitempty"`
	AliasFinal                    string        `json:"aliasFinal,omitempty"`
	AliasWarning                  AliasWarning  `json:"aliasWarning,omitempty"`
	AlwaysRefuseToBuild           bool          `json:"alwaysRefuseToBuild,omitempty"`
	AutoAssignCustomDomains       bool          `json:"autoAssignCustomDomains,omitempty"`
	AutomaticAliases              []string      `json:"automaticAliases,omitempty"`
	BootedAt                      int           `json:"bootedAt,omitempty"`
	Build                         Build         `json:"build,omitempty"`
	BuildErrorAt                  int           `json:"buildErrorAt,omitempty"`
	BuildSkipped                  bool          `json:"buildSkipped,omitempty"`
	BuildingAt                    int           `json:"buildingAt,omitempty"`
	Builds                        []string      `json:"builds,omitempty"`
	CanceledAt                    int           `json:"canceledAt,omitempty"`
	ChecksConclusion              string        `json:"checksConclusion,omitempty"`
	ChecksState                   string        `json:"checksState,omitempty"`
	ConnectBuildsEnabled          bool          `json:"connectBuildsEnabled,omitempty"`
	ConnectConfigurationId        string        `json:"connectConfigurationId,omitempty"`
	CreatedAt                     int           `json:"createdAt,omitempty"`
	CreatedIn                     string        `json:"createdIn,omitempty"`
	Creator                       CreatorInfo   `json:"creator,omitempty"`
	Crons                         []Cron        `json:"crons,omitempty"`
	CustomEnvironment             interface{}   `json:"customEnvironment,omitempty"`
	DeletedAt                     int           `json:"deletedAt,omitempty"`
	Env                           []string      `json:"env,omitempty"`
	ErrorCode                     string        `json:"errorCode,omitempty"`
	ErrorLink                     string        `json:"errorLink,omitempty"`
	ErrorMessage                  string        `json:"errorMessage,omitempty"`
	ErrorStep                     string        `json:"errorStep,omitempty"`
	Flags                         interface{}   `json:"flags,omitempty"`
	Functions                     string        `json:"functions,omitempty"`
	GitRepo                       interface{}   `json:"gitRepo,omitempty"`
	GitSource                     interface{}   `json:"gitSource,omitempty"`
	InitReadyAt                   int           `json:"initReadyAt,omitempty"`
	InspectorUrl                  string        `json:"inspectorUrl,omitempty"`
	IsFirstBranchDeployment       bool          `json:"isFirstBranchDeployment,omitempty"`
	IsInConcurrentBuildsQueue     bool          `json:"isInConcurrentBuildsQueue,omitempty"`
	Lambdas                       []interface{} `json:"lambdas,omitempty"`
	MonorepoManager               string        `json:"monorepoManager,omitempty"`
	Name                          string        `json:"name,omitempty"`
	OidcTokenClaims               interface{}   `json:"oidcTokenClaims,omitempty"`
	OwnerId                       string        `json:"ownerId,omitempty"`
	PassiveConnectConfigurationId string        `json:"passiveConnectConfigurationId,omitempty"`
	PassiveRegions                []string      `json:"passiveRegions,omitempty"`
	Plan                          string        `json:"plan,omitempty"`
	PreviewCommentsEnabled        string        `json:"previewCommentsEnabled,omitempty"`
	Project                       ProjectInfo   `json:"project,omitempty"`
	ProjectId                     string        `json:"projectId,omitempty"`
	ProjectSettings               interface{}   `json:"projectSettings,omitempty"`
	Public                        bool          `json:"public,omitempty"`
	Ready                         int           `json:"ready,omitempty"`
	ReadyState                    string        `json:"readyState,omitempty"`
	ReadyStateReason              string        `json:"readyStateReason,omitempty"`
	ReadySubstate                 string        `json:"readySubstate,omitempty"`
	Regions                       []string      `json:"regions,omitempty"`
	Routes                        interface{}   `json:"routes,omitempty"`
	Source                        string        `json:"source,omitempty"`
	Target                        string        `json:"target,omitempty"`
	Team                          Team          `json:"team,omitempty"`
	TtyBuildLogs                  bool          `json:"ttyBuildLogs,omitempty"`
	Type                          string        `json:"type,omitempty"`
	UndeletedAt                   int           `json:"undeletedAt,omitempty"`
	UserAliases                   []string      `json:"userAliases,omitempty"`
	Version                       int           `json:"version,omitempty"`
	Error                         *VercelError  `json:"error,omitempty"`
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
