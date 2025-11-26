package entity

import "github.com/futugyou/platformservice/domain"

type ProjectEntity struct {
	ID                   string            `bson:"id"`
	PlatformID           string            `bson:"platform_id"`
	Name                 string            `bson:"name"`
	Url                  string            `bson:"url"`
	Description          string            `bson:"description"`
	Properties           map[string]string `bson:"properties"`
	Secrets              []SecretEntity    `bson:"secrets"`
	ImageUrl             string            `bson:"image_url"`
	ProviderProjectID    string            `bson:"provider_project_id"`
	Tags                 []string          `bson:"tags"`
	WebhookID            string            `bson:"webhook_id"`
	WebhookUrl           string            `bson:"webhook_url"`
	WebhookEvents        []string          `bson:"webhook_events"`
	WebhookState         string            `bson:"webhook_state"`
	WebhookSigningSecret string            `bson:"webhook_signing_secret"`
}

type ProjectMapper struct{}

func (*ProjectMapper) ToDomain(p *ProjectEntity) *domain.PlatformProject {
	properties := map[string]domain.Property{}
	for k, v := range p.Properties {
		properties[k] = domain.Property{
			Key:   k,
			Value: v,
		}
	}

	secrets := map[string]domain.Secret{}
	for _, s := range p.Secrets {
		secrets[s.Key] = domain.Secret{
			Key:            s.Key,
			Value:          s.Value,
			VaultKey:       s.VaultKey,
			VaultMaskValue: s.VaultMaskValue,
		}
	}

	var webhook *domain.Webhook
	if p.WebhookID != "" {
		webhook = &domain.Webhook{
			ID:            p.WebhookID,
			Url:           p.WebhookUrl,
			Events:        p.WebhookEvents,
			State:         domain.GetWebhookState(p.WebhookState),
			SigningSecret: p.WebhookSigningSecret,
		}
	}
	return &domain.PlatformProject{
		ID:                p.ID,
		Name:              p.Name,
		Url:               p.Url,
		Description:       p.Description,
		Properties:        properties,
		Secrets:           secrets,
		ImageUrl:          p.ImageUrl,
		ProviderProjectID: p.ProviderProjectID,
		Tags:              p.Tags,
		Webhook:           webhook,
	}
}

func (*ProjectMapper) ToEntity(platformID string, project *domain.PlatformProject) *ProjectEntity {
	p := &ProjectEntity{}
	p.ID = project.ID
	p.PlatformID = platformID
	p.Name = project.Name
	p.Url = project.Url
	p.Description = project.Description

	p.Properties = map[string]string{}
	for k, v := range project.Properties {
		p.Properties[k] = v.Value
	}

	p.Secrets = []SecretEntity{}
	for _, s := range project.Secrets {
		p.Secrets = append(p.Secrets, SecretEntity{
			Key:            s.Key,
			Value:          s.Value,
			VaultKey:       s.VaultKey,
			VaultMaskValue: s.VaultMaskValue,
		})
	}

	p.ImageUrl = project.ImageUrl
	p.ProviderProjectID = project.ProviderProjectID
	p.Tags = project.Tags

	if project.Webhook != nil {
		p.WebhookID = project.Webhook.ID
		p.WebhookUrl = project.Webhook.Url
		p.WebhookEvents = project.Webhook.Events
		p.WebhookState = project.Webhook.State.String()
		p.WebhookSigningSecret = project.Webhook.SigningSecret
	}

	return p
}
