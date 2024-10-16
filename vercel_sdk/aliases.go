package vercel

import (
	"fmt"
	"net/url"
	"strings"
)

type AliasService service

func (v *AliasService) AssignAlias(id string, slug string, teamId string, info AssignAliasRequest) (*AssignAliasResponse, error) {
	path := fmt.Sprintf("/v2/deployments/%s/aliases", id)
	if len(slug) > 0 {
		path += ("?slug=" + slug)
	}
	if len(teamId) > 0 {
		if strings.Contains(path, "?") {
			path += ("&teamId=" + teamId)
		} else {
			path += ("?teamId=" + teamId)
		}
	}
	result := &AssignAliasResponse{}
	err := v.client.http.Post(path, info, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *AliasService) DeleteAlias(id string, slug string, teamId string) (*DeleteAliasResponse, error) {
	path := fmt.Sprintf("/v2/aliases/%s", id)
	if len(slug) > 0 {
		path += ("?slug=" + slug)
	}
	if len(teamId) > 0 {
		if strings.Contains(path, "?") {
			path += ("&teamId=" + teamId)
		} else {
			path += ("?teamId=" + teamId)
		}
	}
	result := &DeleteAliasResponse{}
	err := v.client.http.Delete(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *AliasService) GetAlias(id string, slug string, teamId string, projectId string, since string, until string) ([]AliasInfo, error) {
	path := fmt.Sprintf("/v4/aliases/%s", id)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(projectId) > 0 {
		queryParams.Add("projectId", projectId)
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
	result := []AliasInfo{}
	err := v.client.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *AliasService) ListAlias(slug string, teamId string, projectId string, since string, until string) (*ListAliasResponse, error) {
	path := "/v4/aliases"
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(projectId) > 0 {
		queryParams.Add("projectId", projectId)
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
	result := &ListAliasResponse{}
	err := v.client.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *AliasService) ListdeploymentsAlias(id string, slug string, teamId string) ([]AliasInfo, error) {
	path := fmt.Sprintf("/v2/deployments/%s/aliases", id)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	result := []AliasInfo{}
	err := v.client.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type ListAliasResponse struct {
	Aliases    []AliasInfo  `json:"aliases,omitempty"`
	Pagination Pagination   `json:"pagination,omitempty"`
	Error      *VercelError `json:"error,omitempty"`
}

type AssignAliasRequest struct {
	Alias    string       `json:"alias,omitempty"`
	Redirect string       `json:"redirect,omitempty"`
	Error    *VercelError `json:"error,omitempty"`
}

type AssignAliasResponse struct {
	Alias           string       `json:"alias,omitempty"`
	Created         string       `json:"created,omitempty"`
	Uid             string       `json:"uid,omitempty"`
	OldDeploymentId string       `json:"oldDeploymentId"`
	Error           *VercelError `json:"error,omitempty"`
}

type DeleteAliasResponse struct {
	Status string       `json:"status,omitempty"`
	Error  *VercelError `json:"error,omitempty"`
}

type AliasInfo struct {
	Alias              string         `json:"alias,omitempty"`
	Id                 string         `json:"id,omitempty"`
	Status             string         `json:"status,omitempty"`
	Created            string         `json:"created,omitempty"`
	CreatedAt          int            `json:"createdAt,omitempty"`
	Creator            CreatorInfo    `json:"creator,omitempty"`
	DeploymentId       string         `json:"deploymentId,omitempty"`
	ProjectId          string         `json:"projectId,omitempty"`
	Redirect           string         `json:"redirect,omitempty"`
	RedirectStatusCode int            `json:"redirectStatusCode,omitempty"`
	Uid                string         `json:"uid,omitempty"`
	UpdatedAt          int            `json:"updatedAt,omitempty"`
	Deployment         DeploymentInfo `json:"deployment,omitempty"`
	ProtectionBypass   interface{}    `json:"protectionBypass,omitempty"`
	Error              *VercelError   `json:"error,omitempty"`
}

type CreatorInfo struct {
	Email    string `json:"email,omitempty"`
	Creator  string `json:"creator,omitempty"`
	Uid      string `json:"uid,omitempty"`
	Username string `json:"username,omitempty"`
}

type AliasWarning struct {
	Action  string `json:"action,omitempty"`
	Code    string `json:"code,omitempty"`
	Link    string `json:"link,omitempty"`
	Message string `json:"message,omitempty"`
}
