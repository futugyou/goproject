package platform

import (
	"context"

	"github.com/futugyou/infr-project/domain"
)

type PlatformSearch struct {
	Name      string
	NameFuzzy bool
	Activate  *bool
	Tags      []string
	Page      int
	Size      int
}

type IPlatformRepository interface {
	domain.IRepository[Platform]
	GetPlatformByName(ctx context.Context, name string) (*Platform, error)
	SearchPlatforms(ctx context.Context, filter PlatformSearch) ([]Platform, error)
	GetPlatformByIdOrName(ctx context.Context, name string) ( *Platform,  error)
}