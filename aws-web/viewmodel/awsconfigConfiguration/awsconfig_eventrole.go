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
	Arn string `json:"Arn"`
	ID  string `json:"Id"`
}
