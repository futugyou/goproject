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

func (res *ResourceEventSourcer) Apply(event IResourceEvent) Resource {
	resource := Resource{Id: res.ResourceId}
	return resource
}

func (res *ResourceEventSourcer) GetResourceVersions() []Resource {
	return res.allVersions
}

func (res *ResourceEventSourcer) GetResourceVersion(version int) Resource {
	return Resource{}
}
