package application

import (
	"github.com/futugyou/infr-project/resource"
)

type ResourceQueryService struct {
}

func NewResourceQueryService() *ResourceQueryService {
	return &ResourceQueryService{}
}

// TODO: This is temporary because CQRS is not ready yet
func (s *ResourceQueryService) GetAllResourceSnapshots() ([]resource.Resource, error) {
	return []resource.Resource{}, nil
}
