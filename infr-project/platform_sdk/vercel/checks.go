package vercel

import (
	"fmt"
	"net/url"
)

func (v *VercelClient) CreatesCheck(deploymentId string, slug string, teamId string, req CreateCheckRequest) (*CheckInfo, error) {
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
	err := v.http.Post(path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) ListAllChecks(deploymentId string, slug string, teamId string) ([]CheckInfo, error) {
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
	err := v.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) GetSingleCheck(deploymentId string, checkId string, slug string, teamId string) (*CheckInfo, error) {
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
	err := v.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) RerequestCheck(deploymentId string, checkId string, slug string, teamId string) (*CheckInfo, error) {
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
	err := v.http.Post(path, nil, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) UpdateCheck(deploymentId string, checkId string, slug string, teamId string, req CheckInfo) (*CheckInfo, error) {
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
	err := v.http.Post(path, req, result)

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
