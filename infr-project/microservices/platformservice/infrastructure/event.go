package infrastructure

type CreateProviderProjectTriggeredEvent struct {
	ID          string `json:"id"`
	PlatformID  string `json:"platform_id"`
	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
	Provider    string `json:"provider"`
}

func (e *CreateProviderProjectTriggeredEvent) EventType() string {
	return "create_provider_project"
}

type CreateProviderWebhookTriggeredEvent struct {
	ID                string `json:"id"`
	PlatformID        string `json:"platform_id"`
	ProjectID         string `json:"project_id"`
	ProjectName       string `json:"project_name"`
	Provider          string `json:"provider"`
	ProviderProjectId string `json:"provider_project_id"`
	Url               string `json:"url"`
}

func (e *CreateProviderWebhookTriggeredEvent) EventType() string {
	return "create_provider_webhook"
}

type ProjectScreenshotTriggeredEvent struct {
	ID         string `json:"id"`
	PlatformID string `json:"platform_id"`
	ProjectID  string `json:"project_id"`
	Url        string `json:"url"`
}

func (e *ProjectScreenshotTriggeredEvent) EventType() string {
	return "project_screenshot"
}
