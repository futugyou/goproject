package application

import (
	"context"
	"fmt"
	"strings"

	tool "github.com/futugyou/extensions"

	domain "github.com/futugyou/infr-project/domain"
	vault "github.com/futugyou/infr-project/vault"
	models "github.com/futugyou/infr-project/view_models"
)

type VaultService struct {
	innerService *AppService
	repository   vault.IVaultRepositoryAsync
}

func NewVaultService(
	unitOfWork domain.IUnitOfWork,
	repository vault.IVaultRepositoryAsync,
) *VaultService {
	return &VaultService{
		innerService: NewAppService(unitOfWork),
		repository:   repository,
	}
}

func (s *VaultService) GetAllVault(ctx context.Context, page *int, size *int) ([]models.VaultView, error) {
	src, err := s.repository.SearchVaults(ctx, nil, page, size)
	select {
	case datas := <-src:
		result := make([]models.VaultView, len(datas))
		for i := 0; i < len(datas); i++ {
			result[i] = convertVaultToVaultView(datas[i])
		}
		return result, nil
	case errM := <-err:
		return nil, errM
	case <-ctx.Done():
		return nil, fmt.Errorf("GetAllVault timeout")
	}
}

func (s *VaultService) ShowVaultRawValue(ctx context.Context, vaultId string) (string, error) {
	src, err := s.repository.GetAsync(ctx, vaultId)
	select {
	case data := <-src:
		if data == nil {
			return "", fmt.Errorf("vault with id: %s is not exist", vaultId)
		}
		return data.Value, nil
	case errM := <-err:
		return "", errM
	case <-ctx.Done():
		return "", fmt.Errorf("ShowVaultRawValue timeout")
	}
}

func (s *VaultService) CreateVaults(aux models.CreateVaultsRequest, ctx context.Context) (*models.CreateVaultsResponse, error) {
	entities := make([]vault.Vault, 0)
	for i := 0; i < len(aux.Vaults); i++ {
		va := aux.Vaults[i]
		entities = append(entities,
			*vault.NewVault(
				va.Key,
				va.Value,
				vault.WithStorageMedia(vault.GetStorageMedia(va.StorageMedia)),
				vault.WithTags(va.Tags),
				vault.WithVaultType(vault.GetVaultType(va.VaultType), va.TypeIdentity),
			))
	}

	ids := []string{}
	for _, vault := range entities {
		ids = append(ids, vault.Id)
	}

	checksCh, errCh := s.repository.GetVaultByIdsAsync(ctx, ids)
	select {
	case datas := <-checksCh:
		if len(datas) > 0 {
			ids = []string{}
			for _, vault := range datas {
				ids = append(ids, vault.Id)
			}
			return nil, fmt.Errorf("id %s are already existed", strings.Join(ids, ","))
		}
	case errM := <-errCh:
		return nil, errM
	case <-ctx.Done():
		return nil, fmt.Errorf("CreateVaults timeout")
	}

	if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		return <-s.repository.InsertMultipleVaultAsync(ctx, entities)
	}); err != nil {
		return nil, err
	}

	response := models.CreateVaultsResponse{
		Vaults: []models.VaultView{},
	}
	for _, va := range entities {
		response.Vaults = append(response.Vaults, convertVaultToVaultView(va))
	}
	return &response, nil
}

func convertVaultToVaultView(entity vault.Vault) models.VaultView {
	return models.VaultView{
		Id:           entity.Id,
		Key:          entity.Key,
		MaskValue:    tool.MaskString(entity.Value, 5, 0.5),
		StorageMedia: entity.StorageMedia.String(),
		VaultType:    entity.VaultType.String(),
		TypeIdentity: entity.TypeIdentity,
		Tags:         entity.Tags,
	}
}
