package application

import (
	"context"
	"fmt"
	"strings"

	tool "github.com/futugyou/extensions"

	domain "github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/extensions"
	platform "github.com/futugyou/infr-project/platform"
	models "github.com/futugyou/infr-project/view_models"
)

type PlatformService struct {
	innerService *AppService
	repository   platform.IPlatformRepositoryAsync
}

func NewPlatformService(
	unitOfWork domain.IUnitOfWork,
	repository platform.IPlatformRepositoryAsync,
) *PlatformService {
	return &PlatformService{
		innerService: NewAppService(unitOfWork),
		repository:   repository,
	}
}

func (s *PlatformService) CreatePlatform(aux models.CreatePlatformRequest, ctx context.Context) (*models.PlatformDetailView, error) {
	resCh, errCh := s.repository.GetPlatformByNameAsync(ctx, aux.Name)
	var res *platform.Platform
	select {
	case res = <-resCh:
	case err := <-errCh:
		if err != nil && !strings.HasPrefix(err.Error(), extensions.Data_Not_Found_Message) {
			return nil, err
		}
	case <-ctx.Done():
		return nil, fmt.Errorf("CreatePlatform timeout")
	}

	if res != nil && res.Name == aux.Name {
		return nil, fmt.Errorf("name: %s is existed", aux.Name)
	}

	property := make(map[string]platform.PropertyInfo)
	for _, v := range aux.Property {
		property[v.Key] = platform.PropertyInfo{
			Key:   v.Key,
			Value: v.Value,
		}
	}

	if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		res = platform.NewPlatform(aux.Name, aux.Url, aux.Rest, property, aux.Tags)
		return <-s.repository.InsertAsync(ctx, *res)
	}); err != nil {
		return nil, err
	}

	return convertPlatformEntityToViewModel(res)
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
		if err != nil && !strings.HasPrefix(err.Error(), extensions.Data_Not_Found_Message) {
			return nil, err
		}
	case <-ctx.Done():
		return nil, fmt.Errorf("SearchPlatforms timeout")
	}

	result := make([]models.PlatformView, len(src))
	for i := 0; i < len(src); i++ {
		result[i] = models.PlatformView{
			Id:           src[i].Id,
			Name:         src[i].Name,
			Activate:     src[i].Activate,
			Url:          src[i].Url,
			RestEndpoint: src[i].RestEndpoint,
			Tags:         src[i].Tags,
			IsDeleted:    src[i].IsDeleted,
		}
	}
	return result, nil
}

func (s *PlatformService) GetPlatform(id string, ctx context.Context) (*models.PlatformDetailView, error) {
	srcCh, errCh := s.repository.GetAsync(ctx, id)
	select {
	case src := <-srcCh:
		return convertPlatformEntityToViewModel(src)
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("GetPlatform timeout")
	}
}

func (s *PlatformService) UpsertWebhook(id string, projectId string, hook models.UpdatePlatformWebhookRequest, ctx context.Context) (*models.PlatformDetailView, error) {
	platCh, errCh := s.repository.GetAsync(ctx, id)
	var plat *platform.Platform
	var err error
	select {
	case plat = <-platCh:
	case err = <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("UpsertWebhook timeout")
	}

	if _, exists := plat.Projects[projectId]; !exists {
		return nil, fmt.Errorf("projectId: %s is not existed in %s", projectId, id)
	}

	newhook := platform.NewWebhook(hook.Name, hook.Url,
		platform.WithWebhookProperty(hook.Property),
		platform.WithWebhookActivate(hook.Activate),
		platform.WithWebhookState(platform.GetWebhookState(hook.State)),
	)
	if plat, err = plat.UpdateWebhook(projectId, *newhook); err != nil {
		return nil, err
	}

	if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		return <-s.repository.UpdateAsync(ctx, *plat)
	}); err != nil {
		return nil, err
	}

	return convertPlatformEntityToViewModel(plat)
}

func (s *PlatformService) RemoveWebhook(id string, projectId string, hookName string, ctx context.Context) (*models.PlatformDetailView, error) {
	platCh, errCh := s.repository.GetAsync(ctx, id)
	var plat *platform.Platform
	var err error
	select {
	case plat = <-platCh:
	case err = <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("RemoveWebhook timeout")
	}

	if _, exists := plat.Projects[projectId]; !exists {
		return nil, fmt.Errorf("projectId: %s is not existed in %s", projectId, id)
	}

	if plat, err = plat.RemoveWebhook(projectId, hookName); err != nil {
		return nil, err
	}

	if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		return <-s.repository.UpdateAsync(ctx, *plat)
	}); err != nil {
		return nil, err
	}

	return convertPlatformEntityToViewModel(plat)
}

func (s *PlatformService) DeletePlatform(id string, ctx context.Context) (*models.PlatformDetailView, error) {
	if err := <-s.repository.SoftDeleteAsync(ctx, id); err != nil {
		return nil, err
	}
	srcCh, errCh := s.repository.GetAsync(ctx, id)
	select {
	case src := <-srcCh:
		return convertPlatformEntityToViewModel(src)
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("DeletePlatform timeout")
	}
}

