package viewmodel

import "time"

type IAMDataFilter struct {
	AccountId string `json:"accountId"`
}

type IAMData struct {
	UserName   string       `json:"user_name"`
	CreateDate time.Time    `json:"create_date"`
	LastUsed   *time.Time   `json:"last_used"`
	Keys       []IAMDataKey `json:"keys"`
}

type IAMDataKey struct {
	Key        string    `json:"key"`
	Status     string    `json:"status"`
	CreateDate time.Time `json:"create_date"`
}
