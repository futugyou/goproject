package vercel

import (
	"fmt"
	"net/url"
)

func (v *VercelClient) CreateEdgeConfig(slug string, teamId string, req CreateEdgeConfigRequest) (*CreateEdgeConfigResponse, error) {
	path := "/v1/edge-config"
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
	result := &CreateEdgeConfigResponse{}
	err := v.http.Post(path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) CreateEdgeConfigToken(edgeConfigId string, slug string, teamId string, label string) (*CreateEdgeConfigTokenResponse, error) {
	path := fmt.Sprintf("/v1/edge-config/%s/token", edgeConfigId)
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

	req := struct {
		Label string `json:"label"`
	}{
		Label: label,
	}
	result := &CreateEdgeConfigTokenResponse{}
	err := v.http.Post(path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type CreateEdgeConfigRequest struct {
	Slug  string      `json:"slug"`
	Items interface{} `json:"items"`
}

type CreateEdgeConfigResponse struct {
	CreatedAt   int          `json:"createdAt"`
	Digest      string       `json:"digest"`
	Id          string       `json:"id"`
	ItemCount   int          `json:"itemCount"`
	OwnerId     string       `json:"ownerId"`
	Schema      interface{}  `json:"schema"`
	SizeInBytes int          `json:"sizeInBytes"`
	Slug        string       `json:"slug"`
	UpdatedAt   int          `json:"updatedAt"`
	Transfer    EdgeTransfer `json:"transfer"`
	Error       *VercelError `json:"error,omitempty"`
}

type EdgeTransfer struct {
	DoneAt        int    `json:"doneAt"`
	FromAccountId string `json:"fromAccountId"`
	StartedAt     int    `json:"startedAt"`
}

type CreateEdgeConfigTokenResponse struct {
	Token string       `json:"token"`
	Id    string       `json:"id"`
	Error *VercelError `json:"error,omitempty"`
}
