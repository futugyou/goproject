package assembler

import (
	"github.com/futugyou/platformservice/domain"
	"github.com/futugyou/platformservice/provider"
	"github.com/futugyou/platformservice/viewmodel"
)

type ProjectAssembler struct{}

func (s ProjectAssembler) ToViewModels(projects []domain.PlatformProject) []viewmodel.PlatformProject {
	result := []viewmodel.PlatformProject{}
	for _, project := range projects {
		result = append(result, *s.ToViewModel(&project))
	}
	return result
}

func (s ProjectAssembler) ToViewModel(project *domain.PlatformProject) *viewmodel.PlatformProject {
	return &viewmodel.PlatformProject{
		ID:          project.ID,
		Name:        project.Name,
		Url:         project.Url,
		ImageUrl:    project.ImageUrl,
		Description: project.Description,
		Properties: func() []viewmodel.Property {
			res := []viewmodel.Property{}
			for k, v := range project.Properties {
				res = append(res, viewmodel.Property{
					Key:   k,
					Value: v.Value,
				})
			}
			return res
		}(),
		Secrets: func() []viewmodel.Secret {
			res := []viewmodel.Secret{}
			for _, v := range project.Secrets {
				res = append(res, viewmodel.Secret{
					Key:       v.Key,
					VaultID:   v.Value,
					VaultKey:  v.VaultKey,
					MaskValue: v.VaultMaskValue,
				})
			}
			return res
		}(),
		Webhook: func() *viewmodel.Webhook {
			if project.Webhook == nil {
				return nil
			}
			return &viewmodel.Webhook{
				ID:  project.Webhook.ID,
				Url: project.Webhook.Url,
			}
		}(),
		ProviderProjectID: project.ProviderProjectID,
		Followed: func() bool {
			if len(project.ProviderProjectID) > 0 {
				return true
			}
			return false
		}(),
	}
}

func (s *ProjectAssembler) ToDomain(providerProject *provider.Project) *domain.PlatformProject {
	properties := map[string]domain.Property{}
	for key, v := range providerProject.Properties {
		properties[key] = domain.Property{
			Key:   key,
			Value: v,
		}
	}

	projectModel := domain.NewPlatformProject(providerProject.ID, providerProject.Name, providerProject.Url,
		domain.WithProjectProperties(properties),
		domain.WithProjectDescription(providerProject.Description),
		domain.WithProjectTags(providerProject.Tags),
		domain.WithProviderProjectID(providerProject.ID),
	)

	return projectModel
}
