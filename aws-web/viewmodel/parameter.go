package viewmodel

import "time"

type ParameterViewModel struct {
	Id        string    `json:"id,omitempty"`
	AccountId string    `json:"accountId"`
	Region    string    `json:"region"`
	Key       string    `json:"key"`
	Value     string    `json:"value,omitempty"`
	Version   string    `json:"version"`
	NeedSync  bool      `json:"need_sync,omitempty"`
	OperateAt time.Time `json:"operateAt"`
}
