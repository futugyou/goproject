package awsconfigConfiguration

type NatGatewayConfiguration struct {
	CreateTime          int64               `json:"createTime"`
	NatGatewayAddresses []NatGatewayAddress `json:"natGatewayAddresses"`
	NatGatewayID        string              `json:"natGatewayId"`
	State               interface{}         `json:"state"`
	SubnetID            string              `json:"subnetId"`
	VpcID               string              `json:"vpcId"`
	Tags                []Tag               `json:"tags"`
	ConnectivityType    string              `json:"connectivityType"`
}

type NatGatewayAddress struct {
	AllocationID       string `json:"allocationId"`
	NetworkInterfaceID string `json:"networkInterfaceId"`
	PrivateIP          string `json:"privateIp"`
	PublicIP           string `json:"publicIp"`
	AssociationID      string `json:"associationId"`
	IsPrimary          bool   `json:"isPrimary"`
	Status             string `json:"status"`
	Primary            bool   `json:"primary"`
}
