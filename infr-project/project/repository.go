package project

import (
	"context"

	"github.com/futugyou/infr-project/domain"
)

type IProjectRepository interface {
	domain.IRepository[Project]
	GetProjectByName(ctx context.Context, name string) (*Project, error)
	GetAllProject(ctx context.Context) ([]Project, error)
}
