package platform

import "fmt"

// entity
type PlatformProject struct {
	Id         string              `json:"id" bson:"id"`
	Name       string              `json:"name" bson:"name"`
	Url        string              `json:"url" bson:"url"`
	Properties map[string]Property `json:"properties" bson:"properties"`
	Secrets    map[string]Secret   `json:"secrets" bson:"secrets"`
	Webhooks   []Webhook           `json:"webhooks" bson:"webhooks"`
	Followed   bool                `json:"followed" bson:"followed"`
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

func NewPlatformProject(id string, name string, url string, opts ...ProjectOption) *PlatformProject {
	project := &PlatformProject{
		Id:       id,
		Name:     name,
		Followed: false,
		Url:      url, Properties: make(map[string]Property),
		Secrets: make(map[string]Secret),
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

func (w *PlatformProject) UpdateUrl(url string) *PlatformProject {
	w.Url = url
	return w
}

func (w *PlatformProject) UpdateFollowed(followed bool) *PlatformProject {
	w.Followed = followed
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
