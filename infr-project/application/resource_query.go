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

func (s *ResourceQueryService) GetAllResources() ([]models.ResourceDetail, error) {
	return s.repository.GetAll(context.Background())
}

func (s *ResourceQueryService) CurrentResource(id string) (*models.ResourceDetail, error) {
	return s.repository.Get(context.Background(), id)
}
