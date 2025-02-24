package application

import (
	"context"
	"fmt"
	"log"

	tool "github.com/futugyou/extensions"

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
	if provider, err := s.getPlatformProvider(ctx, *src); err == nil {
		projects, _ := s.getProviderProjects(ctx, provider, *src)
		providerProjects = projects
	} else {
		log.Println(err.Error())
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

func (s *PlatformService) convertToPlatformModelProjects(projects map[string]platform.PlatformProject,
	providerProjects []platformProvider.Project) []models.PlatformProject {
	platformProjects := make([]models.PlatformProject, 0)
	for _, project := range providerProjects {
		dbProject := tool.ArrayFirst(tool.GetMapValues(projects), func(p platform.PlatformProject) bool {
			return p.ProviderProjectId == project.ID
		})
		if dbProject != nil {
			delete(projects, dbProject.Id)
		}

		modelProject := s.mergeProject(&project, dbProject)
		platformProjects = append(platformProjects, modelProject)
	}

	for _, v := range projects {
		modelProject := s.mergeProject(nil, &v)
		platformProjects = append(platformProjects, modelProject)
	}

	return platformProjects
}

func (s *PlatformService) convertToPlatformModelEnvironments(values map[string]platformProvider.EnvironmentVariable) []models.EnvironmentVariable {
	return tool.MapToSlice(values, func(key string, v platformProvider.EnvironmentVariable) models.EnvironmentVariable {
		return models.EnvironmentVariable(v)
	})
}

func (s *PlatformService) convertToPlatformModelWorkflows(values map[string]platformProvider.Workflow) []models.Workflow {
	return tool.MapToSlice(values, func(key string, v platformProvider.Workflow) models.Workflow {
		return models.Workflow(v)
	})
}
func (s *PlatformService) convertToPlatformModelWorkflowRuns(values map[string]platformProvider.WorkflowRun) []models.WorkflowRun {
	return tool.MapToSlice(values, func(key string, v platformProvider.WorkflowRun) models.WorkflowRun {
		return models.WorkflowRun(v)
	})
}

func (s *PlatformService) convertToPlatformModelDeployments(values map[string]platformProvider.Deployment) []models.Deployment {
	return tool.MapToSlice(values, func(key string, v platformProvider.Deployment) models.Deployment {
		return models.Deployment(v)
	})
}

func (s *PlatformService) convertToPlatformModelSecrets(secretMap map[string]platform.Secret) []models.Secret {
	return tool.MapToSlice(secretMap, func(key string, v platform.Secret) models.Secret {
		return models.Secret{
			Key:       v.Key,
			VaultId:   v.Value,
			VaultKey:  v.VaultKey,
			MaskValue: v.VaultMaskValue,
		}
	})
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
					if vaults[i].VaultType == "system" {
						return nil, fmt.Errorf("system vault can not be use")
					}

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
	return tool.MapToSlice(properties, func(key string, env platform.Property) models.Property {
		return models.Property(env)
	})
}

func (s *PlatformService) convertToPlatformProperties(propertyList []models.Property) map[string]platform.Property {
	return tool.SliceToMapWithTransform(propertyList, func(v models.Property) string { return v.Key },
		func(v models.Property) platform.Property {
			return platform.Property{
				Key:   v.Key,
				Value: v.Value,
			}
		})
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

func (s *PlatformService) getPlatformProvider(ctx context.Context, src platform.Platform) (platformProvider.IPlatformProviderAsync, error) {
	vaultId, err := src.ProviderVaultInfo()
	if err != nil {
		return nil, err
	}

	token, err := s.vaultService.ShowVaultRawValue(ctx, vaultId)
	if err != nil {
		return nil, fmt.Errorf("get platform provider token error, vaultId is %s, message %s", vaultId, err.Error())
	}

	return platformProvider.PlatformProviderFactory(src.Provider.String(), token)
}

func (s *PlatformService) getProviderProjects(ctx context.Context, provider platformProvider.IPlatformProviderAsync, src platform.Platform) ([]platformProvider.Project, error) {
	parameters := make(map[string]string)
	for _, v := range src.Properties {
		parameters[v.Key] = v.Value
	}
	filter := platformProvider.ProjectFilter{
		Parameters: parameters,
	}

	resCh, errCh := provider.ListProjectAsync(ctx, filter)
	return tool.HandleAsync(ctx, resCh, errCh)
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
	return tool.HandleAsync(ctx, resCh, errCh)
}

func (s *PlatformService) deleteProviderWebhook(ctx context.Context, provider platformProvider.IPlatformProviderAsync, webhookId string, properties map[string]string) error {
	request := platformProvider.DeleteWebHookRequest{
		WebHookId:  webhookId,
		Parameters: properties,
	}

	errCh := provider.DeleteWebHookAsync(ctx, request)
	return tool.HandleErrorAsync(ctx, errCh)
}

func (s *PlatformService) createProviderWebhook(ctx context.Context, provider platformProvider.IPlatformProviderAsync,
	platformId string, projectId string, url string, name string, secret string) (*platformProvider.WebHook, error) {
	request := platformProvider.CreateWebHookRequest{
		PlatformId: platformId,
		ProjectId:  projectId,
		WebHook: platformProvider.WebHook{
			Name: name,
			Url:  url,
			Parameters: map[string]string{
				"SigningSecret": secret,
			},
		},
	}

	resCh, errCh := provider.CreateWebHookAsync(ctx, request)
	return tool.HandleAsync(ctx, resCh, errCh)
}

func mergePropertiesToMap(propertiesList ...map[string]platform.Property) map[string]string {
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

func (s *PlatformService) getProviderUser(ctx context.Context, provider platformProvider.IPlatformProviderAsync) (*platformProvider.User, error) {
	resCh, errCh := provider.GetUserAsync(ctx)
	return tool.HandleAsync(ctx, resCh, errCh)
}

func (s *PlatformService) determineProviderStatus(ctx context.Context, res *platform.Platform) bool {
	provider, err := s.getPlatformProvider(ctx, *res)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	user, err := s.getProviderUser(ctx, provider)
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

func (s *PlatformService) getProviderProject(ctx context.Context, provider platformProvider.IPlatformProviderAsync, name string, parameters map[string]string) (*platformProvider.Project, error) {
	filter := platformProvider.ProjectFilter{
		Parameters: parameters,
		Name:       name,
	}
	resCh, errCh := provider.GetProjectAsync(ctx, filter)
	return tool.HandleAsync(ctx, resCh, errCh)
}

func (*PlatformService) mapToModelProperty(providerProperties map[string]string) []models.Property {
	properties := []models.Property{}
	for key, v := range providerProperties {
		properties = append(properties, models.Property{
			Key:   key,
			Value: v,
		})
	}

	return properties
}

func (s *PlatformService) mergeProject(providerProject *platformProvider.Project, project *platform.PlatformProject) models.PlatformProject {
	modelProject := models.PlatformProject{
		Properties: []models.Property{},
		Secrets:    []models.Secret{},
		Webhooks:   []models.Webhook{},
		ImageData:  []byte{},
		ProviderProject: models.PlatformProviderProject{
			WebHooks:             []models.Webhook{},
			Properties:           []models.Property{},
			EnvironmentVariables: []models.EnvironmentVariable{},
			Environments:         []string{},
			Workflows:            []models.Workflow{},
			WorkflowRuns:         []models.WorkflowRun{},
			Deployments:          []models.Deployment{},
		},
	}

	if project != nil {
		modelProject.Id = project.Id
		modelProject.Description = project.Description
		modelProject.ImageData = project.ImageData
		modelProject.ImageUrl = project.ImageUrl
		modelProject.Name = project.Name
		modelProject.Url = project.Url
		modelProject.Secrets = s.convertToPlatformModelSecrets(project.Secrets)
		modelProject.Properties = s.convertToPlatformModelProperties(project.Properties)
		webhooks := []models.Webhook{}
		for _, hook := range project.Webhooks {
			followed := false
			if len(hook.ID) > 0 {
				followed = true
			}
			mw := models.Webhook{
				ID:         hook.ID,
				Name:       hook.Name,
				Url:        hook.Url,
				Events:     hook.Events,
				Activate:   hook.Activate,
				State:      hook.State.String(),
				Properties: s.convertToPlatformModelProperties(hook.Properties),
				Secrets:    s.convertToPlatformModelSecrets(hook.Secrets),
				Followed:   followed,
			}

			webhooks = append(webhooks, mw)
		}
		modelProject.Webhooks = webhooks
	}

	if providerProject != nil {
		if len(modelProject.Id) > 0 {
			modelProject.Followed = true
		} else {
			modelProject.Followed = false
		}

		modelProject.ProviderProjectId = providerProject.ID
		modelProject.ProviderProject.Id = providerProject.ID
		modelProject.ProviderProject.Name = providerProject.Name
		modelProject.ProviderProject.Description = providerProject.Description
		modelProject.ProviderProject.Url = providerProject.Url
		modelProject.ProviderProject.EnvironmentVariables = s.convertToPlatformModelEnvironments(providerProject.EnvironmentVariables)
		modelProject.ProviderProject.Environments = providerProject.Environments
		modelProject.ProviderProject.Workflows = s.convertToPlatformModelWorkflows(providerProject.Workflows)
		modelProject.ProviderProject.WorkflowRuns = s.convertToPlatformModelWorkflowRuns(providerProject.WorkflowRuns)
		modelProject.ProviderProject.Deployments = s.convertToPlatformModelDeployments(providerProject.Deployments)
		modelProject.ProviderProject.BadgeURL = providerProject.BadgeURL
		modelProject.ProviderProject.BadgeMarkdown = providerProject.BadgeMarkDown
		modelProject.ProviderProject.Properties = s.mapToModelProperty(providerProject.Properties)
		webhooks := []models.Webhook{}
		for _, prow := range providerProject.WebHooks {
			dbWebhook := tool.ArrayFirst(modelProject.Webhooks, func(p models.Webhook) bool {
				return p.ID == prow.ID
			})
			followed := false
			if dbWebhook != nil && len(dbWebhook.ID) > 0 {
				followed = true
			}

			webhooks = append(webhooks, models.Webhook{
				Name:       prow.Name,
				Url:        prow.Url,
				Activate:   prow.Activate,
				State:      "Ready",
				Properties: s.mapToModelProperty(prow.GetParameters()),
				Secrets:    []models.Secret{},
				Followed:   followed,
				ID:         prow.ID,
			})
		}
		modelProject.ProviderProject.WebHooks = webhooks
	}

	return modelProject
}
