package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type CheckService service

type CreateCheckRequest struct {
	Blocking         bool   `json:"blocking"`
	Name             string `json:"name"`
	DetailsUrl       string `json:"detailsUrl,omitempty"`
	ExternalId       string `json:"externalId,omitempty"`
	Path             string `json:"path,omitempty"`
	Rerequestable    bool   `json:"rerequestable,omitempty"`
	DeploymentId     string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *CheckService) CreatesCheck(ctx context.Context, request CreateCheckRequest) (*CheckInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/deployments/%s/checks", request.DeploymentId),
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

	response := &CheckInfo{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type ListChecksParameter struct {
	DeploymentId     string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *CheckService) ListAllChecks(ctx context.Context, request ListChecksParameter) ([]CheckInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/deployments/%s/checks", request.DeploymentId),
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

	response := []CheckInfo{}
	if err := v.client.http.Get(ctx, path, &response); err != nil {
		return nil, err
	}
	return response, nil
}

type SingleCheckParameter struct {
	DeploymentId     string `json:"-"`
	CheckId          string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *CheckService) GetSingleCheck(ctx context.Context, request SingleCheckParameter) (*CheckInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/deployments/%s/checks/%s", request.DeploymentId, request.CheckId),
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

	response := &CheckInfo{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type RerequestCheckRequest struct {
	DeploymentId     string `json:"-"`
	CheckId          string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *CheckService) RerequestCheck(ctx context.Context, request RerequestCheckRequest) (*CheckInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/deployments/%s/checks/%s/rerequest", request.DeploymentId, request.CheckId),
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

	response := &CheckInfo{}
	if err := v.client.http.Post(ctx, path, nil, response); err != nil {
		return nil, err
	}
	return response, nil
}

type UpdateCheckRequest struct {
	Conclusion       string      `json:"conclusion"`
	DetailsUrl       string      `json:"detailsUrl"`
	ExternalId       string      `json:"externalId"`
	Name             string      `json:"name"`
	Output           CheckOutput `json:"output"`
	Path             string      `json:"path"`
	Status           string      `json:"status"`
	DeploymentId     string      `json:"-"`
	CheckId          string      `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *CheckService) UpdateCheck(ctx context.Context, request UpdateCheckRequest) (*CheckInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/deployments/%s/checks/%s", request.DeploymentId, request.CheckId),
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

	response := &CheckInfo{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type CheckInfo struct {
	Conclusion    string       `json:"conclusion"`
	DetailsUrl    string       `json:"detailsUrl"`
	ExternalId    string       `json:"externalId"`
	Name          string       `json:"name"`
	Output        CheckOutput  `json:"output"`
	Path          string       `json:"path"`
	Status        string       `json:"status"`
	Blocking      bool         `json:"blocking"`
	CompletedAt   int          `json:"completedAt"`
	CreatedAt     int          `json:"createdAt"`
	StartedAt     int          `json:"startedAt"`
	UpdatedAt     int          `json:"updatedAt"`
	DeploymentId  string       `json:"deploymentId"`
	Id            string       `json:"id"`
	IntegrationId string       `json:"integrationId"`
	Rerequestable bool         `json:"rerequestable,omitempty"`
	Error         *VercelError `json:"error"`
}

type CheckOutput struct {
	Metrics map[string]CheckCls `json:"metrics"`
}

type CheckCls struct {
	Value         int64  `json:"value"`
	PreviousValue int64  `json:"previousValue"`
	Source        string `json:"source"`
}
