package application

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	models "github.com/futugyou/infr-project/view_models"
	"github.com/futugyou/infr-project/webhook"
)

type WebhookService struct {
	repository webhook.IWebhookLogRepository
}

func NewWebhookService(repository webhook.IWebhookLogRepository) *WebhookService {
	return &WebhookService{
		repository: repository,
	}
}

func (s *WebhookService) ProviderWebhookCallback(ctx context.Context, data models.WebhookRequestInfo) error {
	webhookLog, err := s.buildWebhookLog(data)
	if err != nil {
		return err
	}

	verifyResult, err := webhookLog.Verify(os.Getenv("TRIGGER_AUTH_KEY"))
	if err != nil {
		return err
	}

	if !verifyResult {
		return fmt.Errorf("signature verification failed")
	}

	return s.repository.Insert(ctx, *webhookLog)
}

func (s *WebhookService) VerifyTesting(ctx context.Context) ([]models.VerifyResponse, error) {
	datas := []models.VerifyResponse{}
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

func (s *WebhookService) buildWebhookLog(data models.WebhookRequestInfo) (*webhook.WebhookLogs, error) {
	signature := s.getProviderWebhookSignature(data.Header)
	if len(signature) == 0 {
		return nil, fmt.Errorf("signature verification failed")
	}

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
	} else if strings.HasPrefix(data.UserAgent, "GitHub-") {
		if len(common.Action) > 0 && common.Action != "completed" {
			return nil, fmt.Errorf("github webhook log when action is completed")
		}

		source = "github"
		if h, ok := data.Header["X-Github-Event"]; ok && len(h) > 0 {
			eventType = h[0]
		} else {
			return nil, fmt.Errorf("github webhook need X-Github-Event in header")
		}

		fulls := strings.Split(common.Repository.FullName, "/")
		if len(fulls) == 2 {
			providerPlatformId = fulls[0]
			providerProjectId = fulls[1]
		} else {
			return nil, fmt.Errorf("github webhook need respository fullname")
		}

		if h, ok := data.Header["X-Github-Hook-Id"]; ok && len(h) > 0 {
			providerWebhookId = h[0]
		} else {
			return nil, fmt.Errorf("github webhook need X-Github-Hook-Id in header")
		}
	} else {
		return nil, fmt.Errorf("unsupport webhook")
	}

	webhookLog := webhook.NewWebhookLogs(source, eventType, providerPlatformId, providerProjectId, providerWebhookId, data.Body, signature)

	return webhookLog, nil
}

func (s *WebhookService) SearchWebhookLogs(ctx context.Context, searcher models.WebhookSearch) ([]models.WebhookLogs, error) {
	filter := webhook.WebhookLogSearch(searcher)
	result := []models.WebhookLogs{}
	logs, err := s.repository.SearchWebhookLogs(ctx, filter)
	if err != nil {
		return nil, err
	}

	for _, log := range logs {
		result = append(result, models.WebhookLogs{
			Source:             log.Source,
			EventType:          log.EventType,
			ProviderPlatformId: log.ProviderPlatformId,
			ProviderProjectId:  log.ProviderProjectId,
			ProviderWebhookId:  log.ProviderWebhookId,
			Data:               log.Data,
			HappenedAt:         log.HappenedAt,
		})
	}

	return result, nil
}

type CommonWebhook struct {
	Type         string          `json:"type"`
	Action       string          `json:"action"`
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
