package viewmodel

import "strings"

type VaultView struct {
	ID           string            `json:"id"`
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
	Key          string  `json:"key" form:"key"`
	StorageMedia string  `json:"storage_media" form:"storage_media"`
	VaultType    string  `json:"vault_type" form:"vault_type"`
	TypeIdentity string  `json:"type_identity" form:"type_identity"`
	Description  string  `json:"description" form:"description"`
	Tags         CSVList `json:"tags" form:"tags"`
	Page         int     `json:"page" form:"page,default=1"`
	Size         int     `json:"size" form:"size,default=100"`
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

type CSVList []string

func (c *CSVList) UnmarshalParam(src string) error {
	if src == "" {
		*c = nil
		return nil
	}

	items := strings.Split(src, ",")
	for i, s := range items {
		items[i] = strings.TrimSpace(s)
	}
	*c = items
	return nil
}
