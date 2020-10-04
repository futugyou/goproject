package model

type Model struct {
	ID          uint32 `gorm:"primary_key" json:"id"`
	CreatedBy   string `json:"created_by"`
	CreatedOn   uint32 `json:"created_on"`
	ModififedBy string `json:"modifed_by"`
	ModififedOn uint32 `json:"modifed_on"`
	Deleted     uint32 `json:"deleted_on"`
	IsDel       uint32 `json:"is_del"`
}
