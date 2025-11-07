package application

import (
	"context"
	"fmt"
	"log"
	"strings"

	coreapp "github.com/futugyou/domaincore/application"
	coredomain "github.com/futugyou/domaincore/domain"
	coreinfr "github.com/futugyou/domaincore/infrastructure"

	"github.com/futugyou/platformservice/application/service"
	"github.com/futugyou/platformservice/assembler"
	"github.com/futugyou/platformservice/domain"
	"github.com/futugyou/platformservice/viewmodel"
)

type PlatformService struct {
	innerService   *coreapp.AppService
	repository     domain.PlatformRepository
	eventPublisher coreinfr.EventDispatcher
	vaultService   service.VaultService
}

func NewPlatformService(
	unitOfWork coredomain.UnitOfWork,
	repository domain.PlatformRepository,
	eventPublisher coreinfr.EventDispatcher,
	vaultService service.VaultService,
) *PlatformService {
	return &PlatformService{
		innerService:   coreapp.NewAppService(unitOfWork),
		repository:     repository,
		eventPublisher: eventPublisher,
		vaultService:   vaultService,
	}
}

func (s *PlatformService) CreatePlatform(ctx context.Context, aux viewmodel.CreatePlatformRequest) (*viewmodel.PlatformDetailView, error) {
	properties := map[string]domain.Property{}
	for _, v := range aux.Properties {
		properties[v.Key] = domain.Property(v)
	}

	secretMapper := assembler.SecretAssembler{}
	secrets, err := secretMapper.ToModel(ctx, s.vaultService, aux.Secrets)
	if err != nil {
		return nil, err
	}

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

func (s *PlatformService) SearchPlatforms(ctx context.Context, request viewmodel.SearchPlatformsRequest) ([]viewmodel.PlatformView, error) {
	filter := domain.PlatformSearch{
		Name:      request.Name,
		NameFuzzy: false,
		Activate:  request.Activate,
		Tags:      request.Tags,
		Page:      request.Page,
		Size:      request.Size,
	}

	platforms, err := s.repository.SearchPlatforms(ctx, filter)

	if err != nil && !strings.HasPrefix(err.Error(), coredomain.DATA_NOT_FOUND_MESSAGE) {
		return nil, err
	}

	mapper := assembler.PlatformAssembler{}

	return mapper.ToPlatformViews(platforms), nil
}

func (s *PlatformService) GetPlatform(ctx context.Context, idOrName string) (*viewmodel.PlatformDetailView, error) {
	res, err := s.repository.GetPlatformByIdOrName(ctx, idOrName)
	if err != nil {
		return nil, err
	}

	mapper := assembler.PlatformAssembler{}
	provider, err := s.getPlatformProvider(ctx, *res)
	if err != nil {
		log.Println(err.Error())
		return mapper.ToPlatformDetailView(res), nil
	}

	projects, err := s.getProviderProjects(ctx, provider, *res)
	if err != nil {
		log.Println(err.Error())
		return mapper.ToPlatformDetailView(res), nil
	}

	return s.toPlatformDetailViewWithProjects(mapper,res, projects), err
}
