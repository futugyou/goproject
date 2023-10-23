package awsconfigConfiguration

type NetworkInterfaceConfiguration struct {
	Association        Association                `json:"association"`
	Attachment         NetworkInterfaceAttachment `json:"attachment"`
	AvailabilityZone   string                     `json:"availabilityZone"`
	Description        string                     `json:"description"`
	Groups             []Group                    `json:"groups"`
	InterfaceType      string                     `json:"interfaceType"`
	Ipv6Addresses      []interface{}              `json:"ipv6Addresses"`
	MACAddress         string                     `json:"macAddress"`
	NetworkInterfaceID string                     `json:"networkInterfaceId"`
	OwnerID            string                     `json:"ownerId"`
	PrivateDNSName     string                     `json:"privateDnsName"`
	PrivateIPAddress   string                     `json:"privateIpAddress"`
	PrivateIPAddresses []PrivateIPAddress         `json:"privateIpAddresses"`
	RequesterID        string                     `json:"requesterId"`
	RequesterManaged   bool                       `json:"requesterManaged"`
	SourceDestCheck    bool                       `json:"sourceDestCheck"`
	Status             string                     `json:"status"`
	SubnetID           string                     `json:"subnetId"`
	TagSet             []Tag                      `json:"tagSet"`
	VpcID              string                     `json:"vpcId"`
}

type NetworkInterfaceAttachment struct {
	AttachTime          string `json:"attachTime"`
	AttachmentID        string `json:"attachmentId"`
	DeleteOnTermination bool   `json:"deleteOnTermination"`
	DeviceIndex         int64  `json:"deviceIndex"`
	NetworkCardIndex    int64  `json:"networkCardIndex"`
	InstanceOwnerID     string `json:"instanceOwnerId"`
	Status              string `json:"status"`
}
