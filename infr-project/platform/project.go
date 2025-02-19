package platform

import "fmt"

// entity
type PlatformProject struct {
	Id                string
	Name              string
	Url               string
	Description       string
	Properties        map[string]Property
	Secrets           map[string]Secret
	ImageData         []byte
	Webhooks          []Webhook
	ProviderProjectId string
}

func (s PlatformProject) GetKey() string {
	return s.Id
}

type ProjectOption func(*PlatformProject)

func WithProjectProperties(properties map[string]Property) ProjectOption {
	return func(w *PlatformProject) {
		w.Properties = properties
	}
}

func WithProjectSecrets(secrets map[string]Secret) ProjectOption {
	return func(w *PlatformProject) {
		w.Secrets = secrets
	}
}

func WithProjectDescription(description string) ProjectOption {
	return func(w *PlatformProject) {
		w.Description = description
	}
}

func NewPlatformProject(id string, name string, url string, opts ...ProjectOption) *PlatformProject {
	project := &PlatformProject{
		Id:         id,
		Name:       name,
		Url:        url,
		Properties: make(map[string]Property),
		Secrets:    make(map[string]Secret),
	}

	for _, opt := range opts {
		opt(project)
	}

	return project
}

func (w *PlatformProject) UpdateName(name string) *PlatformProject {
	w.Name = name
	return w
}

func (w *PlatformProject) UpdateDescription(description string) *PlatformProject {
	w.Description = description
	return w
}

func (w *PlatformProject) UpdateUrl(url string) *PlatformProject {
	w.Url = url
	return w
}

func (w *PlatformProject) UpdateProviderProjectId(id string) *PlatformProject {
	w.ProviderProjectId = id
	return w
}

func (w *PlatformProject) UpdateProperties(properties map[string]Property) *PlatformProject {
	w.Properties = properties
	return w
}

func (w *PlatformProject) UpdateSecrets(secrets map[string]Secret) *PlatformProject {
	w.Secrets = secrets
	return w
}

func (w *PlatformProject) UpdateImageData(imageData []byte) *PlatformProject {
	w.ImageData = imageData
	return w
}

func (w *PlatformProject) UpsertWebhook(hook Webhook) {
	for i := 0; i < len(w.Webhooks); i++ {
		if w.Webhooks[i].Name == hook.Name {
			w.Webhooks[i] = hook
			return
		}
	}

	w.Webhooks = append(w.Webhooks, hook)
}

func (w *PlatformProject) RemoveWebhook(hookName string) error {
	for i := len(w.Webhooks) - 1; i >= 0; i-- {
		if w.Webhooks[i].Name == hookName {
			w.Webhooks = append(w.Webhooks[:i], w.Webhooks[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("webhook name: %s does not exist", hookName)
}

func (w *PlatformProject) ClearWebhooks() {
	if w == nil {
		return
	}

	w.Webhooks = []Webhook{}
}

func (w *PlatformProject) GetWebhook(hookName string) (*Webhook, error) {
	for i := len(w.Webhooks) - 1; i >= 0; i-- {
		if w.Webhooks[i].Name == hookName {
			return &w.Webhooks[i], nil
		}
	}

	return nil, fmt.Errorf("webhook name: %s does not exist", hookName)
}

func (w *PlatformProject) GetId() string {
	if w == nil {
		return ""
	}
	return w.Id
}

func (w *PlatformProject) GetName() string {
	if w == nil {
		return ""
	}
	return w.Name
}

func (w *PlatformProject) GetUrl() string {
	if w == nil {
		return ""
	}
	return w.Url
}

func (w *PlatformProject) GetDescription() string {
	if w == nil {
		return ""
	}
	return w.Description
}

func (w *PlatformProject) GeProviderProjectId() string {
	if w == nil {
		return ""
	}
	return w.ProviderProjectId
}

func (w *PlatformProject) GetProperties() map[string]Property {
	if w == nil {
		return map[string]Property{}
	}
	return w.Properties
}

func (w *PlatformProject) GetSecrets() map[string]Secret {
	if w == nil {
		return map[string]Secret{}
	}
	return w.Secrets
}

func (w *PlatformProject) GetWebhooks() []Webhook {
	if w == nil {
		return []Webhook{}
	}
	return w.Webhooks
}
