package vercel

import (
	"fmt"
	"net/url"
)

type IntegrationService service

func (v *IntegrationService) DeleteIntegrationConfiguration(id string, slug string, teamId string) (*string, error) {
	path := fmt.Sprintf("/v1/integrations/configuration/%s", id)
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

func (v *IntegrationService) RetrieveIntegrationConfiguration(id string, slug string, teamId string) (*IntegrationConfigurationInfo, error) {
	path := fmt.Sprintf("/v1/integrations/configuration/%s", id)
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
	result := &IntegrationConfigurationInfo{}
	err := v.client.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *IntegrationService) GetConfiguration(view string, slug string, teamId string) (*string, error) {
	path := "/v1/integrations/configurations"
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(view) > 0 {
		queryParams.Add("view", view)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := ""
	err := v.client.http.Get(path, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *IntegrationService) ListGitByProvider(host string, provider string, slug string, teamId string) ([]ListGitByProviderResponse, error) {
	path := "/v1/integrations/git-namespaces"
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(host) > 0 {
		queryParams.Add("host", host)
	}
	if len(provider) > 0 {
		queryParams.Add("provider", provider)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := []ListGitByProviderResponse{}
	err := v.client.http.Get(path, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *IntegrationService) ListGitLinkedByProvider(host string, installationId string, namespaceId string, query string,
	provider string, slug string, teamId string) (*ListGitLinkedByProviderResponse, error) {
	path := "/v1/integrations/search-repo"
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(host) > 0 {
		queryParams.Add("host", host)
	}
	if len(provider) > 0 {
		queryParams.Add("provider", provider)
	}
	if len(installationId) > 0 {
		queryParams.Add("installationId", installationId)
	}
	if len(namespaceId) > 0 {
		queryParams.Add("namespaceId", namespaceId)
	}
	if len(query) > 0 {
		queryParams.Add("query", query)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &ListGitLinkedByProviderResponse{}
	err := v.client.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type ListGitLinkedByProviderResponse struct {
	GitAccount GitAccount   `json:"gitAccount"`
	Repos      []Repos      `json:"repos"`
	Error      *VercelError `json:"error,omitempty"`
}

type GitAccount struct {
	NamespaceId string `json:"namespaceId"`
	Provider    string `json:"provider"`
}

type Repos struct {
	Id            string `json:"id"`
	DefaultBranch string `json:"defaultBranch"`
	Name          string `json:"name"`
	Namespace     string `json:"namespace"`
	OwnerType     string `json:"ownerType"`
	Provider      string `json:"provider"`
	Slug          string `json:"slug"`
	Private       bool   `json:"private"`
	UpdatedAt     int    `json:"updatedAt"`
	Url           string `json:"url"`
	Owner         Owner  `json:"owner"`
}

type Owner struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ListGitByProviderResponse struct {
	Id                 string       `json:"id"`
	IsAccessRestricted bool         `json:"isAccessRestricted"`
	InstallationId     int          `json:"installationId"`
	Name               string       `json:"name"`
	OwnerType          string       `json:"ownerType"`
	Provider           string       `json:"provider"`
	RequireReauth      bool         `json:"requireReauth"`
	Slug               string       `json:"slug"`
	Error              *VercelError `json:"error,omitempty"`
}

type IntegrationConfigurationInfo struct {
	BillingTotal              string                     `json:"billingTotal,omitempty"`
	CompletedAt               int                        `json:"completedAt"`
	CreatedAt                 int                        `json:"createdAt"`
	DeletedAt                 int                        `json:"deletedAt"`
	DisabledAt                int                        `json:"disabledAt"`
	DisabledReason            string                     `json:"disabledReason"`
	Id                        string                     `json:"id"`
	InstallationType          string                     `json:"installationType"`
	IntegrationId             string                     `json:"integrationId"`
	NorthstarMigratedAt       int                        `json:"northstarMigratedAt"`
	OwnerId                   string                     `json:"ownerId"`
	PeriodEnd                 string                     `json:"periodEnd"`
	PeriodStart               string                     `json:"periodStart"`
	Projects                  []string                   `json:"projects"`
	RemovedLogDrainsAt        int                        `json:"removedLogDrainsAt"`
	RemovedProjectEnvsAt      int                        `json:"removedProjectEnvsAt"`
	RemovedTokensAt           int                        `json:"removedTokensAt"`
	RemovedWebhooksAt         int                        `json:"removedWebhooksAt"`
	Scopes                    []string                   `json:"scopes"`
	ScopesQueue               IntegrationConfScopesQueue `json:"scopesQueue"`
	Slug                      string                     `json:"slug"`
	Source                    string                     `json:"source"`
	TeamId                    string                     `json:"teamId"`
	Type                      string                     `json:"type"`
	UpdatedAt                 int                        `json:"updatedAt"`
	UserId                    string                     `json:"userId"`
	CanConfigureOpenTelemetry bool                       `json:"canConfigureOpenTelemetry,omitempty"`
	ProjectSelection          string                     `json:"projectSelection,omitempty"`
	Error                     *VercelError               `json:"error,omitempty"`
}

type IntegrationConfScopesQueue struct {
	ConfirmedAt int                    `json:"confirmedAt"`
	Note        string                 `json:"note"`
	RequestedAt int                    `json:"requestedAt"`
	Scopes      []IntegrationConfScope `json:"scopes"`
}

type IntegrationConfScope struct {
	Added    []string `json:"added"`
	Upgraded []string `json:"upgraded"`
}
