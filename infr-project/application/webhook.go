package application

import (
	"context"
	"encoding/json"
	"fmt"
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

func (s *WebhookService) ProviderWebhookCallback(ctx context.Context, data WebhookRequestInfo) error {
	signature := s.getProviderWebhookSignature(data.Header)
	verifyResult, err := tool.VerifySignatureHMAC(os.Getenv("TRIGGER_AUTH_KEY"), signature, data.Body)
	if err != nil {
		return err
	}

	if !verifyResult {
		return fmt.Errorf("signature verification failed")
	}

	callLog := s.buildWebhookLog(data)
	c := s.database.Collection("platform_webhook_logs")
	_, err = c.InsertOne(ctx, callLog)
	return err
}

func (s *WebhookService) VerifyTesting(ctx context.Context) ([]VerifyResponse, error) {
	result := make([]WebhookLogs, 0)
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

func (*WebhookService) buildWebhookLog(data WebhookRequestInfo) WebhookLogs {
	common := CommonWebhook{}
	json.Unmarshal([]byte(data.Body), &common)

	source := ""
	eventType := ""
	providerPlatformId := ""
	providerProjectId := ""
	providerWebhookId := ""
	if strings.HasPrefix(data.UserAgent, "CircleCI-") {
		source = "circleci"
		eventType = common.Type
		providerPlatformId = common.Organization.Name
		providerProjectId = common.Project.Slug
		providerWebhookId = common.Webhook.ID
	}

	if strings.HasPrefix(data.UserAgent, "GitHub-") {
		source = "github"
		if h, ok := data.Header["X-Github-Event"]; ok && len(h) > 0 {
			eventType = h[0]
		}
		fulls := strings.Split(common.Repository.FullName, "/")
		if len(fulls) == 2 {
			providerPlatformId = fulls[0]
			providerProjectId = fulls[1]
		}
		if h, ok := data.Header["X-Github-Hook-Id"]; ok && len(h) > 0 {
			providerWebhookId = h[0]
		}
	}

	callLog := WebhookLogs{
		Id:                 uuid.NewString(),
		Source:             source,
		EventType:          eventType,
		ProviderPlatformId: providerPlatformId,
		ProviderProjectId:  providerProjectId,
		ProviderWebhookId:  providerWebhookId,
		Data:               data.Body,
	}

	return callLog
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

type WebhookLogs struct {
	Id                 string `bson:"_id"`
	Source             string `bson:"source"` // github/vercel/circleci
	EventType          string `bson:"event_type"`
	ProviderPlatformId string `bson:"provider_platform_id"`
	ProviderProjectId  string `bson:"provider_project_id"`
	ProviderWebhookId  string `bson:"provider_webhook_id"`
	Data               string `bson:"data"`
	HappenedAt         string `json:"happened_at"`
}

type CommonWebhook struct {
	Type         string          `json:"type"`
	Webhook      ProviderWebhook `json:"webhook"`
	Repository   Repository      `json:"repository"`
	Project      CircleProject   `json:"project"`
	Organization CircleOrg       `json:"organization"`
}

type CircleOrg struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CircleProject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Repository struct {
	FullName string `json:"full_name"`
}

type ProviderWebhook struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}