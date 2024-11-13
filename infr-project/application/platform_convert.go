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
	if provider, err := s.getPlatfromProvider(ctx, *src); err != nil {
		projects, _ := s.getProviderProjects(ctx, provider)
		providerProjects = projects
	}

	return &models.PlatformDetailView{
		Id:         src.Id,
		Name:       src.Name,
		Activate:   src.Activate,
		Url:        src.Url,
		Properties: s.convertToPlatformModelProperty(src.Properties),
		Secrets:    s.convertToPlatformModelSecret(src.Secrets),
		Projects:   s.convertToPlatformModelProjects(src.Projects, providerProjects),
		Tags:       src.Tags,
		IsDeleted:  src.IsDeleted,
		Provider:   src.Provider.String(),
	}, nil
}

func (s *PlatformService) convertToPlatformModelProjects(projects map[string]platform.PlatformProject, providerProjects []platformProvider.Project) []models.PlatformProject {
	platformProjects := make([]models.PlatformProject, 0)

	if len(providerProjects) > 0 {
		//TODO:
	}

	for _, v := range projects {
		platformProjects = append(platformProjects, models.PlatformProject{
			Id:         v.Id,
			Name:       v.Name,
			Followed:   false,
			Url:        v.Url,
			Properties: s.convertToPlatformModelProperty(v.Properties),
			Secrets:    s.convertToPlatformModelSecret(v.Secrets),
			Webhooks:   s.convertToPlatformModelWebhook(v.Webhooks),
		})
	}

	return platformProjects
}

func (s *PlatformService) convertToPlatformModelWebhook(hooks []platform.Webhook) []models.Webhook {
	webhooks := make([]models.Webhook, 0)
	for i := 0; i < len(hooks); i++ {
		webhooks = append(webhooks, models.Webhook{
			Name:       hooks[i].Name,
			Url:        hooks[i].Url,
			Activate:   hooks[i].Activate,
			State:      hooks[i].State.String(),
			Properties: s.convertToPlatformModelProperty(hooks[i].Properties),
			Secrets:    s.convertToPlatformModelSecret(hooks[i].Secrets),
		})
	}

	return webhooks
}

func (s *PlatformService) convertToPlatformModelProperty(properties map[string]platform.Property) []models.Property {
	wps := []models.Property{}
	for _, v := range properties {
		wps = append(wps, models.Property(v))
	}

	return wps
}

func (s *PlatformService) convertToPlatformModelSecret(secretMap map[string]platform.Secret) []models.Secret {
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

func (s *PlatformService) convertToPlatformView(src []platform.Platform) []models.PlatformView {
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
	var provider string
	var vaultId string
	var token string
	var err error

	switch src.Provider {
	case platform.PlatformProviderCircleci:
		provider = platform.PlatformProviderCircleci.String()
		vaultId = src.Secrets["CIRCLECI_TOKEN"].Value
	case platform.PlatformProviderVercel:
		provider = platform.PlatformProviderVercel.String()
		vaultId = src.Secrets["VERCEL_TOKEN"].Value
	case platform.PlatformProviderGithub:
		provider = platform.PlatformProviderGithub.String()
		vaultId = src.Secrets["GITHUB_TOKEN"].Value
	default:
		return nil, fmt.Errorf("%s not supported", src.Provider.String())
	}

	token, err = s.vaultService.ShowVaultRawValue(ctx, vaultId)
	if err != nil {
		return nil, fmt.Errorf("get platfrom provider token err, vaultId is %s, message %s", vaultId, err.Error())
	}

	return platformProvider.PlatformProviderFatory(provider, token)
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
