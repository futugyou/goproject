package application

import (
	view_models "github.com/futugyou/infr-project/view_models"
)

type ResourceQueryService struct {
}

func NewResourceQueryService() *ResourceQueryService {
	return &ResourceQueryService{}
}

// TODO: This is temporary because CQRS is not ready yet
func (s *ResourceQueryService) GetAllResourceSnapshots() ([]view_models.ResourceDetail, error) {
	return []view_models.ResourceDetail{}, nil
}
