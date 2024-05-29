package application

import (
	"context"

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

func (s *PlatformService) CreateResource(name string, url string, rest string, property map[string]string) (*platform.Platform, error) {
	var res *platform.Platform
	err := s.innerService.withUnitOfWork(context.Background(), func(ctx context.Context) error {
		res = platform.NewPlatform(name, url, rest, property)
		return s.repository.Insert(ctx, res)
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
