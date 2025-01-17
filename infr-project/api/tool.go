package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"reflect"

	_ "github.com/joho/godotenv/autoload"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/futugyou/extensions"

	"github.com/futugyou/infr-project/application"
	"github.com/futugyou/infr-project/controller"
	tool "github.com/futugyou/infr-project/extensions"
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
	case "webhook":
		handleWebhook(ctrl, r, w)
	default:
		w.Write([]byte("system error"))
		w.WriteHeader(500)
		return
	}
}

func redistool(_ *controller.Controller, r *http.Request, w http.ResponseWriter) {
	client, err := tool.RedisClient(os.Getenv("REDIS_URL"))
	if err != nil {
		w.Write([]byte("linkMsg:" + err.Error()))
		w.WriteHeader(500)
		return
	}

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
		return reflect.TypeOf(application.ResourceChangeData{})
	default:
		return nil
	}
}

func handleWebhook(_ *controller.Controller, r *http.Request, w http.ResponseWriter) {
	ctx := r.Context()
	bodyBytes, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	query := r.URL.Query()

	reqInfo := RequestInfo{
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

	resp, _ := json.MarshalIndent(reqInfo, "", "  ")
	config := infra.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	if mongoclient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.ConnectString)); err != nil {
		c := mongoclient.Database(config.DBName).Collection("webhook_testing_logs")
		c.InsertOne(ctx, bson.M{
			"_id":  uuid.New().String(),
			"data": string(resp),
		})
	}

	w.WriteHeader(200)

}

type RequestInfo struct {
	Method     string              `json:"method"`
	URL        string              `json:"url"`
	Proto      string              `json:"proto"`
	Host       string              `json:"host"`
	Header     http.Header         `json:"header"`
	Body       string              `json:"body"`
	Query      map[string][]string `json:"query"`
	RemoteAddr string              `json:"remote_addr"`
	UserAgent  string              `json:"user_agent"`
}
