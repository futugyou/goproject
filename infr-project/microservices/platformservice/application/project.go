package application

import (
	"context"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	coreinfr "github.com/futugyou/domaincore/infrastructure"
	tool "github.com/futugyou/extensions"

	"github.com/futugyou/platformservice/assembler"
	"github.com/futugyou/platformservice/domain"
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
	if len(plat.Projects) == 0 {
		plat.Projects = map[string]domain.PlatformProject{}
		for _, project := range providerProjects {
			plat.Projects[project.ID] = *projectMapper.ToDomain(&project)
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
			newProject := *projectMapper.ToDomain(&project)
			plat.Projects[project.ID] = newProject
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
	secrets, err := secretMapper.ToModel(ctx, s.vaultService, project.Secrets)
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

	s.sendProjectChangeEvent(ctx, project, plat, projectID, screenshot)

	return s.innerService.WithUnitOfWork(ctx, func(ctx context.Context) error {
		return s.repository.Update(ctx, *plat)
	})
}

func (s *PlatformService) sendProjectChangeEvent(ctx context.Context, project viewmodel.UpdatePlatformProjectRequest, plat *domain.Platform, projectID string, screenshot bool) {
	events := []coreinfr.Event{}
	if project.Operate == "sync" {
		events = append(events, &CreateProviderProjectTriggeredEvent{
			PlatformID:  plat.ID,
			ProjectID:   projectID,
			ProjectName: project.Name,
			Provider:    plat.Provider.String(),
		})
	}

	if project.ImportWebhooks {
		events = append(events, &CreateProviderWebhookTriggeredEvent{
			PlatformID:        plat.ID,
			ProjectID:         projectID,
			ProjectName:       project.Name,
			Provider:          plat.Provider.String(),
			ProviderProjectId: project.ProviderProjectID,
		})

	}

	if screenshot {
		events = append(events, &ProjectScreenshotTriggeredEvent{
			PlatformID: plat.ID,
			ProjectID:  projectID,
			Url:        domain.GetWebhookUrl(plat.Name, project.Name),
		})
	}

	if len(events) > 0 {
		if err := s.eventPublisher.DispatchIntegrationEvents(ctx, events); err != nil {
			// Why not just report an error? Because `IntegrationEvents` are not that important.
			log.Println(err.Error())
		}
	}
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

func (s *PlatformService) HandlePlatformProjectUpsert(ctx context.Context, event *viewmodel.PlatformProjectUpsertEvent) error {
	plat, err := s.repository.GetPlatformByIdOrNameWithoutProjects(ctx, event.PlatformID)
	if err != nil {
		return err
	}

	project, err := s.repository.GetPlatformProjectByProjectID(ctx, plat.ID, event.ProjectID)
	if err != nil {
		return err
	}

	if event.CreateProviderProject {
		err1 := s.handleProviderProjectCreate(ctx, plat, project)
		if err1 != nil {
			return err1
		}
	}

	if event.ImportWebhooks {
		//TODO:
	}

	if event.Screenshot {
		s.handleScreenshot(ctx, project)
	}

	plat.UpdateProject(*project)

	return s.innerService.WithUnitOfWork(ctx, func(ctx context.Context) error {
		return s.repository.Update(ctx, *plat)
	})
}

func (s *PlatformService) handleScreenshot(ctx context.Context, project *domain.PlatformProject) error {
	if len(project.Url) > 0 && os.Getenv("SCREENSHOT_ALLOW") != "false" {
		imageUrl, err := s.screenshot.Create(ctx, project.Url)
		if err != nil {
			return err
		}

		project.UpdateImageUrl(*imageUrl)
	}

	return nil
}

func (s *PlatformService) handleProviderProjectCreate(ctx context.Context, plat *domain.Platform, project *domain.PlatformProject) error {
	provider, err := s.getPlatformProvider(ctx, *plat)
	if err != nil {
		return err
	}

	parameters := mergePropertiesToMap(project.Properties, plat.Properties)
	providerProject, err := s.getProviderProject(ctx, provider, project.ProviderProjectID, parameters)
	if err != nil {
		return err
	}

	if providerProject == nil || len(providerProject.ID) == 0 {
		var err error
		if providerProject, err = s.createProviderProject(ctx, provider, project.Name, project.Properties); err != nil {
			return err
		}
	}

	if providerProject != nil && len(providerProject.ID) > 0 {
		project.UpdateProviderProjectID(providerProject.ID)
	}

	return nil
}
