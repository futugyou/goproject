package domain

import (
	"fmt"

	tool "github.com/futugyou/extensions"

	"github.com/futugyou/domaincore/domain"
	"github.com/google/uuid"
)

type Vault struct {
	domain.Aggregate
	Key          string
	Value        string
	Description  string
	Extension    map[string]string
	StorageMedia StorageMedia // local, aws, HCP,...
	VaultType    VaultType    // system, common, project, resource, platform, platform_project, platform_webhook
	TypeIdentity string       // system, common, projectId, resourceId, platformId, platform_project_id, platform_webhook_id
	State        VaultState   // default, changing
	Tags         []string
	hasChange    bool
}

type VaultOption func(*Vault)

func WithExtension(extension map[string]string) VaultOption {
	return func(w *Vault) {
		w.Extension = extension
	}
}

func WithDescription(description string) VaultOption {
	return func(w *Vault) {
		w.Description = description
	}
}

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
			ID: uuid.New().String(),
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

func (v *Vault) UpdateState(state VaultState) error {
	if v.State != state {
		v.State = state
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

func (v *Vault) UpdateDescription(value string) error {
	if v.Description != value {
		v.Description = value
		v.hasChange = true
	}
	return nil
}

func (v *Vault) UpdateExtension(extension map[string]string) error {
	if !tool.MapsCompare(v.Extension, extension) {
		v.Extension = extension
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

func (v *Vault) GetIdentityKey() string {
	t := v.VaultType.String()
	if t == VaultTypeSystem.String() {
		return fmt.Sprintf("system/system/%s", v.Key)
	} else if t == VaultTypeCommon.String() {
		return fmt.Sprintf("common/common/%s", v.Key)
	} else {
		return fmt.Sprintf("%s/%s/%s", v.VaultType.String(), v.TypeIdentity, v.Key)
	}
}
