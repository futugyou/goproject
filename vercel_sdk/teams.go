package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type TeamService service

func (v *TeamService) CreateTeam(ctx context.Context, req CreateTeamRequest) (*TeamInfo, error) {
	path := "/v1/teams"
	response := &TeamInfo{}
	if err := v.client.http.Post(ctx, path, req, response); err != nil {
		return nil, err
	}
	return response, nil
}

type DeleteTeamRequest struct {
	TeamId           string         `json:"-"`
	NewDefaultTeamId *string        `json:"-"`
	Slug             *string        `json:"-"`
	Reasons          []DeleteReason `json:"reasons"`
}

type DeleteReason struct {
	Description string `json:"description"`
	Slug        string `json:"slug"`
}

func (v *TeamService) DeleteTeam(ctx context.Context, request DeleteTeamRequest) (*DeleteTeamResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/teams/%s", request.TeamId),
	}
	params := url.Values{}
	if request.Slug != nil {
		params.Add("slug", *request.Slug)
	}
	if request.NewDefaultTeamId != nil {
		params.Add("newDefaultTeamId", *request.NewDefaultTeamId)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &DeleteTeamResponse{}
	if err := v.client.http.DoRequest(ctx, path, "DELETE", request, response); err != nil {
		return nil, err
	}

	return response, nil
}

type DeleteTeamInviteCodeRequest struct {
	TeamId   string
	InviteId string
}

func (v *TeamService) DeleteTeamInviteCode(ctx context.Context, request DeleteTeamInviteCodeRequest) (*DeleteTeamInviteCodeResponse, error) {
	path := fmt.Sprintf("/v1/teams/%s/invites/%s", request.TeamId, request.InviteId)

	response := &DeleteTeamInviteCodeResponse{}
	if err := v.client.http.Delete(ctx, path, response); err != nil {
		return nil, err
	}

	return response, nil
}

type GetTeamParameter struct {
	TeamId string
	Slug   *string
}

func (v *TeamService) GetTeam(ctx context.Context, request GetTeamParameter) (*TeamInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v2/teams/%s", request.TeamId),
	}
	queryParams := url.Values{}
	if request.Slug != nil {
		queryParams.Add("slug", *request.Slug)
	}
	u.RawQuery = queryParams.Encode()
	path := u.String()

	response := &TeamInfo{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}

	return response, nil
}

type GetAccessRequestStatusParameter struct {
	TeamId string
	UserId string
}

func (v *TeamService) GetAccessRequestStatus(ctx context.Context, request GetAccessRequestStatusParameter) (*AccessRequestStatus, error) {
	path := fmt.Sprintf("/v1/teams/%s/request/%s", request.TeamId, request.UserId)
	response := &AccessRequestStatus{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}

	return response, nil
}

type ListTeamMembersParameter struct {
	TeamId                      string
	EligibleMembersForProjectId *string
	ExcludeProject              *string
	Limit                       *string
	Role                        *string
	Search                      *string
	Since                       *string
	Until                       *string
}

func (v *TeamService) ListTeamMembers(ctx context.Context, request ListTeamMembersParameter) (*ListTeamMembersResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v2/teams/%s/members", request.TeamId),
	}
	queryParams := url.Values{}
	if request.EligibleMembersForProjectId != nil {
		queryParams.Add("eligibleMembersForProjectId", *request.EligibleMembersForProjectId)
	}
	if request.ExcludeProject != nil {
		queryParams.Add("excludeProject", *request.ExcludeProject)
	}
	if request.Limit != nil {
		queryParams.Add("limit", *request.Limit)
	}
	if request.Role != nil {
		queryParams.Add("role", *request.Role)
	}
	if request.Search != nil {
		queryParams.Add("search", *request.Search)
	}
	if request.Since != nil {
		queryParams.Add("since", *request.Since)
	}
	if request.Until != nil {
		queryParams.Add("until", *request.Until)
	}
	u.RawQuery = queryParams.Encode()
	path := u.String()

	response := &ListTeamMembersResponse{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}

	return response, nil
}

type ListTeamParameter struct {
	Limit *string
	Since *string
	Until *string
}

