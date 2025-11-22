package controller

import (
	"context"
	"net/http"

	"github.com/futugyou/vaultservice/application"
	v1 "github.com/futugyou/vaultservice/routes/v1"
	"github.com/futugyou/vaultservice/viewmodel"
)

type VaultController struct {
}

func NewVaultController() *VaultController {
	return &VaultController{}
}

func (c *VaultController) CreateVaults(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreateVaultService, func(ctx context.Context, service *application.VaultService, req viewmodel.CreateVaultsRequest) (any, error) {
		return service.CreateVaults(ctx, req)
	})
}

func (c *VaultController) CreateVault(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreateVaultService, func(ctx context.Context, service *application.VaultService, req viewmodel.CreateVaultRequest) (any, error) {
		return service.CreateVault(ctx, req)
	})
}

func (c *VaultController) ShowVaultRawValue(w http.ResponseWriter, r *http.Request, vaultId string) {
	handleRequest(w, r, v1.CreateVaultService, func(ctx context.Context, service *application.VaultService, _ struct{}) (any, error) {
		return service.ShowVaultRawValue(ctx, vaultId)
	})
}

func (c *VaultController) SearchVaults(w http.ResponseWriter, r *http.Request, aux viewmodel.SearchVaultsRequest) {
	handleRequest(w, r, v1.CreateVaultService, func(ctx context.Context, service *application.VaultService, _ struct{}) (any, error) {
		return service.SearchVaults(ctx, aux)
	})
}

func (c *VaultController) ChangeVault(w http.ResponseWriter, r *http.Request, vaultId string) {
	handleRequest(w, r, v1.CreateVaultService, func(ctx context.Context, service *application.VaultService, req viewmodel.ChangeVaultRequest) (any, error) {
		return nil, service.ChangeVault(ctx, vaultId, req)
	})
}

func (c *VaultController) DeleteVault(w http.ResponseWriter, r *http.Request, vaultId string) {
	handleRequest(w, r, v1.CreateVaultService, func(ctx context.Context, service *application.VaultService, _ struct{}) (any, error) {
		return nil, service.DeleteVault(ctx, vaultId)
	})
}

func (c *VaultController) ImportVaults(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreateVaultService, func(ctx context.Context, service *application.VaultService, req viewmodel.ImportVaultsRequest) (any, error) {
		return service.ImportVaults(ctx, req)
	})
}

func (c *VaultController) GetVaultsByIDs(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, v1.CreateVaultService, func(ctx context.Context, service *application.VaultService, req viewmodel.SearchVaultsByIDsRequest) (any, error) {
		return service.GetVaultsByIDs(ctx, req.IDs)
	})
}
