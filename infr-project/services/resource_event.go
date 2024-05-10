package services

import (
	"time"
)

type IResourceEvent interface {
	EventType() string
}

type ResourceCreatedEvent struct {
	Id        string
	Name      string
	Type      ResourceType
	Data      string
	CreatedAt time.Time
}

func (e ResourceCreatedEvent) EventType() string {
	return "ResourceCreated"
}

type ResourceUpdatedEvent struct {
	Id        string
	Name      string
	Type      ResourceType
	Data      string
	Version   int
	UpdatedAt time.Time
}

func (e ResourceUpdatedEvent) EventType() string {
	return "ResourceUpdated"
}

type ResourceDeletedEvent struct {
	Id string
}

func (e ResourceDeletedEvent) EventType() string {
	return "ResourceDeleted"
}

type ResourceEventSourcer struct {
	ResourceId  string
	events      []IResourceEvent
	allVersions []Resource
}

func (res *ResourceEventSourcer) Add(event IResourceEvent) error {
	res.events = append(res.events, event)
	return nil
}

func (res *ResourceEventSourcer) Save(events []IResourceEvent) error {
	// save to repo
	// res.Events
	return nil
}

func (res *ResourceEventSourcer) Load(id string) ([]IResourceEvent, error) {
	// load from repo
	// res.Events = ....
	return nil, nil
}

func (res *ResourceEventSourcer) Apply(aggregate Resource, event IResourceEvent) Resource {
	resource := aggregate

	switch e := event.(type) {
	case ResourceCreatedEvent:
		resource = Resource{Id: e.Id, Name: e.Name, Type: e.Type, Data: e.Data, Version: 1, CreatedAt: e.CreatedAt}
	case ResourceUpdatedEvent:
		resource.Name = e.Name
		resource.Type = e.Type
		resource.Data = e.Data
		resource.Version = e.Version
		resource.CreatedAt = e.UpdatedAt
	case ResourceDeletedEvent:
		// TODO: how to handle delete
	}

	return resource
}

func (res *ResourceEventSourcer) GetAlltVersions() []Resource {
	if len(res.allVersions) > 0 {
		return res.allVersions
	}

	if len(res.events) == 0 {
		if _, err := res.Load(res.ResourceId); err != nil {
			return res.allVersions
		}
	}

	resource := Resource{}

	for i := 0; i < len(res.events); i++ {
		prevVersion := resource.Version
		resource = res.Apply(resource, res.events[i])

		if resource.Version != prevVersion {
			res.allVersions = append(res.allVersions, resource)
		}
	}

	return res.allVersions
}

func (res *ResourceEventSourcer) GetSpecificVersion(version int) Resource {
	if len(res.allVersions) == 0 {
		res.GetAlltVersions()
	}

	for i := 0; i < len(res.allVersions); i++ {
		if res.allVersions[i].Version == version {
			return res.allVersions[i]
		}
	}
	return Resource{}
}
