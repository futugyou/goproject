package viewmodels

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
