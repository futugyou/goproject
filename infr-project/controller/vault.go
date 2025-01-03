package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/futugyou/infr-project/application"
	"github.com/futugyou/infr-project/extensions"
	infra "github.com/futugyou/infr-project/infrastructure_mongo"
	"github.com/futugyou/infr-project/vault"
	models "github.com/futugyou/infr-project/view_models"
)

type VaultController struct {
}

func NewVaultController() *VaultController {
	return &VaultController{}
}

func (c *VaultController) CreateVaults(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createVaultService(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	var aux models.CreateVaultsRequest
	if err := json.NewDecoder(r.Body).Decode(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	if err := extensions.Validate.Struct(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	res, err := service.CreateVaults(ctx, aux)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *VaultController) SearchVaults(w http.ResponseWriter, r *http.Request, aux models.SearchVaultsRequest) {
	ctx := r.Context()
	service, err := createVaultService(ctx)

	if err != nil {
		handleError(w, err, 500)
		return
	}
	if err := extensions.Validate.Struct(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	query := application.VaultSearchQuery{
		Filters: []vault.VaultSearch{
			{
				Key:          aux.Key,
				KeyFuzzy:     true,
				StorageMedia: aux.StorageMedia,
				VaultType:    aux.VaultType,
				TypeIdentity: aux.TypeIdentity,
				Tags:         aux.Tags,
			},
		},
		Page: aux.Page,
		Size: aux.Size,
	}
	res, err := service.SearchVaults(ctx, query)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *VaultController) ShowVaultRawValue(w http.ResponseWriter, r *http.Request, vaultId string) {
	ctx := r.Context()
	service, err := createVaultService(ctx)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	res, err := service.ShowVaultRawValue(ctx, vaultId)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *VaultController) ChangeVault(w http.ResponseWriter, r *http.Request, vaultId string) {
	ctx := r.Context()
	service, err := createVaultService(ctx)

	if err != nil {
		handleError(w, err, 500)
		return
	}
	var aux models.ChangeVaultRequest
	if err := json.NewDecoder(r.Body).Decode(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	res, err := service.ChangeVault(ctx, vaultId, aux)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *VaultController) DeleteVault(w http.ResponseWriter, r *http.Request, vaultId string) {
	ctx := r.Context()
	service, err := createVaultService(ctx)

	if err != nil {
		handleError(w, err, 500)
		return
	}

	res, err := service.DeleteVault(ctx, vaultId)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *VaultController) ImportVaults(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createVaultService(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	var aux models.ImportVaultsRequest
	if err := json.NewDecoder(r.Body).Decode(&aux); err != nil {
		handleError(w, err, 400)
		return
	}

	res, err := service.ImportVaults(ctx, aux)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
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

	return application.NewVaultService(unitOfWork, repo), nil
}
