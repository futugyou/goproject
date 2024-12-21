package qstash

import (
	"context"
	"fmt"
)

type URLGroups service

func (s *URLGroups) UpsertURLGroupAndEndpoint(ctx context.Context, request UpsertURLGroupRequest) error {
	path := fmt.Sprintf("/topics/%s/endpoints", request.UrlGroupName)
	result := ""
	return s.client.http.Post(ctx, path, request, &result)
}

func (s *URLGroups) GetURLGroup(ctx context.Context, urlGroupName string) (*URLGroup, error) {
	path := fmt.Sprintf("/topics/%s", urlGroupName)
	result := &URLGroup{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *URLGroups) ListURLGroup(ctx context.Context) (*ListURLGroup, error) {
	path := "/topics"
	result := &ListURLGroup{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
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

type ListURLGroup []URLGroup
