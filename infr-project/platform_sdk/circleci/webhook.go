package circleci

import "log"

func (s *CircleciClient) CreateWebhook(name string, url string, projectId string) string {
	path := "/webhook"
	request := CreateWebhookRequest{
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
	result := ""
	err := s.http.Post(path, request, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
}

type CreateWebhookRequest struct {
	Name          string       `json:"name"`
	Events        []string     `json:"events"` // "workflow-completed" "job-completed"
	Url           string       `json:"url"`
	VerifyTLS     bool         `json:"verify-tls"`
	SigningSecret string       `json:"signing-secret"`
	Scope         WebhookScope `json:"scope"`
}

type WebhookScope struct {
	Id   string `json:"id"`
	Type string `json:"type"` //"project"
}
