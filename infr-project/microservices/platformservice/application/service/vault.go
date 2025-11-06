package service

import "context"

type VaultService interface {
	GetVaultsByIDs(ctx context.Context, ids []string) ([]VaultView, error)
	ShowVaultRawValue(ctx context.Context, vaultId string) (string, error)
	CreateVault(ctx context.Context, aux CreateVaultRequest) (*VaultView, error)
}
