package vercel

import (
	"fmt"
	"net/url"
)

func (v *VercelClient) CreateDNSRecord(domain string, slug string, teamId string, req CreateDNSRecordRequest) (*CreateDNSRecordResponse, error) {
	path := fmt.Sprintf("/v2/domains/%s/records", domain)
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
	result := &CreateDNSRecordResponse{}
	err := v.http.Post(path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type CreateDNSRecordResponse struct {
	Uid     string       `json:"uid"`
	Updated int          `json:"updated,omitempty"`
	Error   *VercelError `json:"error,omitempty"`
}

type CreateDNSRecordRequest struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Ttl        int    `json:"ttl"`
	Value      string `json:"value,omitempty"`
	Comment    string `json:"comment"`
	MxPriority int    `json:"mxPriority,omitempty"`
	SRV        *SRV   `json:"srv,omitempty"`
	Https      *Https `json:"https,omitempty"`
}

type SRV struct {
	Port     int    `json:"port,omitempty"`
	Priority int    `json:"priority,omitempty"`
	Target   string `json:"target,omitempty"`
	Weight   int    `json:"weight,omitempty"`
}

type Https struct {
	Params   string `json:"params,omitempty"`
	Priority int    `json:"priority,omitempty"`
	Target   string `json:"target,omitempty"`
}
