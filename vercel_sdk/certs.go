package vercel

import (
	"fmt"
	"net/url"
)

func (v *VercelClient) GetCertById(id string) (*CertInfo, error) {
	path := fmt.Sprintf("/v7/certs/%s", id)
	result := &CertInfo{}
	err := v.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) CreateCert(slug string, teamId string, cns []string) (*CertInfo, error) {
	path := "/v7/certs"
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
	request := struct {
		Cns []string `json:"cns"`
	}{
		Cns: cns,
	}
	result := &CertInfo{}
	err := v.http.Post(path, request, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) RemoveCert(slug string, teamId string, id string) (*string, error) {
	path := fmt.Sprintf("/v7/certs/%s", id)
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
	err := v.http.Delete(path, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *VercelClient) UploadCert(slug string, teamId string, req UploadCertRequest) (*CertInfo, error) {
	path := "/v7/certs"
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
	result := &CertInfo{}
	err := v.http.Put(path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
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
	Ca             string `json:"ca"`
	Cert           string `json:"cert"`
	Key            string `json:"key"`
	SkipValidation bool   `json:"skipValidation"`
}
