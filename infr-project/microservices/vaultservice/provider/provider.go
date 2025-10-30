package provider

import (
	"context"
	"time"
)

type ProviderVault struct {
	Key       string
	Value     string
	CreatedAt time.Time
}

type VaultProvider interface {
	Search(ctx context.Context, key string) (*ProviderVault, error)
	PrefixSearch(ctx context.Context, prefix string) (map[string]ProviderVault, error)
	BatchSearch(ctx context.Context, keys []string) (map[string]ProviderVault, error)
	Upsert(ctx context.Context, key string, value string) (*ProviderVault, error)
	Delete(ctx context.Context, key string) error
}
