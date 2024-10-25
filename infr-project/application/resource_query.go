package application

import (
	"context"
	"time"

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

func (s *ResourceQueryService) HandleResourceChaged(ctx context.Context, data ResourceChangeData) error {
	return nil
}

type ResourceChangeData struct {
	Id              string    `bson:"id" json:"id"`
	ResourceVersion int       `bson:"version" json:"version"`
	EventType       string    `bson:"event_type" json:"event_type"`
	CreatedAt       time.Time `bson:"created_at" json:"created_at"`
	Name            string    `bson:"name" json:"name"`
	Type            string    `bson:"type" json:"type"`
	Data            string    `bson:"data" json:"data"`
	Tags            []string  `bson:"tags" json:"tags"`
}
