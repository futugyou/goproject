package platform

import "fmt"

// entity
type PlatformProject struct {
	Id         string            `json:"id" bson:"id"`
	Name       string            `json:"name" bson:"name"`
	Url        string            `json:"url" bson:"url"`
	Properties map[string]string `json:"property" bson:"properties"`
	Webhooks   []Webhook         `json:"webhooks" bson:"webhooks"`
}

func (s PlatformProject) GetKey() string {
	return s.Id
}

func NewPlatformProject(id string, name string, url string, properties map[string]string) *PlatformProject {
	return &PlatformProject{
		Id:         id,
		Name:       name,
		Url:        url,
		Properties: properties,
	}
}

func (w *PlatformProject) UpdateName(name string) *PlatformProject {
	w.Name = name
	return w
}

func (w *PlatformProject) UpdateUrl(url string) *PlatformProject {
	w.Url = url
	return w
}

func (w *PlatformProject) UpdateProperties(properties map[string]string) *PlatformProject {
	w.Properties = properties
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
