package application

import (
	"context"
	"fmt"

	platform "github.com/futugyou/infr-project/platform"
	vault "github.com/futugyou/infr-project/vault"
	models "github.com/futugyou/infr-project/view_models"
)

func (s *PlatformService) convertPlatformEntityToViewModel(ctx context.Context, src *platform.Platform) (*models.PlatformDetailView, error) {
	if src == nil {
		return nil, fmt.Errorf("no platform data found")
	}

	platformProjects := make([]models.PlatformProject, 0)
	for _, v := range src.Projects {
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

	if provider, err := s.getPlatfromProvider(ctx, *src); err != nil {
		providerProjects, _ := s.getProviderProjects(ctx, provider)
		for _, p := range providerProjects {
			has := false
			for _, pp := range platformProjects {
				if pp.Id == p.Id {
					has = true
					break
				}
			}
			if !has {
				platformProjects = append(platformProjects, p)
			}
		}
	}

	properties := make([]models.Property, 0)
	for _, v := range src.Properties {
		properties = append(properties, models.Property(v))
	}

	return &models.PlatformDetailView{
		Id:         src.Id,
		Name:       src.Name,
		Activate:   src.Activate,
		Url:        src.Url,
		Properties: properties,
		Secrets:    s.convertToPlatformModelSecret(src.Secrets),
		Projects:   platformProjects,
		Tags:       src.Tags,
		IsDeleted:  src.IsDeleted,
		Provider:   src.Provider.String(),
	}, nil
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

func (s *PlatformService) convertPlatformProperties(propertyList []models.Property) map[string]platform.Property {
	properties := make(map[string]platform.Property)
	for _, v := range propertyList {
		properties[v.Key] = platform.Property{
			Key:   v.Key,
			Value: v.Value,
		}
	}

	return properties
}
