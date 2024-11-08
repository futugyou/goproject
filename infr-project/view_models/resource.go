package viewmodels

import "time"

type CreateResourceRequest struct {
	Name string   `json:"name" validate:"required,min=3,max=50"`
	Type string   `json:"type" validate:"oneof=DrawIO Markdown Excalidraw Plate"`
	Data string   `json:"data" validate:"required,min=3"`
	Tags []string `json:"tags"`
}

type CreateResourceResponse struct {
	Id string `json:"id"`
}

type UpdateResourceRequest struct {
	Name string   `json:"name" validate:"required,min=3,max=50"`
	Data string   `json:"data" validate:"required,min=3"`
	Tags []string `json:"tags"`
}

type ResourceView struct {
	Id        string    `json:"id" redis:"id"`
	Name      string    `json:"name" redis:"name"`
	Type      string    `json:"type" redis:"type"`
	Data      string    `json:"data" redis:"data"`
	Version   int       `json:"version" redis:"version"`
	IsDelete  bool      `json:"is_deleted" redis:"is_deleted"`
	CreatedAt time.Time `json:"created_at" redis:"created_at"`
	UpdatedAt time.Time `json:"updated_at" redis:"updated_at"`
	Tags      []string  `json:"tags" redis:"tags"`
}

func (r ResourceView) GetTable() string {
	return "resources_query"
}
