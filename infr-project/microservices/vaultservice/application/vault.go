package application

import (
	"context"
	"fmt"
	"sync"

	coreapp "github.com/futugyou/domaincore/application"
	coredomain "github.com/futugyou/domaincore/domain"
	coreinfr "github.com/futugyou/domaincore/infrastructure"

	tool "github.com/futugyou/extensions"

	"github.com/futugyou/vaultservice/domain"
	"github.com/futugyou/vaultservice/options"
	"github.com/futugyou/vaultservice/provider"
	"github.com/futugyou/vaultservice/viewmodel"
)

type VaultService struct {
	innerService   *coreapp.AppService
	repository     domain.VaultRepository
	eventPublisher coreinfr.EventDispatcher
	opts           *options.Options
}

func NewVaultService(
	unitOfWork coredomain.UnitOfWork,
	repository domain.VaultRepository,
	eventPublisher coreinfr.EventDispatcher,
	opts *options.Options,
) *VaultService {
	return &VaultService{
		innerService:   coreapp.NewAppService(unitOfWork),
		repository:     repository,
		eventPublisher: eventPublisher,
		opts:           opts,
	}
}

func toVaultSearchQuery(aux viewmodel.SearchVaultsRequest) domain.VaultQuery {
	page := aux.Page
	size := aux.Size
	return domain.VaultQuery{
		Filters: []domain.VaultFilter{
			{
				Key:          aux.Key,
				KeyFuzzy:     true,
				StorageMedia: aux.StorageMedia,
				VaultType:    aux.VaultType,
				TypeIdentity: aux.TypeIdentity,
				Description:  aux.Description,
				Tags:         aux.Tags,
			},
		},
		Page: &page,
		Size: &size,
	}
}

func (s *VaultService) SearchVaults(ctx context.Context, request viewmodel.SearchVaultsRequest) ([]viewmodel.VaultView, error) {
	query := toVaultSearchQuery(request)
	datas, err := s.repository.SearchVaults(ctx, query)
	if err != nil {
		return nil, err
	}

	result := make([]viewmodel.VaultView, len(datas))
	for i := range datas {
		result[i] = convertVaultToVaultView(datas[i])
	}
	return result, nil
}

func (s *VaultService) ShowVaultRawValue(ctx context.Context, vaultId string) (string, error) {
	src, err := s.repository.FindByID(ctx, vaultId)
	if err != nil {
		return "", err
	}
	if src == nil {
		return "", fmt.Errorf("vault with id: %s is not exist", vaultId)
	}
	return src.Value, nil
}

