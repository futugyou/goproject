package application

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/futugyou/platformservice/domain"
	"github.com/futugyou/platformservice/options"
	"github.com/futugyou/platformservice/viewmodel"
)

type WebhookLogService struct {
	repository domain.WebhookLogRepository
	opts       *options.Options
}

func NewWebhookLogService(repository domain.WebhookLogRepository, opts *options.Options) *WebhookLogService {
	return &WebhookLogService{
		repository: repository,
		opts:       opts,
	}
}

func (s *WebhookLogService) ProviderWebhookCallback(ctx context.Context, data viewmodel.WebhookRequestInfo) error {
	webhookLog, err := s.buildWebhookLog(data)
	if err != nil {
		return err
	}

	verifyResult, err := webhookLog.Verify(s.opts.TriggerAuthKey)
	if err != nil {
		return err
	}

	if !verifyResult {
		return fmt.Errorf("signature verification failed")
	}

	tenDaysAgo := time.Now().Add(-10 * 24 * time.Hour) // 10 days ago
	return s.repository.InsertAndDeleteOldData(ctx, []domain.WebhookLogs{*webhookLog}, tenDaysAgo)
}

func (*WebhookLogService) getProviderWebhookSignature(header map[string][]string) string {
	signature := ""
	if h, ok := header["Circleci-Signature"]; ok && len(h) > 0 {
		for part := range strings.SplitSeq(h[0], ",") {
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

func getMapString(v map[string][]string, key string) string {
	if v == nil {
		return ""
	}
	vs := v[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

func (s *WebhookLogService) buildWebhookLog(data viewmodel.WebhookRequestInfo) (*domain.WebhookLogs, error) {
	signature := s.getProviderWebhookSignature(data.Header)
	if len(signature) == 0 {
		return nil, fmt.Errorf("signature verification failed")
	}

	common := CommonWebhook{}
	json.Unmarshal([]byte(data.Body), &common)

	source := ""
	eventType := ""
	providerPlatformId := getMapString(data.Query, "platform")
	providerProjectId := getMapString(data.Query, "project")
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
		}

		if h, ok := data.Header["X-Github-Hook-Id"]; ok && len(h) > 0 {
			providerWebhookId = h[0]
		} else {
			return nil, fmt.Errorf("github webhook need X-Github-Hook-Id in header")
		}
	} else {
		return nil, fmt.Errorf("unsupport webhook")
	}

	webhookLog := domain.NewWebhookLogs(source, eventType, providerPlatformId, providerProjectId, providerWebhookId, data.Body, signature)

	return webhookLog, nil
}

func (s *WebhookLogService) SearchWebhookLogs(ctx context.Context, searcher viewmodel.WebhookSearch) ([]viewmodel.WebhookLogs, error) {
	filter := domain.WebhookLogSearch(searcher)
	result := []viewmodel.WebhookLogs{}
	logs, err := s.repository.SearchWebhookLogs(ctx, filter)
	if err != nil {
		return nil, err
	}

	for _, log := range logs {
		result = append(result, viewmodel.WebhookLogs{
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
