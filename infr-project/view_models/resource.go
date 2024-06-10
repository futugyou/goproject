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
	Id        string    `json:"id" bson:"id"`
	Name      string    `json:"name" bson:"name"`
	Type      string    `json:"type" bson:"type"`
	Data      string    `json:"data" bson:"data"`
	IsDelete  bool      `json:"is_deleted" bson:"is_deleted"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (r ResourceDetail) GetTable() string {
	return "resources_query"
}
