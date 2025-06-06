package viewmodels

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

type CreateVaultsRequest struct {
	Vaults      []CreateVaultModel `json:"vaults" validate:"required,gt=0,dive"`
	ForceInsert bool               `json:"force_insert"`
}

type CreateVaultsResponse struct {
	Vaults []VaultView `json:"vaults"`
}

type SearchVaultsRequest struct {
	Key          string   `json:"key"`
	StorageMedia string   `json:"storage_media"`
	VaultType    string   `json:"vault_type"`
	TypeIdentity string   `json:"type_identity"`
	Description  string   `json:"description"`
	Tags         []string `json:"tags"`
	Page         int      `json:"page"`
	Size         int      `json:"size"`
}

type CreateVaultRequest struct {
	CreateVaultModel `json:",inline"`
	ForceInsert      bool `json:"force_insert"`
}

type ChangeVaultRequest struct {
	Data        ChangeVaultItem `json:"vault_data"`
	ForceInsert bool            `json:"force_insert"`
}

type ChangeVaultItem struct {
	Key          *string            `json:"key" validate:"min=3,max=150"`
	Value        *string            `json:"value" validate:"min=3,max=150"`
	StorageMedia *string            `json:"storage_media" validate:"oneof=Local AWS HCP AzureVault"`
	VaultType    *string            `json:"vault_type" validate:"oneof=system common project resource platform platform_project platform_webhook"`
	TypeIdentity *string            `json:"type_identity" validate:"min=3,max=150"`
	Tags         *[]string          `json:"tags"`
	Description  *string            `json:"description" validate:"min=3,max=250"`
	Extension    *map[string]string `json:"extension"`
}

type ImportVaultsRequest struct {
	StorageMedia string  `json:"storage_media" validate:"oneof=AWS HCP AzureVault"`
	VaultType    *string `json:"vault_type" validate:"oneof=system common project resource platform platform_project platform_webhook"`
	TypeIdentity *string `json:"type_identity" validate:"min=3,max=150"`
}

type ImportVaultsResponse struct {
	Vaults []VaultView `json:"vaults"`
}
