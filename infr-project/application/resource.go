package application

import (
	"errors"

	infra "github.com/futugyou/infr-project/infrastructure"
	"github.com/futugyou/infr-project/resource"
)

type ResourceService struct {
	service *ApplicationService[resource.IResourceEvent, *resource.Resource]
}

func NewResourceService(eventStore infra.IEventStore[resource.IResourceEvent],
	snapshotStore infra.ISnapshotStore[*resource.Resource],
) *ResourceService {
	return &ResourceService{
		service: NewApplicationService(eventStore, snapshotStore, resource.ResourceFactory),
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
	res := resource.NewResource(name, resourceType, data)

	es := res.DomainEvents()
	events := make([]resource.IResourceEvent, 0)
	for i := 0; i < len(es); i++ {
		events = append(events, es[i])
	}
	if err := s.service.eventStore.Save(events); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ResourceService) UpdateResourceDate(id string, data string) error {
	allVersions, _ := s.service.RetrieveAllVersions(id)

	if len(allVersions) == 0 {
		return errors.New("no resource id by " + id)
	}

	res := allVersions[len(allVersions)-1]
	res = res.ChangeData(data)
	es := res.DomainEvents()
	events := make([]resource.IResourceEvent, 0)
	for i := 0; i < len(es); i++ {
		events = append(events, es[i])
	}
	return s.service.eventStore.Save(events)
}
