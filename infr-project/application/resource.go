package application

import (
	"context"
	"fmt"

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
) *ResourceService {
	return &ResourceService{
		service:    NewApplicationService(eventStore, snapshotStore, unitOfWork, resource.ResourceFactory, needStoreSnapshot),
		unitOfWork: unitOfWork,
	}
}

// func (s *ResourceService) CurrentResource(id string) (*resource.Resource, error) {
// 	res, err := s.service.RetrieveLatestVersion(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return *res, nil
// }

func (s *ResourceService) CreateResource(aux models.CreateResourceRequest) (*resource.Resource, error) {
	var res *resource.Resource
	resourceType := resource.GetResourceType(aux.Type)
	err := s.service.withUnitOfWork(context.Background(), func(ctx context.Context) error {
		res = resource.NewResource(aux.Name, resourceType, aux.Data, aux.Tags)
		return s.service.SaveSnapshotAndEvent(ctx, res)
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ResourceService) UpdateResourceDate(id string, data string) error {
	res, err := s.service.RetrieveLatestVersion(id)
	if err != nil {
		return err
	}

	aggregate, err := (*res).ChangeData(data)
	if err != nil {
		return err
	}

	return s.service.withUnitOfWork(context.Background(), func(ctx context.Context) error {
		return s.service.SaveSnapshotAndEvent(ctx, aggregate)
	})
}

// show all versions
func (s *ResourceService) AllVersionResource(id string) ([]resource.Resource, error) {
	re, err := s.service.RetrieveAllVersions(id)
	if err != nil {
		return nil, err
	}

	result := make([]resource.Resource, 0)
	for i := 0; i < len(re); i++ {
		result = append(result, *re[i])
	}
	return result, nil
}

func (s *ResourceService) DeleteResource(id string) error {
	res, err := s.service.RetrieveLatestVersion(id)
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

	return s.service.withUnitOfWork(context.Background(), func(ctx context.Context) error {
		return s.service.SaveSnapshotAndEvent(ctx, aggregate)
	})
}
