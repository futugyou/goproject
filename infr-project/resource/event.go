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

type ResourceNameChangedEvent struct {
	ResourceEvent `bson:",inline" json:",inline"`
	Name          string `bson:"name" json:"name"`
}

func (e ResourceNameChangedEvent) EventType() string {
	return "ResourceNameChanged"
}

type ResourceDataChangedEvent struct {
	ResourceEvent `bson:",inline" json:",inline"`
	Data          string `bson:"data" json:"data"`
}

func (e ResourceDataChangedEvent) EventType() string {
	return "ResourceDataChanged"
}

type ResourceTagsChangedEvent struct {
	ResourceEvent `bson:",inline" json:",inline"`
	Tags          []string `bson:"tags" json:"tags"`
}

func (e ResourceTagsChangedEvent) EventType() string {
	return "ResourceTagsChanged"
}

type ResourceTypeChangedEvent struct {
	ResourceEvent `bson:",inline" json:",inline"`
	Type          string `bson:"type" json:"type"`
}

func (e ResourceTypeChangedEvent) EventType() string {
	return "ResourceTypeChanged"
}

func CreateEvent(eventType string) (IResourceEvent, error) {
	switch eventType {
	case "ResourceCreated":
		return &ResourceCreatedEvent{}, nil
	case "ResourceUpdated":
		return &ResourceUpdatedEvent{}, nil
	case "ResourceDeleted":
		return &ResourceDeletedEvent{}, nil
	case "ResourceNameChanged":
		return &ResourceNameChangedEvent{}, nil
	case "ResourceDataChanged":
		return &ResourceDataChangedEvent{}, nil
	case "ResourceTagsChanged":
		return &ResourceTagsChangedEvent{}, nil
	case "ResourceTypeChanged":
		return &ResourceTypeChangedEvent{}, nil
	default:
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}
}
