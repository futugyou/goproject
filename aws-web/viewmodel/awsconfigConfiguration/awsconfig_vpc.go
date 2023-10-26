package awsconfigConfiguration

type VPCConfiguration struct {
	CIDRBlock               string                    `json:"cidrBlock"`
	DHCPOptionsID           string                    `json:"dhcpOptionsId"`
	State                   interface{}               `json:"state"`
	VpcID                   string                    `json:"vpcId"`
	OwnerID                 string                    `json:"ownerId"`
	InstanceTenancy         string                    `json:"instanceTenancy"`
	CIDRBlockAssociationSet []CIDRBlockAssociationSet `json:"cidrBlockAssociationSet"`
}

type CIDRBlockAssociationSet struct {
	AssociationID string `json:"associationId"`
	CIDRBlock     string `json:"cidrBlock"`
}
