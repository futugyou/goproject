package vault

import (
	"context"

	"github.com/futugyou/infr-project/domain"
)

type IVaultRepositoryAsync interface {
	domain.IRepositoryAsync[Vault]
	InsertMultipleVaultAsync(ctx context.Context, vaults []Vault) <-chan error
	GetAllVaultAsync(ctx context.Context, page *int, size *int) (<-chan []Vault, <-chan error)
	GetVaultByStorageMediaAsync(ctx context.Context, media StorageMedia) (<-chan []Vault, <-chan error)
	GetVaultByVaultTypeAsync(ctx context.Context, vType VaultType, identities ...string) (<-chan []Vault, <-chan error)
	GetVaultByTagsAsync(ctx context.Context, tags []string) (<-chan []Vault, <-chan error)
	GetVaultByIdsAsync(ctx context.Context, ids []string) (<-chan []Vault, <-chan error)
}
