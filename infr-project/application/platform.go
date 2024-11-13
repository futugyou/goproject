package application

import (
	"context"
	"fmt"
	"strings"

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
}

func NewPlatformService(
	unitOfWork domain.IUnitOfWork,
	repository platform.IPlatformRepositoryAsync,
	vaultService *VaultService,
) *PlatformService {
	return &PlatformService{
		innerService: NewAppService(unitOfWork),
		repository:   repository,
		vaultService: vaultService,
	}
}

func (s *PlatformService) CreatePlatform(ctx context.Context, aux models.CreatePlatformRequest) (*models.PlatformDetailView, error) {
	properties := s.convertPlatformProperties(aux.Properties)
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
	select {
	case <-resCh:
		return nil, fmt.Errorf("name: %s is existed", aux.Name)
	case err := <-errCh:
		if !strings.HasPrefix(err.Error(), extensions.Data_Not_Found_Message) {
			return nil, err
		}
	case <-ctx.Done():
		return nil, fmt.Errorf("CreatePlatform timeout: %w", ctx.Err())
	}

	if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		errCh := s.repository.InsertAsync(ctx, *res)
		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return fmt.Errorf("CreatePlatform timeout: %w", ctx.Err())
		}
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
	var src []platform.Platform
	select {
	case src = <-srcCh:
	case err := <-errCh:
		if !strings.HasPrefix(err.Error(), extensions.Data_Not_Found_Message) {
			return nil, err
		}
	case <-ctx.Done():
		return nil, fmt.Errorf("SearchPlatforms timeout: %w", ctx.Err())
	}

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
	return result, nil
}

func (s *PlatformService) GetPlatform(ctx context.Context, id string) (*models.PlatformDetailView, error) {
	srcCh, errCh := s.repository.GetAsync(ctx, id)
	select {
	case src := <-srcCh:
		return s.convertPlatformEntityToViewModel(ctx, src)
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("GetPlatform timeout: %w", ctx.Err())
	}
}

func (s *PlatformService) UpdatePlatform(ctx context.Context, id string, data models.UpdatePlatformRequest) (*models.PlatformDetailView, error) {
	return s.updatePlatform(ctx, id, "UpdatePlatform", func(plat *platform.Platform) error {
		if plat.IsDeleted {
			return fmt.Errorf("id: %s was alrealdy deleted", plat.Id)
		}

		if plat.Name != data.Name {
			resCh, errCh := s.repository.GetPlatformByNameAsync(ctx, data.Name)
			var res *platform.Platform
			select {
			case res = <-resCh:
				if res.Id != id {
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

		if data.Activate {
			if _, err := plat.Enable(); err != nil {
				return err
			}
		} else {
			if _, err := plat.Disable(); err != nil {
				return err
			}
		}

		if _, err := plat.UpdateProvider(platform.GetPlatformProvider(data.Provider)); err != nil {
			return err
		}

		newProperty := s.convertPlatformProperties(data.Properties)
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

func (s *PlatformService) DeletePlatform(ctx context.Context, id string) (*models.PlatformDetailView, error) {
	return s.updatePlatform(ctx, id, "DeletePlatform", func(plat *platform.Platform) error {
		if _, err := plat.Delete(); err != nil {
			return err
		}

		return nil
	})
}

func (s *PlatformService) RecoveryPlatform(ctx context.Context, id string) (*models.PlatformDetailView, error) {
	return s.updatePlatform(ctx, id, "RecoveryPlatform", func(plat *platform.Platform) error {
		if _, err := plat.Recovery(); err != nil {
			return err
		}

		return nil
	})
}

func (s *PlatformService) updatePlatform(ctx context.Context, id string, operation string, fn func(*platform.Platform) error) (*models.PlatformDetailView, error) {
	platCh, errCh := s.repository.GetAsync(ctx, id)
	var plat *platform.Platform
	select {
	case plat = <-platCh:
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("%s timeout: %w", operation, ctx.Err())
	}

	if err := fn(plat); err != nil {
		return nil, err
	}

	if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		errCh := s.repository.UpdateAsync(ctx, *plat)
		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return fmt.Errorf("%s timeout: %w", operation, ctx.Err())
		}
	}); err != nil {
		return nil, err
	}

	return s.convertPlatformEntityToViewModel(ctx, plat)
}

func (s *PlatformService) UpsertWebhook(ctx context.Context, id string, projectId string, hook models.UpdatePlatformWebhookRequest) (*models.PlatformDetailView, error) {
	platCh, errCh := s.repository.GetAsync(ctx, id)
	var plat *platform.Platform
	var err error
	select {
	case plat = <-platCh:
	case err = <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("UpsertWebhook timeout: %w", ctx.Err())
	}

	if _, exists := plat.Projects[projectId]; !exists {
		return nil, fmt.Errorf("projectId: %s is not existed in %s", projectId, id)
	}

	properties := make(map[string]platform.Property)
	for key, value := range hook.Properties {
		properties[key] = platform.Property{
			Key:   key,
			Value: value,
		}
	}

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

	if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		errCh := s.repository.UpdateAsync(ctx, *plat)
		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return fmt.Errorf("UpsertWebhook timeout: %w", ctx.Err())
		}
	}); err != nil {
		return nil, err
	}

	return s.convertPlatformEntityToViewModel(ctx, plat)
}

