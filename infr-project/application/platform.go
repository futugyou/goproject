package application

import (
	"context"
	"fmt"
	"strings"

	domain "github.com/futugyou/infr-project/domain"
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

func (s *PlatformService) CreatePlatform(aux models.CreatePlatformRequest) (*platform.Platform, error) {
	var res *platform.Platform
	ctx := context.Background()
	res, err := s.repository.GetPlatformByName(ctx, aux.Name)
	if err != nil && !strings.HasPrefix(err.Error(), "data not found") {
		return nil, err
	}

	if res != nil && res.Name == aux.Name {
		return nil, fmt.Errorf("name: %s is existed", aux.Name)
	}

	err = s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		res = platform.NewPlatform(aux.Name, aux.Url, aux.Rest, aux.Property, aux.Tags)
		return s.repository.Insert(ctx, res)
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *PlatformService) GetAllPlatform() ([]platform.Platform, error) {
	return s.repository.GetAllPlatform(context.Background())
}

func (s *PlatformService) GetPlatform(id string) (*platform.Platform, error) {
	res, err := s.repository.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return *res, nil
}

func (s *PlatformService) AddWebhook(id string, projectId string, hook models.UpdatePlatformWebhookRequest) (*platform.Platform, error) {
	res, err := s.repository.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}

	plat := *res
	if _, exists := plat.Projects[projectId]; !exists {
		return nil, fmt.Errorf("projectId: %s is not existed in %s", projectId, id)
	}

	newhook := platform.NewWebhook(hook.Name, hook.Url, hook.Property)
	newhook.Activate = hook.Activate
	newhook.State = platform.GetWebhookState(hook.State)
	plat.UpdateWebhook(projectId, *newhook)
	err = s.repository.Update(context.Background(), plat)
	if err != nil {
		return nil, err
	}

	return plat, nil
}

func (s *PlatformService) DeletePlatform(id string) error {
	return s.repository.Delete(context.Background(), id)
}

func (s *PlatformService) AddProject(id string, projectId string, project models.UpdatePlatformProjectRequest) (*platform.Platform, error) {
	res, err := s.repository.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}

	if len(projectId) == 0 {
		projectId = project.Name
	}

	plat := *res
	proj := platform.NewPlatformProject(projectId, project.Name, project.Url, project.Property)
	plat.UpdateProject(*proj)
	err = s.repository.Update(context.Background(), plat)
	if err != nil {
		return nil, err
	}

	return plat, nil
}

func (s *PlatformService) DeleteProject(id string, projectId string) (*platform.Platform, error) {
	res, err := s.repository.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}

	plat := *res
	plat.RemoveProject(projectId)
	return plat, s.repository.Update(context.Background(), plat)
}

func (s *PlatformService) UpdatePlatform(id string, data models.UpdatePlatformRequest) (*platform.Platform, error) {
	res, err := s.repository.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}

	plat := *res
	plat.UpdateName(data.Name)
	plat.UpdateUrl(data.Url)
	plat.UpdateRestEndpoint(data.Rest)
	plat.UpdateTags(data.Tags)
	if data.Activate != nil {
		if *data.Activate {
			plat.Enable()
		} else {
			plat.Disable()
		}
	}
	if data.Property != nil {
		plat.UpdateProperty(data.Property)
	}
	err = s.repository.Update(context.Background(), plat)
	if err != nil {
		return nil, err
	}
	return plat, nil
}
