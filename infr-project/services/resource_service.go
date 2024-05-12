package services

import (
	"errors"

	domain "github.com/futugyou/infr-project/domain"
)

type ResourceService struct {
}

func (s *ResourceService) CurrentResource(id string) Resource {
	var sourcer domain.IEventSourcer[IResourceEvent, *Resource] = domain.NewEventSourcer[IResourceEvent, *Resource]()
	allVersions, _ := sourcer.GetAllVersions(id)
	return *allVersions[len(allVersions)-1]
}

func (s *ResourceService) CreateResource(name string, resourceType ResourceType, data string) (*Resource, error) {
	resource := NewResource(name, resourceType, data)
	var sourcer domain.IEventSourcer[IResourceEvent, *Resource] = domain.NewEventSourcer[IResourceEvent, *Resource]()

	if err := sourcer.Save(resource.domainEvents); err != nil {
		return nil, err
	}

	return resource, nil
}

func (s *ResourceService) UpdateResourceDate(id string, data string) error {
	var sourcer domain.IEventSourcer[IResourceEvent, *Resource] = domain.NewEventSourcer[IResourceEvent, *Resource]()
	allVersions, _ := sourcer.GetAllVersions(id)
	if len(allVersions) == 0 {
		return errors.New("no resource id by " + id)
	}

	resource := allVersions[len(allVersions)-1]
	resource = resource.ChangeData(data)

	return sourcer.Save(resource.domainEvents)
}
