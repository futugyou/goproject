package domain

import (
	"context"
	"time"

	"github.com/futugyou/domaincore/domain"
)

type PlatformSearch struct {
	Name      string
	NameFuzzy bool
	Activate  *bool
	Tags      []string
	Page      int
	Size      int
}

type PlatformRepository interface {
	domain.Repository[Platform]
	GetPlatformByName(ctx context.Context, name string) (*Platform, error)
	GetPlatformByIdOrName(ctx context.Context, idOrName string) (*Platform, error)
	GetPlatformByIdOrNameWithoutProjects(ctx context.Context, idOrName string) (*Platform, error)
	SearchPlatforms(ctx context.Context, filter PlatformSearch) ([]Platform, error)
	GetPlatformProjects(ctx context.Context, platformID string) ([]PlatformProject, error)
	GetPlatformProjectByIDOrName(ctx context.Context, platformID string, projectID string) (*PlatformProject, error)
	SyncProjects(ctx context.Context, platformID string, projects []PlatformProject) error
	UpdateProject(ctx context.Context, platformID string, project PlatformProject) error
	DeleteProject(ctx context.Context, platformID string, projectID string) error
}

type WebhookLogSearch struct {
	Source             string
	EventType          string
	ProviderPlatformId string
	ProviderProjectId  string
	ProviderWebhookId  string
}

type WebhookLogRepository interface {
	domain.Repository[WebhookLogs]
	SearchWebhookLogs(ctx context.Context, filter WebhookLogSearch) ([]WebhookLogs, error)
	DeleteWebhookLogsByDate(ctx context.Context, filter time.Time) error
	InsertAndDeleteOldData(ctx context.Context, logs []WebhookLogs, filter time.Time) error
}
