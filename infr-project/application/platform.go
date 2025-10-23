package application

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tool "github.com/futugyou/extensions"

	"github.com/redis/go-redis/v9"

	domain "github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/extensions"
	infra "github.com/futugyou/infr-project/infrastructure"
	platform "github.com/futugyou/infr-project/platform"
	platformProvider "github.com/futugyou/infr-project/platform_provider"
	models "github.com/futugyou/infr-project/view_models"
)

type PlatformService struct {
	innerService   *AppService
	repository     platform.IPlatformRepositoryAsync
	vaultService   *VaultService
	client         *redis.Client
	eventPublisher infra.IEventPublisher
	screenshot     infra.IScreenshot
}

func NewPlatformService(
	unitOfWork domain.IUnitOfWork,
	repository platform.IPlatformRepositoryAsync,
	vaultService *VaultService,
	client *redis.Client,
	eventPublisher infra.IEventPublisher,
	screenshot infra.IScreenshot,
) *PlatformService {
	return &PlatformService{
		innerService:   NewAppService(unitOfWork),
		repository:     repository,
		vaultService:   vaultService,
		client:         client,
		eventPublisher: eventPublisher,
		screenshot:     screenshot,
	}
}

func (s *PlatformService) CreatePlatform(ctx context.Context, aux models.CreatePlatformRequest) (*models.PlatformDetailView, error) {
	properties := s.convertToPlatformProperties(aux.Properties)
	secrets, err := s.convertToPlatformSecrets(ctx, aux.Secrets)
	if err != nil {
		return nil, err
	}

	res, err := platform.NewPlatform(
		aux.Name,
		aux.Url,
		platform.GetPlatformProvider(aux.Provider),
		platform.WithPlatformProperties(properties),
		platform.WithPlatformTags(aux.Tags),
		platform.WithPlatformSecrets(secrets),
	)

	if err != nil {
		return nil, err
	}

	// check name
	resCh, errCh := s.repository.GetPlatformByNameAsync(ctx, aux.Name)
	resdb, err := tool.HandleAsync(ctx, resCh, errCh)
	if resdb != nil {
		return nil, fmt.Errorf("name: %s is existed", aux.Name)
	}

	if err != nil && !strings.HasPrefix(err.Error(), extensions.Data_Not_Found_Message) {
		return nil, err
	}

	if res.Provider != platform.PlatformProviderOther {
		status := s.determineProviderStatus(ctx, res)
		if status {
			res.Enable()
		} else {
			res.Disable()
		}
	}

	if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		errCh := s.repository.InsertAsync(ctx, *res)
		return tool.HandleErrorAsync(ctx, errCh)
	}); err != nil {
		return nil, err
	}

	return s.convertPlatformEntityToViewModel(ctx, res)
}

func (s *PlatformService) SearchPlatforms(ctx context.Context, request models.SearchPlatformsRequest) ([]models.PlatformView, error) {
	filter := platform.PlatformSearch{
		Name:      request.Name,
		NameFuzzy: false,
		Activate:  request.Activate,
		Tags:      request.Tags,
		Page:      request.Page,
		Size:      request.Size,
	}

	srcCh, errCh := s.repository.SearchPlatformsAsync(ctx, filter)
	src, err := tool.HandleAsync(ctx, srcCh, errCh)

	if err != nil && !strings.HasPrefix(err.Error(), extensions.Data_Not_Found_Message) {
		return nil, err
	}

	return s.convertToPlatformViews(src), nil
}

func (s *PlatformService) GetPlatform(ctx context.Context, idOrName string) (*models.PlatformDetailView, error) {
	srcCh, errCh := s.repository.GetPlatformByIdOrNameAsync(ctx, idOrName)
	if src, err := tool.HandleAsync(ctx, srcCh, errCh); err != nil {
		return nil, err
	} else {
		return s.convertPlatformEntityToViewModel(ctx, src)
	}
}

