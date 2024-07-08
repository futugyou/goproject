package vercel

import (
	"fmt"
	"net/url"
)

func (v *VercelClient) DeleteIntegrationConfiguration(id string, slug string, teamId string) (*string, error) {
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
	err := v.http.Delete(path, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *VercelClient) RetrieveIntegrationConfiguration(id string, slug string, teamId string) (*IntegrationConfigurationInfo, error) {
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
	err := v.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) GetConfiguration(view string, slug string, teamId string) (*string, error) {
	path := "/v1/integrations/configurations"
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
	err := v.http.Get(path, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
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
