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

	es := res.DomainEvents()
	events := make([]resource.IResourceEvent, 0)
	for i := 0; i < len(es); i++ {
		events = append(events, es[i])
	}

	if err := s.service.snapshotStore.SaveSnapshot(res); err != nil {
		return nil, err
	}

	if err := s.service.eventStore.Save(events); err != nil {
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

	s.unitOfWork.Start()
	aggregate := (*res)
	aggregate = aggregate.ChangeData(data)

	if aggregate.AggregateVersion()%5 == 1 {
		if err := s.service.snapshotStore.SaveSnapshot(aggregate); err != nil {
			return err
		}
	}

	es := aggregate.DomainEvents()
	events := make([]resource.IResourceEvent, 0)
	for i := 0; i < len(es); i++ {
		events = append(events, es[i])
	}

	if err := s.service.eventStore.Save(events); err != nil {
		s.unitOfWork.Rollback()
		return err
	}

	return s.unitOfWork.Commit()
}
