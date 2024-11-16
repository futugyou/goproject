package application

import (
	"context"
	"fmt"

	platform "github.com/futugyou/infr-project/platform"
	platformProvider "github.com/futugyou/infr-project/platform_provider"
	vault "github.com/futugyou/infr-project/vault"
	models "github.com/futugyou/infr-project/view_models"
)

func (s *PlatformService) convertPlatformEntityToViewModel(ctx context.Context, src *platform.Platform) (*models.PlatformDetailView, error) {
	if src == nil {
		return nil, fmt.Errorf("no platform data found")
	}

	providerProjects := []platformProvider.Project{}
	if provider, err := s.getPlatfromProvider(ctx, *src); err == nil {
		projects, _ := s.getProviderProjects(ctx, provider)
		providerProjects = projects
	}

	return &models.PlatformDetailView{
		Id:         src.Id,
		Name:       src.Name,
		Activate:   src.Activate,
		Url:        src.Url,
		Properties: s.convertToPlatformModelProperties(src.Properties),
		Secrets:    s.convertToPlatformModelSecrets(src.Secrets),
		Projects:   s.convertToPlatformModelProjects(src.Projects, providerProjects),
		Tags:       src.Tags,
		IsDeleted:  src.IsDeleted,
		Provider:   src.Provider.String(),
	}, nil
}

func (s *PlatformService) convertToPlatformModelProjects(projects map[string]platform.PlatformProject, providerProjects []platformProvider.Project) []models.PlatformProject {
	platformProjects := make([]models.PlatformProject, 0)

	providerMap := map[string]models.PlatformProject{}
	if len(providerProjects) > 0 {
		for _, project := range providerProjects {
			providerMap[project.ID] = models.PlatformProject{
				Id:                project.ID,
				Name:              project.Name,
				Url:               project.Url,
				Properties:        []models.Property{},
				Secrets:           []models.Secret{},
				Webhooks:          []models.Webhook{},
				Followed:          false,
				ProviderProjectId: project.ID,
			}
		}
	}

	// convert db project to model project
	for _, v := range projects {
		// followed == false && providerProjectId == "", means need create project to provider
		// followed == true && providerProjectId != "", means need provider project was already followed
		followed := false
		providerProjectId := ""
		if _, ok := providerMap[v.ProviderProjectId]; ok {
			followed = true
			providerProjectId = v.ProviderProjectId
			// remove already followed project
			delete(providerMap, v.ProviderProjectId)
		}

		platformProjects = append(platformProjects, models.PlatformProject{
			Id:                v.Id,
			Name:              v.Name,
			Url:               v.Url,
			Properties:        s.convertToPlatformModelProperties(v.Properties),
			Secrets:           s.convertToPlatformModelSecrets(v.Secrets),
			Webhooks:          s.convertToPlatformModelWebhooks(v.Webhooks),
			Followed:          followed,
			ProviderProjectId: providerProjectId,
		})
	}

	// add project from provider which was not followed
	// followed == false && providerProjectId != "", means need follow the provider project
	for _, v := range providerMap {
		platformProjects = append(platformProjects, v)
	}

	return platformProjects
}

func (s *PlatformService) convertToPlatformModelWebhooks(hooks []platform.Webhook) []models.Webhook {
	webhooks := make([]models.Webhook, 0)
	for i := 0; i < len(hooks); i++ {
		webhooks = append(webhooks, models.Webhook{
			Name:       hooks[i].Name,
			Url:        hooks[i].Url,
			Activate:   hooks[i].Activate,
			State:      hooks[i].State.String(),
			Properties: s.convertToPlatformModelProperties(hooks[i].Properties),
			Secrets:    s.convertToPlatformModelSecrets(hooks[i].Secrets),
		})
	}

	return webhooks
}

func (s *PlatformService) convertToPlatformModelSecrets(secretMap map[string]platform.Secret) []models.Secret {
	secrets := []models.Secret{}
	for _, v := range secretMap {
		secrets = append(secrets, models.Secret{
			Key:       v.Key,
			VaultId:   v.Value,
			VaultKey:  v.VaultKey,
			MaskValue: v.VaultMaskValue,
		})
	}

	return secrets
}

