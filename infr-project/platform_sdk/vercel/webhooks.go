package vercel

import (
	"fmt"
	"net/url"
)

func (v *VercelClient) CreateWebhook(slug string, teamId string, req CreateWebhookRequest) (*WebhookInfo, error) {
	path := "/v1/webhooks"
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
	result := &WebhookInfo{}
	err := v.http.Post(path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *VercelClient) DeleteWebhook(id string, slug string, teamId string) (*string, error) {
	path := fmt.Sprintf("/v1/webhooks%s", id)
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

type CreateWebhookRequest struct {
	Events     []string `json:"events,omitempty"`
	Url        string   `json:"url,omitempty"`
	ProjectIds []string `json:"projectIds,omitempty"`
}

type WebhookInfo struct {
	CreatedAt  int          `json:"createdAt,omitempty"`
	Events     []string     `json:"events,omitempty"`
	Url        string       `json:"url,omitempty"`
	ProjectIds []string     `json:"projectIds,omitempty"`
	Id         string       `json:"id,omitempty"`
	OwnerId    string       `json:"OwnerId,omitempty"`
	Secret     string       `json:"secret,omitempty"`
	UpdatedAt  int          `json:"updatedAt,omitempty"`
	Error      *VercelError `json:"error,omitempty"`
}
