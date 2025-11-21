package controller

import (
	"context"
	"net/http"

	"github.com/futugyou/domaincore/mongoimpl"
	"github.com/futugyou/domaincore/qstashdispatcherimpl"

	"github.com/futugyou/vaultservice/application"
	"github.com/futugyou/vaultservice/infrastructure"
	"github.com/futugyou/vaultservice/options"

	"github.com/futugyou/vaultservice/viewmodel"
)

type VaultController struct {
}

func NewVaultController() *VaultController {
	return &VaultController{}
}

func (c *VaultController) CreateVaults(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createVaultService, func(ctx context.Context, service *application.VaultService, req viewmodel.CreateVaultsRequest) (any, error) {
		return service.CreateVaults(ctx, req)
	})
}

func (c *VaultController) CreateVault(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createVaultService, func(ctx context.Context, service *application.VaultService, req viewmodel.CreateVaultRequest) (any, error) {
		return service.CreateVault(ctx, req)
	})
}

func (c *VaultController) ShowVaultRawValue(w http.ResponseWriter, r *http.Request, vaultId string) {
	handleRequest(w, r, createVaultService, func(ctx context.Context, service *application.VaultService, _ struct{}) (any, error) {
		return service.ShowVaultRawValue(ctx, vaultId)
	})
}

func (c *VaultController) SearchVaults(w http.ResponseWriter, r *http.Request, aux viewmodel.SearchVaultsRequest) {
	handleRequest(w, r, createVaultService, func(ctx context.Context, service *application.VaultService, _ struct{}) (any, error) {
		return service.SearchVaults(ctx, aux)
	})
}

func (c *VaultController) ChangeVault(w http.ResponseWriter, r *http.Request, vaultId string) {
	handleRequest(w, r, createVaultService, func(ctx context.Context, service *application.VaultService, req viewmodel.ChangeVaultRequest) (any, error) {
		return nil, service.ChangeVault(ctx, vaultId, req)
	})
}

func (c *VaultController) DeleteVault(w http.ResponseWriter, r *http.Request, vaultId string) {
	handleRequest(w, r, createVaultService, func(ctx context.Context, service *application.VaultService, _ struct{}) (any, error) {
		return nil, service.DeleteVault(ctx, vaultId)
	})
}

func (c *VaultController) ImportVaults(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createVaultService, func(ctx context.Context, service *application.VaultService, req viewmodel.ImportVaultsRequest) (any, error) {
		return service.ImportVaults(ctx, req)
	})
}

func (c *VaultController) GetVaultsByIDs(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createVaultService, func(ctx context.Context, service *application.VaultService, req viewmodel.SearchVaultsByIDsRequest) (any, error) {
		return service.GetVaultsByIDs(ctx, req.IDs)
	})
}

func createVaultService(ctx context.Context) (*application.VaultService, error) {
	option := options.New()
	mongoclient, err := mongoimpl.CreateMongoDBClient(ctx, option.MongoDBURL)
	config := mongoimpl.DBConfig{
		DBName: option.DBName,
	}

	if err != nil {
		return nil, err
	}

	queryRepo := infrastructure.NewVaultRepository(mongoclient, config)

	unitOfWork, err := mongoimpl.NewMongoUnitOfWork(mongoclient)
	if err != nil {
		return nil, err
	}

	eventPublisher := qstashdispatcherimpl.NewQStashEventDispatcher(option.QstashToken, option.QstashDestination)

	return application.NewVaultService(unitOfWork, queryRepo, eventPublisher, option), nil
}
