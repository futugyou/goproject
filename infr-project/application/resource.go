package application

import (
	"context"
	"fmt"

	tool "github.com/futugyou/extensions"

	domain "github.com/futugyou/infr-project/domain"
	infra "github.com/futugyou/infr-project/infrastructure"
	"github.com/futugyou/infr-project/resource"
	models "github.com/futugyou/infr-project/view_models"
)

type ResourceService struct {
	service    *ApplicationService[resource.IResourceEvent, *resource.Resource]
	unitOfWork domain.IUnitOfWork
}

func needStoreSnapshot(aggregate *resource.Resource) bool {
	return aggregate.AggregateVersion()%5 == 1
}

func NewResourceService(
	eventStore infra.IEventStore[resource.IResourceEvent],
	snapshotStore infra.ISnapshotStore[*resource.Resource],
	unitOfWork domain.IUnitOfWork,
	eventPublisher infra.IEventPublisher,
) *ResourceService {
	return &ResourceService{
		service:    NewApplicationService(eventStore, snapshotStore, unitOfWork, resource.ResourceFactory, needStoreSnapshot, eventPublisher),
		unitOfWork: unitOfWork,
	}
}

func (s *ResourceService) CreateResource(ctx context.Context, aux models.CreateResourceRequest) (*models.CreateResourceResponse, error) {
	var res *resource.Resource
	resourceType := resource.GetResourceType(aux.Type)
	if err := s.service.withUnitOfWork(ctx, func(ctx context.Context) error {
		res = resource.NewResource(aux.Name, resourceType, aux.Data, aux.ImageData, aux.Tags)
		return s.service.SaveSnapshotAndEvent(ctx, res)
	}); err != nil {
		return nil, err
	}

	return &models.CreateResourceResponse{Id: res.Id}, nil
}

func (s *ResourceService) UpdateResource(ctx context.Context, id string, aux models.UpdateResourceRequest) error {
	res, err := s.service.RetrieveLatestVersion(ctx, id)
	if err != nil {
		return err
	}

	source := *res
	oldVersion := source.Version
	var aggregate *resource.Resource
	var aggregates = make([]resource.Resource, 0)
	if source.Data != aux.Data {
		if aggregate, err = source.ChangeData(aux.Data); err != nil {
			return err
		}
		aggregates = append(aggregates, *aggregate)
	}

	if source.Name != aux.Name {
		if aggregate, err = source.ChangeName(aux.Name); err != nil {
			return err
		}
		aggregates = append(aggregates, *aggregate)
	}

	if !tool.StringArrayCompare(aux.Tags, source.Tags) {
		if aggregate, err = source.ChangeTags(aux.Tags); err != nil {
			return err
		}
		aggregates = append(aggregates, *aggregate)
	}

	if aggregate == nil || oldVersion == aggregate.Version {
		return fmt.Errorf("the data in the resource has not changed")
	}

	var aggs = make([]*resource.Resource, 0)
	for i := 0; i < len(aggregates); i++ {
		aggs = append(aggs, &aggregates[i])
	}

	return s.service.withUnitOfWork(ctx, func(ctx context.Context) error {
		return s.service.SaveSnapshotAndEvent2(ctx, aggs)
	})
}

// show all versions
func (s *ResourceService) AllVersionResource(ctx context.Context, id string) ([]models.ResourceView, error) {
	re, err := s.service.RetrieveAllVersions(ctx, id)
	if err != nil {
		return nil, err
	}

	result := make([]models.ResourceView, len(re))
	for i := 0; i < len(re); i++ {
		result[i] = *convertResourceEntityToViewModel(re[i])
	}
	return result, nil
}

func (s *ResourceService) DeleteResource(ctx context.Context, id string) error {
	res, err := s.service.RetrieveLatestVersion(ctx, id)
	if err != nil {
		return err
	}

	if res == nil || (*res).Id == "" {
		return fmt.Errorf("resource: %s not found", id)
	}

	aggregate, err := (*res).DeleteResource()
	if err != nil {
		return err
	}

	return s.service.withUnitOfWork(ctx, func(ctx context.Context) error {
		return s.service.SaveSnapshotAndEvent(ctx, aggregate)
	})
}

func convertResourceEntityToViewModel(src *resource.Resource) *models.ResourceView {
	if src == nil {
		return nil
	}

	return &models.ResourceView{
		Id:        src.Id,
		Name:      src.Name,
		Type:      src.Type.String(),
		Data:      src.Data,
		Version:   src.Version,
		IsDelete:  src.IsDeleted,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
		Tags:      src.Tags,
	}
}
