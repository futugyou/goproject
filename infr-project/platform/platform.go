package platform

import (
	"github.com/google/uuid"

	"github.com/futugyou/infr-project/domain"
)

// aggregate root
type Platform struct {
	domain.Aggregate `json:"-"`
	Name             string                     `json:"name"`
	Activate         bool                       `json:"activate"`
	Url              string                     `json:"url"`
	RestEndpoint     string                     `json:"rest_endpoint"`
	Property         map[string]string          `json:"property"`
	Projects         map[string]PlatformProject `json:"projects"`
}

func NewPlatform(name string, url string, rest string, property map[string]string) *Platform {
	return &Platform{
		Aggregate: domain.Aggregate{
			Id: uuid.New().String(),
		},
		Name:         name,
		Activate:     true,
		Url:          url,
		RestEndpoint: rest,
		Property:     property,
		Projects:     map[string]PlatformProject{},
	}
}

func (w *Platform) Enable() *Platform {
	w.Activate = true
	return w
}

func (w *Platform) Disable() *Platform {
	w.Activate = false
	return w
}

func (w *Platform) UpdateName(name string) *Platform {
	w.Name = name
	return w
}

func (w *Platform) UpdateUrl(url string) *Platform {
	w.Url = url
	return w
}

func (w *Platform) UpdateRestEndpoint(url string) *Platform {
	w.RestEndpoint = url
	return w
}

func (w *Platform) UpdateProperty(property map[string]string) *Platform {
	w.Property = property
	return w
}

func (w *Platform) UpdateWebhook(projectId string, hook Webhook) *Platform {
	if project, exists := w.Projects[projectId]; exists {
		(&project).UpdateWebhook(hook)
	}
	return w
}

func (w *Platform) RemoveWebhook(projectId string, hookName string) *Platform {
	if project, exists := w.Projects[projectId]; exists {
		(&project).RemoveWebhook(hookName)
	}
	return w
}

func (r *Platform) AggregateName() string {
	return "platforms"
}
