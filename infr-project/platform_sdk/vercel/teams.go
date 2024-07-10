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
