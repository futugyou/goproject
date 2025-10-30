package application

import "github.com/futugyou/vaultservice/domain"

type VaultChanged struct {
	ID           string            `json:"id"`
	Key          string            `json:"key"`
	Value        string            `json:"value"`
	Description  string            `json:"description"`
	Extension    map[string]string `json:"extension"`
	StorageMedia string            `json:"storage_media"`
	VaultType    string            `json:"vault_type"`
	TypeIdentity string            `json:"type_identity"`
	State        string            `json:"state"`
	Tags         []string          `json:"tags"`
}

func ToVaultChanged(data *domain.Vault) *VaultChanged {
	return &VaultChanged{
		ID:           data.ID,
		Key:          data.Key,
		Value:        data.Value,
		Description:  data.Description,
		Extension:    data.Extension,
		StorageMedia: data.StorageMedia.String(),
		VaultType:    data.VaultType.String(),
		TypeIdentity: data.TypeIdentity,
		State:        data.State.String(),
		Tags:         data.Tags,
	}
}

func (e *VaultChanged) EventType() string {
	return "vault_changed"
}
