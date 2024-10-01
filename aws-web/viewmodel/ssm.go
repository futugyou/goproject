package viewmodel

import "time"

type SSMDataFilter struct {
	AccountId string `json:"accountId"`
	Name      string `json:"name"`
}

type SSMData struct {
	Key        string    `json:"key"`
	Value      string    `json:"value"`
	CreateDate time.Time `json:"create_date"`
}
