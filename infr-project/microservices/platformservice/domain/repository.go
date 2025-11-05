package domain

import (
	"context"

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
	GetPlatformProjectByProjectID(ctx context.Context, platformID string, projectID string) (*PlatformProject, error)
	UpdateProject(ctx context.Context, platformID string, project PlatformProject) error
}
