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
	Get(ctx context.Context, key string) (*ProviderVault, error)
	Search(ctx context.Context, prefix string) ([]ProviderVault, error)
	Upinsert(ctx context.Context, key string, value string) (*ProviderVault, error)
	Delete(ctx context.Context, key string) error
}
