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
		service:    NewApplicationService(eventStore, snapshotStore, resource.ResourceFactory, needStoreSnapshot),
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
	ctx := context.Background()
	ctx, err := s.unitOfWork.Start(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			s.unitOfWork.Rollback(ctx)
		} else {
			err = s.unitOfWork.Commit(ctx)
		}
		s.unitOfWork.End(ctx)
	}()

	res := resource.NewResource(name, resourceType, data)

	err = s.service.SaveSnapshotAndEvent(ctx, res)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (s *ResourceService) UpdateResourceDate(id string, data string) error {
	res, err := s.service.RetrieveLatestVersion(id)
	if err == nil {
		return err
	}

	aggregate := (*res)
	aggregate = aggregate.ChangeData(data)

	ctx := context.Background()
	ctx, err = s.unitOfWork.Start(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			s.unitOfWork.Rollback(ctx)
		} else {
			err = s.unitOfWork.Commit(ctx)
		}
		s.unitOfWork.End(ctx)
	}()

	err = s.service.SaveSnapshotAndEvent(ctx, aggregate)
	return err
}
