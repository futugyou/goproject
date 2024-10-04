package vault_provider

import (
	"context"
	"time"
)

type ProviderVault struct {
	Key       string
	Value     string
	CreatedAt time.Time
}

type IVaultProvider interface {
	Search(ctx context.Context, key string) (*ProviderVault, error)
	PrefixSearch(ctx context.Context, prefix string) (map[string]ProviderVault, error)
	BatchSearch(ctx context.Context, keys []string) (map[string]ProviderVault, error)
	Upsert(ctx context.Context, key string, value string) (*ProviderVault, error)
	Delete(ctx context.Context, key string) error
}

type IVaultProviderAsync interface {
	SearchAsync(ctx context.Context, key string) (<-chan *ProviderVault, <-chan error)
	PrefixSearchAsync(ctx context.Context, prefix string) (<-chan map[string]ProviderVault, <-chan error)
	BatchSearchAsync(ctx context.Context, keys []string) (<-chan map[string]ProviderVault, <-chan error)
	UpsertAsync(ctx context.Context, key string, value string) (<-chan *ProviderVault, <-chan error)
	DeleteAsync(ctx context.Context, key string) <-chan error
}
