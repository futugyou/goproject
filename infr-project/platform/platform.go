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
	Property         map[string]PropertyInfo    `json:"property"`
	Projects         map[string]PlatformProject `json:"projects"`
	Tags             []string                   `json:"tags"`
	IsDeleted        bool                       `json:"is_deleted"`
}

type PropertyInfo struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	NeedMask bool   `json:"needMask"`
}

func NewPlatform(name string, url string, rest string, property map[string]PropertyInfo, tags []string) *Platform {
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

func (w *Platform) UpdateProperty(property map[string]PropertyInfo) (*Platform, error) {
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
