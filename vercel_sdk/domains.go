package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type DomainService service

type PurchaseDomainRequest struct {
	Address1         string `json:"address1"`
	City             string `json:"city"`
	Country          string `json:"country"`
	Email            string `json:"email"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Name             string `json:"name"`
	Phone            string `json:"phone"`
	PostalCode       string `json:"postalCode"`
	State            string `json:"state"`
	OrgName          string `json:"orgName"`
	ExpectedPrice    int    `json:"expectedPrice"`
	Renew            bool   `json:"renew"`
	BaseUrlParameter `json:"-"`
}

func (v *DomainService) PurchaseDomain(ctx context.Context, request PurchaseDomainRequest) (*PurchaseDomainResponse, error) {
	u := &url.URL{
		Path: "/v5/domains/buy",
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &PurchaseDomainResponse{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type CheckDomainPriceParameter struct {
	Name             string
	Type             *string // allowed value: new	renewal	transfer
	BaseUrlParameter `json:"-"`
}

func (v *DomainService) CheckDomainPrice(ctx context.Context, request CheckDomainPriceParameter) (*CheckDomainPriceResponse, error) {
	u := &url.URL{
		Path: "/v4/domains/price",
	}
	params := request.GetUrlValues()
	params.Add("name", request.Name)
	if request.Type != nil {
		params.Add("type", *request.Type)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &CheckDomainPriceResponse{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type CheckDomainAvailabilityParameter struct {
	Name             string
	BaseUrlParameter `json:"-"`
}

func (v *DomainService) CheckDomainAvailability(ctx context.Context, request CheckDomainAvailabilityParameter) (*CheckDomainAvailabilityResponse, error) {
	u := &url.URL{
		Path: "/v4/domains/status",
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &CheckDomainAvailabilityResponse{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type RegisterDomainRequest struct {
	Name             string `json:"name"`
	CdnEnabled       bool   `json:"cdnEnabled,omitempty"`
	Zone             bool   `json:"zone,omitempty"`
	Mthod            string `json:"method"`
	Token            string `json:"token,omitempty"`
	AuthCode         string `json:"authCode,omitempty"`
	ExpectedPrice    int    `json:"expectedPrice,omitempty"`
	BaseUrlParameter `json:"-"`
}

func (v *DomainService) RegisterDomain(ctx context.Context, request RegisterDomainRequest) (*RegisterDomainResponse, error) {
	u := &url.URL{
		Path: "/v5/domains",
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &RegisterDomainResponse{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type RemoveDomainRequest struct {
	Domain           string
	BaseUrlParameter `json:"-"`
}

func (v *DomainService) RemoveDomain(ctx context.Context, request RemoveDomainRequest) (*RemoveDomainResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v6/domains/%s", request.Domain),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &RemoveDomainResponse{}
	if err := v.client.http.Delete(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type GetDomainParameter struct {
	Domain           string
	BaseUrlParameter `json:"-"`
}

func (v *DomainService) GetDomain(ctx context.Context, request GetDomainParameter) (*DomainInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v5/domains/%s", request.Domain),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &DomainInfo{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type GetDomainConfigurationParameter struct {
	Domain           string
	Strict           *string //Allowed values:	true	false
	BaseUrlParameter `json:"-"`
}

func (v *DomainService) GetDomainConfiguration(ctx context.Context, request GetDomainConfigurationParameter) (*DomainConfiguration, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v6/domains/%s/config", request.Domain),
	}
	params := request.GetUrlValues()
	if request.Strict != nil {
		params.Add("strict", *request.Strict)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &DomainConfiguration{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type GetDomainTransferParameter struct {
	Domain           string
	BaseUrlParameter `json:"-"`
}

func (v *DomainService) GetDomainTransfer(ctx context.Context, request GetDomainTransferParameter) (*DomainTransfer, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/domains/%s/registry", request.Domain),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &DomainTransfer{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type ListDomainsParameter struct {
	Limit            *string
	Since            *string
	Until            *string
	BaseUrlParameter `json:"-"`
}

func (v *DomainService) ListDomains(ctx context.Context, request ListDomainsParameter) (*ListDomainsResponse, error) {
	u := &url.URL{
		Path: "/v5/domains",
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

	response := &ListDomainsResponse{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type UpdateDomainRequest struct {
	Op                string   `json:"op"`
	CustomNameservers []string `json:"customNameservers,omitempty"`
	Renew             bool     `json:"renew,omitempty"`
	Zone              bool     `json:"zone,omitempty"`
	Destination       string   `json:"destination,omitempty"`
	Domain            string   `json:"-"`
	BaseUrlParameter  `json:"-"`
}

func (v *DomainService) UpdateDomain(ctx context.Context, request UpdateDomainRequest) (*UpdateDomainResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v3/domains/%s", request.Domain),
	}

	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &UpdateDomainResponse{}
	if err := v.client.http.Patch(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type UpdateDomainResponse struct {
	Moved             bool         `json:"moved"`
	Token             string       `json:"token,omitempty"`
	CustomNameservers []string     `json:"customNameservers,omitempty"`
	Renew             bool         `json:"renew,omitempty"`
	Zone              bool         `json:"zone,omitempty"`
	Error             *VercelError `json:"error,omitempty"`
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
