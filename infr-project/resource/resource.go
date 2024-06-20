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
	UpdatedAt                         time.Time    `json:"updated_at"`
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

	event := NewResourceCreatedEvent(r)
	r.AddDomainEvent(event)
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

	event := NewResourceNameChangedEvent(r)
	r.AddDomainEvent(event)
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

	event := NewResourceTypeChangedEvent(r)
	r.AddDomainEvent(event)
	return r, nil
}

func (r *Resource) ChangeData(data string) (*Resource, error) {
	if err := r.stateCheck(); err != nil {
		return r, err
	}

	r.Version = r.Version + 1
	r.CreatedAt = time.Now().UTC()
	r.Data = data

	event := NewResourceDataChangedEvent(r)
	r.AddDomainEvent(event)
	return r, nil
}

func (r *Resource) ChangeTags(tags []string) (*Resource, error) {
	if err := r.stateCheck(); err != nil {
		return r, err
	}

	r.Version = r.Version + 1
	r.CreatedAt = time.Now().UTC()
	r.Tags = tags

	event := NewResourceTagsChangedEvent(r)
	r.AddDomainEvent(event)
	return r, nil
}

func (r *Resource) DeleteResource() (*Resource, error) {
	if err := r.stateCheck(); err != nil {
		return r, err
	}

	r.Version = r.Version + 1
	r.IsDelete = true

	event := NewResourceDeletedEvent(r)
	r.AddDomainEvent(event)
	return r, nil
}

func (r *Resource) AggregateName() string {
	return "resources"
}

func (r *Resource) Apply(event domain.IDomainEvent) error {
	switch e := event.(type) {
	case *ResourceCreatedEvent:
		r.Id = e.Id
		r.Version = e.Version()
		r.CreatedAt = e.CreatedAt
		r.Name = e.Name
		r.Type = GetResourceType(e.Type)
		r.Data = e.Data
		r.Tags = e.Tags
	case *ResourceUpdatedEvent:
		r.Id = e.Id
		r.Name = e.Name
		r.Type = GetResourceType(e.Type)
		r.Version = e.Version()
		r.Data = e.Data
		r.Tags = e.Tags
		r.UpdatedAt = e.CreatedAt
	case *ResourceNameChangedEvent:
		r.Id = e.Id
		r.Version = e.Version()
		r.Name = e.Name
		r.UpdatedAt = e.CreatedAt
	case *ResourceDataChangedEvent:
		r.Id = e.Id
		r.Version = e.Version()
		r.Data = e.Data
		r.UpdatedAt = e.CreatedAt
	case *ResourceTagsChangedEvent:
		r.Id = e.Id
		r.Version = e.Version()
		r.Tags = e.Tags
		r.UpdatedAt = e.CreatedAt
	case *ResourceTypeChangedEvent:
		r.Id = e.Id
		r.Type = GetResourceType(e.Type)
		r.Version = e.Version()
		r.UpdatedAt = e.CreatedAt
	case *ResourceDeletedEvent:
		r.Id = e.Id
		r.Version = e.Version()
		r.IsDelete = true
		r.UpdatedAt = e.CreatedAt
	}

	return errors.New("event type not supported")
}
