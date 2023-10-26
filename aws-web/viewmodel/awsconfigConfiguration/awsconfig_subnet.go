package awsconfigConfiguration

type SubnetConfiguration struct {
	AvailabilityZone            string        `json:"availabilityZone"`
	AvailabilityZoneID          string        `json:"availabilityZoneId"`
	AvailableIPAddressCount     int64         `json:"availableIpAddressCount"`
	CIDRBlock                   string        `json:"cidrBlock"`
	DefaultForAz                bool          `json:"defaultForAz"`
	MapPublicIPOnLaunch         bool          `json:"mapPublicIpOnLaunch"`
	MapCustomerOwnedIPOnLaunch  bool          `json:"mapCustomerOwnedIpOnLaunch"`
	State                       interface{}   `json:"state"`
	SubnetID                    string        `json:"subnetId"`
	VpcID                       string        `json:"vpcId"`
	OwnerID                     string        `json:"ownerId"`
	AssignIpv6AddressOnCreation bool          `json:"assignIpv6AddressOnCreation"`
	Ipv6CIDRBlockAssociationSet []interface{} `json:"ipv6CidrBlockAssociationSet"`
	Tags                        []Tag         `json:"tags"`
	SubnetArn                   string        `json:"subnetArn"`
}