func (s *PlatformService) GetProviderProjectList(ctx context.Context, idOrName string) ([]models.PlatformProviderProject, error) {
	srcCh, errCh := s.repository.GetPlatformByIdOrNameAsync(ctx, idOrName)
	src, err := tool.HandleAsync(ctx, srcCh, errCh)
	if err != nil {
		return nil, err
	}

	if src == nil {
		return nil, fmt.Errorf("no platform data found")
	}

	result := []models.PlatformProviderProject{}
	if provider, err := s.getPlatformProvider(ctx, *src); err == nil {
		projects, _ := s.getProviderProjects(ctx, provider, *src)
		for _, pro := range projects {
			project := s.convertProviderProjectToModel(pro)
			result = append(result, project)
		}
	}

	return result, nil
}

func (s *PlatformService) UpdatePlatform(ctx context.Context, idOrName string, data models.UpdatePlatformRequest) (*models.PlatformDetailView, error) {
	return s.updatePlatform(ctx, idOrName, "UpdatePlatform", func(plat *platform.Platform) error {
		if plat.IsDeleted {
			return fmt.Errorf("id: %s was already deleted", plat.Id)
		}

		if plat.Name != data.Name {
			resCh, errCh := s.repository.GetPlatformByNameAsync(ctx, data.Name)
			var res *platform.Platform
			select {
			case res = <-resCh:
				if res.Id != plat.Id {
					return fmt.Errorf("name: %s is existed", data.Name)
				}
			case err := <-errCh:
				if !strings.HasPrefix(err.Error(), extensions.Data_Not_Found_Message) {
					return err
				}
			case <-ctx.Done():
				return fmt.Errorf("UpdatePlatform timeout: %w", ctx.Err())
			}

			if _, err := plat.UpdateName(data.Name); err != nil {
				return err
			}
		}

		if _, err := plat.UpdateUrl(data.Url); err != nil {
			return err
		}

		if _, err := plat.UpdateTags(data.Tags); err != nil {
			return err
		}

		if _, err := plat.UpdateProvider(platform.GetPlatformProvider(data.Provider)); err != nil {
			return err
		}

		if plat.Provider != platform.PlatformProviderOther {
			status := s.determineProviderStatus(ctx, plat)
			if status {
				plat.Enable()
			} else {
				plat.Disable()
			}
		}

		newProperty := s.convertToPlatformProperties(data.Properties)
		if _, err := plat.UpdateProperties(newProperty); err != nil {
			return err
		}

		newSecrets, err := s.convertToPlatformSecrets(ctx, data.Secrets)
		if err != nil {
			return err
		}

		if _, err := plat.UpdateSecrets(newSecrets); err != nil {
			return err
		}

		return nil
	})
}

func (s *PlatformService) DeletePlatform(ctx context.Context, idOrName string) (*models.PlatformDetailView, error) {
	return s.updatePlatform(ctx, idOrName, "DeletePlatform", func(plat *platform.Platform) error {
		if _, err := plat.Delete(); err != nil {
			return err
		}

		return nil
	})
}

func (s *PlatformService) RecoveryPlatform(ctx context.Context, idOrName string) (*models.PlatformDetailView, error) {
	return s.updatePlatform(ctx, idOrName, "RecoveryPlatform", func(plat *platform.Platform) error {
		if _, err := plat.Recovery(); err != nil {
			return err
		}

		return nil
	})
}

func (s *PlatformService) updatePlatform(ctx context.Context, idOrName string, _ string, fn func(*platform.Platform) error) (*models.PlatformDetailView, error) {
	platCh, errCh := s.repository.GetPlatformByIdOrNameAsync(ctx, idOrName)
	plat, err := tool.HandleAsync(ctx, platCh, errCh)
	if err != nil {
		return nil, err
	}

	if err := fn(plat); err != nil {
		return nil, err
	}

	if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		errCh := s.repository.UpdateAsync(ctx, *plat)
		return tool.HandleErrorAsync(ctx, errCh)
	}); err != nil {
		return nil, err
	}

	return s.convertPlatformEntityToViewModel(ctx, plat)
}

