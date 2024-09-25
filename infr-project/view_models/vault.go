package viewmodels

type VaultView struct {
	Id           string   `json:"id"`
	Key          string   `json:"key"`
	MaskValue    string   `json:"mask_value"`
	StorageMedia string   `json:"storage_media"`
	VaultType    string   `json:"vault_type"`
	TypeIdentity string   `json:"type_identity"`
	Tags         []string `json:"tags"`
}

type CreateVaultModel struct {
	Key          string   `json:"key" validate:"required,min=3,max=150"`
	Value        string   `json:"value" validate:"required,min=3,max=150"`
	StorageMedia string   `json:"storage_media" validate:"oneof=local aws HCP"`
	VaultType    string   `json:"vault_type" validate:"oneof=system common project resource platform"`
	TypeIdentity string   `json:"type_identity" validate:"min=3,max=150"`
	Tags         []string `json:"tags"`
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
	Tags         []string `json:"tags"`
}

type ChangeVaultRequest struct {
	Key          *string   `json:"key" validate:"min=3,max=150"`
	Value        *string   `json:"value" validate:"min=3,max=150"`
	StorageMedia *string   `json:"storage_media" validate:"oneof=local aws HCP"`
	VaultType    *string   `json:"vault_type" validate:"oneof=system common project resource platform"`
	TypeIdentity *string   `json:"type_identity" validate:"min=3,max=150"`
	Tags         *[]string `json:"tags"`
}
