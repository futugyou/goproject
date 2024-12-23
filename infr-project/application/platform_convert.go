package application

import (
	"context"
	"fmt"
	"log"

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
		var dbProject *platform.PlatformProject
		for _, v := range projects {
			if v.ProviderProjectId == project.ID {
				dbProject = &v
				delete(projects, v.Id)
				break
			}
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

func (s *PlatformService) convertToPlatformModelEnvironments(values map[string]platformProvider.Env) []models.ProjectEnv {
	wps := []models.ProjectEnv{}
	for _, v := range values {
		wps = append(wps, models.ProjectEnv(v))
	}

	return wps
}

func (s *PlatformService) convertToPlatformModelWorkflows(values map[string]platformProvider.Workflow) []models.Workflow {
	wps := []models.Workflow{}
	for _, v := range values {
		wps = append(wps, models.Workflow(v))
	}

	return wps
}

func (s *PlatformService) convertToPlatformModelDeployments(values map[string]platformProvider.Deployment) []models.Deployment {
	wps := []models.Deployment{}
	for _, v := range values {
		wps = append(wps, models.Deployment(v))
	}

	return wps
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

func (s *PlatformService) getProviderProjects(ctx context.Context, provider platformProvider.IPlatformProviderAsync, src platform.Platform) ([]platformProvider.Project, error) {
	parameters := make(map[string]string)
	for _, v := range src.Properties {
		parameters[v.Key] = v.Value
	}
	filter := platformProvider.ProjectFilter{
		Parameters: parameters,
	}
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
		return nil, fmt.Errorf("createProviderProject timeout: %w", ctx.Err())
	}
}

func (s *PlatformService) deleteProviderWebhook(ctx context.Context, provider platformProvider.IPlatformProviderAsync, webhookId string, properties map[string]string) error {
	request := platformProvider.DeleteWebHookRequest{
		WebHookId:  webhookId,
		Parameters: properties,
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
		return nil, fmt.Errorf("createProviderWebhook timeout: %w", ctx.Err())
	}
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
	select {
	case user := <-resCh:
		return user, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("getProviderUser timeout: %w", ctx.Err())
	}
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
	select {
	case project := <-resCh:
		return project, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("getProviderProject timeout: %w", ctx.Err())
	}
}

// followed == false && providerProjectId == "", means need create project to provider
// followed == true && providerProjectId != "", means need provider project was already followed
// followed == false && providerProjectId != "", means need follow the provider project
func (s *PlatformService) mergeProject(providerProject *platformProvider.Project, project *platform.PlatformProject) models.PlatformProject {
	modelProject := models.PlatformProject{}
	if project != nil {
		modelProject.Id = project.Id
		modelProject.Name = project.Name
		modelProject.Secrets = s.convertToPlatformModelSecrets(project.Secrets)
		modelProject.Webhooks = s.convertToPlatformModelWebhooks(project.Webhooks)
	}

	properties := []models.Property{}
	if providerProject != nil {
		if len(modelProject.Id) == 0 {
			modelProject.Id = providerProject.ID
		}
		if len(modelProject.Name) == 0 {
			modelProject.Name = providerProject.Name
		}
		modelProject.Url = providerProject.Url
		modelProject.ProviderProjectId = providerProject.ID
		modelProject.Environments = s.convertToPlatformModelEnvironments(providerProject.Envs)
		modelProject.Workflows = s.convertToPlatformModelWorkflows(providerProject.Workflows)
		modelProject.Deployments = s.convertToPlatformModelDeployments(providerProject.Deployments)
		modelProject.BadgeURL = providerProject.BadgeURL
		modelProject.BadgeMarkdown = providerProject.BadgeMarkDown
		// TODO redesign hook
		// modelProject.Webhooks = providerProject.Hooks
		for k, v := range providerProject.Properties {
			properties = append(properties, models.Property{Key: k, Value: v})
		}
	}

	if providerProject != nil && project != nil {
		modelProject.Followed = true
		for _, v := range project.Properties {
			if _, ok := providerProject.Properties[v.Key]; !ok {
				properties = append(properties, models.Property{Key: v.Key, Value: v.Value})
			}
		}
	}

	if project != nil && providerProject == nil {
		for _, v := range project.Properties {
			properties = append(properties, models.Property{Key: v.Key, Value: v.Value})
		}
	}

	modelProject.Properties = properties

	return modelProject
}
