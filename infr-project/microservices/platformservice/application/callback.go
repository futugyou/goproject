package application

import (
	"context"
	"fmt"

	tool "github.com/futugyou/extensions"

	"github.com/futugyou/platformservice/domain"
	platformprovider "github.com/futugyou/platformservice/provider"
	"github.com/futugyou/platformservice/viewmodel"
)

func (s *PlatformService) HandleCreateProviderProject(ctx context.Context, event *viewmodel.CreateProviderProjectRequest) error {
	plat, err := s.repository.GetPlatformByIdOrNameWithoutProjects(ctx, event.PlatformID)
	if err != nil {
		return err
	}

	project, err := s.repository.GetPlatformProjectByIDOrName(ctx, plat.ID, event.ProjectID)
	if err != nil {
		return err
	}

	e, err := s.eventHandler.Get(ctx, event.ID)
	if err != nil {
		return err
	}

	ctx = tool.WithJWT(ctx, e.Token)
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

	project, err := s.repository.GetPlatformProjectByIDOrName(ctx, plat.ID, event.ProjectID)
	if err != nil {
		return err
	}

	if project.ProviderProjectID != event.ProviderProjectId {
		return fmt.Errorf("project provider project id not match")
	}

	e, err := s.eventHandler.Get(ctx, event.ID)
	if err != nil {
		return err
	}

	ctx = tool.WithJWT(ctx, e.Token)
	err = s.handleProviderWebhookCreate(ctx, plat, project, event.Url)
	if err != nil {
		return err
	}

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

	project, err := s.repository.GetPlatformProjectByIDOrName(ctx, plat.ID, event.ProjectID)
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
	if len(project.Url) > 0 && s.opts.ScreenshotAllow == "true" {
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

func (s *PlatformService) handleProviderWebhookCreate(ctx context.Context, plat *domain.Platform, project *domain.PlatformProject, url string) error {
	provider, err := s.getPlatformProvider(ctx, *plat)
	if err != nil {
		return err
	}

	parameters := mergePropertiesToMap(project.Properties, plat.Properties)
	hook, err := provider.GetWebHookByUrl(ctx, platformprovider.GetWebHookRequest{
		Parameters: parameters,
		Url:        url,
		ProjectID:  project.ProviderProjectID,
	})
	if err != nil {
		return err
	}

	if hook != nil {
		s.updateProjectWebhook(hook, project,hook.Parameters["SigningSecret"])

		return nil
	}

	secret, _ := tool.GenerateRandomKey(6)
	hook, err = provider.CreateWebHook(ctx, platformprovider.CreateWebHookRequest{
		PlatformID:    plat.ID,
		ProjectID:     project.ProviderProjectID,
		VerifyTLS:     true,
		SigningSecret: secret,
		Name:          project.Name,
		Url:           url,
		Parameters:    parameters,
	})
	if err != nil {
		return err
	}

	// Besides vercel, `hook.SigningSecret should be the same as the `secret`, because the vercel's secret is obtained after creation.
	s.updateProjectWebhook(hook, project,hook.SigningSecret)

	return nil
}

func (*PlatformService) updateProjectWebhook(hook *platformprovider.WebHook, project *domain.PlatformProject,secret string) {
	opts := []domain.WebhookOption{
		domain.WithWebhookEvents(hook.Events),
		domain.WithWebhookID(hook.ID),
		domain.WithWebhookState(domain.WebhookReady),
	}

	if   len(secret) > 0 {
		opts = append(opts, domain.WithWebhookSigningSecret(secret))
	}

	webhook := domain.NewWebhook(hook.Url, opts...)

	project.UpdateWebhook(webhook)
}
