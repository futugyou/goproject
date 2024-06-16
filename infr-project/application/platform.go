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

func (s *PlatformService) CreatePlatform(name string, url string, rest string, property map[string]string) (*platform.Platform, error) {
	var res *platform.Platform
	ctx := context.Background()
	res, err := s.repository.GetPlatformByName(ctx, name)
	if err != nil && !strings.HasPrefix(err.Error(), "data not found") {
		return nil, err
	}

	if res != nil && res.Name == name {
		return nil, fmt.Errorf("name: %s is existed", name)
	}

	err = s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		res = platform.NewPlatform(name, url, rest, property)
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

func (s *PlatformService) AddWebhook(projectId string, hook platform.Webhook) (*platform.Platform, error) {
	res, err := s.repository.GetAllPlatform(context.Background())
	if err != nil {
		return nil, err
	}

	var plat *platform.Platform
	for _, p := range res {
		if _, exists := p.Projects[projectId]; exists {
			plat = &p
			break
		}
	}

	if plat == nil {
		return nil, fmt.Errorf("projectId: %s is no existed", projectId)
	}

	plat.UpdateWebhook(projectId, hook)
	err = s.repository.Update(context.Background(), plat)
	if err != nil {
		return nil, err
	}

	return plat, nil
}

func (s *PlatformService) DeletePlatform(id string) error {
	return s.repository.Delete(context.Background(), id)
}

func (s *PlatformService) UpdatePlatform(id string, data models.UpdatePlatformRequest) (*platform.Platform, error) {
	res, err := s.repository.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}

	plat := *res
	if len(data.Name) > 0 {
		plat.UpdateName(data.Name)
	}
	if len(data.Url) > 0 {
		plat.UpdateUrl(data.Url)
	}
	if len(data.Rest) > 0 {
		plat.UpdateRestEndpoint(data.Rest)
	}
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
