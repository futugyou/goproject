package domain

import (
	"context"

	"github.com/futugyou/domaincore/domain"
)

type ProjectRepository interface {
	domain.Repository[Project]
	GetProjectByName(ctx context.Context, name string) (*Project, error)
	GetAllProject(ctx context.Context, page *int, size *int) ([]Project, error)
}
