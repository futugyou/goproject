package viewmodels

type PropertyInfo struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	NeedMask bool   `json:"needMask"`
}

type CreatePlatformRequest struct {
	Name     string                  `json:"name" validate:"required,min=3,max=50"`
	Url      string                  `json:"url" validate:"required,min=3,max=50"`
	Rest     string                  `json:"rest" validate:"required,min=3,max=50"`
	Tags     []string                `json:"tags"`
	Property map[string]PropertyInfo `json:"property"`
}

type UpdatePlatformRequest struct {
	Name     string                  `json:"name" validate:"required,min=3,max=50"`
	Url      string                  `json:"url" validate:"required,min=3,max=50"`
	Rest     string                  `json:"rest" validate:"required,min=3,max=50"`
	Property map[string]PropertyInfo `json:"property,omitempty"`
	Tags     []string                `json:"tags"`
	Activate *bool                   `json:"activate,omitempty"`
}

type UpdatePlatformProjectRequest struct {
	Name     string            `json:"name" validate:"required,min=3,max=50"`
	Url      string            `json:"url" validate:"required,min=3,max=50"`
	Property map[string]string `json:"property,omitempty"`
}

type UpdatePlatformWebhookRequest struct {
	Name     string            `json:"name" validate:"required,min=3,max=50"`
	Url      string            `json:"url" validate:"required,min=3,max=50"`
	Activate bool              `json:"activate" validate:"required"`
	State    string            `json:"state" validate:"oneof=Init Creating Ready"`
	Property map[string]string `json:"property"`
}
