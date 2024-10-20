package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type AliasService service

func (v *AliasService) AssignAlias(ctx context.Context, request AssignAliasRequest) (*AssignAliasResponse, error) {
	if request.Id == nil {
		return nil, fmt.Errorf("vercel assign alias request need id")
	}

	u := &url.URL{
		Path: fmt.Sprintf("/v2/deployments/%s/aliases", *request.Id),
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

	response := &AssignAliasResponse{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}

	return response, nil
}

type DeleteAliasRequest struct {
	Id               string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *AliasService) DeleteAlias(ctx context.Context, request DeleteAliasRequest) (*DeleteAliasResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v2/aliases/%s/", request.Id),
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

	response := &DeleteAliasResponse{}
	if err := v.client.http.Delete(ctx, path, response); err != nil {
		return nil, err
	}

	return response, nil
}

type GetAliasParameter struct {
	Id               string  `json:"-"`
	ProjectId        *string `json:"-"`
	Since            *string `json:"-"`
	Until            *string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *AliasService) GetAlias(ctx context.Context, request GetAliasParameter) ([]AliasInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v4/aliases/%s/", request.Id),
	}
	params := url.Values{}
	if request.TeamId != nil {
		params.Add("teamId", *request.TeamId)
	}
	if request.TeamSlug != nil {
		params.Add("slug", *request.TeamSlug)
	}
	if request.ProjectId != nil {
		params.Add("projectId", *request.ProjectId)
	}
	if request.Since != nil {
		params.Add("since", *request.Since)
	}
	if request.Until != nil {
		params.Add("until", *request.Until)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := []AliasInfo{}
	if err := v.client.http.Get(ctx, path, &response); err != nil {
		return nil, err
	}

	return response, nil
}

type ListAliasParameter struct {
	ProjectId        *string `json:"-"`
	Since            *string `json:"-"`
	Until            *string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *AliasService) ListAlias(ctx context.Context, request ListAliasParameter) (*ListAliasResponse, error) {
	u := &url.URL{
		Path: "/v4/aliases",
	}
	params := url.Values{}
	if request.TeamId != nil {
		params.Add("teamId", *request.TeamId)
	}
	if request.TeamSlug != nil {
		params.Add("slug", *request.TeamSlug)
	}
	if request.ProjectId != nil {
		params.Add("projectId", *request.ProjectId)
	}
	if request.Since != nil {
		params.Add("since", *request.Since)
	}
	if request.Until != nil {
		params.Add("until", *request.Until)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ListAliasResponse{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}

	return response, nil
}

type ListDeploymentsAliasParameter struct {
	DeploymentId     string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *AliasService) ListDeploymentsAlias(ctx context.Context, request ListDeploymentsAliasParameter) ([]AliasInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v2/deployments/%s/aliases", request.DeploymentId),
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

	response := []AliasInfo{}
	if err := v.client.http.Get(ctx, path, &response); err != nil {
		return nil, err
	}

	return response, nil
}

type ListAliasResponse struct {
	Aliases    []AliasInfo  `json:"aliases,omitempty"`
	Pagination Pagination   `json:"pagination,omitempty"`
	Error      *VercelError `json:"error,omitempty"`
}

type AssignAliasRequest struct {
	Id               *string `json:"-"`
	Alias            string  `json:"alias,omitempty"`
	Redirect         string  `json:"redirect,omitempty"`
	BaseUrlParameter `json:"-"`
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
