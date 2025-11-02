package domain

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/futugyou/domaincore/domain"
)

// platform aggregate root
// The difference between Property and Secrets is like ConfigMap and Secrets in k8s
type Platform struct {
	domain.Aggregate
	Name        string
	Activate    bool
	Url         string
	Description string
	Provider    PlatformProvider
	Properties  map[string]Property
	Secrets     map[string]Secret
	Projects    map[string]PlatformProject
	Tags        []string
	IsDeleted   bool
}

type PlatformOption func(*Platform)

func WithPlatformSecrets(secrets map[string]Secret) PlatformOption {
	return func(w *Platform) {
		w.Secrets = secrets
	}
}

func WithPlatformTags(tags []string) PlatformOption {
	return func(w *Platform) {
		w.Tags = tags
	}
}

func WithPlatformProperties(properties map[string]Property) PlatformOption {
	return func(w *Platform) {
		w.Properties = properties
	}
}

func NewPlatform(name string, url string, provider PlatformProvider, opts ...PlatformOption) (*Platform, error) {
	platform := &Platform{
		Aggregate: domain.Aggregate{
			ID: uuid.New().String(),
		},
		Name:       name,
		Activate:   true,
		Url:        url,
		Properties: map[string]Property{},
		Secrets:    map[string]Secret{},
		Projects:   map[string]PlatformProject{},
		Tags:       []string{},
		IsDeleted:  false,
		Provider:   provider,
	}

	for _, opt := range opts {
		opt(platform)
	}

	// check property
	if err := platform.checkPlatformProviderProperties(platform.Properties); err != nil {
		return nil, err
	}

	// check secret
	if err := platform.checkPlatformProviderSecrets(platform.Secrets); err != nil {
		return nil, err
	}

	return platform, nil
}

func (r *Platform) stateCheck() error {
	if r.IsDeleted {
		return fmt.Errorf("id: %s was alrealdy deleted", r.ID)
	}

	return nil
}

func (w *Platform) Delete() (*Platform, error) {
	w.IsDeleted = true
	return w, nil
}

func (w *Platform) Recovery() (*Platform, error) {
	w.IsDeleted = false
	return w, nil
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

// update secret and propert need after update provider
func (w *Platform) UpdateProvider(provider PlatformProvider) (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}
	w.Provider = provider
	return w, nil
}

func (w *Platform) UpdateTags(tags []string) (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}
	w.Tags = tags
	return w, nil
}

func (w *Platform) UpdateProperties(properties map[string]Property) (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}

	if err := w.checkPlatformProviderProperties(properties); err != nil {
		return nil, err
	}

	w.Properties = properties
	return w, nil
}

func (w *Platform) UpdateSecrets(secrets map[string]Secret) (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}

	if err := w.checkPlatformProviderSecrets(secrets); err != nil {
		return nil, err
	}

	w.Secrets = secrets
	return w, nil
}

func (w *Platform) UpdateProject(project PlatformProject) (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}

	if w.Projects == nil {
		w.Projects = map[string]PlatformProject{}
	}

	if w.Provider == PlatformProviderGithub {
		if _, ok := project.Properties["GITHUB_REPO"]; !ok {
			project.Properties["GITHUB_REPO"] = Property{Key: "GITHUB_REPO", Value: project.Name}
		}
	}

	w.Projects[project.ID] = project
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

func (w *Platform) UpdateWebhook(projectId string, hook *Webhook) (*Platform, error) {
	if err := w.stateCheck(); err != nil {
		return nil, err
	}

	if project, exists := w.Projects[projectId]; exists {
		projectPointer := &project
		projectPointer.UpdateWebhook(hook)
		w.Projects[projectId] = *projectPointer
		return w, nil
	} else {
		return nil, fmt.Errorf("project id: %s does not exist", projectId)
	}
}

func (w *Platform) GetWebhook(projectId string) *Webhook {
	project, exists := w.Projects[projectId]
	if !exists {
		return nil
	}

	return project.GetWebhook()
}

func (w *Platform) ProviderVaultInfo() (string, error) {
	vaultId := ""
	switch w.Provider {
	case PlatformProviderCircleci:
		vaultId = w.Secrets["CIRCLECI_TOKEN"].Value
	case PlatformProviderVercel:
		vaultId = w.Secrets["VERCEL_TOKEN"].Value
	case PlatformProviderGithub:
		vaultId = w.Secrets["GITHUB_TOKEN"].Value
	default:
		return "", fmt.Errorf("%s not supported", w.Provider.String())
	}

	return vaultId, nil
}

func (r Platform) AggregateName() string {
	return "platforms"
}

func (w *Platform) checkPlatformProviderSecrets(secretMap map[string]Secret) error {
	switch w.Provider {
	case PlatformProviderCircleci:
		if _, ok := secretMap["CIRCLECI_TOKEN"]; !ok {
			return fmt.Errorf("%s provider MUST have CIRCLECI_TOKEN in Secret", w.Provider.String())
		}
	case PlatformProviderVercel:
		if _, ok := secretMap["VERCEL_TOKEN"]; !ok {
			return fmt.Errorf("%s provider MUST have VERCEL_TOKEN in Secret", w.Provider.String())
		}
	case PlatformProviderGithub:
		if _, ok := secretMap["GITHUB_TOKEN"]; !ok {
			return fmt.Errorf("%s provider MUST have GITHUB_TOKEN in Secret", w.Provider.String())
		}
	}

	return nil
}

func (w *Platform) checkPlatformProviderProperties(properties map[string]Property) error {
	switch w.Provider {
	case PlatformProviderOther:
	case PlatformProviderVercel:
	case PlatformProviderCircleci:
		if _, ok := properties["org_slug"]; !ok {
			return fmt.Errorf("%s provider MUST have org_slug in Property", w.Provider.String())
		}
	case PlatformProviderGithub:
		if _, ok := properties["GITHUB_OWNER"]; !ok {
			return fmt.Errorf("%s provider MUST have GITHUB_OWNER in Property", w.Provider.String())
		}
	}

	return nil
}