func (s *PlatformService) convertToPlatformSecrets(ctx context.Context, secrets []models.Secret) (map[string]platform.Secret, error) {
	secretInfos := make(map[string]platform.Secret)
	if len(secrets) == 0 {
		return secretInfos, nil
	}

	filter := []vault.VaultSearch{}
	for _, secret := range secrets {
		filter = append(filter, vault.VaultSearch{
			ID: secret.VaultId,
		})
	}

	query := VaultSearchQuery{Filters: filter, Page: 0, Size: 0}
	if vaults, err := s.vaultService.SearchVaults(ctx, query); err == nil {
		for _, secret := range secrets {
			for i := 0; i < len(vaults); i++ {
				if vaults[i].Id == secret.VaultId {
					secretInfos[secret.Key] = platform.Secret{
						Key:            secret.Key,
						Value:          vaults[i].Id,
						VaultKey:       vaults[i].Key,
						VaultMaskValue: vaults[i].MaskValue,
					}
					break
				}
			}
		}
	} else {
		return nil, err
	}

	return secretInfos, nil
}

func (s *PlatformService) convertToPlatformModelProperties(properties map[string]platform.Property) []models.Property {
	wps := []models.Property{}
	for _, v := range properties {
		wps = append(wps, models.Property(v))
	}

	return wps
}

func (s *PlatformService) convertToPlatformProperties(propertyList []models.Property) map[string]platform.Property {
	properties := make(map[string]platform.Property)
	for _, v := range propertyList {
		properties[v.Key] = platform.Property{
			Key:   v.Key,
			Value: v.Value,
		}
	}

	return properties
}

func (s *PlatformService) convertToPlatformViews(src []platform.Platform) []models.PlatformView {
	result := make([]models.PlatformView, len(src))
	for i := 0; i < len(src); i++ {
		result[i] = models.PlatformView{
			Id:        src[i].Id,
			Name:      src[i].Name,
			Activate:  src[i].Activate,
			Url:       src[i].Url,
			Tags:      src[i].Tags,
			IsDeleted: src[i].IsDeleted,
			Provider:  src[i].Provider.String(),
		}
	}

	return result
}

func (s *PlatformService) getPlatfromProvider(ctx context.Context, src platform.Platform) (platformProvider.IPlatformProviderAsync, error) {
	vaultId, err := src.ProviderVaultInfo()
	if err != nil {
		return nil, err
	}

	token, err := s.vaultService.ShowVaultRawValue(ctx, vaultId)
	if err != nil {
		return nil, fmt.Errorf("get platfrom provider token error, vaultId is %s, message %s", vaultId, err.Error())
	}

	return platformProvider.PlatformProviderFatory(src.Provider.String(), token)
}

func (s *PlatformService) getProviderProjects(ctx context.Context, provider platformProvider.IPlatformProviderAsync) ([]platformProvider.Project, error) {
	filter := platformProvider.ProjectFilter{}
	resCh, errCh := provider.ListProjectAsync(ctx, filter)
	select {
	case projects := <-resCh:
		return projects, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("getProviderProjects timeout: %w", ctx.Err())
	}
}

func (s *PlatformService) createProviderProject(ctx context.Context, provider platformProvider.IPlatformProviderAsync, name string, properties map[string]platform.Property) (*platformProvider.Project, error) {
	request := platformProvider.CreateProjectRequest{
		Name:       name,
		Parameters: map[string]string{},
	}

	for _, v := range properties {
		request.Parameters[v.Key] = v.Value
	}

	resCh, errCh := provider.CreateProjectAsync(ctx, request)
	select {
	case project := <-resCh:
		return project, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("getProviderProject timeout: %w", ctx.Err())
	}
}

func (s *PlatformService) deleteProviderWebhook(ctx context.Context, provider platformProvider.IPlatformProviderAsync, webhookId string, properties map[string]platform.Property) error {
	request := platformProvider.DeleteWebHookRequest{
		WebHookId:  webhookId,
		Parameters: map[string]string{},
	}

	for _, v := range properties {
		request.Parameters[v.Key] = v.Value
	}

	errCh := provider.DeleteWebHookAsync(ctx, request)
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return fmt.Errorf("deleteProviderWebhook timeout: %w", ctx.Err())
	}
}

func (s *PlatformService) createProviderWebhook(ctx context.Context, provider platformProvider.IPlatformProviderAsync,
	platformId string, projectId string, url string, name string) (*platformProvider.WebHook, error) {
	request := platformProvider.CreateWebHookRequest{
		PlatformId: platformId,
		ProjectId:  projectId,
		WebHook: platformProvider.WebHook{
			Name: name,
			Url:  url,
		},
	}

	resCh, errCh := provider.CreateWebHookAsync(ctx, request)
	select {
	case webhook := <-resCh:
		return webhook, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("getProviderProject timeout: %w", ctx.Err())
	}
}
