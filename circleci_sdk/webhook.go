package circleci

type WebhookService service

func (s *WebhookService) CreateWebhook(name string, url string, projectId string, signingSecret string) (*WebhookItem, error) {
	path := "/webhook"
	request := WebhookItem{
		Name:          name,
		Events:        []string{"workflow-completed"},
		Url:           url,
		VerifyTLS:     false,
		SigningSecret: signingSecret,
		Scope: WebhookScope{
			Id:   projectId,
			Type: "project",
		},
	}

	result := &WebhookItem{}
	if err := s.client.http.Post(path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *WebhookService) ListWebhook(projectId string) (*ListWebhookResponse, error) {
	path := "/webhook?scope-id=" + projectId + "&scope-type=project"
	result := &ListWebhookResponse{}
	if err := s.client.http.Get(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *WebhookService) GetWebhook(webhookId string) (*WebhookItem, error) {
	path := "/webhook/" + webhookId
	result := &WebhookItem{}
	if err := s.client.http.Get(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *WebhookService) UpdateWebhook(webhookId string, name string, url string, signingSecret string) (*WebhookItem, error) {
	path := "/webhook/" + webhookId
	request := WebhookItem{
		Name:          name,
		Events:        []string{"workflow-completed"},
		Url:           url,
		VerifyTLS:     false,
		SigningSecret: signingSecret,
	}

	result := &WebhookItem{}
	if err := s.client.http.Put(path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *WebhookService) DeleteWebhook(webhookId string) (*BaseResponse, error) {
	path := "/webhook/" + webhookId
	result := &BaseResponse{}
	if err := s.client.http.Delete(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

type ListWebhookResponse struct {
	Items         []WebhookItem `json:"items"`
	NextPageToken string        `json:"next_page_token"`
	Message       *string       `json:"message,omitempty"`
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
	Message       *string      `json:"message,omitempty"`
}

type WebhookScope struct {
	Id   string `json:"id"`
	Type string `json:"type"` //"project"
}
