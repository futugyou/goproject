package viewmodel

import "time"

type EcsClusterFilter struct {
	AccountId string `json:"accountId,omitempty"`
}

type EcsClusterViewModel struct {
	Name         string    `json:"name"`
	ClusterName  string    `json:"cluster_name"`
	Service      string    `json:"service"`
	AccountAlias string    `json:"account_alias"`
	CreatedAt    time.Time `json:"createdAt"`
}
