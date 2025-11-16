package viewmodel

type CreateProviderProjectRequest struct {
	PlatformID  string `json:"platform_id"`
	ProjectID   string `json:"project_id"`
	ProjectName bool   `json:"project_name"`
	Provider    bool   `json:"provider"`
}

type CreateProviderWebhookRequest struct {
	PlatformID        string `json:"platform_id"`
	ProjectID         string `json:"project_id"`
	ProjectName       bool   `json:"project_name"`
	Provider          bool   `json:"provider"`
	ProviderProjectId string `json:"provider_project_id"`
	Url               string `json:"url"`
}

type ProjectScreenshotRequest struct {
	PlatformID string `json:"platform_id"`
	ProjectID  string `json:"project_id"`
	Url        bool   `json:"url"`
}
