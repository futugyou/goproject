package vercel

import (
	"net/url"
)

func (v *VercelClient) PurchaseDomain(slug string, teamId string, req PurchaseDomainRequest) (*PurchaseDomainResponse, error) {
	path := "/v5/domains/buy"
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
	result := &PurchaseDomainResponse{}
	err := v.http.Post(path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) CheckDomainPrice(slug string, teamId string, name string, domainType string) (*CheckDomainPriceResponse, error) {
	path := "/v4/domains/price"
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(domainType) > 0 {
		queryParams.Add("type", domainType)
	}
	if len(name) > 0 {
		queryParams.Add("name", name)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &CheckDomainPriceResponse{}
	err := v.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type PurchaseDomainRequest struct {
	Address1      string `json:"address1"`
	City          string `json:"city"`
	Country       string `json:"country"`
	Email         string `json:"email"`
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	PostalCode    string `json:"postalCode"`
	State         string `json:"state"`
	OrgName       string `json:"orgName"`
	ExpectedPrice int    `json:"expectedPrice"`
	Renew         bool   `json:"renew"`
}

type PurchaseDomainResponse struct {
	Created  int          `json:"created"`
	Ns       []string     `json:"ns"`
	Pending  bool         `json:"pending"`
	Uid      string       `json:"uid"`
	Verified bool         `json:"verified"`
	Error    *VercelError `json:"error,omitempty"`
}

type CheckDomainPriceResponse struct {
	Period int          `json:"period"`
	Price  int          `json:"price"`
	Error  *VercelError `json:"error,omitempty"`
}