func (s *PlatformService) RemoveWebhook(ctx context.Context, id string, projectId string, hookName string) (*models.PlatformDetailView, error) {
	platCh, errCh := s.repository.GetAsync(ctx, id)
	var plat *platform.Platform
	var err error
	select {
	case plat = <-platCh:
	case err = <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("RemoveWebhook timeout: %w", ctx.Err())
	}

	if _, exists := plat.Projects[projectId]; !exists {
		return nil, fmt.Errorf("projectId: %s is not existed in %s", projectId, id)
	}

	if plat, err = plat.RemoveWebhook(projectId, hookName); err != nil {
		return nil, err
	}

	if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		errCh := s.repository.UpdateAsync(ctx, *plat)
		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return fmt.Errorf("RemoveWebhook timeout: %w", ctx.Err())
		}
	}); err != nil {
		return nil, err
	}

	return s.convertPlatformEntityToViewModel(ctx, plat)
}

func (s *PlatformService) UpsertProject(ctx context.Context, id string, projectId string, project models.UpdatePlatformProjectRequest) (*models.PlatformDetailView, error) {
	platCh, errCh := s.repository.GetAsync(ctx, id)
	var plat *platform.Platform
	var err error
	select {
	case plat = <-platCh:
	case err = <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("UpsertProject timeout: %w", ctx.Err())
	}

	if len(projectId) == 0 {
		projectId = project.Name
	}

	properties := make(map[string]platform.Property)
	for key, value := range project.Properties {
		properties[key] = platform.Property{
			Key:   key,
			Value: value,
		}
	}

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
	)

	if _, err = plat.UpdateProject(*proj); err != nil {
		return nil, err
	}

	if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		errCh := s.repository.UpdateAsync(ctx, *plat)
		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return fmt.Errorf("UpsertProject timeout: %w", ctx.Err())
		}
	}); err != nil {
		return nil, err
	}

	return s.convertPlatformEntityToViewModel(ctx, plat)
}

func (s *PlatformService) DeleteProject(ctx context.Context, id string, projectId string) (*models.PlatformDetailView, error) {
	platCh, errCh := s.repository.GetAsync(ctx, id)
	var plat *platform.Platform
	var err error
	select {
	case plat = <-platCh:
	case err = <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("DeleteProject timeout: %w", ctx.Err())
	}

	if _, err := plat.RemoveProject(projectId); err != nil {
		return nil, err
	}
	if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		errCh := s.repository.UpdateAsync(ctx, *plat)
		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return fmt.Errorf("DeleteProject timeout: %w", ctx.Err())
		}
	}); err != nil {
		return nil, err
	}
	return s.convertPlatformEntityToViewModel(ctx, plat)
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

func (s *PlatformService) getProviderProjects(ctx context.Context, provider platformProvider.IPlatformProviderAsync) ([]models.PlatformProject, error) {
	var projects []platformProvider.Project

	filter := platformProvider.ProjectFilter{}
	resCh, errCh := provider.ListProjectAsync(ctx, filter)
	select {
	case projects = <-resCh:
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("getProviderProjects timeout: %w", ctx.Err())
	}

	result := make([]models.PlatformProject, 0)
	for _, project := range projects {
		result = append(result, models.PlatformProject{
			Id:         project.ID,
			Name:       project.Name,
			Url:        project.Url,
			Properties: []models.Property{},
			Secrets:    []models.Secret{},
			Webhooks:   []models.Webhook{},
			Followed:   false,
		})
	}

	return result, nil
}