func (s *VaultService) CreateVaults(ctx context.Context, aux viewmodel.CreateVaultsRequest) (*viewmodel.CreateVaultsResponse, error) {
	if len(aux.Vaults) == 0 {
		return nil, fmt.Errorf("no vaults need to create")
	}

	entities := make([]domain.Vault, 0)
	storageMediaList := make(map[string]struct{})
	storageMedia := ""

	for i := 0; i < len(aux.Vaults); i++ {
		va := aux.Vaults[i]
		entity := domain.NewVault(
			va.Key,
			va.Value,
			domain.WithStorageMedia(domain.GetStorageMedia(va.StorageMedia)),
			domain.WithTags(va.Tags),
			domain.WithVaultType(domain.GetVaultType(va.VaultType), va.TypeIdentity),
			domain.WithExtension(va.Extension),
			domain.WithDescription(va.Description),
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
		filter := []domain.VaultFilter{}
		for i := 0; i < len(aux.Vaults); i++ {
			va := aux.Vaults[i]
			filter = append(filter, domain.VaultFilter{
				Key:          va.Key,
				KeyFuzzy:     false,
				StorageMedia: va.StorageMedia,
				VaultType:    va.VaultType,
				TypeIdentity: va.TypeIdentity,
			})
		}
		query := domain.VaultQuery{
			Filters: filter,
		}
		datas, err := s.repository.SearchVaults(ctx, query)
		if err != nil {
			return nil, err
		}
		if len(datas) > 0 {
			return nil, fmt.Errorf("some vaults are already existed, check again")
		}
	}

	if err := s.innerService.WithUnitOfWork(ctx, func(ctx context.Context) error {
		err := s.repository.InsertMultipleVault(ctx, entities)
		if err != nil {
			return err
		}

		if storageMedia == domain.StorageMediaLocal.String() {
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

	response := viewmodel.CreateVaultsResponse{
		Vaults: []viewmodel.VaultView{},
	}
	for _, va := range entities {
		response.Vaults = append(response.Vaults, convertVaultToVaultView(va))
	}
	return &response, nil
}

func (s *VaultService) CreateVault(ctx context.Context, aux viewmodel.CreateVaultRequest) (*viewmodel.VaultView, error) {
	createVaultsRequest := viewmodel.CreateVaultsRequest{
		Vaults:      []viewmodel.CreateVaultModel{aux.CreateVaultModel},
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

func (s *VaultService) ChangeVault(ctx context.Context, id string, aux viewmodel.ChangeVaultRequest) (*viewmodel.VaultView, error) {
	if tool.IsAllFieldsNil(aux.Data) {
		return nil, fmt.Errorf("no data need change")
	}

	var data *domain.Vault
	query := generateChangeVaultSearchFilter(aux.Data, id)
	datas, err := s.repository.SearchVaults(ctx, query)
	if err != nil {
		return nil, err
	}
	if len(datas) == 0 || (len(datas) == 1 && id != datas[0].ID) {
		return nil, fmt.Errorf("id %s are not existed", id)
	}
	if len(datas) > 1 && !aux.ForceInsert {
		return nil, fmt.Errorf("vaults with 'key+storage_media+vault_type+type_identity' was already existed, check again")
	}
	for _, da := range datas {
		if da.ID == id {
			data = &da
			break
		}
	}

	if data == nil {
		return nil, fmt.Errorf("id %s are not existed", id)
	}

	doVaultChange(data, aux.Data)

	if data.HasChange() {
		if err := s.innerService.WithUnitOfWork(ctx, func(ctx context.Context) error {
			err = s.repository.Update(ctx, *data)
			if err != nil {
				return err
			}

			s.eventPublisher.DispatchIntegrationEvent(ctx, ToVaultChanged(data))

			if data.StorageMedia == domain.StorageMediaLocal {
				return nil
			}

			return s.upsertVaultInProvider(ctx, data.StorageMedia.String(), map[string]string{data.GetIdentityKey(): data.Value})
		}); err != nil {
			return nil, err
		}
	}

	model := convertVaultToVaultView(*data)

	return &model, nil
}

func (s *VaultService) DeleteVault(ctx context.Context, vaultId string) (bool, error) {
	va, err := s.repository.FindByID(ctx, vaultId)
	if err != nil {
		return false, err
	}

	if va == nil {
		return false, fmt.Errorf("vault with id: %s is not exist", vaultId)
	}

	if err := s.innerService.WithUnitOfWork(ctx, func(ctx context.Context) error {
		err := s.repository.Delete(ctx, vaultId)
		if err != nil {
			return err
		}

		if va.StorageMedia == domain.StorageMediaLocal {
			return nil
		}

		return s.deleteVaultInProvider(ctx, va.VaultType.String(), va.GetIdentityKey())
	}); err != nil {
		return false, err
	}

	return true, nil
}

func (s *VaultService) ImportVaults(ctx context.Context, aux viewmodel.ImportVaultsRequest) (*viewmodel.ImportVaultsResponse, error) {
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

	entities := make([]domain.Vault, 0)
	if datas, err := s.searchVaultInProvider(ctx, aux.StorageMedia, fmt.Sprintf("%s/%s", vt, vi)); err != nil {
		return nil, err
	} else {
		for _, data := range datas {
			entities = append(entities, *domain.NewVault(
				data.Key,
				data.Value,
				domain.WithStorageMedia(domain.GetStorageMedia(aux.StorageMedia)),
				domain.WithVaultType(domain.GetVaultType(vt), vi),
			))
		}
	}

	if err := s.innerService.WithUnitOfWork(ctx, func(ctx context.Context) error {
		return s.repository.InsertMultipleVault(ctx, entities)
	}); err != nil {
		return nil, err
	}

	response := viewmodel.ImportVaultsResponse{
		Vaults: []viewmodel.VaultView{},
	}
	for _, va := range entities {
		response.Vaults = append(response.Vaults, convertVaultToVaultView(va))
	}
	return &response, nil
}

func generateChangeVaultSearchFilter(aux viewmodel.ChangeVaultItem, id string) domain.VaultQuery {
	filter := []domain.VaultFilter{{
		ID: id,
	}}

	subFilter := domain.VaultFilter{}
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
	return domain.VaultQuery{
		Filters: filter,
		Page:    nil,
		Size:    nil,
	}
}

func doVaultChange(data *domain.Vault, aux viewmodel.ChangeVaultItem) {
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
		storageMedia := domain.GetStorageMedia(*aux.StorageMedia)
		data.UpdateStorageMedia(storageMedia)
	}

	if aux.VaultType != nil && aux.TypeIdentity != nil {
		vaultType := domain.GetVaultType(*aux.VaultType)
		data.UpdateVaultType(vaultType, *aux.TypeIdentity)
	}

	if aux.Tags != nil {
		data.UpdateTags(*aux.Tags)
	}
}

func convertVaultToVaultView(entity domain.Vault) viewmodel.VaultView {
	return viewmodel.VaultView{
		ID:           entity.ID,
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
	p, err := provider.VaultProviderFactory(provider_type, s.opts)
	if err != nil {
		return err
	}

	return p.Delete(ctx, key)
}

func (s *VaultService) upsertVaultInProvider(ctx context.Context, provider_type string, datas map[string]string) error {
	p, err := provider.VaultProviderFactory(provider_type, s.opts)
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
	p, err := provider.VaultProviderFactory(provider_type, s.opts)
	if err != nil {
		return nil, err
	}

	return p.PrefixSearch(ctx, prefix)
}
