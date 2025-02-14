package controller

import (
	"context"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/futugyou/infr-project/application"
	infra "github.com/futugyou/infr-project/infrastructure_mongo"
	publisher "github.com/futugyou/infr-project/infrastructure_qstash"
	"github.com/futugyou/infr-project/vault"
	models "github.com/futugyou/infr-project/view_models"
)

type VaultController struct {
}

func NewVaultController() *VaultController {
	return &VaultController{}
}

func (c *VaultController) CreateVaults(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createVaultService, func(ctx context.Context, service *application.VaultService, req models.CreateVaultsRequest) (interface{}, error) {
		return service.CreateVaults(ctx, req)
	})
}

func (c *VaultController) CreateVault(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createVaultService, func(ctx context.Context, service *application.VaultService, req models.CreateVaultRequest) (interface{}, error) {
		return service.CreateVault(ctx, req)
	})
}

func (c *VaultController) SearchVaults(w http.ResponseWriter, r *http.Request, aux models.SearchVaultsRequest) {
	handleRequest(w, r, createVaultService, func(ctx context.Context, service *application.VaultService, _ struct{}) (interface{}, error) {
		query := application.VaultSearchQuery{
			Filters: []vault.VaultSearch{
				{
					Key:          aux.Key,
					KeyFuzzy:     true,
					StorageMedia: aux.StorageMedia,
					VaultType:    aux.VaultType,
					TypeIdentity: aux.TypeIdentity,
					Description:  aux.Description,
					Tags:         aux.Tags,
				},
			},
			Page: aux.Page,
			Size: aux.Size,
		}
		return service.SearchVaults(ctx, query)
	})
}

func (c *VaultController) ShowVaultRawValue(w http.ResponseWriter, r *http.Request, vaultId string) {
	handleRequest(w, r, createVaultService, func(ctx context.Context, service *application.VaultService, _ struct{}) (interface{}, error) {
		return service.ShowVaultRawValue(ctx, vaultId)
	})
}

func (c *VaultController) ChangeVault(w http.ResponseWriter, r *http.Request, vaultId string) {
	handleRequest(w, r, createVaultService, func(ctx context.Context, service *application.VaultService, req models.ChangeVaultRequest) (interface{}, error) {
		return service.ChangeVault(ctx, vaultId, req)
	})
}

func (c *VaultController) DeleteVault(w http.ResponseWriter, r *http.Request, vaultId string) {
	handleRequest(w, r, createVaultService, func(ctx context.Context, service *application.VaultService, _ struct{}) (interface{}, error) {
		return service.DeleteVault(ctx, vaultId)
	})
}

func (c *VaultController) ImportVaults(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createVaultService, func(ctx context.Context, service *application.VaultService, req models.ImportVaultsRequest) (interface{}, error) {
		return service.ImportVaults(ctx, req)
	})
}

func createVaultService(ctx context.Context) (*application.VaultService, error) {
	config := infra.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectString))
	if err != nil {
		return nil, err
	}

	repo := infra.NewVaultRepository(client, config)
	unitOfWork, err := infra.NewMongoUnitOfWork(client)
	if err != nil {
		return nil, err
	}

	eventPublisher := publisher.NewQStashEventPulisher(os.Getenv("QSTASH_TOKEN"), os.Getenv("QSTASH_DESTINATION"))
	return application.NewVaultService(unitOfWork, repo, eventPublisher), nil
}
