package application

import (
	"context"
	"fmt"
	"slices"
	"strings"

	coreinfr "github.com/futugyou/domaincore/infrastructure"
	tool "github.com/futugyou/extensions"
	"github.com/google/uuid"

	"github.com/futugyou/platformservice/assembler"
	"github.com/futugyou/platformservice/domain"
	"github.com/futugyou/platformservice/infrastructure"
	platformprovider "github.com/futugyou/platformservice/provider"
	"github.com/futugyou/platformservice/viewmodel"
)

func (s *PlatformService) ImportProjectsFromProvider(ctx context.Context, idOrName string, providerProjectIDs []string) error {
	plat, err := s.repository.GetPlatformByIdOrName(ctx, idOrName)
	if err != nil {
		return err
	}
	provider, err := s.getPlatformProvider(ctx, *plat)
	if err != nil {
		return err
	}
	providerProjects, err := s.getProviderProjects(ctx, provider, *plat)
	if err != nil {
		return err
	}

	providerProjects = tool.ArrayFilter(providerProjects, func(project platformprovider.Project) bool {
		if len(providerProjectIDs) == 0 {
			return true
		}

		return slices.Contains(providerProjectIDs, project.ID)
	})

	if len(providerProjects) == 0 {
		return fmt.Errorf("no provider project found in %s", plat.Provider.String())
	}

	projectMapper := assembler.ProjectAssembler{}
	events := []coreinfr.Event{}

	if len(plat.Projects) == 0 {
		plat.Projects = map[string]domain.PlatformProject{}
		for _, project := range providerProjects {
			plat.Projects[project.ID] = *projectMapper.ToDomain(&project)
			events = s.buildProjectChangeEvent(plat.ID, plat.Name, plat.Provider.String(), project.ID, project.Name, project.ID, project.Url, false, true, true)
		}
	} else {
		for _, project := range providerProjects {
			find := false
			for _, entity := range plat.Projects {
				// The following situations are considered to be `link`
				// 1. already `linked`: ProviderProjectId = ID
				// 2. name is the same, this is a normal situation
				// 3. ID is the same, linked, then cancel(ProviderProjectId is ""), then link again
				if entity.ProviderProjectID == project.ID || entity.Name == project.Name || entity.ID == project.ID {
					find = true
					rawProject := plat.Projects[project.ID]
					rawProject.ProviderProjectID = project.ID
					plat.Projects[project.ID] = rawProject
					break
				}
			}

			if find {
				continue
			}

			newPro := *projectMapper.ToDomain(&project)
			events = s.buildProjectChangeEvent(plat.ID, plat.Name, plat.Provider.String(), newPro.ID, newPro.Name, newPro.ID, newPro.Url, false, true, true)
			plat.Projects[project.ID] = newPro
		}
	}

	bMap := make(map[string]struct{})
	projects := []domain.PlatformProject{}
	for _, v := range plat.Projects {
		if _, ok := bMap[v.ProviderProjectID]; ok {
			return fmt.Errorf("provider project id: %s is duplicated", v.ProviderProjectID)
		}
		bMap[v.ProviderProjectID] = struct{}{}
		projects = append(projects, v)
	}

	return s.innerService.WithUnitOfWork(ctx, func(ctx context.Context) error {
		err := s.eventHandler.DispatchIntegrationEvents(ctx, events)
		if err != nil {
			return err
		}

		return s.repository.SyncProjects(ctx, plat.ID, projects)
	})
}

