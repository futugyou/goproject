package awsconfigConfiguration

type TargetGroupConfiguration struct {
	HealthCheckEnabled         bool     `json:"HealthCheckEnabled"`
	HealthCheckIntervalSeconds int64    `json:"HealthCheckIntervalSeconds"`
	HealthCheckPath            string   `json:"HealthCheckPath"`
	HealthCheckPort            string   `json:"HealthCheckPort"`
	HealthCheckProtocol        string   `json:"HealthCheckProtocol"`
	HealthCheckTimeoutSeconds  int64    `json:"HealthCheckTimeoutSeconds"`
	HealthyThresholdCount      int64    `json:"HealthyThresholdCount"`
	IPAddressType              string   `json:"IpAddressType"`
	LoadBalancerArns           []string `json:"LoadBalancerArns"`
	Matcher                    Matcher  `json:"Matcher"`
	Port                       int64    `json:"Port"`
	Protocol                   string   `json:"Protocol"`
	ProtocolVersion            string   `json:"ProtocolVersion"`
	TargetGroupArn             string   `json:"TargetGroupArn"`
	TargetGroupName            string   `json:"TargetGroupName"`
	TargetType                 string   `json:"TargetType"`
	UnhealthyThresholdCount    int64    `json:"UnhealthyThresholdCount"`
	VpcID                      string   `json:"VpcId"`
}

type Matcher struct {
	GrpcCode interface{} `json:"GrpcCode"`
	HTTPCode string      `json:"HttpCode"`
}