// this method is considered deprecated. we should create a webhook through the project callback.
func (s *PlatformService) UpsertWebhook(ctx context.Context, idOrName string, projectId string, hook models.UpdatePlatformWebhookRequest) (*models.PlatformDetailView, error) {
	platCh, errCh := s.repository.GetPlatformByIdOrNameAsync(ctx, idOrName)
	plat, err := tool.HandleAsync(ctx, platCh, errCh)
	if err != nil {
		return nil, err
	}

	project, exists := plat.Projects[projectId]
	if !exists {
		return nil, fmt.Errorf("projectId: %s is not existed in %s", projectId, idOrName)
	}

	webhook, _ := project.GetWebhook(hook.Name)
	needInsert := webhook == nil && hook.Sync

	properties := s.convertToPlatformProperties(hook.Properties)

	secrets, err := s.convertToPlatformSecrets(ctx, hook.Secrets)
	if err != nil {
		return nil, err
	}

	newhook := platform.NewWebhook(hook.Name, platform.GetWebhookUrl(plat.Name, project.Name),
		platform.WithWebhookProperties(properties),
		platform.WithWebhookActivate(hook.Activate),
		platform.WithWebhookState(platform.GetWebhookState(hook.State)),
		platform.WithWebhookSecrets(secrets),
	)

	if plat, err = plat.UpdateWebhook(projectId, *newhook); err != nil {
		return nil, err
	}

	if needInsert {
		// Regardless of whether sync is successful, the program will continue
		if providerHook, err := s.handlingProviderWebhookCreation(ctx, plat, project, hook.Name); err != nil {
			log.Println(err.Error())
		} else {
			hookSecrets, hookProperties := s.createWebhookVault(ctx, providerHook.GetParameters(), plat.Id, project.ProviderProjectId, hook.Name)
			newhook.UpdateProperties(hookProperties)
			newhook.UpdateSecrets(hookSecrets)
			newhook.UpdateProviderHookId(providerHook.ID)
			plat.UpdateWebhook(project.Id, *newhook)
		}
	}

	if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		errCh := s.repository.UpdateAsync(ctx, *plat)
		return tool.HandleErrorAsync(ctx, errCh)
	}); err != nil {
		return nil, err
	}

	return s.convertPlatformEntityToViewModel(ctx, plat)
}

func (s *PlatformService) createWebhookVault(ctx context.Context, parameters map[string]string,
	platformId string, providerProjectId string, hookName string) (map[string]platform.Secret, map[string]platform.Property) {
	properties := map[string]platform.Property{}
	for k, v := range parameters {
		properties[k] = platform.Property{
			Key:   k,
			Value: v,
		}
	}

	secrets := map[string]platform.Secret{}
	signingSecret := parameters["SigningSecret"]
	if len(signingSecret) == 0 {
		return secrets, properties
	}

	delete(properties, "SigningSecret")

	aux := models.CreateVaultRequest{
		CreateVaultModel: models.CreateVaultModel{
			Key:          "WebHookSecret",
			Value:        signingSecret,
			StorageMedia: "Local",
			VaultType:    "project",
			TypeIdentity: fmt.Sprintf("%s/%s/%s", platformId, providerProjectId, hookName),
		},
		ForceInsert: true,
	}

	vaultReps, err := s.vaultService.CreateVault(ctx, aux)
	if err != nil {
		log.Println(err.Error())
		return secrets, properties
	}

	secrets["WebHookSecret"] = platform.Secret{
		Key:            "WebHookSecret",
		Value:          signingSecret,
		VaultKey:       vaultReps.Key,
		VaultMaskValue: vaultReps.MaskValue,
	}

	return secrets, properties
}

