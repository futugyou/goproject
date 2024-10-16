package vercel

import (
	"fmt"
	"net/url"
)

type SecretService service

func (v *SecretService) CreateSecret(slug string, teamId string, req CreateSecretRequest) (*SecretInfo, error) {
	path := fmt.Sprintf("/v2/secrets/%s", req.Name)
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
	result := &SecretInfo{}
	err := v.client.http.Post(path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *SecretService) DeleteSecret(idOrName string, slug string, teamId string) (*DeleteSecretResponse, error) {
	path := fmt.Sprintf("/v2/secrets/%s", idOrName)
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
	result := &DeleteSecretResponse{}
	err := v.client.http.Delete(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *SecretService) GetSecret(idOrName string, slug string, teamId string, decrypt string) (*SecretInfo, error) {
	path := fmt.Sprintf("/v3/secrets/%s", idOrName)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(decrypt) > 0 {
		queryParams.Add("decrypt", decrypt)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &SecretInfo{}
	err := v.client.http.Delete(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *SecretService) ListSecret(slug string, teamId string) (*ListSecretResponse, error) {
	path := "/v3/secrets"
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
	result := &ListSecretResponse{}
	err := v.client.http.Delete(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *SecretService) ChangeSecretName(name string, slug string, teamId string, req CreateSecretRequest) (*string, error) {
	path := fmt.Sprintf("/v2/secrets/%s", name)
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
	err := v.client.http.Patch(path, req, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

type ListSecretResponse struct {
	Secrets    []SecretInfo `json:"secrets,omitempty"`
	Pagination Pagination   `json:"pagination,omitempty"`
	Error      *VercelError `json:"error,omitempty"`
}

type DeleteSecretResponse struct {
	Created string       `json:"created,omitempty"`
	Name    string       `json:"name,omitempty"`
	Uid     string       `json:"uid,omitempty"`
	Error   *VercelError `json:"error,omitempty"`
}

type CreateSecretRequest struct {
	Name        string `json:"name,omitempty"`
	Value       string `json:"value,omitempty"`
	Decryptable bool   `json:"decryptable,omitempty"`
}

type SecretInfo struct {
	Created     string       `json:"created,omitempty"`
	CreatedAt   int          `json:"createdAt,omitempty"`
	Decryptable bool         `json:"decryptable,omitempty"`
	Name        string       `json:"name,omitempty"`
	ProjectId   string       `json:"projectId,omitempty"`
	TeamId      string       `json:"teamId,omitempty"`
	Uid         string       `json:"uid,omitempty"`
	UserId      string       `json:"userId,omitempty"`
	Value       interface{}  `json:"value,omitempty"`
	Error       *VercelError `json:"error,omitempty"`
}

type SecretValue struct {
	Data []string `json:"data,omitempty"`
	Type string   `json:"type,omitempty"`
}
