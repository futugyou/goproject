package resource

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"

	domain "github.com/futugyou/infr-project/domain"
)

type Resource struct {
	domain.AggregateWithEventSourcing
	Name      string
	Type      ResourceType
	Data      string
	ImageData string
	Tags      []string
	IsDeleted bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewResource(name string, resourceType ResourceType, data string, imageData string,tags []string) *Resource {
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
		ImageData: imageData,
		Tags:      tags,
		CreatedAt: time.Now().UTC(),
	}

	event := NewResourceCreatedEvent(r)
	r.AddDomainEvent(event)
	return r
}

func (r *Resource) stateCheck() error {
	if r.IsDeleted {
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

// Discuss: type should not be changed after creation
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
	r.IsDeleted = true

	event := NewResourceDeletedEvent(r)
	r.AddDomainEvent(event)
	return r, nil
}

func (r *Resource) AggregateName() string {
	return "resources"
}

func (r *Resource) Apply(event domain.IDomainEvent) error {
	if processor, ok := ResourceEventProcessors[reflect.TypeOf(event)]; ok {
		processor.EventApply(r, event)
		return nil
	}

	return errors.New("event type not supported")
}
