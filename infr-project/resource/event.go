package resource

import (
	"time"

	domain "github.com/futugyou/infr-project/domain"
)

type IResourceEvent interface {
	domain.IDomainEvent
}

type ResourceCreatedEvent struct {
	domain.DomainEvent
	Name      string
	Type      ResourceType
	Data      string
	CreatedAt time.Time
}

func (e ResourceCreatedEvent) EventType() string {
	return "ResourceCreated"
}

type ResourceUpdatedEvent struct {
	domain.DomainEvent
	Name      string
	Type      ResourceType
	Data      string
	UpdatedAt time.Time
}

func (e ResourceUpdatedEvent) EventType() string {
	return "ResourceUpdated"
}

type ResourceDeletedEvent struct {
	domain.DomainEvent
}

func (e ResourceDeletedEvent) EventType() string {
	return "ResourceDeleted"
}
