package assembler

import (
	"github.com/futugyou/platformservice/domain"
	"github.com/futugyou/platformservice/provider"
)

type ProjectAssembler struct{}

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
		domain.WithProviderProjectId(providerProject.ID),
	)

	return projectModel
}