func (s *PlatformService) handlingProviderWebhookCreation(ctx context.Context, plat *platform.Platform, project platform.PlatformProject, webhookName string) (*platformProvider.WebHook, error) {
	provider, err := s.getPlatformProvider(ctx, *plat)
	if err != nil {
		return nil, err
	}

	secret, _ := tool.GenerateRandomKey(6)
	platformId := ""
	if plat.Provider == platform.PlatformProviderGithub {
		if prop, ok := plat.Properties["GITHUB_OWNER"]; ok {
			platformId = prop.Value
		}
	}

	return s.createProviderWebhook(
		ctx,
		provider,
		platformId,
		project.ProviderProjectId,
		platform.GetWebhookUrl(plat.Name, project.Name),
		webhookName,
		secret,
	)
}

func (s *PlatformService) RemoveWebhook(ctx context.Context, request models.RemoveWebhookRequest) (*models.PlatformDetailView, error) {
	platCh, errCh := s.repository.GetPlatformByIdOrNameAsync(ctx, request.PlatformId)
	plat, err := tool.HandleAsync(ctx, platCh, errCh)
	if err != nil {
		return nil, err
	}

	project, exists := plat.Projects[request.ProjectId]
	if !exists {
		return nil, fmt.Errorf("projectId: %s is not existed in %s", request.ProjectId, request.PlatformId)
	}

	hook, err := project.GetWebhook(request.HookName)
	if err != nil {
		return nil, err
	}

	if plat, err = plat.RemoveWebhook(request.ProjectId, request.HookName); err != nil {
		return nil, err
	}

	if request.Sync {
		if provider, err := s.getPlatformProvider(ctx, *plat); err == nil {
			parameters := mergePropertiesToMap(plat.Properties, project.Properties)
			if err = s.deleteProviderWebhook(ctx, provider, hook.ID, parameters); err != nil {
				log.Println(err.Error())
			}
		} else {
			log.Println(err.Error())
		}
	}

	if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		errCh := s.repository.UpdateAsync(ctx, *plat)
		return tool.HandleErrorAsync(ctx, errCh)
	}); err != nil {
		return nil, err
	}

	return s.convertPlatformEntityToViewModel(ctx, plat)
}

func (s *PlatformService) UpsertProject(ctx context.Context, idOrName string, projectId string, project models.UpdatePlatformProjectRequest) (*models.PlatformDetailView, error) {
	platCh, errCh := s.repository.GetPlatformByIdOrNameAsync(ctx, idOrName)
	plat, err := tool.HandleAsync(ctx, platCh, errCh)
	if err != nil {
		return nil, err
	}

	if len(projectId) == 0 {
		projectId = tool.Sanitize2String(strings.ToLower(project.Name), "_")
	}

	properties := s.convertToPlatformProperties(project.Properties)

	secrets, err := s.convertToPlatformSecrets(ctx, project.Secrets)
	if err != nil {
		return nil, err
	}

	var projectDb *platform.PlatformProject
	screenshot := false
	if proj, ok := plat.Projects[projectId]; ok {
		projectDb = &proj
		if len(project.Url) > 0 && (projectDb.Url != project.Url || len(projectDb.ImageUrl) == 0) {
			screenshot = true
		}

		projectDb.UpdateName(project.Name)
		projectDb.UpdateDescription(project.Description)
		projectDb.UpdateProperties(properties)
		projectDb.UpdateUrl(project.Url)
		projectDb.UpdateSecrets(secrets)
		projectDb.UpdateProviderProjectId(project.ProviderProjectId)
		projectDb.UpdateTags(project.Tags)
	} else {
		projectDb = platform.NewPlatformProject(
			projectId,
			project.Name,
			project.Url,
			platform.WithProjectProperties(properties),
			platform.WithProjectSecrets(secrets),
			platform.WithProjectDescription(project.Description),
			platform.WithProviderProjectId(project.ProviderProjectId),
			platform.WithProjectTags(project.Tags),
		)

		if len(projectDb.Url) > 0 {
			screenshot = true
		}
	}

	if _, err = plat.UpdateProject(*projectDb); err != nil {
		return nil, err
	}

	if project.Operate == "sync" || project.ImportWebhooks || screenshot {
		event := models.PlatformProjectUpsertEvent{
			PlatformId:            plat.Id,
			ProjectId:             projectId,
			CreateProviderProject: project.Operate == "sync",
			ImportWebhooks:        project.ImportWebhooks,
			EventName:             "upsert_project",
			Screenshot:            screenshot,
		}
		if err := s.eventPublisher.PublishCommon(ctx, event, "upsert_project"); err != nil {
			log.Println(err.Error())
		}
	}

	if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		errCh := s.repository.UpdateAsync(ctx, *plat)
		return tool.HandleErrorAsync(ctx, errCh)
	}); err != nil {
		return nil, err
	}

	return s.convertPlatformEntityToViewModel(ctx, plat)
}

