package application

import (
	"errors"

	services "github.com/futugyou/infr-project/services"
)

type ResourceService struct {
}

func (s *ResourceService) CurrentResource(id string) services.Resource {
	var sourcer IEventSourcingService[services.IResourceEvent, *services.Resource] = NewEventSourcer[services.IResourceEvent, *services.Resource]()
	allVersions, _ := sourcer.RetrieveAllVersions(id)
	return *allVersions[len(allVersions)-1]
}

func (s *ResourceService) CreateResource(name string, resourceType services.ResourceType, data string) (*services.Resource, error) {
	resource := services.NewResource(name, resourceType, data)
	var sourcer IEventSourcingService[services.IResourceEvent, *services.Resource] = NewEventSourcer[services.IResourceEvent, *services.Resource]()

	es := resource.DomainEvents()
	events := make([]services.IResourceEvent, 0)
	for i := 0; i < len(es); i++ {
		events = append(events, es[i])
	}
	if err := sourcer.Save(events); err != nil {
		return nil, err
	}

	return resource, nil
}

func (s *ResourceService) UpdateResourceDate(id string, data string) error {
	var sourcer IEventSourcingService[services.IResourceEvent, *services.Resource] = NewEventSourcer[services.IResourceEvent, *services.Resource]()
	allVersions, _ := sourcer.RetrieveAllVersions(id)
	if len(allVersions) == 0 {
		return errors.New("no resource id by " + id)
	}

	resource := allVersions[len(allVersions)-1]
	resource = resource.ChangeData(data)
	es := resource.DomainEvents()
	events := make([]services.IResourceEvent, 0)
	for i := 0; i < len(es); i++ {
		events = append(events, es[i])
	}
	return sourcer.Save(events)
}
