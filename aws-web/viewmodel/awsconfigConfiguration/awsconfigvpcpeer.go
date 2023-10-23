package awsconfigConfiguration

type VPCPeeringConnectionConfiguration struct {
	AccepterVpcInfo        TerVpcInfo    `json:"accepterVpcInfo"`
	RequesterVpcInfo       TerVpcInfo    `json:"requesterVpcInfo"`
	Tags                   []interface{} `json:"tags"`
	VpcPeeringConnectionID string        `json:"vpcPeeringConnectionId"`
}

type TerVpcInfo struct {
	CIDRBlock        string         `json:"cidrBlock"`
	Ipv6CIDRBlockSet []interface{}  `json:"ipv6CidrBlockSet"`
	CIDRBlockSet     []CIDRBlockSet `json:"cidrBlockSet"`
	OwnerID          string         `json:"ownerId"`
	VpcID            string         `json:"vpcId"`
	Region           string         `json:"region"`
}

type CIDRBlockSet struct {
	CIDRBlock string `json:"cidrBlock"`
}
