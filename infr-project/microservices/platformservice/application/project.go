package application

import (
	"context"
	"fmt"
	"slices"

	tool "github.com/futugyou/extensions"

	"github.com/futugyou/platformservice/assembler"
	"github.com/futugyou/platformservice/domain"
	platformprovider "github.com/futugyou/platformservice/provider"
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
	projects, err := s.getProviderProjects(ctx, provider, *plat)
	if err != nil {
		return err
	}

	if len(projects) == 0 {
		return fmt.Errorf("no provider project found in %s", plat.Provider.String())
	}

	projects = tool.ArrayFilter(projects, func(project platformprovider.Project) bool {
		if len(providerProjectIDs) == 0 {
			return true
		}
		return slices.Contains(providerProjectIDs, project.ID)
	})

	projectMapper := assembler.ProjectAssembler{}
	if len(plat.Projects) == 0 {
		plat.Projects = map[string]domain.PlatformProject{}
		for _, project := range projects {
			plat.Projects[project.ID] = *projectMapper.ToDomain(&project)
		}
	} else {
		for _, project := range projects {
			find := false
			for _, entity := range plat.Projects {
				// The following situations are considered to be `link`
				// 1. already `linked`: ProviderProjectId = ID
				// 2. name is the same, this is a normal situation
				// 3. ID is the same, linked, then cancel(ProviderProjectId is ""), then link again
				if entity.ProviderProjectId == project.ID || entity.Name == project.Name || entity.ID == project.ID {
					find = true
					rawProject := plat.Projects[project.ID]
					rawProject.ProviderProjectId = project.ID
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

	for _, v := range plat.Projects {
		if _, ok := bMap[v.ProviderProjectId]; ok {
			return fmt.Errorf("provider project id: %s is duplicated", v.ProviderProjectId)
		}
		bMap[v.ProviderProjectId] = struct{}{}
	}

	return s.innerService.WithUnitOfWork(ctx, func(ctx context.Context) error {
		return s.repository.Update(ctx, *plat)
	})
}
