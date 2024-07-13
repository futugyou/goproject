package application

import (
	"context"

	models "github.com/futugyou/infr-project/view_models"
)

type ResourceQueryService struct {
	repository IPlatformRepository
}

func NewResourceQueryService(repository IPlatformRepository) *ResourceQueryService {
	return &ResourceQueryService{
		repository: repository,
	}
}

func (s *ResourceQueryService) GetAllResources(ctx context.Context) ([]models.ResourceDetail, error) {
	return s.repository.GetAll(ctx)
}

func (s *ResourceQueryService) CurrentResource(id string, ctx context.Context) (*models.ResourceDetail, error) {
	return s.repository.Get(ctx, id)
}
