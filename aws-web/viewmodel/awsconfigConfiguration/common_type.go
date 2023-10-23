package awsconfigConfiguration

type VpcInfo struct {
	VpcId   string
	Subnets []SubnetInfo
}

type SubnetInfo struct {
	Subnet           string
	AvailabilityZone string
}

type Group struct {
	GroupID   string `json:"groupId"`
	GroupName string `json:"groupName"`
}

type DNSRecord struct {
	Type string `json:"Type"`
	TTL  int64  `json:"TTL"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Subnet struct {
	SubnetIdentifier       string                 `json:"subnetIdentifier"`
	SubnetAvailabilityZone SubnetAvailabilityZone `json:"subnetAvailabilityZone"`
	SubnetStatus           string                 `json:"subnetStatus"`
}

type Association struct {
	IPOwnerID     string `json:"ipOwnerId"`
	PublicDNSName string `json:"publicDnsName"`
	PublicIP      string `json:"publicIp"`
}

type PrivateIPAddress struct {
	Association      Association `json:"association"`
	Primary          bool        `json:"primary"`
	PrivateDNSName   string      `json:"privateDnsName"`
	PrivateIPAddress string      `json:"privateIpAddress"`
}
