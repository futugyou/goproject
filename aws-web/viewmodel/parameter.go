package viewmodel

import "time"

type ParameterViewModel struct {
	Id           string    `json:"id,omitempty"`
	AccountId    string    `json:"accountId,omitempty"`
	AccountAlias string    `json:"alias,omitempty"`
	Region       string    `json:"region,omitempty"`
	Key          string    `json:"key"`
	Value        string    `json:"value,omitempty"`
	Version      string    `json:"version,omitempty"`
	NeedSync     bool      `json:"need_sync,omitempty"`
	OperateAt    time.Time `json:"operateAt"`
}

type ParameterFilter struct {
	AccountAlias string `json:"alias"`
	Region       string `json:"region"`
	Key          string `json:"key"`
}

type CompareViewModel struct {
	Key     string `json:"key"`
	Value   string `json:"value,omitempty"`
	Version string `json:"version,omitempty"`
}
