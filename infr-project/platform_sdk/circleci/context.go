package circleci

func (s *CircleciClient) CreateContext(name string, ownerId string, ownerType string) (*ContextInfo, error) {
	path := "/context/"

	result := &ContextInfo{}
	request := CreateContextInfo{
		Name: name,
		Owner: ContextOwner{
			Id:        ownerId,
			OwnerType: ownerType,
		},
	}

	if err := s.http.Post(path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *CircleciClient) ListContext(ownerId string) (*ListContextResponse, error) {
	path := "/context?owner-id=" + ownerId

	result := &ListContextResponse{}
	if err := s.http.Get(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *CircleciClient) DeleteContext(contextId string) (*BaseResponse, error) {
	path := "/context/" + contextId

	result := &BaseResponse{}
	if err := s.http.Delete(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *CircleciClient) GetContext(contextId string) (*ContextInfo, error) {
	path := "/context/" + contextId

	result := &ContextInfo{}
	if err := s.http.Get(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CircleciClient) ListContextEnvironment(contextId string) (*ListContextEnvResponse, error) {
	path := "/context/" + contextId + "/environment-variable"

	result := &ListContextEnvResponse{}
	if err := s.http.Get(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CircleciClient) UpsertContextEnvironment(contextId string, name string, value string) (*ContextEnvInfo, error) {
	path := "/context/" + contextId + "/environment-variable/" + name
	request := UpsertContextEnv{
		Value: value,
	}
	result := &ContextEnvInfo{}
	if err := s.http.Put(path, request, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CircleciClient) RemoveContextEnvironment(contextId string, name string) (*BaseResponse, error) {
	path := "/context/" + contextId + "/environment-variable/" + name

	result := &BaseResponse{}
	if err := s.http.Delete(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CircleciClient) GetContextRestrictions(contextId string) (*ContextRestrictionsResponse, error) {
	path := "/context/" + contextId + "/restrictions"

	result := &ContextRestrictionsResponse{}
	if err := s.http.Get(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CircleciClient) CreateContextRestriction(contextId string, restriction_type string, restriction_value string) (*ContextRestriction, error) {
	path := "/context/" + contextId + "/restrictions"

	result := &ContextRestriction{}
	request := &struct {
		RestrictionType  string `json:"restriction_type"`
		RestrictionValue string `json:"restriction_value"`
	}{
		RestrictionType:  restriction_type,
		RestrictionValue: restriction_value,
	}
	if err := s.http.Post(path, request, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CircleciClient) DeleteContextRestriction(contextId string, restriction_id string) (*BaseResponse, error) {
	path := "/context/" + contextId + "/restrictions/" + restriction_id

	result := &BaseResponse{}
	if err := s.http.Delete(path, result); err != nil {
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
