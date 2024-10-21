package vercel

import (
	"net/http"
	"net/url"
)

type VercelClient struct {
	token  string
	http   *httpClient
	common service

	AccessGroups   *AccessGroupService
	Aliases        *AliasService
	Artifacts      *ArtifactService
	Authentication *AuthService
	Certs          *CertService
	Checks         *CheckService
	Deployments    *DeploymentService
	DNS            *DNSService
	Domains        *DomainService
	EdgeConfig     *EdgeService
	Integrations   *IntegrationService
	LogDrains      *LogDrainService
	ProjectMembers *MemberService
	Projects       *ProjectService
	Secrets        *SecretService
	Teams          *TeamService
	User           *UserService
	Webhooks       *WebhookService
}

type service struct {
	client *VercelClient
}

const vercle_url string = "https://api.vercel.com"

func NewClient(token string) *VercelClient {
	c := &VercelClient{
		http: newClient(token, vercle_url),
	}
	c.token = token
	c.initialize()
	return c
}

func NewClientWithHeader(customeHeader map[string]string) *VercelClient {
	c := &VercelClient{
		http: newClientWithHeader(vercle_url, customeHeader),
	}
	c.initialize()
	return c
}

func NewClientWithHttp(client *http.Client) *VercelClient {
	c := &VercelClient{
		http: newClientWithHttp(vercle_url, client),
	}
	c.initialize()
	return c
}

func (c *VercelClient) initialize() {
	c.common.client = c
	c.AccessGroups = (*AccessGroupService)(&c.common)
	c.Aliases = (*AliasService)(&c.common)
	c.Artifacts = (*ArtifactService)(&c.common)
	c.Authentication = (*AuthService)(&c.common)
	c.Certs = (*CertService)(&c.common)
	c.Checks = (*CheckService)(&c.common)
	c.Deployments = (*DeploymentService)(&c.common)
	c.DNS = (*DNSService)(&c.common)
	c.Domains = (*DomainService)(&c.common)
	c.EdgeConfig = (*EdgeService)(&c.common)
	c.Integrations = (*IntegrationService)(&c.common)
	c.LogDrains = (*LogDrainService)(&c.common)
	c.ProjectMembers = (*MemberService)(&c.common)
	c.Projects = (*ProjectService)(&c.common)
	c.Secrets = (*SecretService)(&c.common)
	c.Teams = (*TeamService)(&c.common)
	c.User = (*UserService)(&c.common)
	c.Webhooks = (*WebhookService)(&c.common)
}

type VercelError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Pagination struct {
	Count int    `json:"count"`
	Next  string `json:"next"`
}

type BaseUrlParameter struct {
	TeamSlug *string `json:"slug"`
	TeamId   *string `json:"teamId"`
}

func (u *BaseUrlParameter) GetUrlValues() url.Values {
	params := url.Values{}
	if u.TeamId != nil {
		params.Add("teamId", *u.TeamId)
	}
	if u.TeamSlug != nil {
		params.Add("slug", *u.TeamSlug)
	}
	return params
}
