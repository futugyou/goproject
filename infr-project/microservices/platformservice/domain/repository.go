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
	SearchPlatforms(ctx context.Context, filter PlatformSearch) ([]Platform, error)
	GetPlatformByIdOrName(ctx context.Context, name string) (*Platform, error)
}
