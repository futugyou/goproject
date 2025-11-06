package service

type CreateVaultRequest struct {
	CreateVaultModel `json:",inline"`
	ForceInsert      bool `json:"force_insert"`
}

type CreateVaultModel struct {
	Key          string            `json:"key" validate:"required,min=3,max=150"`
	Value        string            `json:"value" validate:"required,min=3,max=150"`
	StorageMedia string            `json:"storage_media" validate:"oneof=Local AWS HCP AzureVault"`
	VaultType    string            `json:"vault_type" validate:"oneof=system common project resource platform platform_project platform_webhook"`
	TypeIdentity string            `json:"type_identity" validate:"min=3,max=150"`
	Tags         []string          `json:"tags"`
	Description  string            `json:"description"`
	Extension    map[string]string `json:"extension"`
}

type VaultView struct {
	Id           string            `json:"id"`
	Key          string            `json:"key"`
	MaskValue    string            `json:"mask_value"`
	StorageMedia string            `json:"storage_media"`
	VaultType    string            `json:"vault_type"`
	TypeIdentity string            `json:"type_identity"`
	Tags         []string          `json:"tags"`
	Description  string            `json:"description"`
	Extension    map[string]string `json:"extension"`
}
