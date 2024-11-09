package resourcequery

import (
	"context"

	"github.com/futugyou/infr-project/domain"
)

type IResourceRepository interface {
	domain.IRepository[Resource]
	GetResourceByName(ctx context.Context, name string) (*Resource, error)
	GetAllResource(ctx context.Context, page *int, size *int) ([]Resource, error)
}
