package entity

import "time"

type EcsServiceSearchFilter struct {
	AccountId string `json:"accountId,omitempty"`
}

type EcsServiceEntity struct {
	Id             string    `bson:"_id,omitempty"`
	AccountId      string    `bson:"account_id"`
	Cluster        string    `bson:"cluster"`
	ClusterArn     string    `bson:"cluster_arn"`
	ServiceName    string    `bson:"service_name"`
	ServiceNameArn string    `bson:"service_name_arn"`
	RoleArn        string    `bson:"role_arn"`
	OperateAt      time.Time `bson:"operate_at"`
}

func (EcsServiceEntity) GetType() string {
	return "ecs_services"
}
