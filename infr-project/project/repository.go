package project

import (
	"context"

	"github.com/futugyou/infr-project/domain"
)

type IProjectRepository interface {
	domain.IRepository[Project]
	GetProjectByName(ctx context.Context, name string) (*Project, error)
	GetAllProject(ctx context.Context, page *int, size *int) ([]Project, error)
}

type IProjectRepositoryAsync interface {
	domain.IRepositoryAsync[Project]
	GetProjectByNameAsync(ctx context.Context, name string) (<-chan *Project, <-chan error)
	GetAllProjectAsync(ctx context.Context, page *int, size *int) (<-chan []Project, <-chan error)
}
