package viewmodel

import "time"

type AwsConfigFileData struct {
	RelatedEvents                []string          `json:"relatedEvents"`
	Relationships                []Relationship    `json:"relationships"`
	Configuration                interface{}       `json:"configuration"`
	Tags                         map[string]string `json:"tags"`
	ConfigurationItemVersion     string            `json:"configurationItemVersion"`
	ConfigurationItemCaptureTime time.Time         `json:"configurationItemCaptureTime"`
	ConfigurationStateID         int64             `json:"configurationStateId"`
	AwsAccountID                 string            `json:"awsAccountId"`
	ConfigurationItemStatus      string            `json:"configurationItemStatus"`
	ResourceType                 string            `json:"resourceType"`
	ResourceID                   string            `json:"resourceId"`
	ResourceName                 string            `json:"resourceName"`
	ARN                          string            `json:"ARN"`
	AwsRegion                    string            `json:"awsRegion"`
	AvailabilityZone             string            `json:"availabilityZone"`
	ConfigurationStateMd5Hash    string            `json:"configurationStateMd5Hash"`
	ResourceCreationTime         time.Time         `json:"resourceCreationTime"`
}

type Relationship struct {
	ResourceID   string `json:"resourceId"`
	ResourceType string `json:"resourceType"`
	Name         string `json:"name"`
}

type Configuration struct {
	CertificateArn          string    `json:"certificateArn"`
	DomainName              string    `json:"domainName"`
	SubjectAlternativeNames []string  `json:"subjectAlternativeNames"`
	Serial                  string    `json:"serial"`
	Subject                 string    `json:"subject"`
	Issuer                  string    `json:"issuer"`
	CreatedAt               time.Time `json:"createdAt"`
	IssuedAt                time.Time `json:"issuedAt"`
	Status                  string    `json:"status"`
	NotBefore               time.Time `json:"notBefore"`
	NotAfter                time.Time `json:"notAfter"`
	KeyAlgorithm            string    `json:"keyAlgorithm"`
	SignatureAlgorithm      string    `json:"signatureAlgorithm"`
	InUseBy                 []string  `json:"inUseBy"`
	Type                    string    `json:"type"`
	SubnetIds               []string  `json:"subnetIds"`
	SecurityGroups          []string  `json:"securityGroups"`
}

type VPCEndpointConfiguration struct {
	VpcEndpointID string   `json:"vpcEndpointId"`
	VpcID         string   `json:"vpcId"`
	ServiceName   string   `json:"serviceName"`
	SubnetIDS     []string `json:"subnetIds"`
	Groups        []Group  `json:"groups"`
}

type Group struct {
	GroupID   string `json:"groupId"`
	GroupName string `json:"groupName"`
}

type VPCConfiguration struct {
	CIDRBlock               string                    `json:"cidrBlock"`
	DHCPOptionsID           string                    `json:"dhcpOptionsId"`
	State                   string                    `json:"state"`
	VpcID                   string                    `json:"vpcId"`
	OwnerID                 string                    `json:"ownerId"`
	InstanceTenancy         string                    `json:"instanceTenancy"`
	CIDRBlockAssociationSet []CIDRBlockAssociationSet `json:"cidrBlockAssociationSet"`
}

type CIDRBlockAssociationSet struct {
	AssociationID string `json:"associationId"`
	CIDRBlock     string `json:"cidrBlock"`
}

type ServiceDiscoveryConfiguration struct {
	ID          string    `json:"Id"`
	NamespaceID string    `json:"NamespaceId"`
	Arn         string    `json:"Arn"`
	Name        string    `json:"Name"`
	Type        string    `json:"Type"`
	Description string    `json:"Description"`
	DNSConfig   DNSConfig `json:"DnsConfig"`
}

type DNSConfig struct {
	NamespaceID   string      `json:"NamespaceId"`
	RoutingPolicy string      `json:"RoutingPolicy"`
	DNSRecords    []DNSRecord `json:"DnsRecords"`
}

type DNSRecord struct {
	Type string `json:"Type"`
	TTL  int64  `json:"TTL"`
}
type SubnetConfiguration struct {
	AvailabilityZone            string        `json:"availabilityZone"`
	AvailabilityZoneID          string        `json:"availabilityZoneId"`
	AvailableIPAddressCount     int64         `json:"availableIpAddressCount"`
	CIDRBlock                   string        `json:"cidrBlock"`
	DefaultForAz                bool          `json:"defaultForAz"`
	MapPublicIPOnLaunch         bool          `json:"mapPublicIpOnLaunch"`
	MapCustomerOwnedIPOnLaunch  bool          `json:"mapCustomerOwnedIpOnLaunch"`
	State                       string        `json:"state"`
	SubnetID                    string        `json:"subnetId"`
	VpcID                       string        `json:"vpcId"`
	OwnerID                     string        `json:"ownerId"`
	AssignIpv6AddressOnCreation bool          `json:"assignIpv6AddressOnCreation"`
	Ipv6CIDRBlockAssociationSet []interface{} `json:"ipv6CidrBlockAssociationSet"`
	Tags                        []Tag         `json:"tags"`
	SubnetArn                   string        `json:"subnetArn"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
