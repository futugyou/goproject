package platform

import (
	"context"

	"github.com/futugyou/infr-project/domain"
)

type IPlatformRepository interface {
	domain.IRepository[Platform]
	GetPlatformByName(ctx context.Context, name string) (*Platform, error)
	GetAllPlatform(ctx context.Context) ([]Platform, error)
}

type IPlatformRepositoryAsync interface {
	domain.IRepositoryAsync[Platform]
	GetPlatformByNameAsync(ctx context.Context, name string) (<-chan *Platform, <-chan error)
	GetAllPlatformAsync(ctx context.Context) (<-chan []Platform, <-chan error)
}
