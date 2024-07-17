package application

import (
	"context"
	"fmt"
	"os"
	"strings"

	tool "github.com/futugyou/extensions"

	domain "github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/extensions"
	platform "github.com/futugyou/infr-project/platform"
	models "github.com/futugyou/infr-project/view_models"
)

type PlatformService struct {
	innerService *AppService
	repository   platform.IPlatformRepository
}

func NewPlatformService(
	unitOfWork domain.IUnitOfWork,
	repository platform.IPlatformRepository,
) *PlatformService {
	return &PlatformService{
		innerService: NewAppService(unitOfWork),
		repository:   repository,
	}
}

func (s *PlatformService) CreatePlatform(aux models.CreatePlatformRequest, ctx context.Context) (*models.PlatformView, error) {
	var res *platform.Platform
	res, err := s.repository.GetPlatformByName(ctx, aux.Name)
	if err != nil && !strings.HasPrefix(err.Error(), extensions.Data_Not_Found_Message) {
		return nil, err
	}

	if res != nil && res.Name == aux.Name {
		return nil, fmt.Errorf("name: %s is existed", aux.Name)
	}

	property := make(map[string]platform.PropertyInfo)
	for _, v := range aux.Property {
		value, err := tool.AesCTREncrypt(v.Value, os.Getenv("Encrypt_Key"))
		if err != nil {
			return nil, err
		}
		property[v.Key] = platform.PropertyInfo{
			Key:      v.Key,
			Value:    value,
			NeedMask: v.NeedMask,
		}
	}
	err = s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		res = platform.NewPlatform(aux.Name, aux.Url, aux.Rest, property, aux.Tags)
		return s.repository.Insert(ctx, *res)
	})
	if err != nil {
		return nil, err
	}
	property = make(map[string]platform.PropertyInfo)
	for _, v := range aux.Property {
		property[v.Key] = platform.PropertyInfo(v)
	}
	res.UpdateProperty(property)
	return convertPlatformEntityToViewModel(res), nil
}

func (s *PlatformService) GetAllPlatform(ctx context.Context) ([]models.PlatformView, error) {
	src, err := s.repository.GetAllPlatform(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]models.PlatformView, len(src))
	for i := 0; i < len(src); i++ {
		result = append(result, *convertPlatformEntityToViewModel(&src[i]))
	}
	return result, nil
}

func (s *PlatformService) GetPlatform(id string, ctx context.Context) (*models.PlatformView, error) {
	src, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return convertPlatformEntityToViewModel(src), nil
}

func (s *PlatformService) AddWebhook(id string, projectId string, hook models.UpdatePlatformWebhookRequest, ctx context.Context) (*models.PlatformView, error) {
	plat, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, exists := plat.Projects[projectId]; !exists {
		return nil, fmt.Errorf("projectId: %s is not existed in %s", projectId, id)
	}

	newhook := platform.NewWebhook(hook.Name, hook.Url, hook.Property)
	newhook.Activate = hook.Activate
	newhook.State = platform.GetWebhookState(hook.State)
	if plat, err = plat.UpdateWebhook(projectId, *newhook); err != nil {
		return nil, err
	}
	if err = s.repository.Update(ctx, *plat); err != nil {
		return nil, err
	}

	return convertPlatformEntityToViewModel(plat), nil
}

func (s *PlatformService) DeletePlatform(id string, ctx context.Context) (*models.PlatformView, error) {
	if err := s.repository.SoftDelete(ctx, id); err != nil {
		return nil, err
	}
	plat, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return convertPlatformEntityToViewModel(plat), nil
}

func (s *PlatformService) AddProject(id string, projectId string, project models.UpdatePlatformProjectRequest, ctx context.Context) (*models.PlatformView, error) {
	plat, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if len(projectId) == 0 {
		projectId = project.Name
	}

	proj := platform.NewPlatformProject(projectId, project.Name, project.Url, project.Property)
	if _, err = plat.UpdateProject(*proj); err != nil {
		return nil, err
	}

	if err = s.repository.Update(ctx, *plat); err != nil {
		return nil, err
	}

	return convertPlatformEntityToViewModel(plat), nil
}

func (s *PlatformService) DeleteProject(id string, projectId string, ctx context.Context) (*models.PlatformView, error) {
	plat, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, err := plat.RemoveProject(projectId); err != nil {
		return nil, err
	}
	if err = s.repository.Update(ctx, *plat); err != nil {
		return nil, err
	}
	return convertPlatformEntityToViewModel(plat), nil
}

func (s *PlatformService) UpdatePlatform(id string, data models.UpdatePlatformRequest, ctx context.Context) (*models.PlatformView, error) {
	plat, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if plat.Name != data.Name {
		res, err := s.repository.GetPlatformByName(ctx, data.Name)
		if err != nil && !strings.HasPrefix(err.Error(), extensions.Data_Not_Found_Message) {
			return nil, err
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
		newProperty[v.Key] = platform.PropertyInfo(v)
	}
	if !tool.MapsCompareCommon(plat.Property, newProperty) {
		if _, err := plat.UpdateProperty(newProperty); err != nil {
			return nil, err
		}
	}

	if err = s.repository.Update(ctx, *plat); err != nil {
		return nil, err
	}

	return convertPlatformEntityToViewModel(plat), nil
}

func convertPlatformEntityToViewModel(src *platform.Platform) *models.PlatformView {
	if src == nil {
		return nil
	}
	propertyInfos := make([]models.PropertyInfo, 0)
	for _, v := range src.Property {
		propertyInfos = append(propertyInfos, models.PropertyInfo{
			Key:      v.Key,
			Value:    v.Value,
			NeedMask: v.NeedMask,
		})
	}
	platformProjects := make([]models.PlatformProject, 0)
	for _, v := range src.Projects {
		webhooks := make([]models.Webhook, len(v.Webhooks))
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
	return &models.PlatformView{
		Id:           src.Id,
		Name:         src.Name,
		Activate:     src.Activate,
		Url:          src.Url,
		RestEndpoint: src.RestEndpoint,
		Property:     propertyInfos,
		Projects:     platformProjects,
		Tags:         src.Tags,
		IsDeleted:    src.IsDeleted,
	}
}
