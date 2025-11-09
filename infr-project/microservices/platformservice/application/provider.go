package application

import (
	"context"
	"fmt"
	"log"

	tool "github.com/futugyou/extensions"

	"github.com/futugyou/platformservice/assembler"
	"github.com/futugyou/platformservice/domain"
	"github.com/futugyou/platformservice/provider"
	"github.com/futugyou/platformservice/viewmodel"
)

func (s *PlatformService) determineProviderStatus(ctx context.Context, res *domain.Platform) bool {
	provider, err := s.getPlatformProvider(ctx, *res)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	user, err := provider.GetUser(ctx)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	if user == nil || len(user.ID) == 0 {
		log.Printf("no user found for %s provider\n", res.Provider.String())
		return false
	}

	return true
}

func (s *PlatformService) getPlatformProvider(ctx context.Context, src domain.Platform) (provider.PlatformProvider, error) {
	vaultId, err := src.ProviderVaultInfo()
	if err != nil {
		return nil, err
	}

	token, err := s.vaultService.ShowVaultRawValue(ctx, vaultId)
	if err != nil {
		return nil, fmt.Errorf("get platform provider token error, vaultId is %s, message %s", vaultId, err.Error())
	}

	return provider.PlatformProviderFactory(ctx, src.Provider.String(), token)
}

func (s *PlatformService) getProviderProjects(ctx context.Context, prov provider.PlatformProvider, src domain.Platform) ([]provider.Project, error) {
	parameters := make(map[string]string)
	for _, v := range src.Properties {
		parameters[v.Key] = v.Value
	}
	filter := provider.ProjectFilter{
		Parameters: parameters,
	}

	return prov.ListProject(ctx, filter)
}

func (a *PlatformService) toPlatformDetailViewWithProviderProjects(mapper assembler.PlatformAssembler, d *domain.Platform, projects []provider.Project) *viewmodel.PlatformDetailView {
	platform := mapper.ToPlatformDetailView(d)
	for i := range platform.Projects {
		for _, v := range projects {
			if platform.Projects[i].ProviderProjectID == v.ID {
				platform.Projects[i].Followed = true
				platform.Projects[i].ProviderProject = a.convertProviderProjectToSimpleModel(&v)
				break
			}
		}
	}

	return platform
}

func (s *PlatformService) convertProviderProjectToModel(providerProject *provider.Project) *viewmodel.PlatformProviderProject {
	projectModel := &viewmodel.PlatformProviderProject{
		ID:                   providerProject.ID,
		Name:                 providerProject.Name,
		Description:          providerProject.Description,
		Url:                  providerProject.Url,
		EnvironmentVariables: s.convertToPlatformModelEnvironments(providerProject.EnvironmentVariables),
		Environments:         providerProject.Environments,
		Workflows:            s.convertToPlatformModelWorkflows(providerProject.Workflows),
		WorkflowRuns:         s.convertToPlatformModelWorkflowRuns(providerProject.WorkflowRuns),
		Deployments:          s.convertToPlatformModelDeployments(providerProject.Deployments),
		BadgeURL:             providerProject.BadgeURL,
		BadgeMarkdown:        providerProject.BadgeMarkDown,
		Properties:           s.mapToModelProperty(providerProject.Properties),
		Tags:                 providerProject.Tags,
		Readme:               providerProject.Readme,
	}

	return projectModel
}

func (s *PlatformService) convertToPlatformModelEnvironments(values map[string]provider.EnvironmentVariable) []viewmodel.EnvironmentVariable {
	return tool.MapToSlice(values, func(key string, v provider.EnvironmentVariable) viewmodel.EnvironmentVariable {
		return viewmodel.EnvironmentVariable(v)
	})
}

func (s *PlatformService) convertToPlatformModelWorkflows(values map[string]provider.Workflow) []viewmodel.Workflow {
	return tool.MapToSlice(values, func(key string, v provider.Workflow) viewmodel.Workflow {
		return viewmodel.Workflow(v)
	})
}

func (s *PlatformService) convertToPlatformModelWorkflowRuns(values map[string]provider.WorkflowRun) []viewmodel.WorkflowRun {
	return tool.MapToSlice(values, func(key string, v provider.WorkflowRun) viewmodel.WorkflowRun {
		return viewmodel.WorkflowRun(v)
	})
}

func (s *PlatformService) convertToPlatformModelDeployments(values map[string]provider.Deployment) []viewmodel.Deployment {
	return tool.MapToSlice(values, func(key string, v provider.Deployment) viewmodel.Deployment {
		return viewmodel.Deployment(v)
	})
}

func (*PlatformService) mapToModelProperty(providerProperties map[string]string) []viewmodel.Property {
	properties := []viewmodel.Property{}
	for key, v := range providerProperties {
		properties = append(properties, viewmodel.Property{
			Key:   key,
			Value: v,
		})
	}

	return properties
}

func (s *PlatformService) convertProviderProjectToSimpleModel(providerProject *provider.Project) *viewmodel.PlatformProviderProject {
	projectModel := &viewmodel.PlatformProviderProject{
		ID:            providerProject.ID,
		Name:          providerProject.Name,
		Description:   providerProject.Description,
		Url:           providerProject.Url,
		BadgeURL:      providerProject.BadgeURL,
		BadgeMarkdown: providerProject.BadgeMarkDown,
		Tags:          providerProject.Tags,
		Readme:        providerProject.Readme,
	}

	return projectModel
}

func (s *PlatformService) getProviderProject(ctx context.Context, platformProvider provider.PlatformProvider, name string, parameters map[string]string) (*provider.Project, error) {
	filter := provider.ProjectFilter{
		Parameters: parameters,
		Name:       name,
	}

	return platformProvider.GetProject(ctx, filter)
}

func mergePropertiesToMap(propertiesList ...map[string]domain.Property) map[string]string {
	properties := make(map[string]string)
	for _, propertyMap := range propertiesList {
		for _, v := range propertyMap {
			if _, ok := properties[v.Key]; !ok {
				properties[v.Key] = v.Value
			}
		}
	}
	return properties
}
