package vault

import (
	"context"

	"github.com/futugyou/infr-project/domain"
)

type IVaultRepositoryAsync interface {
	domain.IRepositoryAsync[Vault]
	InsertMultipleVault(ctx context.Context, vaults []Vault) <-chan error
	GetAllVaultAsync(ctx context.Context, page *int, size *int) (<-chan []Vault, <-chan error)
	GetAllVaultByStorageMediaAsync(ctx context.Context, media StorageMedia) (<-chan []Vault, <-chan error)
	GetAllVaultByVaultTypeAsync(ctx context.Context, vType VaultType, identities ...string) (<-chan []Vault, <-chan error)
	GetAllVaultByTagsAsync(ctx context.Context, tags []string) (<-chan []Vault, <-chan error)
	GetAllVaultByIdsAsync(ctx context.Context, ids []string) (<-chan []Vault, <-chan error)
}
