package application

import (
	"context"

	"github.com/futugyou/infr-project/extensions"
	models "github.com/futugyou/infr-project/view_models"
)

type IRepository[Query models.IQuery] interface {
	Get(ctx context.Context, id string) (*Query, error)
	GetAll(ctx context.Context) ([]Query, error)
	GetWithSearch(ctx context.Context, condition extensions.Search) ([]Query, error)
}

type IPlatformRepository interface {
	IRepository[models.ResourceDetail]
}
