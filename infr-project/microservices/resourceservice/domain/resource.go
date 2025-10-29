package domain

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/futugyou/domaincore/domain"
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
	r.raise(event)
	return r
}

func (r *Resource) Clone() *Resource {
	return &Resource{
		AggregateWithEventSourcing: domain.AggregateWithEventSourcing{
			Version: r.Version,
			Aggregate: domain.Aggregate{
				ID: r.ID,
			},
		},
		Name:      r.Name,
		Type:      r.Type,
		Data:      r.Data,
		ImageData: r.ImageData,
		Tags:      r.Tags,
		IsDeleted: r.IsDeleted,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func (r *Resource) stateCheck() error {
	if r.IsDeleted {
		return fmt.Errorf("id: %s is alrealdy deleted", r.ID)
	}

	return nil
}

func (r *Resource) ChangeResource(name string, resourceType ResourceType, data string, imageData string, tags []string) (*Resource, error) {
	if err := r.stateCheck(); err != nil {
		return r, err
	}

	event := NewResourceUpdatedEvent(r.ID, r.Version, name, resourceType, data, imageData, tags)
	err := r.raise(event)
	return r, err
}

func (r *Resource) DeleteResource() (*Resource, error) {
	if err := r.stateCheck(); err != nil {
		return r, err
	}

	event := NewResourceDeletedEvent(r.ID, r.Version)
	err := r.raise(event)
	return r, err
}

func (r *Resource) AggregateName() string {
	return "resources"
}

func (r *Resource) Apply(event domain.DomainEvent) error {
	if processor, ok := ResourceEventProcessors[reflect.TypeOf(event)]; ok {
		processor.EventApply(r, event)
		return nil
	}

	return errors.New("event type not supported")
}

func (r *Resource) Replay(events []domain.DomainEvent) error {
	for _, event := range events {
		if err := r.Apply(event); err != nil {
			return err
		}
	}

	return nil
}

func (r *Resource) raise(event ResourceEvent) error {
	r.AddDomainEvent(event)
	return r.Apply(event)
}
