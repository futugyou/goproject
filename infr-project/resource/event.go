package resource

import (
	"fmt"

	domain "github.com/futugyou/infr-project/domain"
)

type IResourceEvent interface {
	domain.IDomainEvent
}

type ResourceEvent struct {
	domain.DomainEvent `bson:",inline" json:",inline"`
}

type ResourceCreatedEvent struct {
	ResourceEvent `bson:",inline" json:",inline"`
	Name          string   `bson:"name" json:"name"`
	Type          string   `bson:"type" json:"type"`
	Data          string   `bson:"data" json:"data"`
	Tags          []string `bson:"tags" json:"tags"`
}

func (e ResourceCreatedEvent) EventType() string {
	return "ResourceCreated"
}

type ResourceUpdatedEvent struct {
	ResourceEvent `bson:",inline" json:",inline"`
	Name          string   `bson:"name" json:"name"`
	Type          string   `bson:"type" json:"type"`
	Data          string   `bson:"data" json:"data"`
	Tags          []string `bson:"tags" json:"tags"`
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
