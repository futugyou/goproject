package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"reflect"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/redis/go-redis/v9"

	"github.com/futugyou/extensions"

	"github.com/futugyou/infr-project/application"
	"github.com/futugyou/infr-project/controller"
	infra "github.com/futugyou/infr-project/infrastructure_mongo"
	models "github.com/futugyou/infr-project/view_models"
)

func ToolsDispatch(w http.ResponseWriter, r *http.Request) {
	// cors
	if extensions.Cors(w, r) {
		return
	}

	op := r.URL.Query().Get("optype")
	ctrl := controller.NewController()
	switch op {
	case "redis":
		redistool(ctrl, r, w)
	case "event":
		eventHandler(ctrl, r, w)
	default:
		w.Write([]byte("system error"))
		w.WriteHeader(500)
		return
	}
}

func redistool(_ *controller.Controller, r *http.Request, w http.ResponseWriter) {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		w.Write([]byte("linkMsg:" + err.Error()))
		w.WriteHeader(500)
		return
	}
	opt.MaxRetries = 3
	opt.DialTimeout = 10 * time.Second
	opt.ReadTimeout = -1
	opt.WriteTimeout = -1
	opt.DB = 0

	client := redis.NewClient(opt)

	ctx := r.Context()

	err = client.Set(ctx, "foo", "bar", 0).Err()
	if err != nil {
		w.Write([]byte("WriteMsg:" + err.Error()))
		w.WriteHeader(500)
		return
	}

	val, err := client.Get(ctx, "foo").Result()
	if err != nil {
		w.Write([]byte("ReadMsg:" + err.Error()))
		w.WriteHeader(500)
		return
	}

	w.Write([]byte("ResultMsg:" + val))
	w.WriteHeader(200)
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

	if resourceData, ok := dataInstance.(*application.ResourceChangeData); ok {
		if err = service.HandleResourceChaged(ctx, *resourceData); err != nil {
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
	queryRepo, err := createResourceQueryRepository(ctx)
	if err != nil {
		return nil, err
	}
	return application.NewResourceQueryService(queryRepo), nil
}

func createResourceQueryRepository(ctx context.Context) (*infra.ResourceQueryRepository, error) {
	config := infra.DBConfig{
		DBName:        os.Getenv("query_db_name"),
		ConnectString: os.Getenv("query_mongodb_url"),
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectString))
	if err != nil {
		return nil, err
	}

	return infra.NewResourceQueryRepository(client, config), nil
}

func getDataType(tableName string) reflect.Type {
	switch tableName {
	case "resource_events":
		return reflect.TypeOf(application.ResourceChangeData{})
	default:
		return nil
	}
}
