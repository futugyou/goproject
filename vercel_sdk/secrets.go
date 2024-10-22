package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type SecretService service

type CreateSecretRequest struct {
	Name             string `json:"name,omitempty"`
	Value            string `json:"value,omitempty"`
	Decryptable      bool   `json:"decryptable,omitempty"`
	BaseUrlParameter `json:"-"`
}

func (v *SecretService) CreateSecret(ctx context.Context, request CreateSecretRequest) (*SecretInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v2/secrets/%s", request.Name),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &SecretInfo{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type DeleteSecretRequest struct {
	IdOrName         string
	BaseUrlParameter `json:"-"`
}

func (v *SecretService) DeleteSecret(ctx context.Context, request DeleteSecretRequest) (*DeleteSecretResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v2/secrets/%s", request.IdOrName),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &DeleteSecretResponse{}
	if err := v.client.http.Delete(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type GetSecretParameter struct {
	IdOrName         string
	Decrypt          *string
	BaseUrlParameter `json:"-"`
}

func (v *SecretService) GetSecret(ctx context.Context, request GetSecretParameter) (*SecretInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v3/secrets/%s", request.IdOrName),
	}
	params := request.GetUrlValues()
	if request.Decrypt != nil {
		params.Add("decrypt", *request.Decrypt)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &SecretInfo{}
	if err := v.client.http.Delete(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type ListSecretParameter struct {
	BaseUrlParameter `json:"-"`
}

func (v *SecretService) ListSecret(ctx context.Context, request ListSecretParameter) (*ListSecretResponse, error) {
	u := &url.URL{
		Path: "/v3/secrets",
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ListSecretResponse{}
	if err := v.client.http.Delete(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type ChangeSecretNameRequest struct {
	OldName          string `json:"-"`
	NewName          string `json:"name"`
	BaseUrlParameter `json:"-"`
}

type ChangeSecretNameResponse struct {
	Created string       `json:"created"`
	Name    string       `json:"name"`
	OldName string       `json:"oldName"`
	Uid     string       `json:"uid"`
	Error   *VercelError `json:"error,omitempty"`
}

func (v *SecretService) ChangeSecretName(ctx context.Context, request ChangeSecretNameRequest) (*ChangeSecretNameResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v2/secrets/%s", request.OldName),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ChangeSecretNameResponse{}
	if err := v.client.http.Patch(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
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
