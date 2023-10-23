package awsconfigConfiguration

type LoadBalancerConfiguration struct {
	LoadBalancerArn       string             `json:"loadBalancerArn"`
	DNSName               string             `json:"dNSName"`
	CanonicalHostedZoneID string             `json:"canonicalHostedZoneId"`
	CreatedTime           string             `json:"createdTime"`
	LoadBalancerName      string             `json:"loadBalancerName"`
	Scheme                string             `json:"scheme"`
	VpcID                 string             `json:"vpcId"`
	Type                  string             `json:"type"`
	AvailabilityZones     []AvailabilityZone `json:"availabilityZones"`
	SecurityGroups        []string           `json:"securityGroups"`
	IPAddressType         string             `json:"ipAddressType"`
}

type AvailabilityZone struct {
	ZoneName              string        `json:"zoneName"`
	SubnetID              string        `json:"subnetId"`
	LoadBalancerAddresses []interface{} `json:"loadBalancerAddresses"`
}
