package platform

import (
	"fmt"

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
	Tags             []string                   `json:"tags"`
	IsDeleted        bool                       `json:"is_deleted"`
}

func NewPlatform(name string, url string, rest string, property map[string]string, tags []string) *Platform {
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
		Tags:         tags,
		IsDeleted:    false,
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

func (w *Platform) UpdateRestEndpoint(url string) (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}
	w.RestEndpoint = url
	return w, nil
}

func (w *Platform) UpdateTags(tags []string) (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}
	w.Tags = tags
	return w, nil
}

func (w *Platform) UpdateProperty(property map[string]string) (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}
	w.Property = property
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
	delete(w.Projects, projectId)
	return w, nil
}

func (w *Platform) UpdateWebhook(projectId string, hook Webhook) (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}
	if project, exists := w.Projects[projectId]; exists {
		projectPointer := &project
		projectPointer.UpdateWebhook(hook)
		w.Projects[projectId] = *projectPointer
	}
	return w, nil
}

func (w *Platform) RemoveWebhook(projectId string, hookName string) (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}
	if project, exists := w.Projects[projectId]; exists {
		projectPointer := &project
		projectPointer.RemoveWebhook(hookName)
		w.Projects[projectId] = *projectPointer
	}
	return w, nil
}

func (r *Platform) AggregateName() string {
	return "platforms"
}
