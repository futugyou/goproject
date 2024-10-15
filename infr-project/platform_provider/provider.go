package platform_provider

import "context"

//TODO: fill field

type CreateProjectRequest struct {
	PlatformId string
	Name       string
	Parameters map[string]string
}

type Project struct {
	Name  string
	Url   string
	Hooks []WebHook
}

type WebHook struct {
	Name       string
	Url        string
	Parameters map[string]string
}

type ProjectFilter struct {
	Name       string
	Parameters map[string]string
}

type CreateWebHookRequest struct {
	PlatformId string
	ProjectId  string
	WebHook    WebHook
}

type IPlatformProviderAsync interface {
	CreateProjectAsync(ctx context.Context, request CreateProjectRequest) (<-chan *Project, <-chan error)
	// no webhook info
	ListProjectAsync(ctx context.Context, filter ProjectFilter) (<-chan []Project, <-chan error)
	// include webhook info
	GetProjectAsync(ctx context.Context, filter ProjectFilter) (<-chan *Project, <-chan error)
	CreateWebHookAsync(ctx context.Context, request CreateWebHookRequest) (<-chan *WebHook, <-chan error)
}
