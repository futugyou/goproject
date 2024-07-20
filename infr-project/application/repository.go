package application

import (
	"context"

	models "github.com/futugyou/infr-project/view_models"
)

type IRepository[Query models.IQuery] interface {
	Get(ctx context.Context, id string) (*Query, error)
	GetAll(ctx context.Context) ([]Query, error)
}

type IRepositoryAsync[Query models.IQuery] interface {
	GetAsync(ctx context.Context, id string) (<-chan *Query, <-chan error)
	GetAllAsync(ctx context.Context) (<-chan []Query, <-chan error)
}

type IResourceViewRepository interface {
	IRepository[models.ResourceView]
	GetResourceByName(ctx context.Context, name string) (*models.ResourceView, error)
}

type IResourceViewRepositoryAsync interface {
	IRepositoryAsync[models.ResourceView]
	GetResourceByNameAsync(ctx context.Context, name string) (<-chan *models.ResourceView, <-chan error)
}
