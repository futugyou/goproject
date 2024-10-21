package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type IntegrationService service

type DeleteIntegrationConfigurationRequest struct {
	Id               string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *IntegrationService) DeleteIntegrationConfiguration(ctx context.Context, request DeleteIntegrationConfigurationRequest) (*string, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/integrations/configuration/%s", request.Id),
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

type RetrieveIntegrationConfigurationParameter struct {
	Id               string `json:"-"`
	BaseUrlParameter `json:"-"`
}

func (v *IntegrationService) RetrieveIntegrationConfiguration(ctx context.Context, request RetrieveIntegrationConfigurationParameter) (*IntegrationConfigurationInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/integrations/configuration/%s", request.Id),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &IntegrationConfigurationInfo{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type GetConfigurationParameter struct {
	View                *string `json:"-"`
	InstallationType    *string `json:"-"`
	IntegrationIdOrSlug *string `json:"-"`
	BaseUrlParameter    `json:"-"`
}

func (v *IntegrationService) GetConfiguration(ctx context.Context, request GetConfigurationParameter) (*string, error) {
	u := &url.URL{
		Path: "/v1/integrations/configurations",
	}
	params := request.GetUrlValues()
	if request.View != nil {
		params.Add("view", *request.View)
	}
	if request.InstallationType != nil {
		params.Add("installationType", *request.InstallationType)
	}
	if request.IntegrationIdOrSlug != nil {
		params.Add("integrationIdOrSlug", *request.IntegrationIdOrSlug)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := ""
	if err := v.client.http.Get(ctx, path, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

type ListGitByProviderParameter struct {
	Host     *string
	Provider *string
}

func (v *IntegrationService) ListGitByProvider(ctx context.Context, request ListGitByProviderParameter) ([]ListGitByProviderResponse, error) {
	u := &url.URL{
		Path: "/v1/integrations/git-namespaces",
	}
	params := url.Values{}
	if request.Host != nil {
		params.Add("host", *request.Host)
	}
	if request.Provider != nil {
		params.Add("provider", *request.Provider)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := []ListGitByProviderResponse{}
	if err := v.client.http.Get(ctx, path, &response); err != nil {
		return nil, err
	}
	return response, nil
}

type ListGitLinkedByProviderParameter struct {
	Host             *string
	Provider         *string //Allowed values:	github	github-custom-host	gitlab	bitbucket
	InstallationId   *string
	NamespaceId      *string
	Query            *string
	BaseUrlParameter `json:"-"`
}

func (v *IntegrationService) ListGitLinkedByProvider(ctx context.Context, request ListGitLinkedByProviderParameter) (*ListGitLinkedByProviderResponse, error) {
	u := &url.URL{
		Path: "/v1/integrations/search-repo",
	}
	params := request.GetUrlValues()
	if request.Host != nil {
		params.Add("host", *request.Host)
	}
	if request.Provider != nil {
		params.Add("provider", *request.Provider)
	}
	if request.InstallationId != nil {
		params.Add("installationId", *request.InstallationId)
	}
	if request.NamespaceId != nil {
		params.Add("namespaceId", *request.NamespaceId)
	}
	if request.Query != nil {
		params.Add("query", *request.Query)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := &ListGitLinkedByProviderResponse{}
	if err := v.client.http.Get(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
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
