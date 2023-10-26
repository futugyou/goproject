package awsconfigConfiguration

type ServiceDiscoveryInstanceConfiguration struct {
	InstanceAttributes InstanceAttributes `json:"InstanceAttributes"`
	InstanceID         string             `json:"InstanceId"`
	ServiceID          string             `json:"ServiceId"`
}

type InstanceAttributes struct {
	EcsTaskDefinitionFamily string `json:"ECS_TASK_DEFINITION_FAMILY"`
	EcsClusterName          string `json:"ECS_CLUSTER_NAME"`
	AvailabilityZone        string `json:"AVAILABILITY_ZONE"`
	EcsServiceName          string `json:"ECS_SERVICE_NAME"`
	AwsInstanceIpv4         string `json:"AWS_INSTANCE_IPV4"`
	AwsInitHealthStatus     string `json:"AWS_INIT_HEALTH_STATUS"`
	Region                  string `json:"REGION"`
}
