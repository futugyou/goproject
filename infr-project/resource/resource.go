package resource

import (
	"encoding/json"
	"errors"
	"time"

	domain "github.com/futugyou/infr-project/domain"
	"github.com/google/uuid"
)

type Resource struct {
	domain.BaseEventSourcing
	Id        string       `json:"id"`
	Name      string       `json:"name"`
	Type      ResourceType `json:"type"`
	Version   int          `json:"version"`
	Data      string       `json:"data"`
	CreatedAt time.Time    `json:"created_at"`
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

// MarshalJSON is a custom marshaler for Resource that handles the serialization of ResourceType.
// In this case, we can skip MarshalJSON, only implement UnmarshalJSON
func (r *Resource) MarshalJSON() ([]byte, error) {
	type Alias Resource
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  r.Type.String(),
		Alias: (*Alias)(r),
	})
}

// UnmarshalJSON is a custom unmarshaler for Resource that handles the deserialization of ResourceType.
func (r *Resource) UnmarshalJSON(data []byte) error {
	type Alias Resource
	aux := &struct {
		Type string `json:"type"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch aux.Type {
	case string(DrawIO):
		r.Type = DrawIO
	case string(Markdown):
		r.Type = Markdown
	case string(Excalidraw):
		r.Type = Excalidraw
	case string(Plate):
		r.Type = Plate
	default:
		return json.Unmarshal(data, &r)
	}
	return nil
}

func NewResource(name string, resourceType ResourceType, data string) *Resource {
	r := &Resource{
		Id:        uuid.New().String(),
		Name:      name,
		Type:      resourceType,
		Version:   1,
		Data:      data,
		CreatedAt: time.Now().UTC(),
	}
	r.createCreatedEvent()
	return r
}

func (r *Resource) ChangeName(name string) *Resource {
	r.Version = r.Version + 1
	r.CreatedAt = time.Now().UTC()
	r.Name = name
	r.createUpdatedEvent()
	return r
}

func (r *Resource) ChangeType(resourceType ResourceType, data string) *Resource {
	r.Version = r.Version + 1
	r.CreatedAt = time.Now().UTC()
	r.Type = resourceType
	r.Data = data
	r.createUpdatedEvent()
	return r
}

func (r *Resource) ChangeData(data string) *Resource {
	r.Version = r.Version + 1
	r.CreatedAt = time.Now().UTC()
	r.Data = data
	r.createUpdatedEvent()
	return r
}

func (r *Resource) AggregateName() string {
	return "resources"
}

func (r *Resource) AggregateId() string {
	return r.Id
}

func (r *Resource) AggregateVersion() int {
	return r.Version
}

func (r *Resource) Apply(event domain.IDomainEvent) error {
	switch e := event.(type) {
	case ResourceCreatedEvent:
		r.Id = e.Id
		r.Name = e.Name
		r.Type = e.Type
		r.Version = e.Version()
		r.CreatedAt = e.CreatedAt
		r.Data = e.Data
	case ResourceUpdatedEvent:
		r.Id = e.Id
		r.Name = e.Name
		r.Type = e.Type
		r.Version = e.Version()
		r.CreatedAt = e.UpdatedAt
		r.Data = e.Data
	case ResourceDeletedEvent:
		// TODO: how to handle delete
	}

	return errors.New("event type not supported")
}

func (r *Resource) createCreatedEvent() {
	event := ResourceCreatedEvent{
		DomainEvent: domain.DomainEvent{
			Id:              r.Id,
			ResourceVersion: r.Version,
		},
		Name:      r.Name,
		Type:      r.Type,
		Data:      r.Data,
		CreatedAt: r.CreatedAt,
	}

	r.AddDomainEvent(event)
}

func (r *Resource) createUpdatedEvent() {
	event := ResourceUpdatedEvent{
		DomainEvent: domain.DomainEvent{
			Id:              r.Id,
			ResourceVersion: r.Version,
		},
		Name:      r.Name,
		Type:      r.Type,
		Data:      r.Data,
		UpdatedAt: time.Now().UTC(),
	}
	r.AddDomainEvent(event)
}

// func (r *Resource) createDeletedEvent() {
// 	event := ResourceDeletedEvent{
// 		Id: r.Id,
// 	}
// 	r.AddDomainEvent(event)
// }
