package services

import (
	"encoding/json"
	"errors"
	"time"

	eventsourcing "github.com/futugyou/infr-project/event_sourcing"
	"github.com/google/uuid"
)

type Resource struct {
	eventsourcing.IAggregate `json:"-"`
	Id                       string       `json:"id"`
	Name                     string       `json:"name"`
	Type                     ResourceType `json:"type"`
	Version                  int          `json:"version"`
	Data                     string       `json:"data"`
	CreatedAt                time.Time    `json:"created_at"`
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
	return &Resource{
		Id:        uuid.New().String(),
		Name:      name,
		Type:      resourceType,
		Version:   1,
		Data:      data,
		CreatedAt: time.Now().UTC(),
	}
}

func (r *Resource) ChangeName(name string) *Resource {
	r.Version = r.Version + 1
	r.CreatedAt = time.Now().UTC()
	r.Name = name
	return r
}

func (r *Resource) ChangeType(resourceType ResourceType, data string) *Resource {
	r.Version = r.Version + 1
	r.CreatedAt = time.Now().UTC()
	r.Type = resourceType
	r.Data = data
	return r
}

func (r *Resource) ChangeData(data string) *Resource {
	r.Version = r.Version + 1
	r.CreatedAt = time.Now().UTC()
	r.Data = data
	return r
}

func (r Resource) AggregateName() string {
	return "resources"
}

func (r Resource) AggregateId() string {
	return r.Id
}
func (r Resource) AggregateVersion() int {
	return r.Version
}

func (r Resource) Apply(event eventsourcing.IDomainEvent) (eventsourcing.IEventSourcing, error) {
	switch e := event.(type) {
	case ResourceCreatedEvent:
		return Resource{Id: e.Id, Name: e.Name, Type: e.Type, Data: e.Data, Version: 0, CreatedAt: e.CreatedAt}, nil
	case ResourceUpdatedEvent:
		return Resource{Id: e.Id, Name: e.Name, Type: e.Type, Data: e.Data, Version: e.Version(), CreatedAt: e.UpdatedAt}, nil
	case ResourceDeletedEvent:
		// TODO: how to handle delete
	}

	return r, errors.New("event type not supported")
}

func CreateCreatedEvent(resource Resource) IResourceEvent {
	event := ResourceCreatedEvent{
		Id:        resource.Id,
		Name:      resource.Name,
		Type:      resource.Type,
		Data:      resource.Data,
		CreatedAt: resource.CreatedAt,
	}

	return event
}

func CreateUpdatedEvent(resource Resource) IResourceEvent {
	event := ResourceUpdatedEvent{
		Id:              resource.Id,
		Name:            resource.Name,
		Type:            resource.Type,
		Data:            resource.Data,
		ResourceVersion: resource.Version,
		UpdatedAt:       time.Now().UTC(),
	}
	return event
}

func CreateDeletedEvent(resource Resource) IResourceEvent {
	event := ResourceDeletedEvent{
		Id: resource.Id,
	}
	return event
}

type ResourceService struct {
}

func (s *ResourceService) CurrentResource(id string) Resource {
	var sourcer eventsourcing.IEventSourcer[IResourceEvent, Resource] = eventsourcing.NewEventSourcer[IResourceEvent, Resource]()
	allVersions, _ := sourcer.GetAllVersions(id)
	return allVersions[len(allVersions)-1]
}

func (s *ResourceService) CreateResource(name string, resourceType ResourceType, data string) (*Resource, error) {
	resource := NewResource(name, resourceType, data)
	var sourcer eventsourcing.IEventSourcer[IResourceEvent, Resource] = eventsourcing.NewEventSourcer[IResourceEvent, Resource]()

	evt := CreateCreatedEvent(*resource)

	if err := sourcer.Save([]IResourceEvent{evt}); err != nil {
		return nil, err
	}

	return resource, nil
}

func (s *ResourceService) UpdateResourceDate(id string, data string) error {
	var sourcer eventsourcing.IEventSourcer[IResourceEvent, Resource] = eventsourcing.NewEventSourcer[IResourceEvent, Resource]()
	allVersions, _ := sourcer.GetAllVersions(id)
	if len(allVersions) == 0 {
		return errors.New("no resource id by " + id)
	}

	resource := allVersions[len(allVersions)-1]
	resource = *resource.ChangeData(data)
	evt := CreateUpdatedEvent(resource)

	return sourcer.Save([]IResourceEvent{evt})
}
