package vercel

import (
	"fmt"
)

func (v *VercelClient) CreatesCheck(deploymentId string, slug string, teamId string, req CreateCheckRequest) (*CheckInfo, error) {
	path := fmt.Sprintf("/v1/deployments/%s/checks", deploymentId)
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
	Output        interface{}  `json:"output"`
	Error         *VercelError `json:"error"`
}
