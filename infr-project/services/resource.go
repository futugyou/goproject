package services

import "time"

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
