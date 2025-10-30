package domain

import (
	"context"

	"github.com/futugyou/domaincore/domain"
)

type VaultQuery struct {
	Filters []VaultFilter
	Page    *int
	Size    *int
}

type VaultFilter struct {
	ID           string
	Key          string
	KeyFuzzy     bool
	StorageMedia string
	VaultType    string
	TypeIdentity string
	Description  string
	Tags         []string
}

type VaultRepository interface {
	domain.Repository[Vault]
	InsertMultipleVault(ctx context.Context, vaults []Vault) error
	GetVaultByIds(ctx context.Context, ids []string) ([]Vault, error)
	SearchVaults(ctx context.Context, query VaultQuery) ([]Vault, error)
}
