package services

import (
	"fmt"
	"time"

	eventsourcing "github.com/futugyou/infr-project/event_sourcing"
)

type IResourceEvent interface {
	eventsourcing.IEvent
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
	aggregate.Apply(event)
	return aggregate
}

func (res *ResourceEventSourcer) GetAllVersions(id string) ([]Resource, error) {
	if len(res.allVersions) > 0 {
		return res.allVersions, nil
	}

	if len(res.events) == 0 {
		if _, err := res.Load(id); err != nil {
			return []Resource{}, err
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

	return res.allVersions, nil
}

func (res *ResourceEventSourcer) GetSpecificVersion(id string, version int) (*Resource, error) {
	if len(res.allVersions) == 0 {
		res.GetAllVersions(id)
	}

	for i := 0; i < len(res.allVersions); i++ {
		if res.allVersions[i].Version == version {
			return &res.allVersions[i], nil
		}
	}
	return nil, fmt.Errorf("not found with id:%s version:%d", id, version)
}

type ResourceEventSourcerWithSnapshot struct {
	ResourceEventSourcer
	// other
}

func (res *ResourceEventSourcerWithSnapshot) TakeSnapshot() error {
	// create snapshot
	return nil
}

func (res *ResourceEventSourcerWithSnapshot) RestoreFromSnapshot() error {
	// restore
	return nil
}
