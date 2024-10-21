package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type EdgeService service

type UpsertEdgeConfigRequest struct {
	EdgeConfigId     string      `json:"-"`
	Slug             string      `json:"slug"`
	Items            interface{} `json:"items,omitempty"`
	BaseUrlParameter `json:"-"`
}

func (v *EdgeService) CreateEdgeConfig(ctx context.Context, request UpsertEdgeConfigRequest) (*EdgeConfigInfo, error) {
	u := &url.URL{
		Path: "/v1/edge-config",
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &EdgeConfigInfo{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type CreateEdgeConfigTokenRequest struct {
	Label            string `json:"label"`
	EdgeConfigId     string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *EdgeService) CreateEdgeConfigToken(ctx context.Context, request CreateEdgeConfigTokenRequest) (*CreateEdgeConfigTokenResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/edge-config/%s/token", request.EdgeConfigId),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &CreateEdgeConfigTokenResponse{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type DeleteEdgeConfigRequest struct {
	EdgeConfigId     string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *EdgeService) DeleteEdgeConfig(ctx context.Context, request DeleteEdgeConfigRequest) (*string, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/edge-config/%s", request.EdgeConfigId),
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

type DeleteEdgeConfigSchemaRequest struct {
	EdgeConfigId     string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *EdgeService) DeleteEdgeConfigSchema(ctx context.Context, request DeleteEdgeConfigSchemaRequest) (*string, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/edge-config/%s/schema", request.EdgeConfigId),
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

type DeleteEdgeConfigTokensRequest struct {
	EdgeConfigId     string   `json:"-"`
	Tokens           []string `json:"tokens"`
	BaseUrlParameter `json:"-"`
}

func (v *EdgeService) DeleteEdgeConfigTokens(ctx context.Context, request DeleteEdgeConfigTokensRequest) (*string, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/edge-config/%s/tokens", request.EdgeConfigId),
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

type GetEdgeConfigParameter struct {
	EdgeConfigId     string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *EdgeService) GetEdgeConfig(ctx context.Context, request GetEdgeConfigParameter) (*EdgeConfigInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/edge-config/%s", request.EdgeConfigId),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &EdgeConfigInfo{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type GetEdgeConfigBackupParameter struct {
	EdgeConfigId              string `json:"-"`
	EdgeConfigBackupVersionId string `json:"-"`
	BaseUrlParameter          `json:"-"`
}

func (v *EdgeService) GetEdgeConfigBackup(ctx context.Context, request GetEdgeConfigBackupParameter) (*EdgeConfigBackupInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/edge-config/%s/backups/%s", request.EdgeConfigId, request.EdgeConfigBackupVersionId),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &EdgeConfigBackupInfo{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type ListEdgeConfigBackupsParameter struct {
	EdgeConfigId     string  `json:"-"`
	Limit            *string `json:"-"`
	Metadata         *string `json:"-"`
	Next             *string `json:"-"`
	BaseUrlParameter `json:"-"`
}

type ListEdgeConfigBackupsResponse struct {
	Backups    []ListEdgeConfigBackupItem `json:"backups"`
	Pagination Pagination                 `json:"pagination,omitempty"`
	Error      *VercelError               `json:"error,omitempty"`
}

type ListEdgeConfigBackupItem struct {
	Id           string       `json:"id"`
	LastModified int          `json:"lastModified"`
	Metadata     EdgeMetadata `json:"metadata"`
}

func (v *EdgeService) ListEdgeConfigBackups(ctx context.Context, request ListEdgeConfigBackupsParameter) (*ListEdgeConfigBackupsResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/edge-config/%s/backups", request.EdgeConfigId),
	}
	params := request.GetUrlValues()
	if request.Limit != nil {
		params.Add("limit", *request.Limit)
	}
	if request.Metadata != nil {
		params.Add("metadata", *request.Metadata)
	}
	if request.Next != nil {
		params.Add("next", *request.Next)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ListEdgeConfigBackupsResponse{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type GetEdgeConfigItemParameter struct {
	EdgeConfigId      string `json:"-"`
	EdgeConfigItemKey string `json:"-"`
	BaseUrlParameter  `json:"-"`
}

func (v *EdgeService) GetEdgeConfigItem(ctx context.Context, request GetEdgeConfigItemParameter) (*EdgeConfigItemInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/edge-config/%s/item/%s", request.EdgeConfigId, request.EdgeConfigItemKey),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &EdgeConfigItemInfo{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type ListEdgeConfigItemsParameter struct {
	EdgeConfigId     string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *EdgeService) ListEdgeConfigItems(ctx context.Context, request ListEdgeConfigItemsParameter) ([]EdgeConfigItemInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/edge-config/%s/items", request.EdgeConfigId),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := []EdgeConfigItemInfo{}
	if err := v.client.http.Get(ctx, path, &response); err != nil {
		return nil, err
	}
	return response, nil
}

type GetEdgeConfigSchemaParameter struct {
	EdgeConfigId     string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *EdgeService) GetEdgeConfigSchema(ctx context.Context, request GetEdgeConfigSchemaParameter) (*string, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/edge-config/%s/schema", request.EdgeConfigId),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := ""
	if err := v.client.http.Get(ctx, path, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

type GetEdgeConfigTokenParameter struct {
	EdgeConfigId     string `json:"-"`
	Token            string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *EdgeService) GetEdgeConfigToken(ctx context.Context, request GetEdgeConfigTokenParameter) (*EdgeConfigTokenInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/edge-config/%s/token/%s", request.EdgeConfigId, request.Token),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &EdgeConfigTokenInfo{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type ListEdgeConfigTokensParameter struct {
	EdgeConfigId     string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *EdgeService) ListEdgeConfigTokens(ctx context.Context, request ListEdgeConfigTokensParameter) ([]EdgeConfigTokenInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/edge-config/%s/tokens", request.EdgeConfigId),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := []EdgeConfigTokenInfo{}
	err := v.client.http.Get(ctx, path, &response)

	if err != nil {
		return nil, err
	}
	return response, nil
}

type ListEdgeConfigsParameter struct {
	BaseUrlParameter `json:"-"`
}

// TODO: The document shows that an object is returned, I guess it may be a list
// https://vercel.com/docs/rest-api/endpoints/edge-config#get-edge-configs
func (v *EdgeService) ListEdgeConfigs(ctx context.Context, request ListEdgeConfigsParameter) ([]EdgeConfigInfo, error) {
	u := &url.URL{
		Path: "/v1/edge-config",
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := []EdgeConfigInfo{}
	if err := v.client.http.Get(ctx, path, &response); err != nil {
		return nil, err
	}
	return response, nil
}

type UpdateEdgeConfigItemsRequest struct {
	EdgeConfigId     string                 `json:"-"`
	DryRun           *string                `json:"-"`
	Items            []UpdateEdgeConfigItem `json:"items"`
	BaseUrlParameter `json:"-"`
}

type UpdateEdgeConfigItem struct {
	Description string      `json:"description"`
	Key         string      `json:"key"`
	Operation   string      `json:"operation"`
	Value       interface{} `json:"value"`
}

func (v *EdgeService) UpdateEdgeConfigItems(ctx context.Context, request UpdateEdgeConfigItemsRequest) (*UpdateEdgeConfigItemResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/edge-config/%s/items", request.EdgeConfigId),
	}
	params := request.GetUrlValues()
	if request.DryRun != nil {
		params.Add("dryRun", *request.DryRun)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &UpdateEdgeConfigItemResponse{}
	if err := v.client.http.Patch(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type UpdateEdgeConfigSchemaRequest struct {
	EdgeConfigId     string      `json:"-"`
	DryRun           *string     `json:"-"`
	Definition       interface{} `json:"definition"`
	BaseUrlParameter `json:"-"`
}

func (v *EdgeService) UpdateEdgeConfigSchema(ctx context.Context, request UpdateEdgeConfigSchemaRequest) (*string, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/edge-config/%s/schema", request.EdgeConfigId),
	}
	params := request.GetUrlValues()
	if request.DryRun != nil {
		params.Add("dryRun", *request.DryRun)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := ""
	if err := v.client.http.Post(ctx, path, request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (v *EdgeService) UpdateEdgeConfig(ctx context.Context, request UpsertEdgeConfigRequest) (*EdgeConfigTokenInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/edge-config/%s", request.EdgeConfigId),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &EdgeConfigTokenInfo{}
	if err := v.client.http.Put(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type UpdateEdgeConfigItemResponse struct {
	Status string       `json:"status"`
	Error  *VercelError `json:"error,omitempty"`
}

type EdgeConfigTokenInfo struct {
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

type EdgeConfigBackupInfo struct {
	Backup       int          `json:"backup"`
	Id           string       `json:"id"`
	LastModified int          `json:"lastModified"`
	Metadata     EdgeMetadata `json:"metadata"`
	User         EdgeUser     `json:"user"`
	Error        *VercelError `json:"error,omitempty"`
}

type EdgeUser struct {
	ItemsCount int    `json:"itemsCount"`
	ItemsBytes int    `json:"itemsBytes"`
	UpdatedAt  string `json:"updatedAt"`
	UpdatedBy  string `json:"updatedBy"`
}

type EdgeMetadata struct {
	ItemsCount int    `json:"itemsCount"`
	ItemsBytes int    `json:"itemsBytes"`
	UpdatedAt  string `json:"updatedAt"`
	UpdatedBy  string `json:"updatedBy"`
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
