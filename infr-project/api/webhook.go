package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"reflect"

	_ "github.com/joho/godotenv/autoload"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/futugyou/extensions"

	"github.com/futugyou/infr-project/application"
	"github.com/futugyou/infr-project/controller"
	tool "github.com/futugyou/infr-project/extensions"
	infra "github.com/futugyou/infr-project/infrastructure_mongo"
	publisher "github.com/futugyou/infr-project/infrastructure_qstash"
	screenshot "github.com/futugyou/infr-project/infrastructure_screenshot"
	models "github.com/futugyou/infr-project/view_models"
)

func WebhookDispatch(w http.ResponseWriter, r *http.Request) {
	// cors
	if extensions.Cors(w, r) {
		return
	}

	op := r.URL.Query().Get("optype")
	version := r.URL.Query().Get("version")
	if len(version) == 0 {
		version = "v1"
	}

	ctrl := controller.NewController()
	switch version {
	case "v1":
		switch op {
		case "event":
			eventHandler(ctrl, r, w)
		case "webhook":
			handleWebhook(ctrl, r, w)
		case "qstash":
			handleQstash(ctrl, r, w)
		default:
			w.Write([]byte("page not found"))
			w.WriteHeader(404)
		}
	case "v2":
		switch op {
		case "webhook":
			handleWebhook(ctrl, r, w)
		default:
			w.Write([]byte("page not found"))
			w.WriteHeader(404)
		}
	default:
		w.Write([]byte("page not found"))
		w.WriteHeader(404)
	}
}

func eventHandler(_ *controller.Controller, r *http.Request, w http.ResponseWriter) {
	ctx := r.Context()
	bearer := r.Header.Get("Authorization")
	if bearer != os.Getenv("TRIGGER_AUTH_KEY") {
		w.Write([]byte("Authorization code error"))
		w.WriteHeader(500)
		return
	}
	var event models.TriggerEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	service, err := createResourceQueryService(ctx)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	dataType := getDataType(event.TableName)
	if dataType == nil {
		w.Write([]byte("can not find data type"))
		w.WriteHeader(500)
		return
	}

	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	dataInstance := reflect.New(dataType).Interface()

	if err := json.Unmarshal(dataBytes, dataInstance); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(500)
		return
	}

	if resourceData, ok := dataInstance.(*models.ResourceChangeData); ok {
		if err = service.HandleResourceChanged(ctx, *resourceData); err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(500)
			return
		}
	} else {
		w.Write([]byte("can not find event handler"))
		w.WriteHeader(200)
	}
}

func createResourceQueryService(ctx context.Context) (*application.ResourceQueryService, error) {
	config := infra.DBConfig{
		DBName:        os.Getenv("query_db_name"),
		ConnectString: os.Getenv("query_mongodb_url"),
	}

	mongoclient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectString))
	if err != nil {
		return nil, err
	}

	client, err := tool.RedisClient(os.Getenv("REDIS_URL"))
	if err != nil {
		return nil, err
	}

	queryRepo := infra.NewResourceQueryRepository(mongoclient, config)

	unitOfWork, err := infra.NewMongoUnitOfWork(mongoclient)
	if err != nil {
		return nil, err
	}

	return application.NewResourceQueryService(queryRepo, client, unitOfWork), nil
}

func getDataType(tableName string) reflect.Type {
	switch tableName {
	case "resource_events":
		return reflect.TypeOf(models.ResourceChangeData{})
	default:
		return nil
	}
}

func handleWebhook(_ *controller.Controller, r *http.Request, w http.ResponseWriter) {
	ctrl := controller.NewWebhookController()
	ctrl.ProviderWebhookCallback(w, r)
}

func handleQstash(_ *controller.Controller, r *http.Request, w http.ResponseWriter) {
	query := r.URL.Query()
	if len(query["event"]) == 0 {
		w.WriteHeader(200)
		return
	}

	ctx := r.Context()
	var err error
	bodyBytes, _ := io.ReadAll(r.Body)
	defer r.Body.Close()
	event := query["event"][0]

	switch event {
	case "upsert_project":
		var data models.PlatformProjectUpsertEvent
		if err = json.Unmarshal(bodyBytes, &data); err != nil {
			w.WriteHeader(500)
			return
		}

		service, err := createPlatformService(ctx)
		if err != nil {
			w.WriteHeader(500)
			return
		}

		if err = service.HandlePlatformProjectUpsert(ctx, data); err != nil {
			w.WriteHeader(500)
			return
		}
	case "vault_changed":
	}

	w.WriteHeader(200)
}

func createPlatformService(ctx context.Context) (*application.PlatformService, error) {
	config := infra.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectString))
	if err != nil {
		return nil, err
	}

	unitOfWork, err := infra.NewMongoUnitOfWork(client)
	if err != nil {
		return nil, err
	}

	redisClient, err := tool.RedisClient(os.Getenv("REDIS_URL"))
	if err != nil {
		return nil, err
	}

	repo := infra.NewPlatformRepository(client, config)
	vaultRepo := infra.NewVaultRepository(client, config)
	eventPublisher := publisher.NewQStashEventPulisher(os.Getenv("QSTASH_TOKEN"), os.Getenv("QSTASH_DESTINATION"))
	vaultService := application.NewVaultService(unitOfWork, vaultRepo, eventPublisher)
	ss := screenshot.NewScreenshot()
	return application.NewPlatformService(unitOfWork, repo, vaultService, redisClient, eventPublisher, ss), nil
}
