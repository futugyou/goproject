package application

import (
	"context"
	"fmt"
	"sync"

	tool "github.com/futugyou/extensions"

	domain "github.com/futugyou/infr-project/domain"
	infra "github.com/futugyou/infr-project/infrastructure"
	vault "github.com/futugyou/infr-project/vault"
	provider "github.com/futugyou/infr-project/vault_provider"
	models "github.com/futugyou/infr-project/view_models"
)

type VaultService struct {
	innerService  *AppService
	repository    vault.IVaultRepositoryAsync
	eventPublisher infra.IEventPublisher
}

func NewVaultService(
	unitOfWork domain.IUnitOfWork,
	repository vault.IVaultRepositoryAsync,
	eventPublisher infra.IEventPublisher,
) *VaultService {
	return &VaultService{
		innerService:  NewAppService(unitOfWork),
		repository:    repository,
		eventPublisher: eventPublisher,
	}
}

type VaultSearchQuery struct {
	Filters []vault.VaultSearch
	Page    int
	Size    int
}

func (s *VaultService) SearchVaults(ctx context.Context, query VaultSearchQuery) ([]models.VaultView, error) {
	src, err := s.repository.SearchVaultsAsync(ctx, query.Filters, &query.Page, &query.Size)
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
		return nil, fmt.Errorf("SearchVaults timeout: %w", ctx.Err())
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
		return "", fmt.Errorf("ShowVaultRawValue timeout: %w", ctx.Err())
	}
}

func (s *VaultService) CreateVaults(ctx context.Context, aux models.CreateVaultsRequest) (*models.CreateVaultsResponse, error) {
	if len(aux.Vaults) == 0 {
		return nil, fmt.Errorf("no vaults need to create")
	}

	entities := make([]vault.Vault, 0)
	storageMediaList := make(map[string]struct{})
	storageMedia := ""

	for i := 0; i < len(aux.Vaults); i++ {
		va := aux.Vaults[i]
		entity := vault.NewVault(
			va.Key,
			va.Value,
			vault.WithStorageMedia(vault.GetStorageMedia(va.StorageMedia)),
			vault.WithTags(va.Tags),
			vault.WithVaultType(vault.GetVaultType(va.VaultType), va.TypeIdentity),
			vault.WithExtension(va.Extension),
			vault.WithDescription(va.Description),
		)
		entities = append(entities, *entity)
		if _, ok := storageMediaList[va.StorageMedia]; !ok {
			storageMediaList[va.StorageMedia] = struct{}{}
			storageMedia = va.StorageMedia
		}
	}

	if len(storageMediaList) > 1 {
		return nil, fmt.Errorf("StorageMedia can only contain one type per request")
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
			return nil, fmt.Errorf("CreateVaults timeout: %w", ctx.Err())
		}
	}

	if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		errCh := s.repository.InsertMultipleVaultAsync(ctx, entities)
		select {
		case err := <-errCh:
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return fmt.Errorf("CreateVaults timeout: %w", ctx.Err())
		}

		if storageMedia == vault.StorageMediaLocal.String() {
			return nil
		}

		vaultDatas := map[string]string{}
		for _, item := range entities {
			vaultDatas[item.GetIdentityKey()] = item.Value
		}

		// If an error occurs, you can force an 'ForceInsert' operation
		return s.upsertVaultInProvider(ctx, storageMedia, vaultDatas)
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

func (s *VaultService) CreateVault(ctx context.Context, aux models.CreateVaultRequest) (*models.VaultView, error) {
	createVaultsRequest := models.CreateVaultsRequest{
		Vaults:      []models.CreateVaultModel{aux.CreateVaultModel},
		ForceInsert: aux.ForceInsert,
	}

	result, err := s.CreateVaults(ctx, createVaultsRequest)
	if err != nil {
		return nil, err
	}

	if len(result.Vaults) == 0 {
		return nil, fmt.Errorf("create vault error, check data again")
	}

	return &result.Vaults[0], nil
}

