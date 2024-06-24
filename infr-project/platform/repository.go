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
