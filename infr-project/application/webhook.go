package application

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	tool "github.com/futugyou/extensions"
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

	c := s.database.Collection("platform_webhook_logs")
	_, err := c.InsertOne(ctx, bson.M{
		"_id":  uuid.New().String(),
		"data": string(resp),
	})
	return err
}

func (s *WebhookService) VerifyTesting(ctx context.Context) ([]VerifyResponse, error) {
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

	datas := []VerifyResponse{}
	for _, d := range result {
		data := WebhookRequestInfo{}
		json.Unmarshal([]byte(d.Data), &data)
		ver := VerifyResponse{
			Id:      d.Id,
			Verify:  false,
			Message: "",
		}

		signature := s.getProviderWebhookSignature(data.Header)
		if r, err := tool.VerifySignatureHMAC(os.Getenv("TRIGGER_AUTH_KEY"), signature, data.Body); err != nil {
			ver.Message = err.Error()
		} else {
			ver.Verify = r
		}

		datas = append(datas, ver)
	}
	return datas, nil
}

func (*WebhookService) getProviderWebhookSignature(header map[string][]string) string {
	signature := ""
	if h, ok := header["Circleci-Signature"]; ok && len(h) > 0 {
		for _, part := range strings.Split(h[0], ",") {
			kv := strings.SplitN(part, "=", 2)
			if len(kv) == 2 && kv[0] == "v1" {
				signature = kv[1]
				break
			}
		}
	}

	if h, ok := header["X-Hub-Signature-256"]; ok {
		if strings.HasPrefix(h[0], "sha256=") {
			signature = h[0][7:]
		}
	}
	return signature
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

type VerifyResponse struct {
	Id      string `json:"id"`
	Verify  bool   `json:"verify"`
	Message string `json:"message"`
}
