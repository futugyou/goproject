package vercel

import (
	"fmt"
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

type CertInfo struct {
	AutoRenew bool         `json:"autoRenew"`
	Cns       []string     `json:"cns"`
	Id        string       `json:"id"`
	CreatedAt int          `json:"createdAt"`
	ExpiresAt int          `json:"expiresAt"`
	Error     *VercelError `json:"error"`
}
