package application

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"log"

	tool "github.com/futugyou/extensions"

	"github.com/redis/go-redis/v9"

	domain "github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/extensions"
	platform "github.com/futugyou/infr-project/platform"
	platformProvider "github.com/futugyou/infr-project/platform_provider"
	models "github.com/futugyou/infr-project/view_models"
)

type PlatformService struct {
	innerService *AppService
	repository   platform.IPlatformRepositoryAsync
	vaultService *VaultService
	client       *redis.Client
}

func NewPlatformService(
	unitOfWork domain.IUnitOfWork,
	repository platform.IPlatformRepositoryAsync,
	vaultService *VaultService,
	client *redis.Client,
) *PlatformService {
	return &PlatformService{
		innerService: NewAppService(unitOfWork),
		repository:   repository,
		vaultService: vaultService,
		client:       client,
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

func (s *PlatformService) UpdatePlatform(ctx context.Context, idOrName string, data models.UpdatePlatformRequest) (*models.PlatformDetailView, error) {
	return s.updatePlatform(ctx, idOrName, "UpdatePlatform", func(plat *platform.Platform) error {
		if plat.IsDeleted {
			return fmt.Errorf("id: %s was alrealdy deleted", plat.Id)
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

	newhook := platform.NewWebhook(hook.Name, hook.Url,
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
		if provider, err := s.getPlatfromProvider(ctx, *plat); err == nil {
			platformId := ""
			if plat.Provider == platform.PlatformProviderGithub {
				if prop, ok := plat.Properties["GITHUB_OWNER"]; ok {
					platformId = prop.Value
				}
			}

			if providerHook, err := s.createProviderWebhook(ctx, provider, platformId,
				project.ProviderProjectId, hook.Url, hook.Name); err == nil {
				newhook.UpdateProviderHookId(providerHook.ID)
				plat.UpdateWebhook(projectId, *newhook)
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
		if provider, err := s.getPlatfromProvider(ctx, *plat); err == nil {
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

	proj := platform.NewPlatformProject(
		projectId,
		project.Name,
		project.Url,
		platform.WithProjectProperties(properties),
		platform.WithProjectSecrets(secrets),
		platform.WithProjectDescription(project.Description),
	)

	if _, err = plat.UpdateProject(*proj); err != nil {
		return nil, err
	}

	providerProjectId := project.ProviderProjectId
	// Regardless of whether sync is successful, the program will continue
	if project.Operate == "sync" {
		if provider, err := s.getPlatfromProvider(ctx, *plat); err == nil {
			shouldCreate := len(providerProjectId) == 0
			if !shouldCreate {
				projects, _ := s.getProviderProjects(ctx, provider, *plat)
				shouldCreate = true

				for _, v := range projects {
					if v.ID == providerProjectId {
						shouldCreate = false
						break
					}
				}
			}

			if shouldCreate {
				if p, err := s.createProviderProject(ctx, provider, proj.Name, proj.Properties); err == nil {
					providerProjectId = p.ID
				}
			}
		} else {
			log.Println(err.Error())
		}
	}

	proj.UpdateProviderProjectId(providerProjectId)
	// The status of PlatformProject has been checked before.
	plat.UpdateProject(*proj)

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

func (s *PlatformService) GetPlatformProject(ctx context.Context, platfromIdOrName string, projectId string) (*models.PlatformProject, error) {
	srcCh, errCh := s.repository.GetPlatformByIdOrNameAsync(ctx, platfromIdOrName)
	src, err := tool.HandleAsync(ctx, srcCh, errCh)
	if err != nil {
		return nil, err
	}

	if project, ok := src.Projects[projectId]; ok {
		providerProject := &platformProvider.Project{}
		if provider, err := s.getPlatfromProvider(ctx, *src); err != nil {
			log.Println(err.Error())
		} else {
			// TODO: add limit or cache when get detail provider
			redisKey := fmt.Sprintf("platform_%s_project_%s", src.Id, projectId)
			if p, err := s.client.Get(ctx, redisKey).Result(); err != nil {
				log.Println(err.Error())
			} else {
				if err = json.Unmarshal([]byte(p), providerProject); err != nil {
					log.Println(err.Error())
				}
			}

			if len(providerProject.ID) == 0 {
				parameters := mergePropertiesToMap(project.Properties, src.Properties)
				if project, err := s.getProviderProject(ctx, provider, project.ProviderProjectId, parameters); err != nil {
					log.Println(err.Error())
				} else {
					providerProject = project
				}
				if data, err := json.Marshal(providerProject); err != nil {
					log.Println(err.Error())
				} else {

					if re, err := s.client.Set(ctx, redisKey, string(data), time.Minute*10).Result(); err != nil {
						log.Println(err.Error())
					} else {
						log.Println(re)
					}
				}
			}
		}

		modelProject := s.mergeProject(providerProject, &project)
		return &modelProject, nil
	} else {
		return nil, fmt.Errorf("can not find project with id: %s", projectId)
	}
}