func (s *VaultService) ChangeVault(ctx context.Context, id string, aux models.ChangeVaultRequest) (*models.VaultView, error) {
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
		return nil, fmt.Errorf("ChangeVault timeout: %w", ctx.Err())
	}

	if data == nil {
		return nil, fmt.Errorf("id %s are not existed", id)
	}

	doVaultChange(data, aux.Data)

	if data.HasChange() {
		if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
			errCh = s.repository.UpdateAsync(ctx, *data)
			select {
			case err := <-errCh:
				if err != nil {
					return err
				}
			case <-ctx.Done():
				return fmt.Errorf("ChangeVault timeout: %w", ctx.Err())
			}

			if data.StorageMedia == vault.StorageMediaLocal {
				return nil
			}

			s.eventPublisher.PublishCommon(ctx, data, "vault_changed")

			return s.upsertVaultInProvider(ctx, data.StorageMedia.String(), map[string]string{data.GetIdentityKey(): data.Value})
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
	case <-ctx.Done():
		return false, fmt.Errorf("DeleteVault timeout: %w", ctx.Err())
	case err := <-errCh:
		return false, err
	case va = <-vaCh:
	}

	if va == nil {
		return false, fmt.Errorf("vault with id: %s is not exist", vaultId)
	}

	if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		errCh := s.repository.DeleteAsync(ctx, vaultId)
		select {
		case err := <-errCh:
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return fmt.Errorf("DeleteVault timeout: %w", ctx.Err())
		}

		if va.StorageMedia == vault.StorageMediaLocal {
			return nil
		}

		return s.deleteVaultInProvider(ctx, va.VaultType.String(), va.GetIdentityKey())
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (s *VaultService) ImportVaults(ctx context.Context, aux models.ImportVaultsRequest) (*models.ImportVaultsResponse, error) {
	vt := "system"
	vi := "system"
	if aux.VaultType != nil {
		switch *aux.VaultType {
		case "common":
			vt = "common"
			vi = "common"
		case "project", "resource", "platform":
			if aux.TypeIdentity == nil {
				return nil, fmt.Errorf("when VaultType is not system and common, the TypeIdentity cannot be nil")
			}
			vt = *aux.VaultType
			vi = *aux.TypeIdentity
		}
	}

	entities := make([]vault.Vault, 0)
	if datas, err := s.searchVaultInProvider(ctx, aux.StorageMedia, fmt.Sprintf("%s/%s", vt, vi)); err != nil {
		return nil, err
	} else {
		for _, data := range datas {
			entities = append(entities, *vault.NewVault(
				data.Key,
				data.Value,
				vault.WithStorageMedia(vault.GetStorageMedia(aux.StorageMedia)),
				vault.WithVaultType(vault.GetVaultType(vt), vi),
			))
		}
	}

	if err := s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		errCh := s.repository.InsertMultipleVaultAsync(ctx, entities)
		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return fmt.Errorf("DeleteVault timeout: %w", ctx.Err())
		}
	}); err != nil {
		return nil, err
	}

	response := models.ImportVaultsResponse{
		Vaults: []models.VaultView{},
	}
	for _, va := range entities {
		response.Vaults = append(response.Vaults, convertVaultToVaultView(va))
	}
	return &response, nil
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

	if aux.Description != nil {
		data.UpdateDescription(*aux.Description)
	}

	if aux.Extension != nil {
		data.UpdateExtension(*aux.Extension)
	}

	if aux.Value != nil {
		value := *aux.Value
		maskValue := tool.MaskString(data.Value, 5, 0.5)
		if value != maskValue {
			data.UpdateValue(value)
		}
	}

	if aux.StorageMedia != nil {
		storageMedia := vault.GetStorageMedia(*aux.StorageMedia)
		data.UpdateStorageMedia(storageMedia)
	}

	if aux.VaultType != nil && aux.TypeIdentity != nil {
		vaultType := vault.GetVaultType(*aux.VaultType)
		data.UpdateVaultType(vaultType, *aux.TypeIdentity)
	}

	if aux.Tags != nil {
		data.UpdateTags(*aux.Tags)
	}
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
		Description:  entity.Description,
		Extension:    entity.Extension,
	}
}

func (s *VaultService) deleteVaultInProvider(ctx context.Context, provider_type string, key string) error {
	p, err := provider.VaultProviderFactory(provider_type)
	if err != nil {
		return err
	}

	return p.Delete(ctx, key)
}

func (s *VaultService) upsertVaultInProvider(ctx context.Context, provider_type string, datas map[string]string) error {
	p, err := provider.VaultProviderFactory(provider_type)
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

func (s *VaultService) searchVaultInProvider(ctx context.Context, provider_type string, prefix string) (map[string]provider.ProviderVault, error) {
	p, err := provider.VaultProviderFactory(provider_type)
	if err != nil {
		return nil, err
	}

	return p.PrefixSearch(ctx, prefix)
}
