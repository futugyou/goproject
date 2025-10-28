package domain

import (
	"context"

	"github.com/futugyou/domaincore/domain"
)

type ResourceRepository interface {
	domain.Repository[Resource]
	GetResourceByName(ctx context.Context, name string) (*Resource, error)
	GetAllResource(ctx context.Context, page *int, size *int) ([]Resource, error)
}
