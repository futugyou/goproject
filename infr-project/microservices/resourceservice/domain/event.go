package domain

import (
	"fmt"
	"reflect"
	"time"

	"github.com/futugyou/domaincore/domain"
	"github.com/google/uuid"
)

type ResourceEvent interface {
	domain.DomainEvent
}

type BaseResourceEvent struct {
	domain.BaseDomainEvent `bson:",inline"`
}

type ResourceCreatedEvent struct {
	BaseResourceEvent `bson:",inline"`
	Name              string   `bson:"name"`
	Type              string   `bson:"type"`
	Data              string   `bson:"data"`
	Tags              []string `bson:"tags"`
	ImageData         string   `bson:"imageData"`
}

func (e ResourceCreatedEvent) EventType() string {
	return "ResourceCreated"
}

func NewResourceCreatedEvent(name string, resourceType ResourceType, data string, imageData string, tags []string) *ResourceCreatedEvent {
	return &ResourceCreatedEvent{
		BaseResourceEvent: BaseResourceEvent{
			BaseDomainEvent: domain.BaseDomainEvent{
				ID:              uuid.New().String(),
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

type ResourceUpdatedEvent struct {
	BaseResourceEvent `bson:",inline"`
	Name              string `bson:"name"`
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
		BaseResourceEvent: BaseResourceEvent{BaseDomainEvent: domain.BaseDomainEvent{ID: id, ResourceVersion: version + 1, CreatedAt: time.Now().UTC()}},
		Name:              name,
		Type:              resourceType.String(),
		Data:              data,
		ImageData:         imageData,
		Tags:              tags,
	}
}

type ResourceDeletedEvent struct {
	BaseResourceEvent `bson:",inline"`
}

func (e ResourceDeletedEvent) EventType() string {
	return "ResourceDeleted"
}

func NewResourceDeletedEvent(id string, version int) *ResourceDeletedEvent {
	return &ResourceDeletedEvent{
		BaseResourceEvent: BaseResourceEvent{
			BaseDomainEvent: domain.BaseDomainEvent{
				ID:              id,
				ResourceVersion: version + 1,
				CreatedAt:       time.Now().UTC(),
			},
		},
	}
}

func CreateEvent(eventType string) (ResourceEvent, error) {
	eventTypeReflect, ok := eventNameMappingEventTypes[eventType]
	if !ok {
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}

	eventValue := reflect.New(eventTypeReflect.Elem()).Interface().(ResourceEvent)
	return eventValue, nil
}
