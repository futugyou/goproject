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
	Id        string    `json:"id" bson:"id"`
	Name      string    `json:"name" bson:"name"`
	Type      string    `json:"type" bson:"type"`
	Data      string    `json:"data" bson:"data"`
	Version   int       `json:"version" bson:"version"`
	IsDelete  bool      `json:"is_deleted" bson:"is_deleted"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	Tags      []string  `json:"tags" bson:"tags"`
}

func (r ResourceView) GetTable() string {
	return "resources_query"
}
