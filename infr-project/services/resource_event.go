package services

import (
	"time"

	eventsourcing "github.com/futugyou/infr-project/event_sourcing"
)

type IResourceEvent interface {
	eventsourcing.IDomainEvent
}

type ResourceCreatedEvent struct {
	Id              string
	Name            string
	Type            ResourceType
	Data            string
	CreatedAt       time.Time
	ResourceVersion int
}

func (e ResourceCreatedEvent) EventType() string {
	return "ResourceCreated"
}

func (e ResourceCreatedEvent) Version() int {
	return e.ResourceVersion
}

type ResourceUpdatedEvent struct {
	Id              string
	Name            string
	Type            ResourceType
	Data            string
	UpdatedAt       time.Time
	ResourceVersion int
}

func (e ResourceUpdatedEvent) EventType() string {
	return "ResourceUpdated"
}

func (e ResourceUpdatedEvent) Version() int {
	return e.ResourceVersion
}

type ResourceDeletedEvent struct {
	Id              string
	ResourceVersion int
}

func (e ResourceDeletedEvent) EventType() string {
	return "ResourceDeleted"
}

func (e ResourceDeletedEvent) Version() int {
	return e.ResourceVersion
}