func (s *PlatformService) DeleteProject(ctx context.Context, idOrName string, projectId string) (*models.PlatformDetailView, error) {
	platCh, errCh := s.repository.GetPlatformByIdOrNameAsync(ctx, idOrName)
	plat, err := tool.HandleAsync(ctx, platCh, errCh)
	if err != nil {
		return nil, err
	}

	if _, err := plat.RemoveProject(projectId); err != nil {
		return nil, err
	}

	if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		errCh := s.repository.UpdateAsync(ctx, *plat)
		return tool.HandleErrorAsync(ctx, errCh)
	}); err != nil {
		return nil, err
	}

	return s.convertPlatformEntityToViewModel(ctx, plat)
}

func (s *PlatformService) GetPlatformProject(ctx context.Context, platformIdOrName string, projectId string) (*models.PlatformProject, error) {
	srcCh, errCh := s.repository.GetPlatformByIdOrNameAsync(ctx, platformIdOrName)
	src, err := tool.HandleAsync(ctx, srcCh, errCh)
	if err != nil {
		return nil, err
	}

	if project, ok := src.Projects[projectId]; ok {
		providerProject := &platformProvider.Project{}
		if provider, err := s.getPlatformProvider(ctx, *src); err != nil {
			log.Println(err.Error())
		} else {
			providerProject = s.getProviderProjectWithCache(ctx, *src, project, provider)
		}

		modelProject := s.mergeProject(providerProject, &project)
		return &modelProject, nil
	} else {
		return nil, fmt.Errorf("can not find project with id: %s", projectId)
	}
}

func (s *PlatformService) getProviderProjectWithCache(ctx context.Context, src platform.Platform, project platform.PlatformProject, provider platformProvider.IPlatformProviderAsync) *platformProvider.Project {
	providerProject := &platformProvider.Project{}
	redisKey := fmt.Sprintf("platform_%s_project_%s", src.Id, project.Id)
	if data, err := s.client.Get(ctx, redisKey).Result(); err != nil {
		log.Println(err.Error())
	} else {
		if err = json.Unmarshal([]byte(data), providerProject); err != nil {
			log.Println(err.Error())
		}
	}

	if len(providerProject.ID) > 0 {
		return providerProject
	}

	parameters := mergePropertiesToMap(project.Properties, src.Properties)
	if project, err := s.getProviderProject(ctx, provider, project.ProviderProjectId, parameters); err != nil {
		log.Println(err.Error())
	} else {
		providerProject = project
	}

	if data, err := json.Marshal(providerProject); err != nil {
		log.Println(err.Error())
	} else {
		// 10Minute maybe configurable
		if _, err := s.client.Set(ctx, redisKey, string(data), time.Minute*10).Result(); err != nil {
			log.Println(err.Error())
		}
	}

	return providerProject
}

