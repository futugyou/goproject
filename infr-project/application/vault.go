package application

import (
	"context"
	"fmt"
	"os"
	"sync"

	tool "github.com/futugyou/extensions"

	domain "github.com/futugyou/infr-project/domain"
	vault "github.com/futugyou/infr-project/vault"
	provider "github.com/futugyou/infr-project/vault_provider"
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

func (s *VaultService) SearchVaults(ctx context.Context, request models.SearchVaultsRequest) ([]models.VaultView, error) {
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
	src, err := s.repository.SearchVaultsAsync(ctx, filter, &request.Page, &request.Size)
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
		decryptVaultValue(data)
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
		entity := vault.NewVault(
			va.Key,
			va.Value,
			vault.WithStorageMedia(vault.GetStorageMedia(va.StorageMedia)),
			vault.WithTags(va.Tags),
			vault.WithVaultType(vault.GetVaultType(va.VaultType), va.TypeIdentity),
		)
		encryptVaultValue(entity)
		entities = append(entities, *entity)
	}

	if !aux.ForceInsert {
		filter := []vault.VaultSearch{}
		for i := 0; i < len(aux.Vaults); i++ {
			va := aux.Vaults[i]
			filter = append(filter, vault.VaultSearch{
				Key:          va.Key,
				KeyFuzzy:     false,
				StorageMedia: va.StorageMedia,
				VaultType:    va.VaultType,
				TypeIdentity: va.TypeIdentity,
			})
		}

		checksCh, errCh := s.repository.SearchVaultsAsync(ctx, filter, nil, nil)
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
	if tool.IsAllFieldsNil(aux.Data) {
		return nil, fmt.Errorf("no data need change")
	}

	var data *vault.Vault
	filter := generateChangeVaultSearchFilter(aux.Data, id)
	vaultCh, errCh := s.repository.SearchVaultsAsync(ctx, filter, nil, nil)
	select {
	case datas := <-vaultCh:
		if len(datas) == 0 || (len(datas) == 1 && id != datas[0].Id) {
			return nil, fmt.Errorf("id %s are not existed", id)
		}
		if len(datas) > 1 && !aux.ForceInsert {
			return nil, fmt.Errorf("vaults with 'key+storage_media+vault_type+type_identity' was already existed, check again")
		}
		for _, da := range datas {
			if da.Id == id {
				data = &da
				break
			}
		}
	case errM := <-errCh:
		return nil, errM
	case <-ctx.Done():
		return nil, fmt.Errorf("CreateVaults timeout")
	}

	if data == nil {
		return nil, fmt.Errorf("id %s are not existed", id)
	}

	doVaultChange(data, aux.Data)

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

func (s *VaultService) DeleteVault(ctx context.Context, vaultId string) (bool, error) {
	vaCh, errCh := s.repository.GetAsync(ctx, vaultId)
	var va *vault.Vault
	select {
	case err := <-errCh:
		if err != nil {
			return false, err
		}
	case va = <-vaCh:
	case <-ctx.Done():
		return false, fmt.Errorf("DeleteVault timeout")
	}

	if va == nil {
		return false, fmt.Errorf("vault with id: %s is not exist", vaultId)
	}

	errCh = s.repository.DeleteAsync(ctx, vaultId)
	select {
	case err := <-errCh:
		if err == nil {
			if err = s.deleteVaultInProvider(ctx, va.VaultType.String(), va.Key); err != nil {
				return true, nil
			} else {
				return false, err
			}
		} else {
			return false, err
		}
	case <-ctx.Done():
		return false, fmt.Errorf("DeleteVault timeout")
	}
}

func generateChangeVaultSearchFilter(aux models.ChangeVaultItem, id string) []vault.VaultSearch {
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

func doVaultChange(data *vault.Vault, aux models.ChangeVaultItem) {
	if aux.Key != nil {
		data.UpdateKey(*aux.Key)
	}

	if aux.Value != nil {
		data.UpdateValue(*aux.Value)
		encryptVaultValue(data)
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
	decryptVaultValue(&entity)

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

func encryptVaultValue(entity *vault.Vault) error {
	if entity == nil {
		return fmt.Errorf("vault can not be nil")
	}
	value, err := tool.AesCTREncrypt(entity.Value, os.Getenv("Encrypt_Key"))
	if err != nil {
		return err
	}
	return entity.UpdateValue(value)
}

func decryptVaultValue(entity *vault.Vault) error {
	if entity == nil {
		return fmt.Errorf("vault can not be nil")
	}
	value, err := tool.AesCTRDecrypt(entity.Value, os.Getenv("Encrypt_Key"))
	if err != nil {
		return err
	}
	return entity.UpdateValue(value)
}

func (s *VaultService) deleteVaultInProvider(ctx context.Context, provider_type string, key string) error {
	p, err := provider.VaultProviderFatory(provider_type)
	if err != nil {
		return err
	}

	return p.Delete(ctx, key)
}

func (s *VaultService) upsertVaultInProvider(ctx context.Context, provider_type string, datas map[string]string) error {
	p, err := provider.VaultProviderFatory(provider_type)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	concurrencyLimit := 5
	sem := make(chan struct{}, concurrencyLimit)

	errCh := make(chan error, len(datas))
	defer close(errCh)

	for key, value := range datas {
		wg.Add(1)

		go func(key string, value string) {
			defer wg.Done()

			sem <- struct{}{}

			_, err := p.Upsert(ctx, key, value)
			if err != nil {
				errCh <- err
			}
			<-sem
		}(key, value)
	}

	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

func (s *VaultService) searchVaultInProvider(ctx context.Context, provider_type string, keys []string) (map[string]provider.ProviderVault, error) {
	p, err := provider.VaultProviderFatory(provider_type)
	if err != nil {
		return nil, err
	}

	return p.BatchSearch(ctx, keys)
}
