package domain

// entity
type PlatformProject struct {
	ID                string
	Name              string
	Url               string
	Description       string
	Properties        map[string]Property
	Secrets           map[string]Secret
	ImageUrl          string
	Webhook           *Webhook
	ProviderProjectID string
	Tags              []string
}

func (s PlatformProject) GetKey() string {
	return s.ID
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

func WithProviderProjectID(providerProjectId string) ProjectOption {
	return func(w *PlatformProject) {
		w.ProviderProjectID = providerProjectId
	}
}

func WithProjectTags(tags []string) ProjectOption {
	return func(w *PlatformProject) {
		w.Tags = tags
	}
}

func NewPlatformProject(id string, name string, url string, opts ...ProjectOption) *PlatformProject {
	project := &PlatformProject{
		ID:         id,
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

func (w *PlatformProject) UpdateProviderProjectID(id string) *PlatformProject {
	w.ProviderProjectID = id
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

func (w *PlatformProject) UpdateImageUrl(url string) *PlatformProject {
	w.ImageUrl = url
	return w
}

func (w *PlatformProject) UpdateTags(tags []string) *PlatformProject {
	w.Tags = tags
	return w
}

func (w *PlatformProject) UpdateWebhook(hook *Webhook) {
	w.Webhook = hook
}

func (w *PlatformProject) GetId() string {
	if w == nil {
		return ""
	}
	return w.ID
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

func (w *PlatformProject) GetProviderProjectID() string {
	if w == nil {
		return ""
	}
	return w.ProviderProjectID
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

func (w *PlatformProject) GetWebhook() *Webhook {
	if w == nil {
		return nil
	}
	return w.Webhook
}