func (s *PlatformService) HandlePlatformProjectUpsert(ctx context.Context, event models.PlatformProjectUpsertEvent) error {
	platCh, errCh := s.repository.GetPlatformByIdOrNameAsync(ctx, event.PlatformId)
	plat, err := tool.HandleAsync(ctx, platCh, errCh)
	if err != nil {
		return err
	}

	var project *platform.PlatformProject
	if proj, ok := plat.Projects[event.ProjectId]; ok {
		project = &proj
	} else {
		return fmt.Errorf("can not find project with id: %s", event.ProjectId)
	}

	if event.CreateProviderProject || event.ImportWebhooks {
		if provider, err := s.getPlatformProvider(ctx, *plat); err != nil {
			log.Println(err.Error())
		} else {
			// 1. providerProject
			providerProject, err := s.createProviderProjectIfNeeded(ctx, event, *plat, *project, provider)
			if err != nil {
				log.Println(err.Error())
			}

			if providerProject != nil && len(providerProject.ID) > 0 {
				project.UpdateProviderProjectId(providerProject.ID)
			}

			// 2. webhook
			if err := s.handleWebhooks(ctx, event, providerProject, plat, project); err != nil {
				log.Println(err.Error())
			}
		}
	}

	if event.Screenshot && len(project.Url) > 0 && os.Getenv("SCREENSHOT_ALLOW") != "false" {
		if imageUrl, err := s.screenshot.Create(ctx, project.Url); err != nil {
			log.Println(err.Error())
		} else {
			project.UpdateImageUrl(*imageUrl)
		}
	}

	plat.UpdateProject(*project)

	return s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		errCh := s.repository.UpdateAsync(ctx, *plat)
		return tool.HandleErrorAsync(ctx, errCh)
	})
}

func (s *PlatformService) setupWebhookWithSecrets(ctx context.Context, providerHook platformProvider.WebHook, plat *platform.Platform, project *platform.PlatformProject) {
	newhook := platform.NewWebhook(providerHook.Name, platform.GetWebhookUrl(plat.Name, project.Name))
	hookSecrets, hookProperties := s.createWebhookVault(ctx, providerHook.GetParameters(), plat.Id, project.ProviderProjectId, newhook.Name)
	newhook.UpdateProperties(hookProperties)
	newhook.UpdateSecrets(hookSecrets)
	newhook.UpdateProviderHookId(providerHook.ID)
	project.UpsertWebhook(*newhook)
}

func (s *PlatformService) createProviderProjectIfNeeded(ctx context.Context, event models.PlatformProjectUpsertEvent,
	plat platform.Platform, project platform.PlatformProject, provider platformProvider.IPlatformProviderAsync) (*platformProvider.Project, error) {

	providerProject := s.getProviderProjectWithCache(ctx, plat, project, provider)
	if providerProject == nil || len(providerProject.ID) == 0 {
		if event.CreateProviderProject {
			var err error
			if providerProject, err = s.createProviderProject(ctx, provider, project.Name, project.Properties); err != nil {
				return nil, err
			}
		}

		if providerProject == nil || len(providerProject.ID) == 0 {
			log.Printf("no corresponding provider information with project: %s", event.ProjectId)
		}
	}

	return providerProject, nil
}

func (s *PlatformService) handleWebhooks(ctx context.Context, event models.PlatformProjectUpsertEvent,
	providerProject *platformProvider.Project, plat *platform.Platform, project *platform.PlatformProject) error {

	if providerProject == nil || len(providerProject.ID) == 0 && !event.ImportWebhooks {
		return nil
	}

	project.ClearWebhooks()
	if len(providerProject.WebHooks) == 0 {
		providerHook, err := s.handlingProviderWebhookCreation(ctx, plat, *project, "infr-project-webhook")
		if err != nil {
			return err
		}

		if providerHook != nil {
			s.setupWebhookWithSecrets(ctx, *providerHook, plat, project)
		}
	} else {
		for _, providerHook := range providerProject.WebHooks {
			if !strings.HasPrefix(providerHook.Url, os.Getenv("PROJECT_URL")) {
				continue
			}

			s.setupWebhookWithSecrets(ctx, providerHook, plat, project)
		}
	}

	return nil
}
