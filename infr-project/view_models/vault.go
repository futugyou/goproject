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
