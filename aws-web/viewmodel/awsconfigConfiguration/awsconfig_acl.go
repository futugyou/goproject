package awsconfigConfiguration

type NetworkAclConfiguration struct {
	Associations []NetworkAclAssociation `json:"associations"`
	Entries      []Entry                 `json:"entries"`
	IsDefault    bool                    `json:"isDefault"`
	NetworkACLID string                  `json:"networkAclId"`
	Tags         []interface{}           `json:"tags"`
	VpcID        string                  `json:"vpcId"`
	OwnerID      string                  `json:"ownerId"`
}

type NetworkAclAssociation struct {
	NetworkACLAssociationID string `json:"networkAclAssociationId"`
	NetworkACLID            string `json:"networkAclId"`
	SubnetID                string `json:"subnetId"`
}

type Entry struct {
	CIDRBlock  string `json:"cidrBlock"`
	Egress     bool   `json:"egress"`
	Protocol   string `json:"protocol"`
	RuleAction string `json:"ruleAction"`
	RuleNumber int64  `json:"ruleNumber"`
}
