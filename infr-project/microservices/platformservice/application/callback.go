package application

import (
	"context"
	"os"

	"github.com/futugyou/platformservice/domain"
	"github.com/futugyou/platformservice/viewmodel"
)

func (s *PlatformService) HandleCreateProviderProject(ctx context.Context, event *viewmodel.CreateProviderProjectRequest) error {
	plat, err := s.repository.GetPlatformByIdOrNameWithoutProjects(ctx, event.PlatformID)
	if err != nil {
		return err
	}

	project, err := s.repository.GetPlatformProjectByProjectID(ctx, plat.ID, event.ProjectID)
	if err != nil {
		return err
	}

	err = s.handleProviderProjectCreate(ctx, plat, project)
	if err != nil {
		return err
	}

	plat.UpdateProject(*project)

	return s.innerService.WithUnitOfWork(ctx, func(ctx context.Context) error {
		return s.repository.Update(ctx, *plat)
	})
}

func (s *PlatformService) HandleCreateProviderWebhook(ctx context.Context, event *viewmodel.CreateProviderWebhookRequest) error {
	plat, err := s.repository.GetPlatformByIdOrNameWithoutProjects(ctx, event.PlatformID)
	if err != nil {
		return err
	}

	project, err := s.repository.GetPlatformProjectByProjectID(ctx, plat.ID, event.ProjectID)
	if err != nil {
		return err
	}

	//TODO: create or link webhook

	plat.UpdateProject(*project)

	return s.innerService.WithUnitOfWork(ctx, func(ctx context.Context) error {
		return s.repository.Update(ctx, *plat)
	})
}

func (s *PlatformService) HandleProjectScreenshot(ctx context.Context, event *viewmodel.ProjectScreenshotRequest) error {
	plat, err := s.repository.GetPlatformByIdOrNameWithoutProjects(ctx, event.PlatformID)
	if err != nil {
		return err
	}

	project, err := s.repository.GetPlatformProjectByProjectID(ctx, plat.ID, event.ProjectID)
	if err != nil {
		return err
	}

	err = s.handleScreenshot(ctx, project)
	if err != nil {
		return err
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
