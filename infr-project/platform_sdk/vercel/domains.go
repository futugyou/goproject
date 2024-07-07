package vercel

import (
	"fmt"
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

func (v *VercelClient) CheckDomainAvailability(slug string, teamId string, name string) (*CheckDomainAvailabilityResponse, error) {
	path := "/v4/domains/status"
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(name) > 0 {
		queryParams.Add("name", name)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &CheckDomainAvailabilityResponse{}
	err := v.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) RegisterDomain(slug string, teamId string, req RegisterDomainRequest) (*RegisterDomainResponse, error) {
	path := "/v5/domains"
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
	result := &RegisterDomainResponse{}
	err := v.http.Post(path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) RemoveDomain(slug string, teamId string, domain string) (*RemoveDomainResponse, error) {
	path := fmt.Sprintf("/v6/domains/%s", domain)
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
	result := &RemoveDomainResponse{}
	err := v.http.Delete(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) GetDomain(slug string, teamId string, domain string) (*DomainInfo, error) {
	path := fmt.Sprintf("/v5/domains/%s", domain)
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
	result := &DomainInfo{}
	err := v.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) GetDomainConfiguration(slug string, teamId string, domain string, strict string) (*DomainConfiguration, error) {
	path := fmt.Sprintf("/v6/domains/%s/config", domain)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(strict) > 0 {
		queryParams.Add("strict", strict)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &DomainConfiguration{}
	err := v.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) GetDomainTransfer(slug string, teamId string, domain string) (*DomainTransfer, error) {
	path := fmt.Sprintf("/v1/domains/%s/registry", domain)
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
	result := &DomainTransfer{}
	err := v.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) ListDomains(slug string, teamId string, limit string, since string, until string) (*ListDomainsResponse, error) {
	path := "/v5/domains"
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
	result := &ListDomainsResponse{}
	err := v.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type ListDomainsResponse struct {
	Domains    []DomainInfo `json:"domains"`
	Pagination Pagination   `json:"pagination,omitempty"`
	Error      *VercelError `json:"error,omitempty"`
}

type DomainTransfer struct {
	Reason         string       `json:"reason"`
	Status         string       `json:"status"`
	TransferPolicy string       `json:"transferPolicy"`
	Transferable   bool         `json:"transferable"`
	Error          *VercelError `json:"error,omitempty"`
}

type DomainConfiguration struct {
	AcceptedChallenges []string     `json:"acceptedChallenges"`
	ConfiguredBy       string       `json:"configuredBy"`
	Misconfigured      bool         `json:"misconfigured"`
	Error              *VercelError `json:"error,omitempty"`
}

type RemoveDomainResponse struct {
	Uid   string       `json:"uid"`
	Error *VercelError `json:"error,omitempty"`
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

type CheckDomainAvailabilityResponse struct {
	Available bool         `json:"available"`
	Error     *VercelError `json:"error,omitempty"`
}

type RegisterDomainRequest struct {
	Name          string `json:"name"`
	CdnEnabled    bool   `json:"cdnEnabled,omitempty"`
	Zone          bool   `json:"zone,omitempty"`
	Mthod         string `json:"method"`
	Token         string `json:"token,omitempty"`
	AuthCode      string `json:"authCode,omitempty"`
	ExpectedPrice int    `json:"expectedPrice,omitempty"`
}

type RegisterDomainResponse struct {
	Domain []DomainInfo `json:"domain"`
	Error  *VercelError `json:"error,omitempty"`
}

type DomainInfo struct {
	BoughtAt            int           `json:"boughtAt"`
	CreatedAt           int           `json:"createdAt"`
	Creator             DomainCreator `json:"creator"`
	CustomNameservers   []string      `json:"customNameservers"`
	ExpiresAt           int           `json:"expiresAt"`
	Id                  string        `json:"id"`
	IntendedNameservers []string      `json:"intendedNameservers"`
	Name                string        `json:"name"`
	Nameservers         []string      `json:"nameservers"`
	OrderedAt           int           `json:"orderedAt"`
	Renew               bool          `json:"renew"`
	ServiceType         string        `json:"serviceType"`
	TeamId              string        `json:"teamId"`
	TransferStartedAt   int           `json:"transferStartedAt"`
	TransferredAt       int           `json:"transferredAt"`
	UserId              string        `json:"userId"`
	Verified            bool          `json:"verified"`
	Error               *VercelError  `json:"error,omitempty"`
}

type DomainCreator struct {
	CustomerId       string `json:"CustomerId"`
	Email            string `json:"email"`
	Id               string `json:"id"`
	IsDomainReseller bool   `json:"isDomainReseller"`
	Username         string `json:"username"`
}
