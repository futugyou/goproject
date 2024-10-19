package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type CheckService service

func (v *CheckService) CreatesCheck(ctx context.Context, deploymentId string, slug string, teamId string, req CreateCheckRequest) (*CheckInfo, error) {
	path := fmt.Sprintf("/v1/deployments/%s/checks", deploymentId)
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
	result := &CheckInfo{}
	err := v.client.http.Post(ctx, path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *CheckService) ListAllChecks(ctx context.Context, deploymentId string, slug string, teamId string) ([]CheckInfo, error) {
	path := fmt.Sprintf("/v1/deployments/%s/checks", deploymentId)
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
	result := []CheckInfo{}
	err := v.client.http.Get(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *CheckService) GetSingleCheck(ctx context.Context, deploymentId string, checkId string, slug string, teamId string) (*CheckInfo, error) {
	path := fmt.Sprintf("/v1/deployments/%s/checks/%s", deploymentId, checkId)
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
	result := &CheckInfo{}
	err := v.client.http.Get(ctx, path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *CheckService) RerequestCheck(ctx context.Context, deploymentId string, checkId string, slug string, teamId string) (*CheckInfo, error) {
	path := fmt.Sprintf("/v1/deployments/%s/checks/%s/rerequest", deploymentId, checkId)
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
	result := &CheckInfo{}
	err := v.client.http.Post(ctx, path, nil, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *CheckService) UpdateCheck(ctx context.Context, deploymentId string, checkId string, slug string, teamId string, req CheckInfo) (*CheckInfo, error) {
	path := fmt.Sprintf("/v1/deployments/%s/checks/%s", deploymentId, checkId)
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
	result := &CheckInfo{}
	err := v.client.http.Post(ctx, path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type CreateCheckRequest struct {
	Blocking      bool   `json:"blocking"`
	Name          string `json:"name"`
	DetailsUrl    string `json:"detailsUrl,omitempty"`
	ExternalId    string `json:"externalId,omitempty"`
	Path          string `json:"path,omitempty"`
	Rerequestable bool   `json:"rerequestable,omitempty"`
}

type CheckInfo struct {
	Blocking      bool         `json:"blocking"`
	CompletedAt   int          `json:"completedAt"`
	Conclusion    string       `json:"conclusion"`
	CreatedAt     int          `json:"createdAt"`
	DeploymentId  string       `json:"deploymentId"`
	DetailsUrl    string       `json:"detailsUrl"`
	ExternalId    string       `json:"externalId"`
	Id            string       `json:"id"`
	IntegrationId string       `json:"integrationId"`
	Name          string       `json:"name"`
	Path          string       `json:"path"`
	Output        CheckOutput  `json:"output"`
	Error         *VercelError `json:"error"`
}

type CheckOutput struct {
	Metrics CheckMetrics `json:"metrics"`
}

type CheckMetrics struct {
	Fcp                    CheckCls `json:"FCP"`
	Lcp                    CheckCls `json:"LCP"`
	Cls                    CheckCls `json:"CLS"`
	Tbt                    CheckCls `json:"TBT"`
	VirtualExperienceScore CheckCls `json:"virtualExperienceScore"`
}

type CheckCls struct {
	Value         int64  `json:"value"`
	PreviousValue int64  `json:"previousValue"`
	Source        string `json:"source"`
}
