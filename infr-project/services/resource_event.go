package services

import (
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
