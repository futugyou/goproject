package application

import (
	"context"

	models "github.com/futugyou/infr-project/view_models"
)

type IRepository[Query models.IQuery] interface {
	Get(ctx context.Context, id string) (*Query, error)
	GetAll(ctx context.Context) ([]Query, error)
}

type IPlatformRepository interface {
	IRepository[models.ResourceView]
	GetResourceByName(ctx context.Context, name string) (*models.ResourceView, error)
}
