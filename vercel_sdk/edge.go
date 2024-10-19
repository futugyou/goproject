package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type EdgeService service

func (v *EdgeService) CreateEdgeConfig(ctx context.Context, slug string, teamId string, req UpsertEdgeConfigRequest) (*EdgeConfigInfo, error) {
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
	err := v.client.http.Post(ctx, path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *EdgeService) CreateEdgeConfigToken(ctx context.Context, edgeConfigId string, slug string, teamId string, label string) (*CreateEdgeConfigTokenResponse, error) {
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
	err := v.client.http.Post(ctx, path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *EdgeService) DeleteEdgeConfig(ctx context.Context, edgeConfigId string, slug string, teamId string) (*string, error) {
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
	err := v.client.http.Delete(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *EdgeService) DeleteEdgeConfigSchema(ctx context.Context, edgeConfigId string, slug string, teamId string) (*string, error) {
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
	err := v.client.http.Delete(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *EdgeService) DeleteEdgeConfigTokens(ctx context.Context, edgeConfigId string, slug string, teamId string, tokens []string) (*string, error) {
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
	err := v.client.http.DeleteWithBody(ctx, path, req, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *EdgeService) GetEdgeConfig(ctx context.Context, edgeConfigId string, slug string, teamId string) (*EdgeConfigInfo, error) {
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
	err := v.client.http.Get(ctx, path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *EdgeService) GetEdgeConfigItems(ctx context.Context, edgeConfigId string, slug string, teamId string) ([]EdgeConfigItemInfo, error) {
	path := fmt.Sprintf("/v1/edge-config/%s/items", edgeConfigId)
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
	result := []EdgeConfigItemInfo{}
	err := v.client.http.Get(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *EdgeService) GetEdgeConfigSchema(ctx context.Context, edgeConfigId string, slug string, teamId string) (*string, error) {
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
	err := v.client.http.Get(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *EdgeService) GetEdgeConfigToken(ctx context.Context, edgeConfigId string, token string, slug string, teamId string) (*EdgeConfigTokennfo, error) {
	path := fmt.Sprintf("/v1/edge-config/%s/token/%s", edgeConfigId, token)
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
	result := &EdgeConfigTokennfo{}
	err := v.client.http.Get(ctx, path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *EdgeService) GetEdgeConfigTokens(ctx context.Context, edgeConfigId string, slug string, teamId string) ([]EdgeConfigTokennfo, error) {
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
	result := []EdgeConfigTokennfo{}
	err := v.client.http.Get(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *EdgeService) GetEdgeConfigs(ctx context.Context, slug string, teamId string) ([]EdgeConfigInfo, error) {
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
	result := []EdgeConfigInfo{}
	err := v.client.http.Get(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *EdgeService) UpdateEdgeConfigItems(ctx context.Context, edgeConfigId string, slug string, teamId string, items []UpdateEdgeConfigItemRequest) ([]UpdateEdgeConfigItemResponse, error) {
	path := fmt.Sprintf("/v1/edge-config/%s/items", edgeConfigId)
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
	result := []UpdateEdgeConfigItemResponse{}
	err := v.client.http.Patch(ctx, path, items, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *EdgeService) UpdateEdgeConfigSchema(ctx context.Context, edgeConfigId string, slug string, teamId string, req interface{}) (*string, error) {
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
	err := v.client.http.Post(ctx, path, req, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *EdgeService) UpdateEdgeConfig(ctx context.Context, edgeConfigId string, slug string, teamId string, req UpsertEdgeConfigRequest) (*EdgeConfigTokennfo, error) {
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
	result := &EdgeConfigTokennfo{}
	err := v.client.http.Put(ctx, path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type UpdateEdgeConfigItemResponse struct {
	Status string       `json:"status"`
	Error  *VercelError `json:"error,omitempty"`
}

type UpdateEdgeConfigItemRequest struct {
	Description string      `json:"description"`
	Key         string      `json:"key"`
	Operation   string      `json:"operation"`
	Value       interface{} `json:"value"`
}

type UpsertEdgeConfigRequest struct {
	Slug  string      `json:"slug"`
	Items interface{} `json:"items,omitempty"`
}

type EdgeConfigTokennfo struct {
	Token        string       `json:"token"`
	Id           string       `json:"id"`
	Label        string       `json:"label"`
	CreatedAt    int          `json:"createdAt"`
	EdgeConfigId string       `json:"edgeConfigId"`
	Error        *VercelError `json:"error,omitempty"`
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

type EdgeConfigItemInfo struct {
	CreatedAt    int          `json:"createdAt"`
	UpdatedAt    int          `json:"updatedAt"`
	Key          string       `json:"key"`
	Description  string       `json:"description"`
	EdgeConfigId string       `json:"edgeConfigId"`
	Value        interface{}  `json:"value"`
	Error        *VercelError `json:"error,omitempty"`
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
