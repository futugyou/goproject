package awsconfigConfiguration

type VPCEndpointConfiguration struct {
	VpcEndpointID string   `json:"vpcEndpointId"`
	VpcID         string   `json:"vpcId"`
	ServiceName   string   `json:"serviceName"`
	SubnetIDS     []string `json:"subnetIds"`
	Groups        []Group  `json:"groups"`
}
