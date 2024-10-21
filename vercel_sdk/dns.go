package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type DNSService service

type UpsertDNSRecordRequest struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	Ttl              int    `json:"ttl"`
	Value            string `json:"value,omitempty"`
	Comment          string `json:"comment"`
	MxPriority       int    `json:"mxPriority,omitempty"`
	SRV              *SRV   `json:"srv,omitempty"`
	Https            *Https `json:"https,omitempty"`
	DomainOrRecordId string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *DNSService) CreateDNSRecord(ctx context.Context, request UpsertDNSRecordRequest) (*CreateDNSRecordResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v2/domains/%s/records", request.DomainOrRecordId),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &CreateDNSRecordResponse{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type ListDNSRecordsParameter struct {
	Domain           string  `json:"-"`
	Limit            *string `json:"-"`
	Since            *string `json:"-"`
	Until            *string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *DNSService) ListDNSRecords(ctx context.Context, request ListDNSRecordsParameter) (*ListDNSRecordResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v4/domains/%s/records", request.Domain),
	}
	params := request.GetUrlValues()
	if request.Limit != nil {
		params.Add("limit", *request.Limit)
	}
	if request.Since != nil {
		params.Add("since", *request.Since)
	}
	if request.Until != nil {
		params.Add("until", *request.Until)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ListDNSRecordResponse{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type DeleteDNSRecordsRequest struct {
	Domain           string `json:"-"`
	RecordId         string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *DNSService) DeleteDNSRecords(ctx context.Context, request DeleteDNSRecordsRequest) (*string, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v2/domains/%s/records/%s", request.Domain, request.RecordId),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := ""
	if err := v.client.http.Delete(ctx, path, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (v *DNSService) UpdateDNSRecords(ctx context.Context, request UpsertDNSRecordRequest) (*UpdateDNSRecordResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/domains/records/%s", request.DomainOrRecordId),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &UpdateDNSRecordResponse{}
	if err := v.client.http.Patch(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
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
