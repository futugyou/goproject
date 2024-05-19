package application

import (
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
	s.unitOfWork.Start()
	res := resource.NewResource(name, resourceType, data)

	if err := s.service.SaveSnapshotAndEvent(res); err != nil {
		s.unitOfWork.Rollback()
		return nil, err
	}

	return res, s.unitOfWork.Commit()
}

func (s *ResourceService) UpdateResourceDate(id string, data string) error {
	res, err := s.service.RetrieveLatestVersion(id)
	if err == nil {
		return err
	}

	aggregate := (*res)
	aggregate = aggregate.ChangeData(data)
	s.unitOfWork.Start()
	if err := s.service.SaveSnapshotAndEvent(aggregate); err != nil {
		s.unitOfWork.Rollback()
		return err
	}

	return s.unitOfWork.Commit()
}
