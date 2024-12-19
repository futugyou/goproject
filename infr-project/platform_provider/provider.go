package platform_provider

import "context"

//TODO: fill field

const CommonProjectBadge = "https://img.shields.io/badge/%s-%s-%s?logo=%s&link=%s"

type CreateProjectRequest struct {
	PlatformId string
	Name       string
	Parameters map[string]string
}

type Project struct {
	ID          string
	Name        string
	Url         string
	Hooks       []WebHook
	Properties  map[string]string
	Envs        map[string]Env
	Workflows   map[string]Workflow
	Deployments map[string]Deployment
	BadgeURL    string
}

type Env struct {
	ID        string
	Key       string
	CreatedAt string
	UpdatedAt string
	Type      string
	Value     string
}

type Workflow struct {
	ID        string
	Name      string
	Status    string
	CreatedAt string
	BadgeURL  string
}

type Deployment struct {
	ID            string
	Name          string
	Plan          string
	ReadyState    string
	ReadySubstate string
	CreatedAt     string
	BadgeURL      string
}

type WebHook struct {
	ID         string
	Name       string
	Url        string
	Events     []string
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

type DeleteWebHookRequest struct {
	Parameters map[string]string
	WebHookId  string
}

type User struct {
	Name string
	ID   string
}

// Although the CreateProject method is provided, it is best not to use it.
// The DeleteProject method is not provided because it is more dangerous.
// The DeleteWebHook method is provided because it is less dangerous
type IPlatformProviderAsync interface {
	CreateProjectAsync(ctx context.Context, request CreateProjectRequest) (<-chan *Project, <-chan error)
	// no webhook info
	ListProjectAsync(ctx context.Context, filter ProjectFilter) (<-chan []Project, <-chan error)
	// include webhook info
	GetProjectAsync(ctx context.Context, filter ProjectFilter) (<-chan *Project, <-chan error)
	CreateWebHookAsync(ctx context.Context, request CreateWebHookRequest) (<-chan *WebHook, <-chan error)
	DeleteWebHookAsync(ctx context.Context, request DeleteWebHookRequest) <-chan error
	GetUserAsync(ctx context.Context) (<-chan *User, <-chan error)
}

func Intersect(setA, setB []string) []string {
	bMap := make(map[string]struct{})

	for _, b := range setB {
		bMap[b] = struct{}{}
	}

	var intersection []string

	for _, a := range setA {
		if _, found := bMap[a]; found {
			intersection = append(intersection, a)
		}
	}

	return intersection
}
