package awsconfigConfiguration

type RedshiftConfiguration struct {
	ClusterSubnetGroupName string   `json:"clusterSubnetGroupName"`
	Description            string   `json:"description"`
	VpcID                  string   `json:"vpcId"`
	SubnetGroupStatus      string   `json:"subnetGroupStatus"`
	Subnets                []Subnet `json:"subnets"`
	Tags                   []Tag    `json:"tags"`
}
