package vault

import (
	"fmt"

	"github.com/futugyou/infr-project/domain"
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

func WithVaultType(vType VaultType, identities ...string) VaultOption {
	return func(w *Vault) {
		type_identity := ""
		if vType == VaultTypeSystem {
			type_identity = "system"
		} else if vType == VaultTypeCommon {
			type_identity = "common"
		} else if len(identities) > 0 {
			type_identity = identities[0]
		}

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
		Aggregate:    domain.Aggregate{},
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

	vault.Id = fmt.Sprintf("%s-%s-%s", vault.VaultType.String(), vault.TypeIdentity, vault.Key)

	return vault
}

func (r Vault) AggregateName() string {
	return "vaults"
}
