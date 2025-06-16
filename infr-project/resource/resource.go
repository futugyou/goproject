package resource

import (
	"errors"
	"fmt"
	"reflect"
	"time"

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

func NewResource(name string, resourceType ResourceType, data string, imageData string, tags []string) *Resource {
	r := &Resource{}
	event := NewResourceCreatedEvent(name, resourceType, data, imageData, tags)
	r.AddDomainEvent(event)
	r.Apply(event)
	return r
}

func (r *Resource) stateCheck() error {
	if r.IsDeleted {
		return fmt.Errorf("id: %s is alrealdy deleted", r.Id)
	}

	return nil
}

// Deprecated: Use ChangeResource.
func (r *Resource) ChangeName(name string) (*Resource, error) {
	if err := r.stateCheck(); err != nil {
		return r, err
	}

	event := NewResourceNameChangedEvent(r.Id, name, r.Version)
	r.raise(event)
	return r, nil
}

// Discuss: type should not be changed after creation
func (r *Resource) ChangeType(resourceType ResourceType, data string) (*Resource, error) {
	if err := r.stateCheck(); err != nil {
		return r, err
	}

	event := NewResourceTypeChangedEvent(r.Id, r.Version, resourceType, data)
	r.raise(event)
	return r, nil
}

// Deprecated: Use ChangeResource.
func (r *Resource) ChangeData(data string) (*Resource, error) {
	if err := r.stateCheck(); err != nil {
		return r, err
	}

	event := NewResourceDataChangedEvent(r.Id, r.Version, data)
	r.raise(event)
	return r, nil
}

// Deprecated: Use ChangeResource.
func (r *Resource) ChangeTags(tags []string) (*Resource, error) {
	if err := r.stateCheck(); err != nil {
		return r, err
	}

	event := NewResourceTagsChangedEvent(r.Id, r.Version, tags)
	r.raise(event)
	return r, nil
}

func (r *Resource) ChangeResource(name string, resourceType ResourceType, data string, imageData string, tags []string) (*Resource, error) {
	if err := r.stateCheck(); err != nil {
		return r, err
	}

	event := NewResourceUpdatedEvent(r.Id, r.Version, name, resourceType, data, imageData, tags)
	r.raise(event)
	return r, nil
}

func (r *Resource) DeleteResource() (*Resource, error) {
	if err := r.stateCheck(); err != nil {
		return r, err
	}

	event := NewResourceDeletedEvent(r.Id, r.Version)
	r.raise(event)
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

func (r *Resource) raise(event IResourceEvent) {
	r.AddDomainEvent(event)
	r.Apply(event)
}
