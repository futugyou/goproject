package sdk

import "log"

type CircleciClient struct {
	http IHttpClient
}

const circleci_url string = "https://circleci.com/api/v2"

func NewCircleciClient(token string) *CircleciClient {
	header := make(map[string]string)
	header["Circle-Token"] = token

	c := &CircleciClient{
		http: newHttpClientWithHeader(circleci_url, header),
	}
	return c
}

func (s *CircleciClient) Pipelines(org_slug string) string {
	path := "/pipeline?org-slug=" + org_slug
	result := "[]"
	err := s.http.Get(path, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
}

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