func (s *PlatformService) UpsertProject(ctx context.Context, idOrName string, projectID string, project viewmodel.UpdatePlatformProjectRequest) error {
	plat, err := s.repository.GetPlatformByIdOrName(ctx, idOrName)
	if err != nil {
		return err
	}

	if len(projectID) == 0 {
		projectID = tool.Sanitize2String(strings.ToLower(project.Name), "_")
	}

	properties := map[string]domain.Property{}
	for _, v := range project.Properties {
		properties[v.Key] = domain.Property(v)
	}

	secretMapper := assembler.SecretAssembler{}
	secrets, err := secretMapper.ToModel(ctx, project.Secrets)
	if err != nil {
		return err
	}

	var projectDb *domain.PlatformProject
	screenshot := false
	if proj, ok := plat.Projects[projectID]; ok {
		projectDb = &proj
		if len(project.Url) > 0 && (projectDb.Url != project.Url || len(projectDb.ImageUrl) == 0) {
			screenshot = true
		}

		projectDb.UpdateName(project.Name)
		projectDb.UpdateDescription(project.Description)
		projectDb.UpdateProperties(properties)
		projectDb.UpdateUrl(project.Url)
		projectDb.UpdateSecrets(secrets)
		projectDb.UpdateProviderProjectID(project.ProviderProjectID)
		projectDb.UpdateTags(project.Tags)
	} else {
		projectDb = domain.NewPlatformProject(
			projectID,
			project.Name,
			project.Url,
			domain.WithProjectProperties(properties),
			domain.WithProjectSecrets(secrets),
			domain.WithProjectDescription(project.Description),
			domain.WithProviderProjectID(project.ProviderProjectID),
			domain.WithProjectTags(project.Tags),
		)

		if len(projectDb.Url) > 0 {
			screenshot = true
		}
	}

	if _, err = plat.UpdateProject(*projectDb); err != nil {
		return err
	}

	return s.innerService.WithUnitOfWork(ctx, func(ctx context.Context) error {
		events := s.buildProjectChangeEvent(plat.ID, plat.Name, plat.Provider.String(), projectID, project.Name, project.ProviderProjectID, project.Url, project.Operate == "sync", project.ImportWebhooks, screenshot)
		err := s.eventHandler.DispatchIntegrationEvents(ctx, events)
		if err != nil {
			return err
		}

		return s.repository.Update(ctx, *plat)
	})
}

func (s *PlatformService) buildProjectChangeEvent(platID, platName, provider, projectID, projectName, providerProjectID, projectUrl string, createProviderProject, importWebhook, screenshot bool) []coreinfr.Event {
	events := []coreinfr.Event{}
	if createProviderProject {
		events = append(events, &infrastructure.CreateProviderProjectTriggeredEvent{
			ID:          uuid.NewString(),
			PlatformID:  platID,
			ProjectID:   projectID,
			ProjectName: projectName,
			Provider:    provider,
		})
	}

	if importWebhook {
		events = append(events, &infrastructure.CreateProviderWebhookTriggeredEvent{
			ID:                uuid.NewString(),
			PlatformID:        platID,
			ProjectID:         projectID,
			ProjectName:       projectName,
			Provider:          provider,
			ProviderProjectId: providerProjectID,
			Url:               domain.GetWebhookUrl(platName, projectName, s.opts.ProjectWebhookUrl),
		})
	}

	if screenshot {
		events = append(events, &infrastructure.ProjectScreenshotTriggeredEvent{
			ID:         uuid.NewString(),
			PlatformID: platID,
			ProjectID:  projectID,
			Url:        projectUrl,
		})
	}

	return events
}

func (s *PlatformService) DeleteProject(ctx context.Context, idOrName string, projectId string) error {
	plat, err := s.repository.GetPlatformByIdOrName(ctx, idOrName)
	if err != nil {
		return err
	}

	if _, err := plat.RemoveProject(projectId); err != nil {
		return err
	}

	return s.innerService.WithUnitOfWork(ctx, func(ctx context.Context) error {
		return s.repository.DeleteProject(ctx, plat.ID, projectId)
	})
}

func (s *PlatformService) GetPlatformProject(ctx context.Context, idOrName string, projectId string) (*viewmodel.PlatformProject, error) {
	src, err := s.repository.GetPlatformByIdOrNameWithoutProjects(ctx, idOrName)
	if err != nil {
		return nil, err
	}

	project, err := s.repository.GetPlatformProjectByProjectID(ctx, src.ID, projectId)
	if err != nil {
		return nil, err
	}

	projectMapper := assembler.ProjectAssembler{}
	projectModel := projectMapper.ToViewModel(project)

	if len(project.ProviderProjectID) > 0 {
		platformProvider, err := s.getPlatformProvider(ctx, *src)
		if err != nil {
			return nil, err
		}

		parameters := mergePropertiesToMap(project.Properties, src.Properties)
		providerProject, err := s.getProviderProject(ctx, platformProvider, project.ProviderProjectID, parameters)
		if err != nil {
			return nil, err
		}

		projectModel.ProviderProject = s.convertProviderProjectToModel(providerProject)
	}

	return projectModel, nil
}
