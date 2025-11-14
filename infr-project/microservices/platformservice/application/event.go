package application

type CreateProviderProjectTriggeredEvent struct {
	PlatformID  string `json:"platform_id"`
	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
	Provider    string `json:"provider"`
}

func (e *CreateProviderProjectTriggeredEvent) EventType() string {
	return "create_provider_project"
}

type CreateProviderWebhookTriggeredEvent struct {
	PlatformID        string `json:"platform_id"`
	ProjectID         string `json:"project_id"`
	ProjectName       string `json:"project_name"`
	Provider          string `json:"provider"`
	ProviderProjectId string `json:"provider_project_id"`
}

func (e *CreateProviderWebhookTriggeredEvent) EventType() string {
	return "create_provider_webhook"
}

type ProjectScreenshotTriggeredEvent struct {
	PlatformID string `json:"platform_id"`
	ProjectID  string `json:"project_id"`
	Url        string `json:"url"`
}

func (e *ProjectScreenshotTriggeredEvent) EventType() string {
	return "project_screenshot"
}
