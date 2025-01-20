package resource

import (
	"fmt"
	"reflect"
	"time"

	domain "github.com/futugyou/infr-project/domain"
)

type IResourceEvent interface {
	domain.IDomainEvent
}

type ResourceEvent struct {
	domain.DomainEvent `bson:",inline"`
}

type ResourceCreatedEvent struct {
	ResourceEvent `bson:",inline"`
	Name          string   `bson:"name"`
	Type          string   `bson:"type"`
	Data          string   `bson:"data"`
	Tags          []string `bson:"tags"`
}

func (e ResourceCreatedEvent) EventType() string {
	return "ResourceCreated"
}

func NewResourceCreatedEvent(r *Resource) *ResourceCreatedEvent {
	return &ResourceCreatedEvent{
		ResourceEvent: ResourceEvent{
			DomainEvent: domain.DomainEvent{
				Id:              r.Id,
				ResourceVersion: r.Version,
				CreatedAt:       r.CreatedAt,
			},
		},
		Name: r.Name,
		Type: r.Type.String(),
		Data: r.Data,
		Tags: r.Tags,
	}
}

// Deprecated: Use a specific resource event type, cannot delete because data already exists
type ResourceUpdatedEvent struct {
	ResourceEvent `bson:",inline"`
	Name          string   `bson:"name"`
	Type          string   `bson:"type"`
	Data          string   `bson:"data"`
	Tags          []string `bson:"tags"`
}

func (e ResourceUpdatedEvent) EventType() string {
	return "ResourceUpdated"
}

type ResourceDeletedEvent struct {
	ResourceEvent `bson:",inline"`
}

func (e ResourceDeletedEvent) EventType() string {
	return "ResourceDeleted"
}

func NewResourceDeletedEvent(r *Resource) *ResourceDeletedEvent {
	return &ResourceDeletedEvent{
		ResourceEvent: ResourceEvent{
			DomainEvent: domain.DomainEvent{
				Id:              r.Id,
				ResourceVersion: r.Version,
				CreatedAt:       time.Now().UTC(),
			},
		},
	}
}

type ResourceNameChangedEvent struct {
	ResourceEvent `bson:",inline"`
	Name          string `bson:"name"`
}

func (e ResourceNameChangedEvent) EventType() string {
	return "ResourceNameChanged"
}

func NewResourceNameChangedEvent(r *Resource) *ResourceNameChangedEvent {
	return &ResourceNameChangedEvent{
		ResourceEvent: ResourceEvent{
			DomainEvent: domain.DomainEvent{
				Id:              r.Id,
				ResourceVersion: r.Version,
				CreatedAt:       time.Now().UTC(),
			},
		},
		Name: r.Name,
	}
}

type ResourceDataChangedEvent struct {
	ResourceEvent `bson:",inline"`
	Data          string `bson:"data"`
}

func (e ResourceDataChangedEvent) EventType() string {
	return "ResourceDataChanged"
}

func NewResourceDataChangedEvent(r *Resource) *ResourceDataChangedEvent {
	return &ResourceDataChangedEvent{
		ResourceEvent: ResourceEvent{
			DomainEvent: domain.DomainEvent{
				Id:              r.Id,
				ResourceVersion: r.Version,
				CreatedAt:       time.Now().UTC(),
			},
		},
		Data: r.Data,
	}
}

type ResourceTagsChangedEvent struct {
	ResourceEvent `bson:",inline"`
	Tags          []string `bson:"tags"`
}

func (e ResourceTagsChangedEvent) EventType() string {
	return "ResourceTagsChanged"
}

func NewResourceTagsChangedEvent(r *Resource) *ResourceTagsChangedEvent {
	return &ResourceTagsChangedEvent{
		ResourceEvent: ResourceEvent{
			DomainEvent: domain.DomainEvent{
				Id:              r.Id,
				ResourceVersion: r.Version,
				CreatedAt:       time.Now().UTC(),
			},
		},
		Tags: r.Tags,
	}
}

type ResourceTypeChangedEvent struct {
	ResourceEvent `bson:",inline"`
	Type          string `bson:"type"`
}

func (e ResourceTypeChangedEvent) EventType() string {
	return "ResourceTypeChanged"
}

func NewResourceTypeChangedEvent(r *Resource) *ResourceTypeChangedEvent {
	return &ResourceTypeChangedEvent{
		ResourceEvent: ResourceEvent{
			DomainEvent: domain.DomainEvent{
				Id:              r.Id,
				ResourceVersion: r.Version,
				CreatedAt:       time.Now().UTC(),
			},
		},
		Type: r.Type.String(),
	}
}

func CreateEvent(eventType string) (IResourceEvent, error) {
	eventTypeReflect, ok := eventNameMappingEventTypes[eventType]
	if !ok {
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}

	eventValue := reflect.New(eventTypeReflect.Elem()).Interface().(IResourceEvent)
	return eventValue, nil
}
