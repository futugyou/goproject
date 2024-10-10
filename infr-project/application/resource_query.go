package application

import (
	"context"

	models "github.com/futugyou/infr-project/view_models"
)

type ResourceQueryService struct {
	repository IResourceViewRepository
}

func NewResourceQueryService(repository IResourceViewRepository) *ResourceQueryService {
	return &ResourceQueryService{
		repository: repository,
	}
}

func (s *ResourceQueryService) GetAllResources(ctx context.Context) ([]models.ResourceView, error) {
	return s.repository.GetAll(ctx)
}

func (s *ResourceQueryService) CurrentResource(ctx context.Context, id string) (*models.ResourceView, error) {
	return s.repository.Get(ctx, id)
}
