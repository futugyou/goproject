package entity

import (
	coredomain "github.com/futugyou/domaincore/domain"

	"github.com/futugyou/platformservice/domain"
)

type PlatformEntity struct {
	ID          string            `bson:"id"`
	Name        string            `bson:"name"`
	Activate    bool              `bson:"activate"`
	Url         string            `bson:"url"`
	Description string            `bson:"description"`
	Provider    string            `bson:"provider"`
	Properties  map[string]string `bson:"properties"`
	Secrets     []SecretEntity    `bson:"secrets"`
	Tags        []string          `bson:"tags"`
	IsDeleted   bool              `bson:"is_deleted"`
}

type SecretEntity struct {
	Key            string `bson:"key"`
	Value          string `bson:"value"`
	VaultKey       string `bson:"vault_key"`
	VaultMaskValue string `bson:"vault_mask_value"`
}

type PlatformMapper struct{}

func (*PlatformMapper) ToDomain(p *PlatformEntity) *domain.Platform {
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

	return &domain.Platform{
		Aggregate: coredomain.Aggregate{
			ID: p.ID,
		},
		Name:        p.Name,
		Activate:    p.Activate,
		Url:         p.Url,
		Description: p.Description,
		Provider:    domain.GetPlatformProvider(p.Provider),
		Properties:  properties,
		Secrets:     secrets,
		Projects:    map[string]domain.PlatformProject{},
		Tags:        p.Tags,
		IsDeleted:   p.IsDeleted,
	}
}

func (*PlatformMapper) ToEntity(platform *domain.Platform) *PlatformEntity {
	properties := map[string]string{}
	for k, v := range platform.Properties {
		properties[k] = v.Value
	}

	secrets := []SecretEntity{}
	for _, s := range platform.Secrets {
		secrets = append(secrets, SecretEntity{
			Key:            s.Key,
			Value:          s.Value,
			VaultKey:       s.VaultKey,
			VaultMaskValue: s.VaultMaskValue,
		})
	}

	return &PlatformEntity{
		ID:          platform.ID,
		Name:        platform.Name,
		Activate:    platform.Activate,
		Url:         platform.Url,
		Description: platform.Description,
		Provider:    platform.Provider.String(),
		Properties:  properties,
		Secrets:     secrets,
		Tags:        platform.Tags,
		IsDeleted:   platform.IsDeleted,
	}
}
