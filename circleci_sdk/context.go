package circleci

import (
	"context"
	"fmt"
)

type ContextService service

func (s *ContextService) CreateContext(ctx context.Context, name string, ownerId string, ownerType string) (*ContextInfo, error) {
	path := "/context/"

	result := &ContextInfo{}
	request := CreateContextInfo{
		Name: name,
		Owner: ContextOwner{
			Id:        ownerId,
			OwnerType: ownerType,
		},
	}

	if err := s.client.http.Post(ctx, path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *ContextService) ListContext(ctx context.Context, ownerId string) (*ListContextResponse, error) {
	path := fmt.Sprintf("/context?owner-id=%s", ownerId)

	result := &ListContextResponse{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *ContextService) DeleteContext(ctx context.Context, contextId string) (*BaseResponse, error) {
	path := fmt.Sprintf("/context/%s", contextId)

	result := &BaseResponse{}
	if err := s.client.http.Delete(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *ContextService) GetContext(ctx context.Context, contextId string) (*ContextInfo, error) {
	path := fmt.Sprintf("/context/%s", contextId)

	result := &ContextInfo{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ContextService) ListContextEnvironment(ctx context.Context, contextId string) (*ListContextEnvResponse, error) {
	path := fmt.Sprintf("/context/%s/environment-variable", contextId)

	result := &ListContextEnvResponse{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ContextService) UpsertContextEnvironment(ctx context.Context, contextId string, name string, value string) (*ContextEnvInfo, error) {
	path := fmt.Sprintf("/context/%s/environment-variable/%s", contextId, name)
	request := UpsertContextEnv{
		Value: value,
	}
	result := &ContextEnvInfo{}
	if err := s.client.http.Put(ctx, path, request, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ContextService) RemoveContextEnvironment(ctx context.Context, contextId string, name string) (*BaseResponse, error) {
	path := fmt.Sprintf("/context/%s/environment-variable/%s", contextId, name)

	result := &BaseResponse{}
	if err := s.client.http.Delete(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ContextService) GetContextRestrictions(ctx context.Context, contextId string) (*ContextRestrictionsResponse, error) {
	path := fmt.Sprintf("/context/%s/restrictions", contextId)

	result := &ContextRestrictionsResponse{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ContextService) CreateContextRestriction(ctx context.Context, contextId string, restriction_type string, restriction_value string) (*ContextRestriction, error) {
	path := fmt.Sprintf("/context/%s/restrictions", contextId)

	result := &ContextRestriction{}
	request := &struct {
		RestrictionType  string `json:"restriction_type"`
		RestrictionValue string `json:"restriction_value"`
	}{
		RestrictionType:  restriction_type,
		RestrictionValue: restriction_value,
	}
	if err := s.client.http.Post(ctx, path, request, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ContextService) DeleteContextRestriction(ctx context.Context, contextId string, restriction_id string) (*BaseResponse, error) {
	path := fmt.Sprintf("/context/%s/restrictions/%s", contextId, restriction_id)

	result := &BaseResponse{}
	if err := s.client.http.Delete(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

type ContextRestrictionsResponse struct {
	Items         []ContextRestriction `json:"items"`
	NextPageToken string               `json:"next_page_token"`
	Message       *string              `json:"message,omitempty"`
}

type ContextRestriction struct {
	ContextID        string  `json:"context_id"`
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	RestrictionType  string  `json:"restriction_type"` //"project" "expression"
	RestrictionValue string  `json:"restriction_value"`
	Message          *string `json:"message,omitempty"`
}

type UpsertContextEnv struct {
	Value string `json:"value"`
}

type ListContextEnvResponse struct {
	Items         []ContextEnvInfo `json:"items"`
	NextPageToken string           `json:"next_page_token"`
	Message       *string          `json:"message,omitempty"`
}

type ContextEnvInfo struct {
	Variable  string  `json:"variable"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	ContextID string  `json:"context_id"`
	Message   *string `json:"message,omitempty"`
}

type ListContextResponse struct {
	Items         []ContextInfo `json:"items"`
	NextPageToken string        `json:"next_page_token"`
	Message       *string       `json:"message,omitempty"`
}

type CreateContextInfo struct {
	Name  string       `json:"name"`
	Owner ContextOwner `json:"owner"`
}

type ContextOwner struct {
	Id        string `json:"id"`
	OwnerType string `json:"type"` //"account" "organization"
}

type ContextInfo struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	CreatedAt string  `json:"created_at"`
	Message   *string `json:"message,omitempty"`
}
