package application

import (
	"context"
	"fmt"

	domain "github.com/futugyou/infr-project/domain"
	platform "github.com/futugyou/infr-project/platform"
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
	if err != nil {
		return nil, err
	}

	if res != nil {
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

type CreatePlatformRequest struct {
	Name     string            `json:"name"`
	Url      string            `json:"url"`
	Rest     string            `json:"rest"`
	Property map[string]string `json:"property"`
}
