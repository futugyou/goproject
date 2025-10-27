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

type IVaultRepository interface {
	domain.IRepository[Vault]
	InsertMultipleVault(ctx context.Context, vaults []Vault) error
	GetVaultByIds(ctx context.Context, ids []string) ([]Vault, error)
	SearchVaults(ctx context.Context, filter []VaultSearch, page *int, size *int) ([]Vault, error)
}
