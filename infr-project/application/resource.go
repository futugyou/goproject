package application

import (
	"context"

	domain "github.com/futugyou/infr-project/domain"
	infra "github.com/futugyou/infr-project/infrastructure"
	"github.com/futugyou/infr-project/resource"
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

func (s *ResourceService) CurrentResource(id string) (*resource.Resource, error) {
	res, err := s.service.RetrieveLatestVersion(id)
	if err != nil {
		return nil, err
	}
	return *res, nil
}

func (s *ResourceService) CreateResource(name string, resourceType resource.ResourceType, data string) (*resource.Resource, error) {
	var res *resource.Resource
	err := s.service.withUnitOfWork(context.Background(), func(ctx context.Context) error {
		res = resource.NewResource(name, resourceType, data)
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

	aggregate := (*res).ChangeData(data)

	return s.service.withUnitOfWork(context.Background(), func(ctx context.Context) error {
		return s.service.SaveSnapshotAndEvent(ctx, aggregate)
	})
}
