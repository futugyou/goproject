package services

import (
	"time"

	"github.com/google/uuid"
)

type Resource struct {
	Id       string       `json:"id"`
	Name     string       `json:"name"`
	Type     ResourceType `json:"type"`
	Version  int          `json:"version"`
	Data     string       `json:"data"`
	CreateAt time.Time    `json:"create_at"`
}

type ResourceType interface {
	privateResourceType()
	String() string
}

type resourceType string

func (c resourceType) privateResourceType() {}
func (c resourceType) String() string {
	return string(c)
}

const DrawIO resourceType = "DrawIO"
const Markdown resourceType = "Markdown"
const Excalidraw resourceType = "Excalidraw"
const Plate resourceType = "Plate"

func NewResource(name string, resourceType ResourceType, data string) *Resource {
	return &Resource{
		Id:       uuid.New().String(),
		Name:     name,
		Type:     resourceType,
		Version:  0,
		Data:     data,
		CreateAt: time.Now().UTC(),
	}
}

func (r *Resource) ChangeName(name string) *Resource {
	r.Version = r.Version + 1
	r.CreateAt = time.Now().UTC()
	r.Name = name
	return r
}

func (r *Resource) ChangeType(resourceType ResourceType, data string) *Resource {
	r.Version = r.Version + 1
	r.CreateAt = time.Now().UTC()
	r.Type = resourceType
	r.Data = data
	return r
}

func (r *Resource) ChangeData(data string) *Resource {
	r.Version = r.Version + 1
	r.CreateAt = time.Now().UTC()
	r.Data = data
	return r
}