func (s *PlatformService) AddProject(id string, projectId string, project models.UpdatePlatformProjectRequest, ctx context.Context) (*models.PlatformDetailView, error) {
	platCh, errCh := s.repository.GetAsync(ctx, id)
	var plat *platform.Platform
	var err error
	select {
	case plat = <-platCh:
	case err = <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("AddProject timeout")
	}

	if len(projectId) == 0 {
		projectId = project.Name
	}

	proj := platform.NewPlatformProject(projectId, project.Name, project.Url, project.Property)
	if _, err = plat.UpdateProject(*proj); err != nil {
		return nil, err
	}

	if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		return <-s.repository.UpdateAsync(ctx, *plat)
	}); err != nil {
		return nil, err
	}

	return convertPlatformEntityToViewModel(plat)
}

func (s *PlatformService) DeleteProject(id string, projectId string, ctx context.Context) (*models.PlatformDetailView, error) {
	platCh, errCh := s.repository.GetAsync(ctx, id)
	var plat *platform.Platform
	var err error
	select {
	case plat = <-platCh:
	case err = <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("DeleteProject timeout")
	}

	if _, err := plat.RemoveProject(projectId); err != nil {
		return nil, err
	}
	if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		return <-s.repository.UpdateAsync(ctx, *plat)
	}); err != nil {
		return nil, err
	}
	return convertPlatformEntityToViewModel(plat)
}

func (s *PlatformService) UpdatePlatform(id string, data models.UpdatePlatformRequest, ctx context.Context) (*models.PlatformDetailView, error) {
	platCh, errCh := s.repository.GetAsync(ctx, id)
	var plat *platform.Platform
	var err error
	select {
	case plat = <-platCh:
	case err = <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("UpdatePlatform timeout")
	}

	if plat.Name != data.Name {
		resCh, errCh := s.repository.GetPlatformByNameAsync(ctx, data.Name)
		var res *platform.Platform
		select {
		case res = <-resCh:
		case err := <-errCh:
			if err != nil && !strings.HasPrefix(err.Error(), extensions.Data_Not_Found_Message) {
				return nil, err
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("UpdatePlatform timeout")
		}

		if res != nil && len(res.Id) > 0 && res.Id != id {
			return nil, fmt.Errorf("name: %s is existed", data.Name)
		}

		if _, err := plat.UpdateName(data.Name); err != nil {
			return nil, err
		}
	}

	if plat.Url != data.Url {
		if _, err := plat.UpdateUrl(data.Url); err != nil {
			return nil, err
		}
	}

	if plat.RestEndpoint != data.Rest {
		if _, err := plat.UpdateRestEndpoint(data.Rest); err != nil {
			return nil, err
		}
	}

	if !tool.StringArrayCompare(plat.Tags, data.Tags) {
		if _, err := plat.UpdateTags(data.Tags); err != nil {
			return nil, err
		}
	}

	if data.Activate != nil && plat.Activate != *data.Activate {
		if *data.Activate {
			if _, err := plat.Enable(); err != nil {
				return nil, err
			}
		} else {
			if _, err := plat.Disable(); err != nil {
				return nil, err
			}
		}
	}

	newProperty := make(map[string]platform.PropertyInfo)
	for _, v := range data.Property {
		newProperty[v.Key] = platform.PropertyInfo{
			Key:   v.Key,
			Value: v.Value,
		}
	}
	if !tool.MapsCompareCommon(plat.Property, newProperty) {
		if _, err := plat.UpdateProperty(newProperty); err != nil {
			return nil, err
		}
	}

	if err = s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		return <-s.repository.UpdateAsync(ctx, *plat)
	}); err != nil {
		return nil, err
	}

	return convertPlatformEntityToViewModel(plat)
}

func convertPlatformEntityToViewModel(src *platform.Platform) (*models.PlatformDetailView, error) {
	if src == nil {
		return nil, nil
	}

	propertyInfos, err := convertProperty(src.Property)
	if err != nil {
		return nil, nil
	}

	platformProjects := make([]models.PlatformProject, 0)
	for _, v := range src.Projects {
		webhooks := make([]models.Webhook, 0)
		for i := 0; i < len(v.Webhooks); i++ {
			webhooks = append(webhooks, models.Webhook{
				Name:     v.Webhooks[i].Name,
				Url:      v.Webhooks[i].Url,
				Activate: v.Webhooks[i].Activate,
				State:    v.Webhooks[i].State.String(),
				Property: v.Webhooks[i].Property,
			})
		}
		platformProjects = append(platformProjects, models.PlatformProject{
			Id:       v.Id,
			Name:     v.Name,
			Url:      v.Url,
			Property: v.Property,
			Webhooks: webhooks,
		})
	}
	return &models.PlatformDetailView{
		Id:           src.Id,
		Name:         src.Name,
		Activate:     src.Activate,
		Url:          src.Url,
		RestEndpoint: src.RestEndpoint,
		Property:     propertyInfos,
		Projects:     platformProjects,
		Tags:         src.Tags,
		IsDeleted:    src.IsDeleted,
	}, nil
}

func convertProperty(properties map[string]platform.PropertyInfo) ([]models.PropertyInfo, error) {
	propertyInfos := make([]models.PropertyInfo, 0)
	for _, v := range properties {
		propertyInfos = append(propertyInfos, models.PropertyInfo{
			Key:   v.Key,
			Value: v.Value,
		})
	}
	return propertyInfos, nil
}
