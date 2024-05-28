package resource

import (
	"errors"
	"fmt"
	"time"

	domain "github.com/futugyou/infr-project/domain"
	"github.com/google/uuid"
)

type Resource struct {
	domain.AggregateWithEventSourcing `json:"-"`
	Name                              string       `json:"name"`
	Type                              ResourceType `json:"type"`
	Data                              string       `json:"data"`
	IsDelete                          bool         `json:"is_delete"`
	CreatedAt                         time.Time    `json:"created_at"`
}

// ResourceType is the interface for resource types.
type ResourceType interface {
	privateResourceType() // Prevents external implementation
	String() string
}

// resourceType is the underlying implementation for ResourceType.
type resourceType string

// privateResourceType makes resourceType implement ResourceType.
func (c resourceType) privateResourceType() {}

// String makes resourceType implement ResourceType.
func (c resourceType) String() string {
	return string(c)
}

// Constants for the different resource types.
const (
	DrawIO     resourceType = "DrawIO"
	Markdown   resourceType = "Markdown"
	Excalidraw resourceType = "Excalidraw"
	Plate      resourceType = "Plate"
)

func GetResourceType(rType string) ResourceType {
	switch rType {
	case "DrawIO":
		return DrawIO
	case "Markdown":
		return Markdown
	case "Excalidraw":
		return Excalidraw
	case "Plate":
		return Plate
	default:
		return Markdown
	}
}

func NewResource(name string, resourceType ResourceType, data string) *Resource {
	r := &Resource{
		AggregateWithEventSourcing: domain.AggregateWithEventSourcing{
			Aggregate: domain.Aggregate{
				Id:   uuid.New().String(),
			},
			Version: 1,
		},
		Name: name,
		Type:      resourceType,
		Data:      data,
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

func (r *Resource) DeleteResource() (*Resource, error) {
	if err := r.stateCheck(); err != nil {
		return r, err
	}
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
	case *ResourceUpdatedEvent:
		r.Id = e.Id
		r.Name = e.Name
		r.Type = GetResourceType(e.Type)
		r.Version = e.Version()
		r.Data = e.Data
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
	}

	r.AddDomainEvent(event)
}

func (r *Resource) createUpdatedEvent() {
	event := ResourceUpdatedEvent{
		ResourceEvent: ResourceEvent{
			DomainEvent: domain.DomainEvent{
				Id:              r.Id,
				ResourceVersion: r.Version,
			},
		},
		Name: r.Name,
		Type: r.Type.String(),
		Data: r.Data,
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
