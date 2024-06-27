package circleci

import "log"

func (s *CircleciClient) CreateWebhook(name string, url string, projectId string) WebhookItem {
	path := "/webhook"
	request := WebhookItem{
		Name:          name,
		Events:        []string{"workflow-completed"},
		Url:           url,
		VerifyTLS:     false,
		SigningSecret: "",
		Scope: WebhookScope{
			Id:   projectId,
			Type: "project",
		},
	}
	result := WebhookItem{}
	err := s.http.Post(path, request, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
}

func (s *CircleciClient) ListWebhook(projectId string) ListWebhookResponse {
	path := "/webhook?scope-id=" + projectId + "&scope-type=project"

	result := ListWebhookResponse{}
	err := s.http.Get(path, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
}

type WebhookScope struct {
	Id   string `json:"id"`
	Type string `json:"type"` //"project"
}

type ListWebhookResponse struct {
	Items         []WebhookItem `json:"items"`
	NextPageToken string        `json:"next_page_token"`
}

type WebhookItem struct {
	Name          string       `json:"name,omitempty"`
	Events        []string     `json:"events,omitempty"` // "workflow-completed" "job-completed"
	Url           string       `json:"url,omitempty"`
	VerifyTLS     bool         `json:"verify-tls,omitempty"`
	SigningSecret string       `json:"signing-secret,omitempty"`
	Scope         WebhookScope `json:"scope,omitempty"`
	UpdatedAt     string       `json:"updated-at,omitempty"`
	Id            string       `json:"id,omitempty"`
	CreatedAt     string       `json:"created-at,omitempty"`
}
