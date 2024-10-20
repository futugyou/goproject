package vercel

import "net/http"

type VercelClient struct {
	token  string
	http   *httpClient
	common service

	AccessGroups *AccessGroupService
	Alias        *AliasService
	Artifact     *ArtifactService
	Auth         *AuthService
	Cert         *CertService
	Check        *CheckService
	Deployment   *DeploymentService
	DNS          *DNSService
	Domain       *DomainService
	Edge         *EdgeService
	Integration  *IntegrationService
	LogDrain     *LogDrainService
	Member       *MemberService
	Project      *ProjectService
	Secret       *SecretService
	Team         *TeamService
	User         *UserService
	Webhook      *WebhookService
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
	c.Alias = (*AliasService)(&c.common)
	c.Artifact = (*ArtifactService)(&c.common)
	c.Auth = (*AuthService)(&c.common)
	c.Cert = (*CertService)(&c.common)
	c.Check = (*CheckService)(&c.common)
	c.Deployment = (*DeploymentService)(&c.common)
	c.DNS = (*DNSService)(&c.common)
	c.Domain = (*DomainService)(&c.common)
	c.Edge = (*EdgeService)(&c.common)
	c.Integration = (*IntegrationService)(&c.common)
	c.LogDrain = (*LogDrainService)(&c.common)
	c.Member = (*MemberService)(&c.common)
	c.Project = (*ProjectService)(&c.common)
	c.Secret = (*SecretService)(&c.common)
	c.Team = (*TeamService)(&c.common)
	c.User = (*UserService)(&c.common)
	c.Webhook = (*WebhookService)(&c.common)
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
