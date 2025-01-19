package application

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"strings"

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
		source := ""
		signatureHeader := ""
		if h, ok := data.Header["Circleci-Signature"]; ok && len(h) > 0 {
			source = "circleci"
			signatureHeader = h[0]
		}

		if h, ok := data.Header["X-Hub-Signature-256"]; ok {
			source = "github"
			signatureHeader = h[0]
		}

		if r, err := VerifySignature(os.Getenv("TRIGGER_AUTH_KEY"), signatureHeader, data.Body, source); err != nil {
			ver.Message = err.Error()
		} else {
			ver.Verify = r
		}

		datas = append(datas, ver)
	}
	return datas, nil
}

func VerifySignature(secret, signatureHeader, payload, source string) (bool, error) {
	var signature string

	switch source {
	case "github":
		if !strings.HasPrefix(signatureHeader, "sha256=") {
			return false, errors.New("invalid GitHub signature format")
		}
		signature = signatureHeader[7:]
	case "circleci":
		for _, part := range strings.Split(signatureHeader, ",") {
			kv := strings.SplitN(part, "=", 2)
			if len(kv) == 2 && kv[0] == "v1" {
				signature = kv[1]
				break
			}
		}
		if signature == "" {
			return false, errors.New("no CircleCI v1 signature found")
		}
	default:
		return false, errors.New("unsupported source")
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(expectedSignature), []byte(signature)), nil
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
