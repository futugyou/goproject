package platform

import (
	"github.com/google/uuid"

	"github.com/futugyou/infr-project/domain"
)

type Platform struct {
	domain.Aggregate `json:"-"`
	Name             string            `json:"name"`
	Activate         bool              `json:"activate"`
	Url              string            `json:"url"`
	RestEndpoint     string            `json:"rest_endpoint"`
	Property         map[string]string `json:"property"`
	Webhooks         []Webhook         `json:"webhooks"`
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
		Webhooks:     []Webhook{},
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

func (w *Platform) UpdateWebhook(hook Webhook) *Platform {
	f := false
	for i := 0; i < len(w.Webhooks); i++ {
		if w.Webhooks[i].Name == hook.Name {
			w.Webhooks[i] = hook
			f = true
			break
		}
	}

	if !f {
		w.Webhooks = append(w.Webhooks, hook)
	}
	return w
}

func (w *Platform) RemoveWebhook(hookName string) *Platform {
	for i := len(w.Webhooks) - 1; i >= 0; i-- {
		if w.Webhooks[i].Name == hookName {
			w.Webhooks = append(w.Webhooks[:i], w.Webhooks[i+1:]...)
		}
	}
	return w
}

func (r *Platform) AggregateName() string {
	return "platforms"
}
