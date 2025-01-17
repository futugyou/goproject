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
	if provider, err := s.getPlatfromProvider(ctx, *src); err == nil {
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
	provider, err := s.getPlatfromProvider(ctx, *res)
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

// followed == false && providerProjectId == "", means need create project to provider
// followed == true && providerProjectId != "", means need provider project was already followed
// followed == false && providerProjectId != "", means need follow the provider project
func (s *PlatformService) mergeProject(providerProject *platformProvider.Project, project *platform.PlatformProject) models.PlatformProject {
	modelProject := models.PlatformProject{
		Properties:           []models.Property{},
		Secrets:              []models.Secret{},
		Webhooks:             []models.Webhook{},
		ProviderProjectId:    "",
		EnvironmentVariables: []models.EnvironmentVariable{},
		Environments:         []string{},
		Workflows:            []models.Workflow{},
		WorkflowRuns:         []models.WorkflowRun{},
		Deployments:          []models.Deployment{},
	}

	if project != nil {
		modelProject.Id = project.Id
		modelProject.Description = project.Description
		modelProject.Name = project.Name
		modelProject.Url = project.Url
		modelProject.Secrets = s.convertToPlatformModelSecrets(project.Secrets)
	}

	if providerProject != nil {
		if len(modelProject.Id) == 0 {
			modelProject.Id = providerProject.ID
		}
		if len(modelProject.Name) == 0 {
			modelProject.Name = providerProject.Name
		}
		if len(modelProject.Description) == 0 {
			modelProject.Description = providerProject.Description
		}
		if len(modelProject.Url) == 0 {
			modelProject.Url = providerProject.Url
		}
		if len(modelProject.Secrets) == 0 {
			modelProject.Secrets = []models.Secret{}
		}

		modelProject.ProviderProjectId = providerProject.ID
		modelProject.EnvironmentVariables = s.convertToPlatformModelEnvironments(providerProject.EnvironmentVariables)
		modelProject.Environments = providerProject.Environments
		modelProject.Workflows = s.convertToPlatformModelWorkflows(providerProject.Workflows)
		modelProject.WorkflowRuns = s.convertToPlatformModelWorkflowRuns(providerProject.WorkflowRuns)
		modelProject.Deployments = s.convertToPlatformModelDeployments(providerProject.Deployments)
		modelProject.BadgeURL = providerProject.BadgeURL
		modelProject.BadgeMarkdown = providerProject.BadgeMarkDown
	}

	modelProject.Webhooks = s.mergeWebhooks(project.GetWebhooks(), providerProject.GetWebhooks())

	if providerProject != nil && project != nil {
		modelProject.Followed = true
	}

	modelProject.Properties = s.mergeProperties(providerProject.GetProperties(), project.GetProperties())

	return modelProject
}

func (s *PlatformService) mergeProperties(providerProperty map[string]string, projectProperty map[string]platform.Property) []models.Property {
	properties := []models.Property{}
	// 1. add provider project property
	for k, v := range providerProperty {
		properties = append(properties, models.Property{Key: k, Value: v})
	}

	// 2. add property which is in project and not in provider
	for _, v := range projectProperty {
		if _, ok := providerProperty[v.Key]; !ok {
			properties = append(properties, models.Property{Key: v.Key, Value: v.Value})
		}
	}

	return properties
}

func (s *PlatformService) mergeWebhooks(platformWebhooks []platform.Webhook, providerWebHooks []platformProvider.WebHook) []models.Webhook {
	webhooks := []models.Webhook{}
	for _, hook := range platformWebhooks {
		mw := models.Webhook{
			ID:         hook.ID,
			Name:       hook.Name,
			Url:        hook.Url,
			Events:     hook.Events,
			Activate:   hook.Activate,
			State:      hook.State.String(),
			Properties: s.convertToPlatformModelProperties(hook.Properties),
			Secrets:    s.convertToPlatformModelSecrets(hook.Secrets),
			Followed:   false,
		}
		prow := tool.ArrayFirst(providerWebHooks, func(t platformProvider.WebHook) bool { return t.ID == hook.ID && t.Url == platform.GetWebhookUrl() })
		if prow != nil {
			mw.Name = prow.Name
			mw.Url = prow.Url
			mw.Activate = prow.Activate
			mw.Followed = true
			mw.Properties = s.mergeWebhookProperty(prow.GetParameters(), mw.Properties)
			mw.Events = prow.Events
		}

		webhooks = append(webhooks, mw)
	}

	for _, prow := range providerWebHooks {
		if prow.Url != platform.GetWebhookUrl() {
			continue
		}

		hook := tool.ArrayFirst(platformWebhooks, func(t platform.Webhook) bool { return t.ID == prow.ID })
		if hook == nil {
			webhooks = append(webhooks, models.Webhook{
				Name:       prow.Name,
				Url:        prow.Url,
				Activate:   prow.Activate,
				State:      "Ready",
				Properties: s.mergeWebhookProperty(prow.GetParameters(), nil),
				Secrets:    []models.Secret{},
				Followed:   true,
				ID:         prow.ID,
			})
		}
	}

	return webhooks
}

func (*PlatformService) mergeWebhookProperty(providerProperties map[string]string, properties []models.Property) []models.Property {
	if properties == nil {
		properties = []models.Property{}
	}

	for key, v := range providerProperties {
		if f := tool.ArrayFirst(properties, func(p models.Property) bool {
			return p.Key == key
		}); f == nil {
			properties = append(properties, models.Property{
				Key:   key,
				Value: v,
			})
		}
	}
	return properties
}
