package viewmodels

type TriggerEvent struct {
	Platform     string      `json:"platform"`
	Operate      string      `json:"operate"`
	DataBaseName string      `json:"db"`
	TableName    string      `json:"table"`
	Data         interface{} `json:"data"`
}
