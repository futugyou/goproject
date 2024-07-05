package vercel

import (
	"fmt"
	"net/url"
	"strings"
)

func (v *VercelClient) AssignAlias(id string, slug string, teamId string, info AssignAliasRequest) (*AssignAliasResponse, error) {
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
	err := v.http.Post(path, info, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) DeleteAlias(id string, slug string, teamId string) (*DeleteAliasResponse, error) {
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
	err := v.http.Delete(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) GetAlias(id string, slug string, teamId string, projectId string, since string, until string) ([]AliasInfo, error) {
	path := fmt.Sprintf("/v2/aliases/%s", id)
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
	err := v.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type AssignAliasRequest struct {
	Alias    string       `json:"alias"`
	Redirect string       `json:"redirect,omitempty"`
	Error    *VercelError `json:"error"`
}

type AssignAliasResponse struct {
	Alias           string       `json:"alias"`
	Created         string       `json:"created"`
	Uid             string       `json:"uid"`
	OldDeploymentId string       `json:"oldDeploymentId"`
	Error           *VercelError `json:"error"`
}

type DeleteAliasResponse struct {
	Status string       `json:"status"`
	Error  *VercelError `json:"error"`
}

type AliasInfo struct {
	Alias              string         `json:"alias"`
	Created            string         `json:"created"`
	CreatedAt          int            `json:"createdAt"`
	Creator            CreatorInfo    `json:"creator"`
	DeploymentId       string         `json:"deploymentId"`
	ProjectId          string         `json:"projectId"`
	Redirect           string         `json:"redirect"`
	RedirectStatusCode int            `json:"redirectStatusCode"`
	Uid                string         `json:"uid"`
	UpdatedAt          int            `json:"updatedAt"`
	Deployment         DeploymentInfo `json:"deployment"`
	ProtectionBypass   interface{}    `json:"protectionBypass"`
	Error              *VercelError   `json:"error"`
}

type DeploymentInfo struct {
	Id   string `json:"id"`
	Meta string `json:"meta"`
	Url  string `json:"url"`
}

type CreatorInfo struct {
	Email    string `json:"email"`
	Uid      string `json:"uid"`
	Username string `json:"username"`
}
