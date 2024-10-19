package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type LogDrainService service

func (v *LogDrainService) CreateLogDrain(ctx context.Context, slug string, teamId string, req CreateLogDrainRequest) (*LogDrainInfo, error) {
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
	err := v.client.http.Post(ctx, path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *LogDrainService) CreateIntegrationLogDrain(ctx context.Context, slug string, teamId string, req CreateLogDrainRequest) (*LogDrainInfo, error) {
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
	err := v.client.http.Post(ctx, path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *LogDrainService) DeletesConfigurableLogDrain(ctx context.Context, id string, slug string, teamId string) (*string, error) {
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
	err := v.client.http.Delete(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *LogDrainService) DeletesIntegrationConfigurableLogDrain(ctx context.Context, id string, slug string, teamId string) (*string, error) {
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
	err := v.client.http.Delete(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *LogDrainService) RetrievesListLogDrain(ctx context.Context, projectId string, slug string, teamId string) ([]LogDrainInfo, error) {
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
	err := v.client.http.Delete(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *LogDrainService) RetrievesConfigurableLogDrain(ctx context.Context, id string, slug string, teamId string) (*LogDrainInfo, error) {
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
	err := v.client.http.Delete(ctx, path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *LogDrainService) RetrievesIntegrationLogDrain(ctx context.Context, slug string, teamId string) ([]LogDrainInfo, error) {
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
	result := []LogDrainInfo{}
	err := v.client.http.Delete(ctx, path, &result)

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
