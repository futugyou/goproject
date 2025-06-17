package resource

import (
	"fmt"
	"reflect"
	"time"

	domain "github.com/futugyou/infr-project/domain"
	"github.com/google/uuid"
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
	ImageData     string   `bson:"imageData"`
}

func (e ResourceCreatedEvent) EventType() string {
	return "ResourceCreated"
}

func NewResourceCreatedEvent(name string, resourceType ResourceType, data string, imageData string, tags []string) *ResourceCreatedEvent {
	return &ResourceCreatedEvent{
		ResourceEvent: ResourceEvent{
			DomainEvent: domain.DomainEvent{
				Id:              uuid.New().String(),
				ResourceVersion: 1,
				CreatedAt:       time.Now().UTC(),
			},
		},
		Name:      name,
		Type:      resourceType.String(),
		Data:      data,
		Tags:      tags,
		ImageData: imageData,
	}
}

// Deprecated: Use a specific resource event type, cannot delete because data already exists
type ResourceUpdatedEvent struct {
	ResourceEvent `bson:",inline"`
	Name          string `bson:"name"`
	// It is usually not recommended to modify the Type
	Type      string   `bson:"type"`
	Data      string   `bson:"data"`
	ImageData string   `bson:"imageData"`
	Tags      []string `bson:"tags"`
}

func (e ResourceUpdatedEvent) EventType() string {
	return "ResourceUpdated"
}

func NewResourceUpdatedEvent(id string, version int, name string, resourceType ResourceType, data string, imageData string, tags []string) *ResourceUpdatedEvent {
	return &ResourceUpdatedEvent{
		ResourceEvent: ResourceEvent{DomainEvent: domain.DomainEvent{Id: id, ResourceVersion: version + 1, CreatedAt: time.Now().UTC()}},
		Name:          name,
		Type:          resourceType.String(),
		Data:          data,
		ImageData:     imageData,
		Tags:          tags,
	}
}

type ResourceDeletedEvent struct {
	ResourceEvent `bson:",inline"`
}

func (e ResourceDeletedEvent) EventType() string {
	return "ResourceDeleted"
}

func NewResourceDeletedEvent(id string, version int) *ResourceDeletedEvent {
	return &ResourceDeletedEvent{
		ResourceEvent: ResourceEvent{
			DomainEvent: domain.DomainEvent{
				Id:              id,
				ResourceVersion: version + 1,
				CreatedAt:       time.Now().UTC(),
			},
		},
	}
}

// Deprecated: Use ResourceUpdatedEvent.
type ResourceNameChangedEvent struct {
	ResourceEvent `bson:",inline"`
	Name          string `bson:"name"`
}

func (e ResourceNameChangedEvent) EventType() string {
	return "ResourceNameChanged"
}

func NewResourceNameChangedEvent(id string, name string, version int) *ResourceNameChangedEvent {
	return &ResourceNameChangedEvent{
		ResourceEvent: ResourceEvent{
			DomainEvent: domain.DomainEvent{
				Id:              id,
				ResourceVersion: version + 1,
				CreatedAt:       time.Now().UTC(),
			},
		},
		Name: name,
	}
}

// Deprecated: Use ResourceUpdatedEvent.
type ResourceDataChangedEvent struct {
	ResourceEvent `bson:",inline"`
	Data          string `bson:"data"`
}

func (e ResourceDataChangedEvent) EventType() string {
	return "ResourceDataChanged"
}

func NewResourceDataChangedEvent(id string, version int, data string) *ResourceDataChangedEvent {
	return &ResourceDataChangedEvent{
		ResourceEvent: ResourceEvent{
			DomainEvent: domain.DomainEvent{
				Id:              id,
				ResourceVersion: version + 1,
				CreatedAt:       time.Now().UTC(),
			},
		},
		Data: data,
	}
}

// Deprecated: Use ResourceUpdatedEvent.
type ResourceTagsChangedEvent struct {
	ResourceEvent `bson:",inline"`
	Tags          []string `bson:"tags"`
}

func (e ResourceTagsChangedEvent) EventType() string {
	return "ResourceTagsChanged"
}

func NewResourceTagsChangedEvent(id string, version int, tags []string) *ResourceTagsChangedEvent {
	return &ResourceTagsChangedEvent{
		ResourceEvent: ResourceEvent{
			DomainEvent: domain.DomainEvent{
				Id:              id,
				ResourceVersion: version + 1,
				CreatedAt:       time.Now().UTC(),
			},
		},
		Tags: tags,
	}
}

type ResourceTypeChangedEvent struct {
	ResourceEvent `bson:",inline"`
	Type          string `bson:"type"`
	Data          string `bson:"data"`
}

func (e ResourceTypeChangedEvent) EventType() string {
	return "ResourceTypeChanged"
}

func NewResourceTypeChangedEvent(id string, version int, typee ResourceType, data string) *ResourceTypeChangedEvent {
	return &ResourceTypeChangedEvent{
		ResourceEvent: ResourceEvent{
			DomainEvent: domain.DomainEvent{
				Id:              id,
				ResourceVersion: version + 1,
				CreatedAt:       time.Now().UTC(),
			},
		},
		Type: typee.String(),
		Data: data,
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
