package awsconfigConfiguration

type Ec2Configuration struct {
	AmiLaunchIndex                          int64                            `json:"amiLaunchIndex"`
	ImageID                                 string                           `json:"imageId"`
	InstanceID                              string                           `json:"instanceId"`
	InstanceType                            string                           `json:"instanceType"`
	KeyName                                 string                           `json:"keyName"`
	LaunchTime                              string                           `json:"launchTime"`
	Monitoring                              Monitoring                       `json:"monitoring"`
	Placement                               Placement                        `json:"placement"`
	PrivateDNSName                          string                           `json:"privateDnsName"`
	PrivateIPAddress                        string                           `json:"privateIpAddress"`
	ProductCodes                            []interface{}                    `json:"productCodes"`
	PublicDNSName                           string                           `json:"publicDnsName"`
	PublicIPAddress                         string                           `json:"publicIpAddress"`
	State                                   State                            `json:"state"`
	StateTransitionReason                   string                           `json:"stateTransitionReason"`
	SubnetID                                string                           `json:"subnetId"`
	VpcID                                   string                           `json:"vpcId"`
	Architecture                            string                           `json:"architecture"`
	BlockDeviceMappings                     []BlockDeviceMapping             `json:"blockDeviceMappings"`
	ClientToken                             string                           `json:"clientToken"`
	EbsOptimized                            bool                             `json:"ebsOptimized"`
	EnaSupport                              bool                             `json:"enaSupport"`
	Hypervisor                              string                           `json:"hypervisor"`
	ElasticGPUAssociations                  []interface{}                    `json:"elasticGpuAssociations"`
	ElasticInferenceAcceleratorAssociations []interface{}                    `json:"elasticInferenceAcceleratorAssociations"`
	NetworkInterfaces                       []NetworkInterface               `json:"networkInterfaces"`
	RootDeviceName                          string                           `json:"rootDeviceName"`
	RootDeviceType                          string                           `json:"rootDeviceType"`
	SecurityGroups                          []Group                          `json:"securityGroups"`
	SourceDestCheck                         bool                             `json:"sourceDestCheck"`
	Tags                                    []Tag                            `json:"tags"`
	VirtualizationType                      string                           `json:"virtualizationType"`
	CPUOptions                              CPUOptions                       `json:"cpuOptions"`
	CapacityReservationSpecification        CapacityReservationSpecification `json:"capacityReservationSpecification"`
	HibernationOptions                      HibernationOptions               `json:"hibernationOptions"`
	Licenses                                []interface{}                    `json:"licenses"`
	MetadataOptions                         MetadataOptions                  `json:"metadataOptions"`
	EnclaveOptions                          EnclaveOptions                   `json:"enclaveOptions"`
}

type BlockDeviceMapping struct {
	DeviceName string `json:"deviceName"`
	Ebs        Ebs    `json:"ebs"`
}

type Ebs struct {
	AttachTime          string `json:"attachTime"`
	DeleteOnTermination bool   `json:"deleteOnTermination"`
	Status              string `json:"status"`
	VolumeID            string `json:"volumeId"`
}

type CPUOptions struct {
	CoreCount      int64 `json:"coreCount"`
	ThreadsPerCore int64 `json:"threadsPerCore"`
}

type CapacityReservationSpecification struct {
	CapacityReservationPreference string `json:"capacityReservationPreference"`
}

type EnclaveOptions struct {
	Enabled bool `json:"enabled"`
}

type HibernationOptions struct {
	Configured bool `json:"configured"`
}

type MetadataOptions struct {
	State                   string `json:"state"`
	HTTPTokens              string `json:"httpTokens"`
	HTTPPutResponseHopLimit int64  `json:"httpPutResponseHopLimit"`
	HTTPEndpoint            string `json:"httpEndpoint"`
}

type Monitoring struct {
	State string `json:"state"`
}

type NetworkInterface struct {
	Association        Association        `json:"association"`
	Attachment         Ec2Attachment      `json:"attachment"`
	Description        string             `json:"description"`
	Groups             []Group            `json:"groups"`
	Ipv6Addresses      []interface{}      `json:"ipv6Addresses"`
	MACAddress         string             `json:"macAddress"`
	NetworkInterfaceID string             `json:"networkInterfaceId"`
	OwnerID            string             `json:"ownerId"`
	PrivateDNSName     string             `json:"privateDnsName"`
	PrivateIPAddress   string             `json:"privateIpAddress"`
	PrivateIPAddresses []PrivateIPAddress `json:"privateIpAddresses"`
	SourceDestCheck    bool               `json:"sourceDestCheck"`
	Status             string             `json:"status"`
	SubnetID           string             `json:"subnetId"`
	VpcID              string             `json:"vpcId"`
	InterfaceType      string             `json:"interfaceType"`
}

type Ec2Attachment struct {
	AttachTime          string `json:"attachTime"`
	AttachmentID        string `json:"attachmentId"`
	DeleteOnTermination bool   `json:"deleteOnTermination"`
	DeviceIndex         int64  `json:"deviceIndex"`
	Status              string `json:"status"`
	NetworkCardIndex    int64  `json:"networkCardIndex"`
}

type Placement struct {
	AvailabilityZone string `json:"availabilityZone"`
	GroupName        string `json:"groupName"`
	Tenancy          string `json:"tenancy"`
}

type State struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
