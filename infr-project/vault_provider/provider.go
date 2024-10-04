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
	IVaultProvider
}

type AsyncWrapper struct {
	IVaultProvider
}

func NewAsyncWrapper(provider IVaultProvider) *AsyncWrapper {
	return &AsyncWrapper{provider}
}

func (s *AsyncWrapper) SearchAsync(ctx context.Context, key string) (<-chan *ProviderVault, <-chan error) {
	resultChan := make(chan *ProviderVault, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := s.Search(ctx, key)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()
	return resultChan, errorChan
}

func (s *AsyncWrapper) PrefixSearchAsync(ctx context.Context, prefix string) (<-chan map[string]ProviderVault, <-chan error) {
	resultChan := make(chan map[string]ProviderVault, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := s.PrefixSearch(ctx, prefix)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()
	return resultChan, errorChan
}

func (s *AsyncWrapper) BatchSearchAsync(ctx context.Context, keys []string) (<-chan map[string]ProviderVault, <-chan error) {
	resultChan := make(chan map[string]ProviderVault, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := s.BatchSearch(ctx, keys)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()
	return resultChan, errorChan
}

func (s *AsyncWrapper) UpsertAsync(ctx context.Context, key string, value string) (<-chan *ProviderVault, <-chan error) {
	resultChan := make(chan *ProviderVault, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := s.Upsert(ctx, key, value)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()
	return resultChan, errorChan
}

func (s *AsyncWrapper) DeleteAsync(ctx context.Context, key string) <-chan error {
	errorChan := make(chan error, 1)
	go func() {
		defer close(errorChan)
		errorChan <- s.Delete(ctx, key)
	}()
	return errorChan
}
