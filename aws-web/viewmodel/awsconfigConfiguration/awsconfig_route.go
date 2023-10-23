package awsconfigConfiguration

type RouteTableConfiguration struct {
	RouteTableID string  `json:"routeTableId"`
	Routes       []Route `json:"routes"`
	Tags         []Tag   `json:"tags"`
	VpcID        string  `json:"vpcId"`
	OwnerID      string  `json:"ownerId"`
}

type Route struct {
	DestinationCIDRBlock string `json:"destinationCidrBlock"`
	GatewayID            string `json:"gatewayId"`
	Origin               string `json:"origin"`
	State                string `json:"state"`
}
