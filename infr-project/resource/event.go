package resource

import (
	"fmt"
	"time"

	domain "github.com/futugyou/infr-project/domain"
)

type IResourceEvent interface {
	domain.IDomainEvent
}

type ResourceEvent struct {
	domain.DomainEvent `bson:",inline" json:",inline"`
}

func (d ResourceEvent) AggregateEventName() string {
	return "resource_events"
}

type ResourceCreatedEvent struct {
	ResourceEvent `bson:",inline" json:",inline"`
	Name          string    `bson:"name" json:"name"`
	Type          string    `bson:"type" json:"type"`
	Data          string    `bson:"data" json:"data"`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
}

func (e ResourceCreatedEvent) EventType() string {
	return "ResourceCreated"
}

type ResourceUpdatedEvent struct {
	ResourceEvent `bson:",inline" json:",inline"`
	Name          string    `bson:"name" json:"name"`
	Type          string    `bson:"type" json:"type"`
	Data          string    `bson:"data" json:"data"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
}

func (e ResourceUpdatedEvent) EventType() string {
	return "ResourceUpdated"
}

type ResourceDeletedEvent struct {
	ResourceEvent `bson:",inline" json:",inline"`
}

func (e ResourceDeletedEvent) EventType() string {
	return "ResourceDeleted"
}

func CreateEvent(eventType string) (IResourceEvent, error) {
	switch eventType {
	case "ResourceCreated":
		return &ResourceCreatedEvent{}, nil
	case "ResourceUpdated":
		return &ResourceUpdatedEvent{}, nil
	case "ResourceDeleted":
		return &ResourceDeletedEvent{}, nil
	default:
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}
}
