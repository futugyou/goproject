package application

import (
	"context"
	"fmt"
	"strings"

	tool "github.com/futugyou/extensions"

	domain "github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/extensions"
	resourcequery "github.com/futugyou/infr-project/resource_query"
	models "github.com/futugyou/infr-project/view_models"
	"github.com/redis/go-redis/v9"
)

type ResourceQueryService struct {
	innerService *AppService
	repository   resourcequery.IResourceRepository
	client       *redis.Client
}

func NewResourceQueryService(repository resourcequery.IResourceRepository, client *redis.Client,
	unitOfWork domain.IUnitOfWork) *ResourceQueryService {
	return &ResourceQueryService{
		repository:   repository,
		client:       client,
		innerService: NewAppService(unitOfWork),
	}
}

func (s *ResourceQueryService) GetAllResources(ctx context.Context) ([]models.ResourceView, error) {
	// ignore error
	resourceViews, _ := tool.RedisListHashWithLua[models.ResourceView](ctx, s.client, "ResourceView:", 100)
	if len(resourceViews) > 0 {
		for i := 0; i < len(resourceViews); i++ {
			resourceViews[i].Tags = strings.Split(resourceViews[i].TagString, ",")
		}
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

func (s *ResourceQueryService) GetResource(ctx context.Context, id string) (*models.ResourceView, error) {
	var viewData models.ResourceView
	s.client.HGetAll(ctx, "ResourceView:"+id).Scan(&viewData)
	if len(viewData.Id) > 0 {
		if len(viewData.TagString) > 0 {
			if strings.Contains(viewData.TagString, ",") {
				viewData.Tags = strings.Split(viewData.TagString, ",")
			} else {
				viewData.Tags = []string{viewData.TagString}
			}
		}
		return &viewData, nil
	}

	data, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	viewData = s.convertData(*data)
	// ignore error
	s.client.HSet(ctx, "ResourceView:"+id, viewData).Result()

	return &viewData, nil
}

func (s *ResourceQueryService) convertData(data resourcequery.Resource) models.ResourceView {
	v := models.ResourceView{
		Id:        data.Id,
		Name:      data.Name,
		Type:      data.Type,
		Data:      data.Data,
		Version:   data.Version,
		IsDelete:  data.IsDelete,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Tags:      data.Tags,
		ImageData: data.ImageData,
	}
	if len(data.Tags) > 0 {
		v.TagString = strings.Join(data.Tags, ",")
	}
	return v
}

func (s *ResourceQueryService) HandleResourceChanged(ctx context.Context, data models.ResourceChangeData) error {
	res, err := s.repository.Get(ctx, data.Id)
	if err != nil && !strings.HasPrefix(err.Error(), extensions.Data_Not_Found_Message) {
		return err
	}

	if res == nil {
		if data.EventType == "ResourceCreated" {
			res = &resourcequery.Resource{
				Aggregate: domain.Aggregate{
					Id: data.Id,
				},
				Name:      data.Name,
				Type:      data.Type,
				Data:      data.Data,
				ImageData: data.ImageData,
				Version:   data.ResourceVersion,
				IsDelete:  false,
				CreatedAt: data.CreatedAt,
				UpdatedAt: data.CreatedAt,
				Tags:      data.Tags,
			}
		}
	} else if res.Version < data.ResourceVersion {
		res.Version = data.ResourceVersion
		res.UpdatedAt = data.CreatedAt
		switch data.EventType {
		case "ResourceDeleted":
			res.IsDelete = true
		case "ResourceUpdated":
			res.Name = data.Name
			res.Type = data.Type
			res.Data = data.Data
			res.ImageData = data.ImageData
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
	}

	if res == nil {
		return fmt.Errorf("resource can not find, ID is %s", data.Id)
	}

	return s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		switch data.EventType {
		case "ResourceCreated":
			err = s.repository.Insert(ctx, *res)
		default:
			err = s.repository.Update(ctx, *res)
		}

		if err != nil {
			return err
		}

		viewData := s.convertData(*res)
		// ignore error
		s.client.HSet(ctx, "ResourceView:"+viewData.Id, viewData).Result()
		return err
	})
}
