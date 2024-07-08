package vercel

import (
	"fmt"
	"net/url"
)

func (v *VercelClient) CreateEdgeConfig(slug string, teamId string, req CreateEdgeConfigRequest) (*EdgeConfigInfo, error) {
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
	result := &EdgeConfigInfo{}
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

func (v *VercelClient) DeleteEdgeConfig(edgeConfigId string, slug string, teamId string) (*string, error) {
	path := fmt.Sprintf("/v1/edge-config/%s", edgeConfigId)
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

func (v *VercelClient) DeleteEdgeConfigSchema(edgeConfigId string, slug string, teamId string) (*string, error) {
	path := fmt.Sprintf("/v1/edge-config/%s/schema", edgeConfigId)
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

func (v *VercelClient) DeleteEdgeConfigTokens(edgeConfigId string, slug string, teamId string, tokens []string) (*string, error) {
	path := fmt.Sprintf("/v1/edge-config/%s/tokens", edgeConfigId)
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
		Tokens []string `json:"tokens"`
	}{
		Tokens: tokens,
	}
	result := ""
	err := v.http.DeleteWithBody(path, req, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *VercelClient) GetEdgeConfig(edgeConfigId string, slug string, teamId string) (*EdgeConfigInfo, error) {
	path := fmt.Sprintf("/v1/edge-config/%s", edgeConfigId)
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
	result := &EdgeConfigInfo{}
	err := v.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type CreateEdgeConfigRequest struct {
	Slug  string      `json:"slug"`
	Items interface{} `json:"items"`
}

type EdgeConfigInfo struct {
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
