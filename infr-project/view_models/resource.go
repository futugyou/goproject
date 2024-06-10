package viewmodels

import "time"

type CreateResourceRequest struct {
	Name string `json:"name" validate:"required,min=3,max=50"`
	Type string `json:"type" validate:"oneof=DrawIO Markdown Excalidraw Plate"`
	Data string `json:"data" validate:"required,min=3"`
}

type UpdateResourceRequest struct {
	Id   string `json:"id" validate:"required"`
	Data string `json:"data" validate:"required,min=3"`
}

type ResourceDetail struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Data      string    `json:"data"`
	IsDelete  bool      `json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