func (v *TeamService) ListTeam(ctx context.Context, request ListTeamParameter) (*ListTeamResponse, error) {
	u := &url.URL{
		Path: "/v2/teams",
	}
	queryParams := url.Values{}
	if request.Limit != nil {
		queryParams.Add("limit", *request.Limit)
	}
	if request.Since != nil {
		queryParams.Add("since", *request.Since)
	}
	if request.Until != nil {
		queryParams.Add("until", *request.Until)
	}
	u.RawQuery = queryParams.Encode()
	path := u.String()

	response := &ListTeamResponse{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (v *TeamService) InviteUser(ctx context.Context, req InviteUserRequest) (*InviteUserResponse, error) {
	path := fmt.Sprintf("/v1/teams/%s/members", req.TeamId)
	response := &InviteUserResponse{}
	if err := v.client.http.Post(ctx, path, req, response); err != nil {
		return nil, err
	}

	return response, nil
}

type JoinTeamRequest struct {
	TeamId     string `json:"-"`
	InviteCode string `json:"inviteCode,omitempty"`
}

func (v *TeamService) JoinTeam(ctx context.Context, request JoinTeamRequest) (*JoinTeamResponse, error) {
	path := fmt.Sprintf("/v1/teams/%s/members/teams/join", request.TeamId)

	response := &JoinTeamResponse{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (v *TeamService) UpdateTeam(ctx context.Context, request UpdateTeamInfo) (*TeamInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v2/teams/%s", request.TeamId),
	}
	queryParams := url.Values{}
	queryParams.Add("slug", request.Slug)
	u.RawQuery = queryParams.Encode()
	path := u.String()

	response := &TeamInfo{}
	if err := v.client.http.Patch(ctx, path, request, response); err != nil {
		return nil, err
	}

	return response, nil
}

type RemoveTeamMemberRequest struct {
	TeamId           string  `json:"-"`
	Uid              string  `json:"-"`
	NewDefaultTeamId *string `json:"-"`
}

func (v *TeamService) RemoveTeamMember(ctx context.Context, request RemoveTeamMemberRequest) (*RemoveTeamMemberResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/teams/%s/members/%s", request.TeamId, request.Uid),
	}
	queryParams := url.Values{}
	if request.NewDefaultTeamId != nil {
		queryParams.Add("newDefaultTeamId", *request.NewDefaultTeamId)
	}
	u.RawQuery = queryParams.Encode()
	path := u.String()

	response := &RemoveTeamMemberResponse{}
	if err := v.client.http.Delete(ctx, path, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (v *TeamService) RequestAccessToTeam(ctx context.Context, req JoinedFrom) (*AccessRequestStatus, error) {
	path := fmt.Sprintf("/v1/teams/%s/request", req.TeamId)
	response := &AccessRequestStatus{}
	if err := v.client.http.Post(ctx, path, req, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (v *TeamService) UpdateTeamMember(ctx context.Context, req UpdateTeamMemberRequest) (*UpdateTeamMemberResponse, error) {
	path := fmt.Sprintf("/v1/teams/%s/members/%s", req.TeamId, req.Uid)
	response := &UpdateTeamMemberResponse{}
	if err := v.client.http.Patch(ctx, path, req, response); err != nil {
		return nil, err
	}

	return response, nil
}

type UpdateTeamMemberResponse struct {
	Id    string       `json:"id,omitempty"`
	Error *VercelError `json:"error,omitempty"`
}

type UpdateTeamMemberRequest struct {
	TeamId     string          `json:"-"`
	Uid        string          `json:"-"`
	Confirmed  bool            `json:"confirmed,omitempty"`
	JoinedFrom JoinedFrom      `json:"joinedFrom,omitempty"`
	Role       string          `json:"role,omitempty"`
	Projects   []InviteProject `json:"projects,omitempty"`
}

type RemoveTeamMemberResponse struct {
	Id    string       `json:"id,omitempty"`
	Error *VercelError `json:"error,omitempty"`
}

type UpdateTeamInfo struct {
	TeamId                             string        `json:"-"`
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
	TeamId   string          `json:"-"`
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
	Id        string       `json:"id,omitempty"`
	Slug      string       `json:"slug,omitempty"`
	Name      string       `json:"name,omitempty"`
	CreatorId string       `json:"creatorId,omitempty"`
	Error     *VercelError `json:"error,omitempty"`
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
	TeamId           string `json:"-"`
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
