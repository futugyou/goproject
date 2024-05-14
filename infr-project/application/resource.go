package application

import (
	"errors"

	"github.com/futugyou/infr-project/resource"
)

type ResourceService struct {
	sourcer IEventSourcingService[resource.IResourceEvent, *resource.Resource]
}

func NewResourceService(sourcer IEventSourcingService[resource.IResourceEvent, *resource.Resource]) *ResourceService {
	return &ResourceService{
		sourcer: sourcer,
	}
}

func (s *ResourceService) CurrentResource(id string) (*resource.Resource, error) {
	res, err := s.sourcer.RetrieveLatestVersion(id)
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
	if err := s.sourcer.Save(events); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ResourceService) UpdateResourceDate(id string, data string) error {
	allVersions, _ := s.sourcer.RetrieveAllVersions(id)

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
	return s.sourcer.Save(events)
}
