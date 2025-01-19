package application

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type WebhookService struct {
	database *mongo.Database
}

func NewWebhookService(database *mongo.Database) *WebhookService {
	return &WebhookService{
		database: database,
	}
}

func (s *WebhookService) PlatformCallback(ctx context.Context, data WebhookRequestInfo) error {
	resp, _ := json.MarshalIndent(data, "", "  ")

	c := s.database.Collection("webhook_testing_logs")
	_, err := c.InsertOne(ctx, bson.M{
		"_id":  uuid.New().String(),
		"data": string(resp),
	})
	return err
}

func (s *WebhookService) VerifyTesting(ctx context.Context) ([]WebhookRequestInfo, error) {
	result := make([]WebhookInfo, 0)
	c := s.database.Collection("webhook_testing_logs")
	filter := bson.D{}
	op := options.Find()
	cursor, err := c.Find(ctx, filter, op)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	for _, data := range result {
		cursor.Decode(&data)
	}

	datas := make([]WebhookRequestInfo, 0)
	for _, d := range result {
		data := WebhookRequestInfo{}
		json.Unmarshal([]byte(d.Data), &data)
		// if h,ok:=data.Header["Circleci-Signature"];ok {

		// }
		datas = append(datas, data)
	}
	return datas, nil
}

type WebhookInfo struct {
	Id   string `bson:"_id"`
	Data string `bson:"data"`
}

type WebhookRequestInfo struct {
	Method     string              `json:"method"`
	URL        string              `json:"url"`
	Proto      string              `json:"proto"`
	Host       string              `json:"host"`
	Header     map[string][]string `json:"header"`
	Body       string              `json:"body"`
	Query      map[string][]string `json:"query"`
	RemoteAddr string              `json:"remote_addr"`
	UserAgent  string              `json:"user_agent"`
}
