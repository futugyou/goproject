package awsconfigConfiguration

type EventRoleConfiguration struct {
	EventBusName       string   `json:"EventBusName"`
	ScheduleExpression string   `json:"ScheduleExpression"`
	State              string   `json:"State"`
	Targets            []Target `json:"Targets"`
	ID                 string   `json:"Id"`
	Arn                string   `json:"Arn"`
	Name               string   `json:"Name"`
}

type Target struct {
	Arn           string        `json:"Arn"`
	ID            string        `json:"Id"`
	RoleArn       string        `json:"RoleArn"`
	EcsParameters EcsParameters `json:"EcsParameters"`
}

type EcsParameters struct {
	TaskDefinitionArn string `json:"TaskDefinitionArn"`
}
