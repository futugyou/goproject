package vercel

import (
	"fmt"
	"net/url"
)

type DNSService service

func (v *DNSService) CreateDNSRecord(domain string, slug string, teamId string, req UpsertDNSRecordRequest) (*CreateDNSRecordResponse, error) {
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
	err := v.client.http.Post(path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *DNSService) ListDNSRecords(domain string, slug string, teamId string, limit string, since string, until string) (*ListDNSRecordResponse, error) {
	path := fmt.Sprintf("/v4/domains/%s/records", domain)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(limit) > 0 {
		queryParams.Add("limit", limit)
	}
	if len(since) > 0 {
		queryParams.Add("since", since)
	}
	if len(until) > 0 {
		queryParams.Add("until", until)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &ListDNSRecordResponse{}
	err := v.client.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *DNSService) DeleteDNSRecords(domain string, recordId string, slug string, teamId string) (*string, error) {
	path := fmt.Sprintf("/v2/domains/%s/records/%s", domain, recordId)
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
	err := v.client.http.Delete(path, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *DNSService) UpdatesDNSRecords(recordId string, slug string, teamId string, req UpsertDNSRecordRequest) (*UpdateDNSRecordResponse, error) {
	path := fmt.Sprintf("/v1/domains/records/%s", recordId)
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
	result := &UpdateDNSRecordResponse{}
	err := v.client.http.Patch(path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type UpdateDNSRecordResponse struct {
	Comment    string       `json:"comment"`
	Created    int          `json:"created"`
	CreatedAt  int          `json:"createdAt"`
	Domain     string       `json:"domain"`
	Name       string       `json:"name"`
	Type       string       `json:"type"`
	Ttl        int          `json:"ttl"`
	Value      string       `json:"value,omitempty"`
	Id         string       `json:"id"`
	RecordType string       `json:"recordType"`
	Error      *VercelError `json:"error,omitempty"`
}

type ListDNSRecordResponse struct {
	Records    []DNSRecordInfo `json:"records"`
	Pagination Pagination      `json:"pagination,omitempty"`
	Error      *VercelError    `json:"error,omitempty"`
}

type CreateDNSRecordResponse struct {
	Uid     string       `json:"uid"`
	Updated int          `json:"updated,omitempty"`
	Error   *VercelError `json:"error,omitempty"`
}

type UpsertDNSRecordRequest struct {
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

type DNSRecordInfo struct {
	Created    int    `json:"created"`
	CreatedAt  int    `json:"createdAt"`
	Creator    int    `json:"creator"`
	Id         int    `json:"id"`
	MxPriority int    `json:"mxPriority,omitempty"`
	Priority   int    `json:"priority,omitempty"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Value      string `json:"value,omitempty"`
	Slug       string `json:"slug"`
	Updated    int    `json:"updated"`
	UpdatedAt  int    `json:"updatedAt"`
}
