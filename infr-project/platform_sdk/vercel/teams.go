package vercel

import (
	"fmt"
	"net/url"
)

func (v *VercelClient) CreateTeam(req CreateTeamRequest) (*TeamInfo, error) {
	path := "/v1/teams"
	result := &TeamInfo{}
	if err := v.http.Post(path, req, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) DeleteTeam(teamId string, slug string, newDefaultTeamId string) (*DeleteTeamResponse, error) {
	path := fmt.Sprintf("/v1/teams/%s", teamId)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(newDefaultTeamId) > 0 {
		queryParams.Add("newDefaultTeamId", newDefaultTeamId)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &DeleteTeamResponse{}
	if err := v.http.Delete(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (v *VercelClient) DeleteTeamInviteCode(teamId string, inviteId string) (*DeleteTeamInviteCodeResponse, error) {
	path := fmt.Sprintf("/v1/teams/%s/invites/%s", teamId, inviteId)

	result := &DeleteTeamInviteCodeResponse{}
	if err := v.http.Delete(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (v *VercelClient) GetTeam(teamId string, slug string) (*TeamInfo, error) {
	path := fmt.Sprintf("/v2/teams/%s", teamId)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &TeamInfo{}
	if err := v.http.Get(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (v *VercelClient) GetAccessRequestStatus(teamId string, userId string) (*AccessRequestStatus, error) {
	path := fmt.Sprintf("/v1/teams/%s/request/%s", teamId, userId)
	result := &AccessRequestStatus{}
	if err := v.http.Get(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (v *VercelClient) ListTeamMembers(teamId string, eligibleMembersForProjectId string, excludeProject string,
	limit string, role string, search string, since string, until string) (*ListTeamMembersResponse, error) {
	path := fmt.Sprintf("/v2/teams/%s/members", teamId)
	queryParams := url.Values{}
	if len(eligibleMembersForProjectId) > 0 {
		queryParams.Add("eligibleMembersForProjectId", eligibleMembersForProjectId)
	}
	if len(excludeProject) > 0 {
		queryParams.Add("excludeProject", excludeProject)
	}
	if len(limit) > 0 {
		queryParams.Add("limit", limit)
	}
	if len(role) > 0 {
		queryParams.Add("role", role)
	}
	if len(search) > 0 {
		queryParams.Add("search", search)
	}
	if len(since) > 0 {
		queryParams.Add("since", since)
	}
	if len(until) > 0 {
		queryParams.Add("until", until)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}

	result := &ListTeamMembersResponse{}
	if err := v.http.Get(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (v *VercelClient) ListTeam(limit string, since string, until string) (*ListTeamResponse, error) {
	path := "/v2/teams"
	queryParams := url.Values{}
	if len(limit) > 0 {
		queryParams.Add("limit", limit)
	}
	if len(since) > 0 {
		queryParams.Add("since", since)
	}
	if len(until) > 0 {
		queryParams.Add("until", until)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &ListTeamResponse{}
	if err := v.http.Get(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (v *VercelClient) InviteUser(teamId string, req InviteUserRequest) (*InviteUserResponse, error) {
	path := fmt.Sprintf("/v1/teams/%s/members", teamId)
	result := &InviteUserResponse{}
	if err := v.http.Post(path, req, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (v *VercelClient) JoinTeam(teamId string, inviteCode string) (*JoinTeamResponse, error) {
	path := fmt.Sprintf("/v1/teams/%s/members/teams/join", teamId)
	req := struct {
		InviteCode string `json:"inviteCode,omitempty"`
	}{
		InviteCode: inviteCode,
	}
	result := &JoinTeamResponse{}
	if err := v.http.Post(path, req, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (v *VercelClient) UpdateTeam(teamId string, slug string, req UpdateTeamInfo) (*TeamInfo, error) {
	path := fmt.Sprintf("/v2/teams/%s", teamId)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &TeamInfo{}
	if err := v.http.Patch(path, req, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (v *VercelClient) RemoveTeamMember(teamId string, uid string, newDefaultTeamId string) (*RemoveTeamMemberResponse, error) {
	path := fmt.Sprintf("/v1/teams/%s/members/%s", teamId, uid)
	queryParams := url.Values{}
	if len(newDefaultTeamId) > 0 {
		queryParams.Add("newDefaultTeamId", newDefaultTeamId)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &RemoveTeamMemberResponse{}
	if err := v.http.Delete(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (v *VercelClient) RequestAccessToTeam(teamId string, req JoinedFrom) (*AccessRequestStatus, error) {
	path := fmt.Sprintf("/v1/teams/%s/request", teamId)
	result := &AccessRequestStatus{}
	if err := v.http.Post(path, req, result); err != nil {
		return nil, err
	}

	return result, nil
}

type RemoveTeamMemberResponse struct {
	Id    string       `json:"id,omitempty"`
	Error *VercelError `json:"error,omitempty"`
}

type UpdateTeamInfo struct {
	Avatar                             string        `json:"avatar,omitempty"`
	Description                        string        `json:"description,omitempty"`
	EmailDomain                        string        `json:"emailDomain,omitempty"`
	EnablePreviewFeedback              string        `json:"enablePreviewFeedback,omitempty"`
	HideIPAddresses                    bool          `json:"hideIpAddresses,omitempty"`
	Name                               string        `json:"name,omitempty"`
	PreviewDeploymentSuffix            string        `json:"previewDeploymentSuffix,omitempty"`
	RegenerateInviteCode               bool          `json:"regenerateInviteCode,omitempty"`
	RemoteCaching                      RemoteCaching `json:"remoteCaching,omitempty"`
	Saml                               Saml          `json:"saml,omitempty"`
	SensitiveEnvironmentVariablePolicy string        `json:"sensitiveEnvironmentVariablePolicy,omitempty"`
	Slug                               string        `json:"slug,omitempty"`
}

type RemoteCaching struct {
	Enabled bool `json:"enabled,omitempty"`
}

type Saml struct {
	Enforced bool   `json:"enforced,omitempty"`
	Roles    string `json:"roles,omitempty"`
}

type JoinTeamResponse struct {
	From   string       `json:"from,omitempty"`
	Name   string       `json:"name,omitempty"`
	Slug   string       `json:"slug,omitempty"`
	TeamId string       `json:"teamId,omitempty"`
	Error  *VercelError `json:"error,omitempty"`
}

type InviteUserRequest struct {
	Email    string          `json:"email,omitempty"`
	Role     string          `json:"role,omitempty"`
	Uid      string          `json:"uid,omitempty"`
	Projects []InviteProject `json:"projects,omitempty"`
}

type InviteProject struct {
	ProjectId string `json:"projectId,omitempty"`
	Role      string `json:"role,omitempty"`
}

type InviteUserResponse struct {
	Email    string       `json:"email,omitempty"`
	Role     string       `json:"role,omitempty"`
	Uid      string       `json:"uid,omitempty"`
	Username string       `json:"username,omitempty"`
	Error    *VercelError `json:"error,omitempty"`
}

type ListTeamResponse struct {
	Teams      []TeamInfo   `json:"teams,omitempty"`
	Pagination Pagination   `json:"pagination,omitempty"`
	Error      *VercelError `json:"error,omitempty"`
}

type ListTeamMembersResponse struct {
	EmailInviteCodes []EmailInviteCode `json:"emailInviteCodes,omitempty"`
	Members          []MemberInfo      `json:"members,omitempty"`
	Pagination       Pagination        `json:"pagination,omitempty"`
	Error            *VercelError      `json:"error,omitempty"`
}

type AccessRequestStatus struct {
	AccessRequestedAt int          `json:"accessRequestedAt,omitempty"`
	Bitbucket         RepoState    `json:"bitbucket,omitempty"`
	Confirmed         bool         `json:"confirmed,omitempty"`
	Github            RepoState    `json:"github,omitempty"`
	Gitlab            RepoState    `json:"gitlab,omitempty"`
	TeamName          string       `json:"teamName,omitempty"`
	TeamSlug          string       `json:"teamSlug,omitempty"`
	JoinedFrom        JoinedFrom   `json:"joinedFrom,omitempty"`
	Error             *VercelError `json:"error,omitempty"`
}

type DeleteTeamInviteCodeResponse struct {
	Id    string       `json:"id,omitempty"`
	Error *VercelError `json:"error,omitempty"`
}

type DeleteTeamResponse struct {
	Id                    string       `json:"id,omitempty"`
	NewDefaultTeamIdError bool         `json:"newDefaultTeamIdError,omitempty"`
	Error                 *VercelError `json:"error,omitempty"`
}

type CreateTeamRequest struct {
	Slug        string      `json:"slug,omitempty"`
	Attribution Attribution `json:"attribution,omitempty"`
	Name        string      `json:"name,omitempty"`
}

type TeamInfo struct {
	Billing interface{}  `json:"billing,omitempty"`
	Id      string       `json:"id,omitempty"`
	Slug    string       `json:"slug,omitempty"`
	Name    string       `json:"name,omitempty"`
	Avatar  string       `json:"avatar,omitempty"`
	Error   *VercelError `json:"error,omitempty"`
}

type Attribution struct {
	LandingPage              string `json:"landingPage,omitempty"`
	PageBeforeConversionPage string `json:"pageBeforeConversionPage,omitempty"`
	SessionReferrer          string `json:"sessionReferrer,omitempty"`
	Utm                      Utm    `json:"utm,omitempty"`
}

type Utm struct {
	UtmCampaign string `json:"utmCampaign,omitempty"`
	UtmMedium   string `json:"utmMedium,omitempty"`
	UtmSource   string `json:"utmSource,omitempty"`
	UtmTerm     string `json:"utmTerm,omitempty"`
}

type RepoState struct {
	Login string `json:"login,omitempty"`
}

type JoinedFrom struct {
	CommitId         string `json:"commitId,omitempty"`
	DsyncConnectedAt int    `json:"dsyncConnectedAt,omitempty"`
	DsyncUserId      string `json:"dsyncUserId,omitempty"`
	GitUserId        string `json:"gitUserId,omitempty"`
	GitUserLogin     string `json:"gitUserLogin,omitempty"`
	IdpUserId        string `json:"idpUserId,omitempty"`
	Origin           string `json:"origin,omitempty"`
	RepoId           string `json:"repoId,omitempty"`
	RepoPath         string `json:"repoPath,omitempty"`
	SsoConnectedAt   int    `json:"ssoConnectedAt,omitempty"`
	SsoUserId        string `json:"ssoUserId,omitempty"`
}

type EmailInviteCode struct {
	AccessGroups []string    `json:"accessGroups,omitempty"`
	CreatedAt    int         `json:"createdAt,omitempty"`
	Email        string      `json:"email,omitempty"`
	Expired      bool        `json:"expired,omitempty"`
	Id           string      `json:"id,omitempty"`
	IsDSyncUser  bool        `json:"isDSyncUser,omitempty"`
	Role         string      `json:"role,omitempty"`
	Projects     interface{} `json:"projects,omitempty"`
}
