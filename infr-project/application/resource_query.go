package application

import (
	"context"
	"strings"
	"time"

	domain "github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/extensions"
	resourcequery "github.com/futugyou/infr-project/resource_query"
	models "github.com/futugyou/infr-project/view_models"
	"github.com/redis/go-redis/v9"
)

type ResourceQueryService struct {
	repository resourcequery.IResourceRepository
	client     *redis.Client
}

func NewResourceQueryService(repository resourcequery.IResourceRepository, client *redis.Client) *ResourceQueryService {
	return &ResourceQueryService{
		repository: repository,
		client:     client,
	}
}

func (s *ResourceQueryService) GetAllResources(ctx context.Context) ([]models.ResourceView, error) {
	resourceViews, _ := extensions.RedisListHashWithLua[models.ResourceView](ctx, s.client, "ResourceView:", 100)

	if len(resourceViews) > 0 {
		return resourceViews, nil
	}

	datas, err := s.repository.GetAllResource(ctx, nil, nil)
	if err != nil {
		return nil, err
	}

	result := make([]models.ResourceView, 0)
	for _, data := range datas {
		result = append(result, s.convertData(data))
	}

	return result, nil
}

func (s *ResourceQueryService) CurrentResource(ctx context.Context, id string) (*models.ResourceView, error) {
	var viewData models.ResourceView
	s.client.HGetAll(ctx, "ResourceView:"+id).Scan(&viewData)
	if len(viewData.Id) > 0 {
		return &viewData, nil
	}

	data, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	viewData = s.convertData(*data)
	s.client.HSet(ctx, "ResourceView:"+id, viewData).Result()

	return &viewData, nil
}

func (s *ResourceQueryService) convertData(data resourcequery.Resource) models.ResourceView {
	return models.ResourceView{
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

func (s *ResourceQueryService) HandleResourceChanged(ctx context.Context, data ResourceChangeData) error {
	res, err := s.repository.Get(ctx, data.Id)
	if err != nil && !strings.HasPrefix(err.Error(), extensions.Data_Not_Found_Message) {
		return err
	}

	if res == nil {
		if data.EventType == "ResourceCreated" {
			aggregate := resourcequery.Resource{
				Aggregate: domain.Aggregate{
					Id: data.Id,
				},
				Name:      data.Name,
				Type:      data.Type,
				Data:      data.Data,
				Version:   data.ResourceVersion,
				IsDelete:  false,
				CreatedAt: data.CreatedAt,
				UpdatedAt: data.CreatedAt,
				Tags:      data.Tags,
			}
			return s.repository.Insert(ctx, aggregate)
		}
	} else if res.Version < data.ResourceVersion {
		res.Version = data.ResourceVersion
		res.UpdatedAt = data.CreatedAt
		switch data.EventType {
		case "ResourceCreated":
			res.IsDelete = true
		case "ResourceUpdated":
			res.Name = data.Name
			res.Type = data.Type
			res.Data = data.Data
			res.Tags = data.Tags
		case "ResourceNameChanged":
			res.Name = data.Name
		case "ResourceDataChanged":
			res.Data = data.Data
		case "ResourceTypeChanged":
			res.Tags = data.Tags
		case "ResourceTagsChanged":
			res.Type = data.Type
		}

		s.client.Del(ctx, "ResourceView:"+data.Id).Result()

		return s.repository.Update(ctx, *res)
	}

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
