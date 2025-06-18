package api

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

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
	"github.com/futugyou/infr-project/resource"
	models "github.com/futugyou/infr-project/view_models"
)

func WebhookDispatch(w http.ResponseWriter, r *http.Request) {
	// Handle CORS preflight
	if extensions.Cors(w, r) {
		return
	}

	op := r.URL.Query().Get("optype")
	version := r.URL.Query().Get("version")
	if version == "" {
		version = "v1"
	}

	ctrl := controller.NewController()

	switch version {
	case "v1":
		handleV1(op, ctrl, r, w)
	case "v2":
		handleV2(op, ctrl, r, w)
	default:
		writeNotFound(w)
	}
}

func handleV1(op string, ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	switch op {
	case "event":
		eventHandler(ctrl, r, w)
	case "webhook":
		handleWebhook(ctrl, r, w)
	case "qstash":
		handleQstash(ctrl, r, w)
	default:
		writeNotFound(w)
	}
}

func handleV2(op string, ctrl *controller.Controller, r *http.Request, w http.ResponseWriter) {
	switch op {
	case "webhook":
		handleWebhook(ctrl, r, w)
	default:
		writeNotFound(w)
	}
}

func writeNotFound(w http.ResponseWriter) {
	http.Error(w, "page not found", http.StatusNotFound)
}

func writeInternalServerError(w http.ResponseWriter, msg string, err error) {
	log.Printf("%s: %v", msg, err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func eventHandler(_ *controller.Controller, r *http.Request, w http.ResponseWriter) {
	ctx := r.Context()
	if r.Header.Get("Authorization") != os.Getenv("TRIGGER_AUTH_KEY") {
		http.Error(w, "Authorization code error", http.StatusInternalServerError)
		return
	}

	var event models.TriggerEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeInternalServerError(w, "decode error", err)
		return
	}

	service, err := createResourceQueryService(ctx)
	if err != nil {
		writeInternalServerError(w, "service creation error", err)
		return
	}

	dataType := getDataType(event.TableName)
	if dataType == nil {
		http.Error(w, "cannot find data type", http.StatusInternalServerError)
		return
	}

	// Decode `event.Data` dynamically
	dataBytes, _ := json.Marshal(event.Data)
	dataInstance := reflect.New(dataType).Interface()
	if err := json.Unmarshal(dataBytes, dataInstance); err != nil {
		writeInternalServerError(w, "unmarshal dynamic data error", err)
		return
	}

	if resourceData, ok := dataInstance.(*models.ResourceChangeData); ok {
		if err := service.HandleResourceChanged(ctx, *resourceData); err != nil {
			writeInternalServerError(w, "handle resource changed error", err)
		}
	} else {
		// Unknown event data structure
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("can not find event handler"))
	}
}

func handleWebhook(_ *controller.Controller, r *http.Request, w http.ResponseWriter) {
	ctrl := controller.NewWebhookController()
	ctrl.ProviderWebhookCallback(w, r)
}

func handleQstash(_ *controller.Controller, r *http.Request, w http.ResponseWriter) {
	query := r.URL.Query()
	event := query.Get("event")
	if event == "" {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()
	bodyBytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		writeInternalServerError(w, "read body error", err)
		return
	}

	// Handle events via typed switch
	switch event {
	case "upsert_project":
		handleEvent(bodyBytes, w, func(data []byte) error {
			var evt models.PlatformProjectUpsertEvent
			if err := json.Unmarshal(data, &evt); err != nil {
				return err
			}
			service, err := createPlatformService(ctx)
			if err != nil {
				return err
			}
			return service.HandlePlatformProjectUpsert(ctx, evt)
		})

	case "ResourceCreated":
		handleEvent(bodyBytes, w, func(data []byte) error {
			var evt resource.ResourceCreatedEvent
			if err := json.Unmarshal(data, &evt); err != nil {
				return err
			}
			service, err := createResourceQueryService(ctx)
			if err != nil {
				return err
			}
			changeData := convertToChangeData(evt.Id, evt.ResourceVersion, event, evt.CreatedAt, evt.Name, evt.Type, evt.Data, evt.ImageData, evt.Tags)
			return service.HandleResourceChanged(ctx, changeData)
		})

	case "ResourceUpdated":
		handleEvent(bodyBytes, w, func(data []byte) error {
			var evt resource.ResourceUpdatedEvent
			if err := json.Unmarshal(data, &evt); err != nil {
				return err
			}
			service, err := createResourceQueryService(ctx)
			if err != nil {
				return err
			}
			changeData := convertToChangeData(evt.Id, evt.ResourceVersion, event, evt.CreatedAt, evt.Name, evt.Type, evt.Data, evt.ImageData, evt.Tags)
			return service.HandleResourceChanged(ctx, changeData)
		})

	case "ResourceDeleted":
		handleEvent(bodyBytes, w, func(data []byte) error {
			var evt resource.ResourceDeletedEvent
			if err := json.Unmarshal(data, &evt); err != nil {
				return err
			}
			service, err := createResourceQueryService(ctx)
			if err != nil {
				return err
			}
			return service.HandleResourceChanged(ctx, models.ResourceChangeData{
				Id:              evt.Id,
				ResourceVersion: evt.ResourceVersion,
				EventType:       event,
				CreatedAt:       evt.CreatedAt,
			})
		})

	case "vault_changed":
		// Reserved for future use
	default:
		http.Error(w, "unknown event", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleEvent(data []byte, w http.ResponseWriter, handler func([]byte) error) {
	if err := handler(data); err != nil {
		writeInternalServerError(w, "event handler error", err)
	}
}

// createResourceQueryService constructs ResourceQueryService
func createResourceQueryService(ctx context.Context) (*application.ResourceQueryService, error) {
	config := infra.DBConfig{
		DBName:        os.Getenv("query_db_name"),
		ConnectString: os.Getenv("query_mongodb_url"),
	}

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectString))
	if err != nil {
		return nil, err
	}

	redisClient, err := tool.RedisClient(os.Getenv("REDIS_URL"))
	if err != nil {
		return nil, err
	}

	queryRepo := infra.NewResourceQueryRepository(mongoClient, config)
	unitOfWork, err := infra.NewMongoUnitOfWork(mongoClient)
	if err != nil {
		return nil, err
	}

	return application.NewResourceQueryService(queryRepo, redisClient, unitOfWork), nil
}

func getDataType(tableName string) reflect.Type {
	switch tableName {
	case "resource_events":
		return reflect.TypeOf(models.ResourceChangeData{})
	default:
		return nil
	}
}

// convertToChangeData creates ResourceChangeData from fields
func convertToChangeData(id string, version int, eventType string, createdAt time.Time, name, typ, data, imageData string, tags []string) models.ResourceChangeData {
	return models.ResourceChangeData{
		Id:              id,
		ResourceVersion: version,
		EventType:       eventType,
		CreatedAt:       createdAt,
		Name:            name,
		Type:            typ,
		Data:            data,
		ImageData:       imageData,
		Tags:            tags,
	}
}

// createPlatformService constructs PlatformService
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
	eventPublisher := publisher.NewQStashEventPulisher(
		os.Getenv("QSTASH_TOKEN"),
		os.Getenv("QSTASH_DESTINATION"),
	)

	vaultService := application.NewVaultService(unitOfWork, vaultRepo, eventPublisher)
	ss := screenshot.NewScreenshot()

	return application.NewPlatformService(unitOfWork, repo, vaultService, redisClient, eventPublisher, ss), nil
}
