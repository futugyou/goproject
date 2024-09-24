package application

import (
	"context"
	"fmt"
	"os"

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

func (s *VaultService) SearchVaults(ctx context.Context, request models.SearchVaultsRequest, page *int, size *int) ([]models.VaultView, error) {
	filter := []vault.VaultSearch{
		{
			Key:          request.Key,
			KeyFuzzy:     true,
			StorageMedia: request.StorageMedia,
			VaultType:    request.VaultType,
			TypeIdentity: request.TypeIdentity,
			Tags:         request.Tags,
		},
	}
	src, err := s.repository.SearchVaults(ctx, filter, page, size)
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
		return nil, fmt.Errorf("SearchVaults timeout")
	}
}

func (s *VaultService) ShowVaultRawValue(ctx context.Context, vaultId string) (string, error) {
	src, err := s.repository.GetAsync(ctx, vaultId)
	select {
	case data := <-src:
		if data == nil {
			return "", fmt.Errorf("vault with id: %s is not exist", vaultId)
		}
		DecryptVaultValue(data)
		return data.Value, nil
	case errM := <-err:
		return "", errM
	case <-ctx.Done():
		return "", fmt.Errorf("ShowVaultRawValue timeout")
	}
}

func (s *VaultService) CreateVaults(aux models.CreateVaultsRequest, ctx context.Context) (*models.CreateVaultsResponse, error) {
	entities := make([]vault.Vault, 0)
	filter := []vault.VaultSearch{}

	for i := 0; i < len(aux.Vaults); i++ {
		va := aux.Vaults[i]
		entity := vault.NewVault(
			va.Key,
			va.Value,
			vault.WithStorageMedia(vault.GetStorageMedia(va.StorageMedia)),
			vault.WithTags(va.Tags),
			vault.WithVaultType(vault.GetVaultType(va.VaultType), va.TypeIdentity),
		)
		EncryptVaultValue(entity)
		entities = append(entities, *entity)

		filter = append(filter, vault.VaultSearch{
			Key:          va.Key,
			KeyFuzzy:     false,
			StorageMedia: va.StorageMedia,
			VaultType:    va.VaultType,
			TypeIdentity: va.TypeIdentity,
		})
	}

	checksCh, errCh := s.repository.SearchVaults(ctx, filter, nil, nil)
	select {
	case datas := <-checksCh:
		if len(datas) > 0 {
			return nil, fmt.Errorf("some vaults are already existed, check again")
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

func (s *VaultService) ChangeVault(id string, aux models.ChangeVaultRequest, ctx context.Context) (*models.VaultView, error) {
	if tool.IsAllFieldsNil(aux) {
		return nil, fmt.Errorf("no data need change")
	}

	var data *vault.Vault
	filter := generateChangeVaultSearchFilter(aux, id)
	vaultCh, errCh := s.repository.SearchVaults(ctx, filter, nil, nil)
	select {
	case datas := <-vaultCh:
		if len(datas) == 0 || (len(datas) == 1 && id != datas[0].Id) {
			return nil, fmt.Errorf("id %s are not existed", id)
		}
		if len(datas) > 1 {
			return nil, fmt.Errorf("vaults with 'key+storage_media+vault_type+type_identity' was already existed, check again")
		}
		data = &datas[0]
	case errM := <-errCh:
		return nil, errM
	case <-ctx.Done():
		return nil, fmt.Errorf("CreateVaults timeout")
	}

	if data == nil {
		return nil, fmt.Errorf("id %s are not existed", id)
	}

	doVaultChange(data, aux)

	if data.HasChange() {
		if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
			return <-s.repository.UpdateAsync(ctx, *data)
		}); err != nil {
			return nil, err
		}
	}

	model := convertVaultToVaultView(*data)
	return &model, nil
}

func generateChangeVaultSearchFilter(aux models.ChangeVaultRequest, id string) []vault.VaultSearch {
	filter := []vault.VaultSearch{{
		ID: id,
	}}

	subFilter := vault.VaultSearch{}
	if aux.Key != nil {
		subFilter.Key = *aux.Key
	}
	if aux.StorageMedia != nil {
		subFilter.StorageMedia = *aux.StorageMedia
	}
	if aux.VaultType != nil {
		subFilter.VaultType = *aux.VaultType
	}
	if aux.TypeIdentity != nil {
		subFilter.TypeIdentity = *aux.TypeIdentity
	}
	if !tool.IsAllFieldsNil(aux) {
		filter = append(filter, subFilter)
	}
	return filter
}

func doVaultChange(data *vault.Vault, aux models.ChangeVaultRequest) {
	if aux.Key != nil {
		data.UpdateKey(*aux.Key)
	}

	if aux.Value != nil {
		data.UpdateValue(*aux.Value)
		EncryptVaultValue(data)
	}

	if aux.StorageMedia != nil {
		storageMedia := vault.GetStorageMedia(*aux.StorageMedia)
		data.UpdateStorageMedia(storageMedia)

	}

	if aux.VaultType != nil || aux.TypeIdentity != nil {
		vaultType := vault.GetVaultType(*aux.VaultType)
		data.UpdateVaultType(vaultType, *aux.TypeIdentity)
	}

	if aux.Tags != nil {
		data.UpdateTags(*aux.Tags)
	}
}

func convertVaultToVaultView(entity vault.Vault) models.VaultView {
	DecryptVaultValue(&entity)

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

func EncryptVaultValue(entity *vault.Vault) {
	if entity == nil {
		return
	}
	value := entity.Value
	if entity.StorageMedia == vault.StorageMediaLocal {
		value, _ = tool.AesCTREncrypt(entity.Value, os.Getenv("Encrypt_Key"))
	}
	entity.Value = value
}

func DecryptVaultValue(entity *vault.Vault) {
	if entity == nil {
		return
	}
	value := entity.Value
	if entity.StorageMedia == vault.StorageMediaLocal {
		value, _ = tool.AesCTRDecrypt(entity.Value, os.Getenv("Encrypt_Key"))
	}
	entity.Value = value
}
