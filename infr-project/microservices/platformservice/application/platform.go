package application

import (
	"context"
	"fmt"
	"strings"

	coreapp "github.com/futugyou/domaincore/application"
	coredomain "github.com/futugyou/domaincore/domain"
	coreinfr "github.com/futugyou/domaincore/infrastructure"

	"github.com/futugyou/platformservice/assembler"
	"github.com/futugyou/platformservice/domain"

	"github.com/futugyou/platformservice/viewmodel"
)

type PlatformService struct {
	innerService   *coreapp.AppService
	repository     domain.PlatformRepository
	eventPublisher coreinfr.EventDispatcher
}

func NewPlatformService(
	unitOfWork coredomain.UnitOfWork,
	repository domain.PlatformRepository,
	eventPublisher coreinfr.EventDispatcher,
) *PlatformService {
	return &PlatformService{
		innerService:   coreapp.NewAppService(unitOfWork),
		repository:     repository,
		eventPublisher: eventPublisher,
	}
}

func (s *PlatformService) CreatePlatform(ctx context.Context, aux viewmodel.CreatePlatformRequest) (*viewmodel.PlatformDetailView, error) {
	properties := map[string]domain.Property{}
	for _, v := range aux.Properties {
		properties[v.Key] = domain.Property(v)
	}

	// TODO: get secrets from vault
	secrets := map[string]domain.Secret{}

	res, err := domain.NewPlatform(
		aux.Name,
		aux.Url,
		domain.GetPlatformProvider(aux.Provider),
		domain.WithPlatformProperties(properties),
		domain.WithPlatformTags(aux.Tags),
		domain.WithPlatformSecrets(secrets),
	)

	if err != nil {
		return nil, err
	}

	// check name
	resdb, err := s.repository.GetPlatformByName(ctx, aux.Name)
	if resdb != nil {
		return nil, fmt.Errorf("name: %s is existed", aux.Name)
	}

	if err != nil && !strings.HasPrefix(err.Error(), coredomain.DATA_NOT_FOUND_MESSAGE) {
		return nil, err
	}

	if res.Provider != domain.PlatformProviderOther {
		status := s.determineProviderStatus(ctx, res)
		if status {
			res.Enable()
		} else {
			res.Disable()
		}
	}

	if err := s.innerService.WithUnitOfWork(ctx, func(ctx context.Context) error {
		return s.repository.Insert(ctx, *res)
	}); err != nil {
		return nil, err
	}

	mapper := assembler.PlatformAssembler{}

	return mapper.ToPlatformDetailView(res), err
}
