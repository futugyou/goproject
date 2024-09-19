package vault

import (
	"github.com/futugyou/infr-project/domain"
	"github.com/google/uuid"
)

type Vault struct {
	domain.Aggregate `json:"-"`
	Key              string       `json:"key"`
	Value            string       `json:"value"`
	StorageMedia     StorageMedia `json:"storage_media"` // local,aws,HCP,...
	VaultType        VaultType    `json:"vault_type"`    // system,common,project,resource,platform
	TypeIdentity     string       `json:"type_identity"` // system-system-key,common-common-key,project-projectId-key,...
	Tags             []string     `json:"tags"`
}

type VaultOption func(*Vault)

func WithStorageMedia(media StorageMedia) VaultOption {
	return func(w *Vault) {
		w.StorageMedia = media
	}
}

func WithVaultType(vType VaultType, type_identity string) VaultOption {
	return func(w *Vault) {
		w.VaultType = vType
		w.TypeIdentity = type_identity
	}
}

func WithTags(tags []string) VaultOption {
	return func(w *Vault) {
		w.Tags = tags
	}
}

func NewVault(key string, value string, opts ...VaultOption) *Vault {
	vault := &Vault{
		Aggregate: domain.Aggregate{
			Id: uuid.New().String()},
		Key:          key,
		Value:        value,
		StorageMedia: StorageMediaLocal,
		VaultType:    VaultTypeCommon,
		TypeIdentity: "common",
		Tags:         []string{},
	}

	for _, opt := range opts {
		opt(vault)
	}

	return vault
}

func (r Vault) AggregateName() string {
	return "vaults"
}
