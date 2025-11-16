package viewmodel

type WebhookRequestInfo struct {
	Method     string              `json:"method"`
	URL        string              `json:"url"`
	Proto      string              `json:"proto"`
	Host       string              `json:"host"`
	Header     map[string][]string `json:"header"`
	Body       string              `json:"body"`
	Query      map[string][]string `json:"query"`
	RemoteAddr string              `json:"remote_addr"`
	UserAgent  string              `json:"user_agent"`
}

type WebhookLogs struct {
	Source             string `json:"source"`
	EventType          string `json:"event_type"`
	ProviderPlatformId string `json:"provider_platform_id"`
	ProviderProjectId  string `json:"provider_project_id"`
	ProviderWebhookId  string `json:"provider_webhook_id"`
	Data               string `json:"data"`
	HappenedAt         string `json:"happened_at"`
}

type WebhookSearch struct {
	Source             string `json:"source"`
	EventType          string `json:"event_type"`
	ProviderPlatformId string `json:"provider_platform_id"`
	ProviderProjectId  string `json:"provider_project_id"`
	ProviderWebhookId  string `json:"provider_webhook_id"`
}
