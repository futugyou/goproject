package viewmodel

import "time"

type EcsClusterFilter struct {
	AccountId string `json:"accountId,omitempty"`
}

type EcsClusterViewModel struct {
	ClusterName  string    `json:"cluster_name"`
	ClusterArn   string    `json:"cluster_arn"`
	Service      string    `json:"service"`
	ServiceArn   string    `json:"service_arn"`
	RoleArn      string    `json:"role_arn"`
	AccountAlias string    `json:"account_alias"`
	OperateAt    time.Time `json:"operate_At"`
}
