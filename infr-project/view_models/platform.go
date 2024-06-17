package viewmodels

type CreatePlatformRequest struct {
	Name     string            `json:"name" validate:"required,min=3,max=50"`
	Url      string            `json:"url" validate:"required,min=3,max=50"`
	Rest     string            `json:"rest" validate:"required,min=3,max=50"`
	Property map[string]string `json:"property"`
}

type UpdatePlatformRequest struct {
	Name     string            `json:"name" validate:"required,min=3,max=50"`
	Url      string            `json:"url" validate:"required,min=3,max=50"`
	Rest     string            `json:"rest" validate:"required,min=3,max=50"`
	Property map[string]string `json:"property,omitempty"`
	Activate *bool             `json:"activate,omitempty"`
}

type UpdatePlatformProjectRequest struct {
	Name     string            `json:"name" validate:"required,min=3,max=50"`
	Url      string            `json:"url" validate:"required,min=3,max=50"`
	Property map[string]string `json:"property,omitempty"`
}
