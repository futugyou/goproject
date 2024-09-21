package application

import (
	"context"
	"fmt"

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
	src, err := s.repository.GetAllVaultAsync(ctx, page, size)
	select {
	case datas := <-src:
		result := make([]models.VaultView, len(datas))
		for i := 0; i < len(datas); i++ {
			result[i] = models.VaultView{
				Id:           datas[i].Id,
				Key:          datas[i].Key,
				MaskValue:    tool.MaskString(datas[i].Value, 5, 0.5),
				StorageMedia: datas[i].StorageMedia.String(),
				VaultType:    datas[i].VaultType.String(),
				TypeIdentity: datas[i].TypeIdentity,
				Tags:         datas[i].Tags,
			}
		}
		return result, nil
	case errM := <-err:
		return nil, errM
	case <-ctx.Done():
		return nil, fmt.Errorf("GetAllVault timeout")
	}
}
