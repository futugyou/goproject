package controller

import (
	"context"
	"io"
	"net/http"

	"github.com/futugyou/domaincore/mongoimpl"
	"github.com/futugyou/platformservice/application"
	"github.com/futugyou/platformservice/infrastructure"
	"github.com/futugyou/platformservice/options"
	"github.com/futugyou/platformservice/viewmodel"
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
	reqInfo := viewmodel.WebhookRequestInfo{
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

	handleRequest(w, r, createWebhookService, func(ctx context.Context, service *application.WebhookLogService, _ struct{}) (any, error) {
		return nil, service.ProviderWebhookCallback(ctx, reqInfo)
	})
}

func (c *WebhookController) SearchWebhookLogs(w http.ResponseWriter, r *http.Request, filter viewmodel.WebhookSearch) {
	handleRequest(w, r, createWebhookService, func(ctx context.Context, service *application.WebhookLogService, _ struct{}) (any, error) {
		return service.SearchWebhookLogs(ctx, filter)
	})
}

func createWebhookService(ctx context.Context) (*application.WebhookLogService, error) {
	option := options.New()
	mongoclient, err := mongoimpl.CreateMongoDBClient(ctx, option.MongoDBURL)
	config := mongoimpl.DBConfig{
		DBName: option.DBName,
	}

	if err != nil {
		return nil, err
	}

	repo := infrastructure.NewWebhookLogRepository(mongoclient, config)

	return application.NewWebhookLogService(repo, option), nil
}
