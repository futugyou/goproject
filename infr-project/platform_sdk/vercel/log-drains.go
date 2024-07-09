package vercel

import (
	"fmt"
	"net/url"
)

func (v *VercelClient) CreateLogDrain(slug string, teamId string, req CreateLogDrainRequest) (*LogDrainInfo, error) {
	path := "/v1/log-drains"
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
	result := &LogDrainInfo{}
	err := v.http.Post(path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) CreateIntegrationLogDrain(slug string, teamId string, req CreateLogDrainRequest) (*LogDrainInfo, error) {
	path := "/v2/integrations/log-drains"
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
	result := &LogDrainInfo{}
	err := v.http.Post(path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) DeletesConfigurableLogDrain(id string, slug string, teamId string) (*string, error) {
	path := fmt.Sprintf("/v1/log-drains/%s", id)
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
	result := ""
	err := v.http.Delete(path, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *VercelClient) DeletesIntegrationConfigurableLogDrain(id string, slug string, teamId string) (*string, error) {
	path := fmt.Sprintf("/v1/integrations/log-drains/%s", id)
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
	result := ""
	err := v.http.Delete(path, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *VercelClient) RetrievesListLogDrain(projectId string, slug string, teamId string) ([]LogDrainInfo, error) {
	path := "/v1/log-drains"
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
	result := []LogDrainInfo{}
	err := v.http.Delete(path, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) RetrievesConfigurableLogDrain(id string, slug string, teamId string) (*LogDrainInfo, error) {
	path := fmt.Sprintf("/v1/log-drains/%s", id)
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
	result := &LogDrainInfo{}
	err := v.http.Delete(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type CreateLogDrainRequest struct {
	DeliveryFormat string      `json:"deliveryFormat,omitempty"`
	Sources        []string    `json:"sources,omitempty"`
	Url            string      `json:"url,omitempty"`
	Environments   []string    `json:"environments,omitempty"`
	Headers        interface{} `json:"headers,omitempty"`
	ProjectIds     []string    `json:"projectIds,omitempty"`
	SamplingRate   int         `json:"samplingRate,omitempty"`
	Secret         string      `json:"secret,omitempty"`
	Name           string      `json:"name,omitempty"`
}

type LogDrainInfo struct {
	Branch              string       `json:"branch,omitempty"`
	ClientId            string       `json:"clientId,omitempty"`
	Compression         string       `json:"compression,omitempty"`
	ConfigurationId     string       `json:"configurationId,omitempty"`
	CreatedFrom         string       `json:"createdFrom,omitempty"`
	DeliveryFormat      string       `json:"deliveryFormat,omitempty"`
	DisabledBy          string       `json:"disabledBy,omitempty"`
	CreatedAt           int          `json:"createdAt,omitempty"`
	DeletedAt           int          `json:"deletedAt,omitempty"`
	DisabledAt          int          `json:"disabledAt,omitempty"`
	DisabledReason      string       `json:"disabledReason,omitempty"`
	Environments        []string     `json:"environments,omitempty"`
	FirstErrorTimestamp int          `json:"firstErrorTimestamp,omitempty"`
	Headers             interface{}  `json:"headers,omitempty"`
	Id                  string       `json:"id,omitempty"`
	Name                string       `json:"name,omitempty"`
	OwnerId             string       `json:"ownerId,omitempty"`
	ProjectIds          []string     `json:"projectIds,omitempty"`
	SamplingRate        int          `json:"samplingRate,omitempty"`
	Secret              string       `json:"secret,omitempty"`
	Sources             []string     `json:"sources,omitempty"`
	Status              string       `json:"status,omitempty"`
	TeamId              string       `json:"teamId,omitempty"`
	Url                 string       `json:"url,omitempty"`
	UpdatedAt           int          `json:"updatedAt,omitempty"`
	Error               *VercelError `json:"error,omitempty"`
}
