package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type CertService service

type GetCertParameter struct {
	Id               string
	BaseUrlParameter `json:"-"`
}

func (v *CertService) GetCertById(ctx context.Context, request GetCertParameter) (*CertInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v7/certs/%s", request.Id),
	}
	params := url.Values{}
	if request.TeamId != nil {
		params.Add("teamId", *request.TeamId)
	}
	if request.TeamSlug != nil {
		params.Add("slug", *request.TeamSlug)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &CertInfo{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type CreateCertRequest struct {
	Cns              []string `json:"cns"`
	BaseUrlParameter `json:"-"`
}

func (v *CertService) CreateCert(ctx context.Context, request CreateCertRequest) (*CertInfo, error) {
	u := &url.URL{
		Path: "/v7/certs",
	}
	params := url.Values{}
	if request.TeamId != nil {
		params.Add("teamId", *request.TeamId)
	}
	if request.TeamSlug != nil {
		params.Add("slug", *request.TeamSlug)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &CertInfo{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type RemoveCertRequest struct {
	Id               string
	BaseUrlParameter `json:"-"`
}

func (v *CertService) RemoveCert(ctx context.Context, request RemoveCertRequest) (*string, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v7/certs/%s", request.Id),
	}
	params := url.Values{}
	if request.TeamId != nil {
		params.Add("teamId", *request.TeamId)
	}
	if request.TeamSlug != nil {
		params.Add("slug", *request.TeamSlug)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := ""
	if err := v.client.http.Delete(ctx, path, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (v *CertService) UploadCert(ctx context.Context, request UploadCertRequest) (*CertInfo, error) {
	u := &url.URL{
		Path: "/v7/certs",
	}
	params := url.Values{}
	if request.TeamId != nil {
		params.Add("teamId", *request.TeamId)
	}
	if request.TeamSlug != nil {
		params.Add("slug", *request.TeamSlug)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &CertInfo{}
	if err := v.client.http.Put(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type CertInfo struct {
	AutoRenew bool         `json:"autoRenew"`
	Cns       []string     `json:"cns"`
	Id        string       `json:"id"`
	CreatedAt int          `json:"createdAt"`
	ExpiresAt int          `json:"expiresAt"`
	Error     *VercelError `json:"error"`
}

type UploadCertRequest struct {
	Ca               string `json:"ca"`
	Cert             string `json:"cert"`
	Key              string `json:"key"`
	SkipValidation   bool   `json:"skipValidation"`
	BaseUrlParameter `json:"-"`
}
