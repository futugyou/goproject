package application

import (
	"context"
	"time"

	resourcequery "github.com/futugyou/infr-project/resource_query"
	models "github.com/futugyou/infr-project/view_models"
)

type ResourceQueryService struct {
	repository resourcequery.IResourceRepository
}

func NewResourceQueryService(repository resourcequery.IResourceRepository) *ResourceQueryService {
	return &ResourceQueryService{
		repository: repository,
	}
}

func (s *ResourceQueryService) GetAllResources(ctx context.Context) ([]models.ResourceView, error) {
	datas, err := s.repository.GetAllResource(ctx, nil, nil)
	if err != nil {
		return nil, err
	}

	result := make([]models.ResourceView, 0)
	for _, data := range datas {
		t := s.convertData(&data)
		result = append(result, *t)
	}

	return result, nil
}

func (s *ResourceQueryService) CurrentResource(ctx context.Context, id string) (*models.ResourceView, error) {
	data, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.convertData(data), nil
}

func (s *ResourceQueryService) convertData(data *resourcequery.Resource) *models.ResourceView {
	if data == nil {
		return nil
	}

	return &models.ResourceView{
		Id:        data.Id,
		Name:      data.Name,
		Type:      data.Type,
		Data:      data.Data,
		Version:   data.Version,
		IsDelete:  data.IsDelete,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Tags:      data.Tags,
	}
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
