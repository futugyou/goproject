package awsconfigConfiguration

type ECSServiceConfiguration struct {
	ServiceArn               string                  `json:"ServiceArn"`
	CapacityProviderStrategy []interface{}           `json:"CapacityProviderStrategy"`
	Cluster                  string                  `json:"Cluster"`
	DeploymentConfiguration  DeploymentConfiguration `json:"DeploymentConfiguration"`
	DeploymentController     DeploymentController    `json:"DeploymentController"`
	DesiredCount             int64                   `json:"DesiredCount"`
	EnableECSManagedTags     bool                    `json:"EnableECSManagedTags"`
	LaunchType               string                  `json:"LaunchType"`
	LoadBalancers            []interface{}           `json:"LoadBalancers"`
	Name                     string                  `json:"Name"`
	NetworkConfiguration     NetworkConfiguration    `json:"NetworkConfiguration"`
	PlacementConstraints     []interface{}           `json:"PlacementConstraints"`
	PlacementStrategies      []interface{}           `json:"PlacementStrategies"`
	PlatformVersion          string                  `json:"PlatformVersion"`
	Role                     string                  `json:"Role"`
	SchedulingStrategy       string                  `json:"SchedulingStrategy"`
	ServiceName              string                  `json:"ServiceName"`
	ServiceRegistries        []ServiceRegistry       `json:"ServiceRegistries"`
	Tags                     []interface{}           `json:"Tags"`
	TaskDefinition           string                  `json:"TaskDefinition"`
}

type DeploymentConfiguration struct {
	DeploymentCircuitBreaker DeploymentCircuitBreaker `json:"DeploymentCircuitBreaker"`
	MaximumPercent           int64                    `json:"MaximumPercent"`
	MinimumHealthyPercent    int64                    `json:"MinimumHealthyPercent"`
}

type DeploymentCircuitBreaker struct {
	Enable   bool `json:"Enable"`
	Rollback bool `json:"Rollback"`
}

type DeploymentController struct {
	Type string `json:"Type"`
}

type NetworkConfiguration struct {
	AwsvpcConfiguration AwsvpcConfiguration `json:"AwsvpcConfiguration"`
}

type AwsvpcConfiguration struct {
	Subnets        []string `json:"Subnets"`
	SecurityGroups []string `json:"SecurityGroups"`
	AssignPublicIP string   `json:"AssignPublicIp"`
}

type ServiceRegistry struct {
	RegistryArn string `json:"RegistryArn"`
}
