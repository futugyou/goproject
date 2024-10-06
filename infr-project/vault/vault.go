package vault

import (
	"fmt"

	tool "github.com/futugyou/extensions"

	"github.com/futugyou/infr-project/domain"
	"github.com/google/uuid"
)

type Vault struct {
	domain.Aggregate `json:"-"`
	Key              string       `json:"key"`
	Value            string       `json:"value"`
	StorageMedia     StorageMedia `json:"storage_media"` // local,aws,HCP,...
	VaultType        VaultType    `json:"vault_type"`    // system,common,project,resource,platform
	TypeIdentity     string       `json:"type_identity"` // system,common,projectId,...
	State            VaultState   `json:"state"`         // default,changing
	Tags             []string     `json:"tags"`
	hasChange        bool         `json:"-"`
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
	if v.Key != key {
		v.Key = key
		v.hasChange = true
	}
	return nil
}

func (v *Vault) UpdateValue(value string) error {
	if v.Value != value {
		v.Value = value
		v.hasChange = true
	}
	return nil
}

func (v *Vault) UpdateStorageMedia(media StorageMedia) error {
	if v.StorageMedia != media {
		v.StorageMedia = media
		v.hasChange = true
	}
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

	if v.VaultType != vType || v.TypeIdentity != type_identity {
		v.VaultType = vType
		v.TypeIdentity = type_identity
		v.hasChange = true
	}
	return nil
}

func (v *Vault) UpdateTags(tags []string) error {
	if !tool.StringArrayCompare(v.Tags, tags) {
		v.Tags = tags
		v.hasChange = true
	}
	return nil
}

func (v *Vault) HasChange() bool {
	return v.hasChange
}
