package viewmodel

import "time"

type ParameterViewModel struct {
	Id           string    `json:"id,omitempty"`
	AccountId    string    `json:"accountId,omitempty"`
	AccountAlias string    `json:"alias,omitempty"`
	Region       string    `json:"region,omitempty"`
	Key          string    `json:"key,omitempty"`
	Value        string    `json:"value,omitempty"`
	Version      string    `json:"version,omitempty"`
	OperateAt    time.Time `json:"operateAt,omitempty"`
}

type ParameterDetailViewModel struct {
	Id           string               `json:"id,omitempty"`
	AccountId    string               `json:"accountId,omitempty"`
	AccountAlias string               `json:"alias,omitempty"`
	Region       string               `json:"region,omitempty"`
	Key          string               `json:"key,omitempty"`
	Value        string               `json:"value,omitempty"`
	Version      string               `json:"version,omitempty"`
	OperateAt    time.Time            `json:"operateAt,omitempty"`
	Current      *ParameterViewModel  `json:"current,omitempty"`
	History      []ParameterViewModel `json:"history,omitempty"`
}

type ParameterFilter struct {
	AccountAlias string `json:"alias,omitempty"`
	Region       string `json:"region,omitempty"`
	Key          string `json:"key,omitempty"`
}

type CompareViewModel struct {
	Key     string `json:"key,omitempty"`
	Value   string `json:"value,omitempty"`
	Version string `json:"version,omitempty"`
}

type SyncModel struct {
	Id string `json:"id,omitempty"`
}
