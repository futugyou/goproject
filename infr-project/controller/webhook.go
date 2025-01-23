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
	models "github.com/futugyou/infr-project/view_models"
)

type WebhookController struct {
}

func NewWebhookController() *WebhookController {
	return &WebhookController{}
}

func (c *WebhookController) ProviderWebhookCallback(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	query := r.URL.Query()
	reqInfo := models.WebhookRequestInfo{
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

	handleRequest(w, r, createWebhookService, func(ctx context.Context, service *application.WebhookService, _ struct{}) (interface{}, error) {
		return nil, service.ProviderWebhookCallback(ctx, reqInfo)
	})
}

func (c *WebhookController) VerifyTesting(w http.ResponseWriter, r *http.Request) {
	handleRequest(w, r, createWebhookService, func(ctx context.Context, service *application.WebhookService, _ struct{}) (interface{}, error) {
		return service.VerifyTesting(ctx)
	})
}

func (c *WebhookController) SearchWebhookLogs(w http.ResponseWriter, r *http.Request, filter models.WebhookSearch) {
	handleRequest(w, r, createWebhookService, func(ctx context.Context, service *application.WebhookService, _ struct{}) (interface{}, error) {
		return service.SearchWebhookLogs(ctx, filter)
	})
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

	repo := infra.NewWebhookLogRepository(client, config)
	return application.NewWebhookService(repo), nil
}
