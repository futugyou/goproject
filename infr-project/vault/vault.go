package vault

import (
	"fmt"

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
		Aggregate: domain.Aggregate{
			Id: uuid.New().String(),
		},
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

func (v *Vault) UpdateKey(key string) error {
	v.Key = key
	return nil
}

func (v *Vault) UpdateValue(value string) error {
	v.Value = value
	return nil
}

func (v *Vault) UpdateStorageMedia(media StorageMedia) error {
	v.StorageMedia = media
	return nil
}

func (v *Vault) UpdateVaultType(vType VaultType, identities ...string) error {
	var type_identity string

	if vType == VaultTypeSystem {
		type_identity = "system"
	} else if vType == VaultTypeCommon {
		type_identity = "common"
	} else if len(identities) == 0 {
		return fmt.Errorf("type identity can not be empty when VaultType is not 'system' or 'common'")
	} else {
		type_identity = identities[0]
	}

	v.VaultType = vType
	v.TypeIdentity = type_identity
	return nil
}

func (v *Vault) UpdateTags(tags []string) error {
	v.Tags = tags
	return nil
}
