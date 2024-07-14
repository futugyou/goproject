package application

import (
	"context"
	"fmt"
	"strings"

	domain "github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/extensions"
	infra "github.com/futugyou/infr-project/infrastructure"
	"github.com/futugyou/infr-project/resource"
	models "github.com/futugyou/infr-project/view_models"
)

type ResourceService struct {
	service    *ApplicationService[resource.IResourceEvent, *resource.Resource]
	unitOfWork domain.IUnitOfWork
	queryRepo  IPlatformRepository
}

func needStoreSnapshot(aggregate *resource.Resource) bool {
	return aggregate.AggregateVersion()%5 == 1
}

func NewResourceService(
	eventStore infra.IEventStore[resource.IResourceEvent],
	snapshotStore infra.ISnapshotStore[*resource.Resource],
	unitOfWork domain.IUnitOfWork,
	queryRepo IPlatformRepository,
) *ResourceService {
	return &ResourceService{
		service:    NewApplicationService(eventStore, snapshotStore, unitOfWork, resource.ResourceFactory, needStoreSnapshot),
		unitOfWork: unitOfWork,
		queryRepo:  queryRepo,
	}
}

func (s *ResourceService) CreateResource(aux models.CreateResourceRequest, ctx context.Context) (*resource.Resource, error) {
	var res *resource.Resource
	resourceType := resource.GetResourceType(aux.Type)
	err := s.service.withUnitOfWork(ctx, func(ctx context.Context) error {
		res = resource.NewResource(aux.Name, resourceType, aux.Data, aux.Tags)
		return s.service.SaveSnapshotAndEvent(ctx, res)
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ResourceService) UpdateResource(id string, aux models.UpdateResourceRequest, ctx context.Context) error {
	res, err := s.service.RetrieveLatestVersion(id, ctx)
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
		res, err := s.queryRepo.GetResourceByName(ctx, aux.Name)
		if err != nil && !strings.HasPrefix(err.Error(), "data not found") {
			return err
		}

		if res != nil && len(res.Id) > 0 && res.Id != id {
			return fmt.Errorf("name: %s is existed", aux.Name)
		}

		if aggregate, err = source.ChangeName(aux.Name); err != nil {
			return err
		}
		aggregates = append(aggregates, *aggregate)
	}

	if !extensions.StringArrayCompare(aux.Tags, source.Tags) {
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
func (s *ResourceService) AllVersionResource(id string, ctx context.Context) ([]resource.Resource, error) {
	re, err := s.service.RetrieveAllVersions(id, ctx)
	if err != nil {
		return nil, err
	}

	result := make([]resource.Resource, 0)
	for i := 0; i < len(re); i++ {
		result = append(result, *re[i])
	}
	return result, nil
}

func (s *ResourceService) DeleteResource(id string, ctx context.Context) error {
	res, err := s.service.RetrieveLatestVersion(id, ctx)
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
