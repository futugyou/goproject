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
	models "github.com/futugyou/infr-project/view_models"
)

func (c *Controller) CreateVaults(w http.ResponseWriter, r *http.Request) {
	service, err := createVaultService(r.Context())
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

	res, err := service.CreateVaults(aux, r.Context())
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *Controller) SearchVaults(w http.ResponseWriter, r *http.Request, aux models.SearchVaultsRequest, page *int, size *int) {
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
	res, err := service.SearchVaults(ctx, aux, page, size)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *Controller) ShowVaultRawValue(w http.ResponseWriter, r *http.Request, vaultId string) {
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

func (c *Controller) ChangeVault(w http.ResponseWriter, r *http.Request, vaultId string) {
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

	res, err := service.ChangeVault(vaultId, aux, ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	writeJSONResponse(w, res, 200)
}

func (c *Controller) DeleteVault(w http.ResponseWriter, r *http.Request, vaultId string) {
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
