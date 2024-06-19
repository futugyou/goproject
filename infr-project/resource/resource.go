package resource

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	domain "github.com/futugyou/infr-project/domain"
)

type Resource struct {
	domain.AggregateWithEventSourcing `json:"-"`
	Name                              string       `json:"name"`
	Type                              ResourceType `json:"type"`
	Data                              string       `json:"data"`
	Tags                              []string     `json:"tags"`
	IsDelete                          bool         `json:"is_deleted"`
	CreatedAt                         time.Time    `json:"created_at"`
}

func NewResource(name string, resourceType ResourceType, data string, tags []string) *Resource {
	r := &Resource{
		AggregateWithEventSourcing: domain.AggregateWithEventSourcing{
			Aggregate: domain.Aggregate{
				Id: uuid.New().String(),
			},
			Version: 1,
		},
		Name:      name,
		Type:      resourceType,
		Data:      data,
		Tags:      tags,
		CreatedAt: time.Now().UTC(),
	}
	r.createCreatedEvent()
	return r
}

func (r *Resource) stateCheck() error {
	if r.IsDelete {
		return fmt.Errorf("id: %s is alrealdy deleted", r.Id)
	}

	return nil
}

func (r *Resource) ChangeName(name string) (*Resource, error) {
	if err := r.stateCheck(); err != nil {
		return r, err
	}
	r.Version = r.Version + 1
	r.CreatedAt = time.Now().UTC()
	r.Name = name
	r.createUpdatedEvent()
	return r, nil
}

func (r *Resource) ChangeType(resourceType ResourceType, data string) (*Resource, error) {
	if err := r.stateCheck(); err != nil {
		return r, err
	}
	r.Version = r.Version + 1
	r.CreatedAt = time.Now().UTC()
	r.Type = resourceType
	r.Data = data
	r.createUpdatedEvent()
	return r, nil
}

func (r *Resource) ChangeData(data string) (*Resource, error) {
	if err := r.stateCheck(); err != nil {
		return r, err
	}
	r.Version = r.Version + 1
	r.CreatedAt = time.Now().UTC()
	r.Data = data
	r.createUpdatedEvent()
	return r, nil
}

func (r *Resource) ChangeTags(tags []string) (*Resource, error) {
	if err := r.stateCheck(); err != nil {
		return r, err
	}
	r.Version = r.Version + 1
	r.CreatedAt = time.Now().UTC()
	r.Tags = tags
	r.createUpdatedEvent()
	return r, nil
}

func (r *Resource) DeleteResource() (*Resource, error) {
	if err := r.stateCheck(); err != nil {
		return r, err
	}
	r.Version = r.Version + 1
	r.IsDelete = true
	r.createDeletedEvent()
	return r, nil
}

func (r *Resource) AggregateName() string {
	return "resources"
}

func (r *Resource) Apply(event domain.IDomainEvent) error {
	switch e := event.(type) {
	case *ResourceCreatedEvent:
		r.Id = e.Id
		r.Name = e.Name
		r.Type = GetResourceType(e.Type)
		r.Version = e.Version()
		r.CreatedAt = e.CreatedAt
		r.Data = e.Data
		r.Tags = e.Tags
	case *ResourceUpdatedEvent:
		r.Id = e.Id
		r.Name = e.Name
		r.Type = GetResourceType(e.Type)
		r.Version = e.Version()
		r.Data = e.Data
		r.Tags = e.Tags
	case *ResourceDeletedEvent:
		r.IsDelete = true
	}

	return errors.New("event type not supported")
}

func (r *Resource) createCreatedEvent() {
	event := ResourceCreatedEvent{
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

	r.AddDomainEvent(event)
}

func (r *Resource) createUpdatedEvent() {
	event := ResourceUpdatedEvent{
		ResourceEvent: ResourceEvent{
			DomainEvent: domain.DomainEvent{
				Id:              r.Id,
				ResourceVersion: r.Version,
				CreatedAt:       time.Now().UTC(),
			},
		},
		Name: r.Name,
		Type: r.Type.String(),
		Data: r.Data,
		Tags: r.Tags,
	}
	r.AddDomainEvent(event)
}

func (r *Resource) createDeletedEvent() {
	event := ResourceDeletedEvent{
		ResourceEvent: ResourceEvent{
			DomainEvent: domain.DomainEvent{
				Id:              r.Id,
				ResourceVersion: r.Version,
				CreatedAt:       time.Now().UTC(),
			},
		},
	}
	r.AddDomainEvent(event)
}
