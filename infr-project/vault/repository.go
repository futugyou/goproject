package vault

import (
	"context"

	"github.com/futugyou/infr-project/domain"
)

type VaultSearch struct {
	ID           string
	Key          string
	KeyFuzzy     bool
	StorageMedia string
	VaultType    string
	TypeIdentity string
	Description  string
	Tags         []string
}

type IVaultRepositoryAsync interface {
	domain.IRepositoryAsync[Vault]
	InsertMultipleVaultAsync(ctx context.Context, vaults []Vault) <-chan error
	GetVaultByIdsAsync(ctx context.Context, ids []string) (<-chan []Vault, <-chan error)
	SearchVaultsAsync(ctx context.Context, filter []VaultSearch, page *int, size *int) (<-chan []Vault, <-chan error)
}
