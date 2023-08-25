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

type EcsClusterDetailViewModel struct {
	ClusterName       string    `json:"cluster_name"`
	ClusterArn        string    `json:"cluster_arn"`
	Service           string    `json:"service"`
	ServiceArn        string    `json:"service_arn"`
	RoleArn           string    `json:"role_arn"`
	AccountAlias      string    `json:"account_alias"`
	LoadBalancers     []string  `json:"loadBalancers,omitempty"`
	SecurityGroups    []string  `json:"security_groups,omitempty"`
	Subnets           []string  `json:"subnets,omitempty"`
	ServiceRegistries []string  `json:"service_registries,omitempty"`
	TaskDefinitions   []string  `json:"task_definitions,omitempty"`
	OperateAt         time.Time `json:"operate_At"`
}

type EcsTaskCompare struct {
	Id            string `json:"Id"`
	SourceTaskArn string `json:"source_task_arn"`
	DestTaskArn   string `json:"dest_task_arn"`
}
