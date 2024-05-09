package services

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Resource struct {
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
func (r Resource) MarshalJSON() ([]byte, error) {
	type Alias Resource
	return json.Marshal(&struct {
		Type string `json:"type"`
		*Alias
	}{
		Type:  r.Type.String(),
		Alias: (*Alias)(&r),
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
		Version:   0,
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
