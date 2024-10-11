package platform

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/futugyou/infr-project/domain"
)

// platform aggregate root
// The difference between Property and Secrets is like ConfigMap and Secrets in k8s
type Platform struct {
	domain.Aggregate `json:"-"`
	Name             string                     `json:"name"`
	Activate         bool                       `json:"activate"`
	Url              string                     `json:"url"`
	Properties       map[string]PropertyInfo    `json:"properties"`
	Secrets          map[string]Secret          `json:"secrets"`
	Projects         map[string]PlatformProject `json:"projects"`
	Tags             []string                   `json:"tags"`
	IsDeleted        bool                       `json:"is_deleted"`
}

func NewPlatform(name string, url string, properties map[string]PropertyInfo, tags []string) *Platform {
	return &Platform{
		Aggregate: domain.Aggregate{
			Id: uuid.New().String(),
		},
		Name:       name,
		Activate:   true,
		Url:        url,
		Properties: properties,
		Projects:   map[string]PlatformProject{},
		Tags:       tags,
		IsDeleted:  false,
	}
}

func (r *Platform) stateCheck() error {
	if r.IsDeleted {
		return fmt.Errorf("id: %s is alrealdy deleted", r.Id)
	}

	return nil
}

func (w *Platform) Enable() (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}
	w.Activate = true
	return w, nil
}

func (w *Platform) Disable() (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}
	w.Activate = false
	return w, nil
}

func (w *Platform) UpdateName(name string) (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}
	w.Name = name
	return w, nil
}

func (w *Platform) UpdateUrl(url string) (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}
	w.Url = url
	return w, nil
}

func (w *Platform) UpdateTags(tags []string) (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}
	w.Tags = tags
	return w, nil
}

func (w *Platform) UpdateProperties(properties map[string]PropertyInfo) (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}
	w.Properties = properties
	return w, nil
}

func (w *Platform) UpdateSecret(secrets map[string]Secret) (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}
	w.Secrets = secrets
	return w, nil
}

// this update not include webhook
func (w *Platform) UpdateProject(project PlatformProject) (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}
	if w.Projects == nil {
		w.Projects = map[string]PlatformProject{}
	}
	if pro, exists := w.Projects[project.Id]; exists {
		project.Webhooks = pro.Webhooks
	}
	w.Projects[project.Id] = project
	return w, nil
}

func (w *Platform) RemoveProject(projectId string) (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}
	if _, ok := w.Projects[projectId]; ok {
		delete(w.Projects, projectId)
		return w, nil
	} else {
		return nil, fmt.Errorf("project id: %s does not exist", projectId)
	}
}

func (w *Platform) UpdateWebhook(projectId string, hook Webhook) (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}

	if project, exists := w.Projects[projectId]; exists {
		projectPointer := &project
		projectPointer.UpsertWebhook(hook)
		w.Projects[projectId] = *projectPointer
		return w, nil
	} else {
		return nil, fmt.Errorf("project id: %s does not exist", projectId)
	}
}

func (w *Platform) RemoveWebhook(projectId string, hookName string) (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}

	if project, exists := w.Projects[projectId]; exists {
		projectPointer := &project
		if err := projectPointer.RemoveWebhook(hookName); err != nil {
			return nil, err
		}
		w.Projects[projectId] = *projectPointer
		return w, nil
	} else {
		return nil, fmt.Errorf("project id: %s does not exist", projectId)
	}
}

func (r Platform) AggregateName() string {
	return "platforms"
}
