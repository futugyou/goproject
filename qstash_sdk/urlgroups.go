package qstash

import (
	"context"
	"fmt"
)

type URLGroupsService service

func (s *URLGroupsService) UpsertURLGroupAndEndpoint(ctx context.Context, request UpsertURLGroupRequest) error {
	path := fmt.Sprintf("/v2/topics/%s/endpoints", request.UrlGroupName)
	result := ""
	return s.client.http.Post(ctx, path, request, &result)
}

func (s *URLGroupsService) GetURLGroup(ctx context.Context, urlGroupName string) (*URLGroup, error) {
	path := fmt.Sprintf("/v2/topics/%s", urlGroupName)
	result := &URLGroup{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *URLGroupsService) ListURLGroup(ctx context.Context) (*URLGroupList, error) {
	path := "/v2/topics"
	result := &URLGroupList{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *URLGroupsService) RemoveEndpoints(ctx context.Context, request RemoveEndpointsRequest) error {
	path := fmt.Sprintf("/v2/topics/%s", request.UrlGroupName)
	result := ""
	return s.client.http.Delete(ctx, path, request, &result)
}

func (s *URLGroupsService) RemoveURLGroup(ctx context.Context, urlGroupName string) error {
	path := fmt.Sprintf("/v2/topics/%s", urlGroupName)
	result := ""
	return s.client.http.Delete(ctx, path, nil, &result)
}

type UpsertURLGroupRequest struct {
	UrlGroupName string     `json:"-"`
	Endpoints    []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type URLGroup struct {
	Name      string     `json:"name"`
	Endpoints []Endpoint `json:"endpoints"`
	CreatedAt int        `json:"createdAt"`
	UpdatedAt int        `json:"updatedAt"`
}

type URLGroupList []URLGroup

type RemoveEndpointsRequest struct {
	UrlGroupName string     `json:"-"`
	Endpoints    []Endpoint `json:"endpoints"`
}
