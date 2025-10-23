package platform_provider

import (
	"context"
	"fmt"
	"strings"
)

const CommonProjectBadge = "https://img.shields.io/badge/%s-%s-%s?logo=%s&link=%s"
const CommonProjectBadgeWithoutUrl = "https://img.shields.io/badge/%s-%s-%s?logo=%s"

func buildCommonBadge(lable string, state string, okState string, logo string, url *string) (badgeUrl string, badgeMarkDown string) {
	lable = strings.ReplaceAll(lable, "-", "%20")
	if state == okState {
		badgeUrl = fmt.Sprintf(CommonProjectBadgeWithoutUrl, lable, state, "brightgreen", logo)
	} else {
		badgeUrl = fmt.Sprintf(CommonProjectBadgeWithoutUrl, lable, state, "red", logo)
	}

	if url != nil {
		badgeUrl += fmt.Sprintf("&link=%s", *url)
	}
	badgeMarkDown = fmt.Sprintf("![%s](%s)", lable, badgeUrl)
	return
}

type CreateProjectRequest struct {
	PlatformId string
	Name       string
	Parameters map[string]string
}

type Project struct {
	ID                   string
	Name                 string
	Url                  string
	Description          string
	WebHooks             []WebHook
	Properties           map[string]string
	EnvironmentVariables map[string]EnvironmentVariable // github use action repositroy secrets
	Environments         []string                       // circleci don't have
	Workflows            map[string]Workflow            // circleci have pipeline, but it's a record of WorkflowRun. vercel don't have
	WorkflowRuns         map[string]WorkflowRun         // circleci use 'pipeline'. vercel don't have
	Deployments          map[string]Deployment          // circleci don't have
	BadgeURL             string
	BadgeMarkDown        string
	Tags                 []string
}

func (w *Project) GetProperties() map[string]string {
	if w == nil {
		return map[string]string{}
	}
	return w.Properties
}

func (w *Project) GetWebhooks() []WebHook {
	if w == nil {
		return []WebHook{}
	}
	return w.WebHooks
}

type EnvironmentVariable struct {
	ID        string
	Key       string
	CreatedAt string
	UpdatedAt string
	Type      string
	Value     string
}

type Workflow struct {
	ID            string
	Name          string
	Status        string
	CreatedAt     string
	BadgeURL      string
	BadgeMarkdown string
}

type WorkflowRun struct {
	ID            string
	Name          string
	Description   string
	Status        string
	CreatedAt     string
	BadgeURL      string
	BadgeMarkdown string
}

type Deployment struct {
	ID            string
	Name          string
	Environment   string
	ReadyState    string
	ReadySubstate string
	CreatedAt     string
	BadgeURL      string
	BadgeMarkdown string
	Description   string
}

type WebHook struct {
	ID         string
	Name       string
	Url        string
	Events     []string
	Activate   bool
	Parameters map[string]string
}

func (w *WebHook) GetParameters() map[string]string {
	if w == nil {
		return map[string]string{}
	}
	return w.Parameters
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
