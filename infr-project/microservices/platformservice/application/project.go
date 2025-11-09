package application

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"

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

	if project.Operate == "sync" || project.ImportWebhooks || screenshot {
		event := &viewmodel.PlatformProjectUpsertEvent{
			PlatformID:            plat.ID,
			ProjectID:             projectID,
			CreateProviderProject: project.Operate == "sync",
			ImportWebhooks:        project.ImportWebhooks,
			EventName:             "upsert_project",
			Screenshot:            screenshot,
		}
		if err := s.eventPublisher.DispatchIntegrationEvent(ctx, event); err != nil {
			log.Println(err.Error())
		}
	}

	return s.innerService.WithUnitOfWork(ctx, func(ctx context.Context) error {
		return s.repository.Update(ctx, *plat)
	})
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
