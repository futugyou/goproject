package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/futugyou/extensions"
	platv1 "github.com/futugyou/platformservice/routes/v1"
	platviewmodel "github.com/futugyou/platformservice/viewmodel"
	v1 "github.com/futugyou/resourcequeryservice/routes/v1"
	"github.com/futugyou/resourcequeryservice/viewmodel"

	resource "github.com/futugyou/resourceservice/domain"

	"github.com/futugyou/infr-project/controller"
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

	var event TriggerEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeInternalServerError(w, "decode error", err)
		return
	}

	service, err := v1.CreateResourceQueryService(ctx)
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

	if resourceData, ok := dataInstance.(*viewmodel.ResourceChangeData); ok {
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
	case "create_provider_project":
		handleEvent(bodyBytes, w, func(data []byte) error {
			var evt platviewmodel.CreateProviderProjectRequest
			if err := json.Unmarshal(data, &evt); err != nil {
				return err
			}
			service, err := platv1.CreatePlatformService(ctx)
			if err != nil {
				return err
			}
			return service.HandleCreateProviderProject(ctx, &evt)
		})
	case "create_provider_webhook":
		handleEvent(bodyBytes, w, func(data []byte) error {
			var evt platviewmodel.CreateProviderWebhookRequest
			if err := json.Unmarshal(data, &evt); err != nil {
				return err
			}
			service, err := platv1.CreatePlatformService(ctx)
			if err != nil {
				return err
			}
			return service.HandleCreateProviderWebhook(ctx, &evt)
		})
	case "project_screenshot":
		handleEvent(bodyBytes, w, func(data []byte) error {
			var evt platviewmodel.ProjectScreenshotRequest
			if err := json.Unmarshal(data, &evt); err != nil {
				return err
			}
			service, err := platv1.CreatePlatformService(ctx)
			if err != nil {
				return err
			}
			return service.HandleProjectScreenshot(ctx, &evt)
		})

	case "ResourceCreated":
		handleEvent(bodyBytes, w, func(data []byte) error {
			var evt resource.ResourceCreatedEvent
			if err := json.Unmarshal(data, &evt); err != nil {
				return err
			}
			service, err := v1.CreateResourceQueryService(ctx)
			if err != nil {
				return err
			}
			changeData := convertToChangeData(evt.ID, evt.ResourceVersion, event, evt.CreatedAt, evt.Name, evt.Type, evt.Data, evt.ImageData, evt.Tags)
			return service.HandleResourceChanged(ctx, changeData)
		})

	case "ResourceUpdated":
		handleEvent(bodyBytes, w, func(data []byte) error {
			var evt resource.ResourceUpdatedEvent
			if err := json.Unmarshal(data, &evt); err != nil {
				return err
			}
			service, err := v1.CreateResourceQueryService(ctx)
			if err != nil {
				return err
			}
			changeData := convertToChangeData(evt.ID, evt.ResourceVersion, event, evt.CreatedAt, evt.Name, evt.Type, evt.Data, evt.ImageData, evt.Tags)
			return service.HandleResourceChanged(ctx, changeData)
		})

	case "ResourceDeleted":
		handleEvent(bodyBytes, w, func(data []byte) error {
			var evt resource.ResourceDeletedEvent
			if err := json.Unmarshal(data, &evt); err != nil {
				return err
			}
			service, err := v1.CreateResourceQueryService(ctx)
			if err != nil {
				return err
			}
			return service.HandleResourceChanged(ctx, viewmodel.ResourceChangeData{
				ID:              evt.ID,
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

func getDataType(tableName string) reflect.Type {
	switch tableName {
	case "resource_events":
		return reflect.TypeOf(viewmodel.ResourceChangeData{})
	default:
		return nil
	}
}

// convertToChangeData creates ResourceChangeData from fields
func convertToChangeData(id string, version int, eventType string, createdAt time.Time, name, typ, data, imageData string, tags []string) viewmodel.ResourceChangeData {
	return viewmodel.ResourceChangeData{
		ID:              id,
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

type TriggerEvent struct {
	Platform     string      `json:"platform"`
	Operate      string      `json:"operate"`
	DataBaseName string      `json:"db"`
	TableName    string      `json:"table"`
	Data         interface{} `json:"data"`
}
