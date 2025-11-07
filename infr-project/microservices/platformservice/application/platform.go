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

	return s.toPlatformDetailViewWithProjects(mapper, res, projects), err
}

func (s *PlatformService) GetProviderProjectList(ctx context.Context, idOrName string) ([]viewmodel.PlatformProviderProject, error) {
	src, err := s.repository.GetPlatformByIdOrName(ctx, idOrName)
	if err != nil {
		return nil, err
	}

	if src == nil {
		return nil, fmt.Errorf("no platform data found")
	}

	result := []viewmodel.PlatformProviderProject{}
	provider, err := s.getPlatformProvider(ctx, *src)
	if err != nil {
		return nil, err
	}

	projects, err := s.getProviderProjects(ctx, provider, *src)
	if err != nil {
		return nil, err
	}

	for _, pro := range projects {
		project := s.convertProviderProjectToSimpleModel(pro)
		result = append(result, project)
	}

	return result, nil
}

func (s *PlatformService) UpdatePlatform(ctx context.Context, idOrName string, data viewmodel.UpdatePlatformRequest) (*viewmodel.PlatformDetailView, error) {
	return s.updatePlatform(ctx, idOrName, "UpdatePlatform", func(plat *domain.Platform) error {
		if plat.IsDeleted {
			return fmt.Errorf("id: %s was already deleted", plat.ID)
		}

		if plat.Name != data.Name {
			res, err := s.repository.GetPlatformByName(ctx, data.Name)
			if err != nil && !strings.HasPrefix(err.Error(), coredomain.DATA_NOT_FOUND_MESSAGE) {
				return err
			}
			if res.ID != plat.ID {
				return fmt.Errorf("name: %s is existed", data.Name)
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

		if _, err := plat.UpdateProvider(domain.GetPlatformProvider(data.Provider)); err != nil {
			return err
		}

		if plat.Provider != domain.PlatformProviderOther {
			status := s.determineProviderStatus(ctx, plat)
			if status {
				plat.Enable()
			} else {
				plat.Disable()
			}
		}

		newProperty := map[string]domain.Property{}
		for _, v := range data.Properties {
			newProperty[v.Key] = domain.Property{
				Key:   v.Key,
				Value: v.Value,
			}
		}
		if _, err := plat.UpdateProperties(newProperty); err != nil {
			return err
		}

		secretMapper := assembler.SecretAssembler{}
		newSecrets, err := secretMapper.ToModel(ctx, s.vaultService, data.Secrets)
		if err != nil {
			return err
		}

		if _, err := plat.UpdateSecrets(newSecrets); err != nil {
			return err
		}

		return nil
	})
}

func (s *PlatformService) DeletePlatform(ctx context.Context, idOrName string) (*viewmodel.PlatformDetailView, error) {
	return s.updatePlatform(ctx, idOrName, "DeletePlatform", func(plat *domain.Platform) error {
		if _, err := plat.Delete(); err != nil {
			return err
		}

		return nil
	})
}

func (s *PlatformService) RecoveryPlatform(ctx context.Context, idOrName string) (*viewmodel.PlatformDetailView, error) {
	return s.updatePlatform(ctx, idOrName, "RecoveryPlatform", func(plat *domain.Platform) error {
		if _, err := plat.Recovery(); err != nil {
			return err
		}

		return nil
	})
}

func (s *PlatformService) updatePlatform(ctx context.Context, idOrName string, _ string, fn func(*domain.Platform) error) (*viewmodel.PlatformDetailView, error) {
	plat, err := s.repository.GetPlatformByIdOrName(ctx, idOrName)
	if err != nil {
		return nil, err
	}

	if err := fn(plat); err != nil {
		return nil, err
	}

	if err := s.innerService.WithUnitOfWork(ctx, func(ctx context.Context) error {
		return s.repository.Update(ctx, *plat)
	}); err != nil {
		return nil, err
	}

	mapper := assembler.PlatformAssembler{}

	return mapper.ToPlatformDetailView(plat), err
}
