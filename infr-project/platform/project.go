package platform

// entity
type PlatformProject struct {
	Id       string            `json:"id" bson:"id"`
	Name     string            `json:"name" bson:"name"`
	Url      string            `json:"url" bson:"url"`
	Property map[string]string `json:"property" bson:"property"`
	Webhooks []Webhook         `json:"webhooks" bson:"webhooks"`
}

func NewPlatformProject(id string, name string, url string, property map[string]string) *PlatformProject {
	return &PlatformProject{
		Id:       id,
		Name:     name,
		Url:      url,
		Property: property,
		Webhooks: []Webhook{},
	}
}

func (w *PlatformProject) UpdateWebhook(hook Webhook) {
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
}

func (w *PlatformProject) RemoveWebhook(hookName string) {
	for i := len(w.Webhooks) - 1; i >= 0; i-- {
		if w.Webhooks[i].Name == hookName {
			w.Webhooks = append(w.Webhooks[:i], w.Webhooks[i+1:]...)
		}
	}
}
