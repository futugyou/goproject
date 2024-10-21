package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type LogDrainService service

type CreateLogDrainRequest struct {
	DeliveryFormat   string      `json:"deliveryFormat,omitempty"`
	Sources          []string    `json:"sources,omitempty"`
	Url              string      `json:"url,omitempty"`
	Environments     []string    `json:"environments,omitempty"`
	Headers          interface{} `json:"headers,omitempty"`
	ProjectIds       []string    `json:"projectIds,omitempty"`
	SamplingRate     int         `json:"samplingRate,omitempty"`
	Secret           string      `json:"secret,omitempty"`
	Name             string      `json:"name,omitempty"`
	BaseUrlParameter `json:"-"`
}

func (v *LogDrainService) CreateLogDrain(ctx context.Context, request CreateLogDrainRequest) (*LogDrainInfo, error) {
	u := &url.URL{
		Path: "/v1/log-drains",
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &LogDrainInfo{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

func (v *LogDrainService) CreateIntegrationLogDrain(ctx context.Context, request CreateLogDrainRequest) (*LogDrainInfo, error) {
	u := &url.URL{
		Path: "/v2/integrations/log-drains",
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &LogDrainInfo{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type DeletesConfigurableLogDrainRequest struct {
	Id               string
	BaseUrlParameter `json:"-"`
}

func (v *LogDrainService) DeletesConfigurableLogDrain(ctx context.Context, request DeletesConfigurableLogDrainRequest) (*string, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/log-drains/%s", request.Id),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := ""
	if err := v.client.http.Delete(ctx, path, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

type DeletesIntegrationConfigurableLogDrainRequest struct {
	Id               string
	BaseUrlParameter `json:"-"`
}

func (v *LogDrainService) DeletesIntegrationConfigurableLogDrain(ctx context.Context, request DeletesIntegrationConfigurableLogDrainRequest) (*string, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/integrations/log-drains/%s", request.Id),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := ""
	if err := v.client.http.Delete(ctx, path, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

type RetrievesListLogDrainParameter struct {
	ProjectId        *string
	BaseUrlParameter `json:"-"`
}

func (v *LogDrainService) RetrievesListLogDrain(ctx context.Context, request RetrievesListLogDrainParameter) ([]LogDrainInfo, error) {
	u := &url.URL{
		Path: "/v1/log-drains",
	}
	params := request.GetUrlValues()
	if request.ProjectId != nil {
		params.Add("projectId", *request.ProjectId)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := []LogDrainInfo{}
	if err := v.client.http.Delete(ctx, path, &response); err != nil {
		return nil, err
	}
	return response, nil
}

type RetrievesConfigurableLogDrainParameter struct {
	Id               string
	BaseUrlParameter `json:"-"`
}

func (v *LogDrainService) RetrievesConfigurableLogDrain(ctx context.Context, request RetrievesConfigurableLogDrainParameter) (*LogDrainInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/log-drains/%s", request.Id),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &LogDrainInfo{}
	if err := v.client.http.Delete(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type RetrievesIntegrationLogDrainParameter struct {
	BaseUrlParameter `json:"-"`
}

func (v *LogDrainService) RetrievesIntegrationLogDrain(ctx context.Context, request RetrievesIntegrationLogDrainParameter) ([]LogDrainInfo, error) {
	u := &url.URL{
		Path: "/v2/integrations/log-drains",
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := []LogDrainInfo{}
	if err := v.client.http.Delete(ctx, path, &response); err != nil {
		return nil, err
	}
	return response, nil
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
