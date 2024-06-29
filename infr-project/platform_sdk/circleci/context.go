package circleci

import "log"

func (s *CircleciClient) CreateContext(name string, ownerId string, ownerType string) ContextInfo {
	path := "/context/"

	result := ContextInfo{}
	request := CreateContextInfo{
		Name: name,
		Owner: ContextOwner{
			Id:        ownerId,
			OwnerType: ownerType,
		},
	}
	err := s.http.Post(path, request, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
}

func (s *CircleciClient) ListContext(ownerId string) ListContextResponse {
	path := "/context?owner-id=" + ownerId

	result := ListContextResponse{}
	err := s.http.Get(path, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
}

func (s *CircleciClient) DeleteContext(contextId string) BaseResponse {
	path := "/context/" + contextId

	result := BaseResponse{}
	err := s.http.Delete(path, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
}

func (s *CircleciClient) GetContext(contextId string) ContextInfo {
	path := "/context/" + contextId

	result := ContextInfo{}
	err := s.http.Get(path, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
}

func (s *CircleciClient) ListContextEnvironment(contextId string) ListContextEnvResponse {
	path := "/context/" + contextId + "/environment-variable"

	result := ListContextEnvResponse{}
	err := s.http.Get(path, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
}

func (s *CircleciClient) UpsertContextEnvironment(contextId string, name string, value string) ContextEnvInfo {
	path := "/context/" + contextId + "/environment-variable/" + name
	request := UpsertContextEnv{
		Value: value,
	}
	result := ContextEnvInfo{}
	err := s.http.Put(path, request, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
}

func (s *CircleciClient) RemoveContextEnvironment(contextId string, name string) BaseResponse {
	path := "/context/" + contextId + "/environment-variable/" + name

	result := BaseResponse{}
	err := s.http.Delete(path, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
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
