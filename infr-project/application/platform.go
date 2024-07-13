package application

import (
	"context"
	"fmt"
	"strings"

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

func (s *PlatformService) CreatePlatform(aux models.CreatePlatformRequest, ctx context.Context) (*platform.Platform, error) {
	var res *platform.Platform
	res, err := s.repository.GetPlatformByName(ctx, aux.Name)
	if err != nil && !strings.HasPrefix(err.Error(), "data not found") {
		return nil, err
	}

	if res != nil && res.Name == aux.Name {
		return nil, fmt.Errorf("name: %s is existed", aux.Name)
	}

	err = s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		res = platform.NewPlatform(aux.Name, aux.Url, aux.Rest, aux.Property, aux.Tags)
		return s.repository.Insert(ctx, *res)
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *PlatformService) GetAllPlatform(ctx context.Context) ([]platform.Platform, error) {
	return s.repository.GetAllPlatform(ctx)
}

func (s *PlatformService) GetPlatform(id string, ctx context.Context) (*platform.Platform, error) {
	return s.repository.Get(ctx, id)
}

func (s *PlatformService) AddWebhook(id string, projectId string, hook models.UpdatePlatformWebhookRequest, ctx context.Context) (*platform.Platform, error) {
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
	plat.UpdateWebhook(projectId, *newhook)
	err = s.repository.Update(ctx, *plat)
	if err != nil {
		return nil, err
	}

	return plat, nil
}

func (s *PlatformService) DeletePlatform(id string, ctx context.Context) (*platform.Platform, error) {
	if err := s.repository.SoftDelete(ctx, id); err != nil {
		return nil, err
	}

	return s.repository.Get(ctx, id)
}

func (s *PlatformService) AddProject(id string, projectId string, project models.UpdatePlatformProjectRequest, ctx context.Context) (*platform.Platform, error) {
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
	err = s.repository.Update(ctx, *plat)
	if err != nil {
		return nil, err
	}

	return plat, nil
}

func (s *PlatformService) DeleteProject(id string, projectId string, ctx context.Context) (*platform.Platform, error) {
	plat, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if _, err := plat.RemoveProject(projectId); err != nil {
		return nil, err
	}

	return plat, s.repository.Update(ctx, *plat)
}

func (s *PlatformService) UpdatePlatform(id string, data models.UpdatePlatformRequest, ctx context.Context) (*platform.Platform, error) {
	plat, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if plat.Name != data.Name {
		res, err := s.repository.GetPlatformByName(ctx, data.Name)
		if err != nil && !strings.HasPrefix(err.Error(), "no data found") {
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

	if !extensions.StringArrayCompare(plat.Tags, data.Tags) {
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

	if !extensions.MapsCompare(data.Property, data.Property) {
		if _, err := plat.UpdateProperty(data.Property); err != nil {
			return nil, err
		}
	}

	err = s.repository.Update(ctx, *plat)
	if err != nil {
		return nil, err
	}

	return plat, nil
}
