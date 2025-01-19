package controller

import (
	"context"
	"io"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/futugyou/infr-project/application"
	infra "github.com/futugyou/infr-project/infrastructure_mongo"
)

type WebhookController struct {
}

func NewWebhookController() *WebhookController {
	return &WebhookController{}
}

func (c *WebhookController) ProviderWebhookCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createWebhookService(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	bodyBytes, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	query := r.URL.Query()

	reqInfo := application.WebhookRequestInfo{
		Method:     r.Method,
		URL:        r.URL.String(),
		Proto:      r.Proto,
		Host:       r.Host,
		Header:     r.Header,
		Body:       string(bodyBytes),
		Query:      query,
		RemoteAddr: r.RemoteAddr,
		UserAgent:  r.UserAgent(),
	}

	service.ProviderWebhookCallback(ctx, reqInfo)
	writeJSONResponse(w, nil, 200)
}

func (c *WebhookController) VerifyTesting(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	service, err := createWebhookService(ctx)
	if err != nil {
		handleError(w, err, 500)
		return
	}

	if data, err := service.VerifyTesting(ctx); err != nil {
		handleError(w, err, 500)
		return
	} else {
		writeJSONResponse(w, data, 200)
	}

}

func createWebhookService(ctx context.Context) (*application.WebhookService, error) {
	config := infra.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectString))
	if err != nil {
		return nil, err
	}

	return application.NewWebhookService(client.Database(config.DBName)), nil
}
