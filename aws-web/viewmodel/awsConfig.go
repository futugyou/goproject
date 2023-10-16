package viewmodel

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/futugyousuzu/goproject/awsgolang/entity"
)

type ResourceGraph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

type Edge struct {
	ID     string   `json:"id"`
	Label  string   `json:"label"`
	Source EdgeItem `json:"source"`
	Target EdgeItem `json:"target"`
}

type EdgeItem struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type Node struct {
	ID         string     `json:"id"`
	Label      string     `json:"label"`
	Properties Properties `json:"properties"`
}

type Properties struct {
	AccountID        string `json:"accountId"`
	Arn              string `json:"arn"`
	AvailabilityZone string `json:"availabilityZone"`
	AwsRegion        string `json:"awsRegion"`
	// Configuration                string    `json:"configuration"`
	ConfigurationItemCaptureTime time.Time `json:"configurationItemCaptureTime"`
	ConfigurationItemStatus      string    `json:"configurationItemStatus"`
	ConfigurationStateID         int64     `json:"configurationStateId"`
	ResourceCreationTime         time.Time `json:"resourceCreationTime"`
	ResourceID                   string    `json:"resourceId"`
	ResourceName                 string    `json:"resourceName"`
	ResourceType                 string    `json:"resourceType"`
	Tags                         string    `json:"tags"`
	Version                      string    `json:"version"`
	VpcID                        string    `json:"vpcId"`
	SubnetID                     string    `json:"subnetId"`
	SubnetIDS                    []string  `json:"subnetIds"`
	Title                        string    `json:"title"`
	SecurityGroups               []string  `json:"securityGroups"`
}

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

func getId(arn string, resourceID string) string {
	if len(arn) == 0 {
		return resourceID
	}
	return arn
}

func getName(resourceName string, tags map[string]string) string {
	if len(tags) > 0 {
		if n, ok := tags["Name"]; ok && len(n) > 0 {
			return n
		}
	}

	if len(resourceName) != 0 {
		return resourceName
	}

	return ""
}

func getDataString(con interface{}) string {
	if con == nil {
		return "{}"
	} else {
		d, _ := json.Marshal(con)
		return string(d)
	}
}

func (data AwsConfigFileData) CreateAwsConfigEntity() entity.AwsConfigEntity {
	configuration := getDataString(data.Configuration)
	name := getName(data.ResourceName, data.Tags)
	vpcid, subnetId, subnetIds, securityGroups := getVpcInfo(data.ResourceType, configuration)

	config := entity.AwsConfigEntity{
		ID:                           getId(data.ARN, data.ResourceID),
		Label:                        name,
		AccountID:                    data.AwsAccountID,
		Arn:                          data.ARN,
		AvailabilityZone:             data.AvailabilityZone,
		AwsRegion:                    data.AwsRegion,
		Configuration:                configuration,
		ConfigurationItemCaptureTime: data.ConfigurationItemCaptureTime,
		ConfigurationItemStatus:      data.ConfigurationItemStatus,
		ConfigurationStateID:         data.ConfigurationStateID,
		ResourceCreationTime:         data.ResourceCreationTime,
		ResourceID:                   data.ResourceID,
		ResourceName:                 data.ResourceName,
		ResourceType:                 data.ResourceType,
		Tags:                         getDataString(data.Tags),
		Version:                      data.ConfigurationItemVersion,
		Title:                        name,
		VpcID:                        vpcid,
		SubnetID:                     subnetId,
		SubnetIds:                    subnetIds,
		SecurityGroups:               securityGroups,
	}
	return config
}

func (data AwsConfigFileData) CreateAwsConfigRelationshipEntity(configs []entity.AwsConfigEntity) []entity.AwsConfigRelationshipEntity {
	lists := make([]entity.AwsConfigRelationshipEntity, 0)
	name := getName(data.ResourceName, data.Tags)
	if len(name) == 0 {
		return lists
	}

	for _, ship := range data.Relationships {
		var id string
		for i := 0; i < len(configs); i++ {
			if configs[i].ResourceID == ship.ResourceID && configs[i].ResourceName == ship.ResourceName {
				id = getId(configs[i].Arn, configs[i].ResourceID)
				break
			}
		}

		if len(id) == 0 {
			continue
		}

		if strings.HasPrefix(ship.Name, "Is") {
			relationship := entity.AwsConfigRelationshipEntity{
				ID:          data.ResourceID + "-" + ship.ResourceID,
				SourceID:    getId(data.ARN, data.ResourceID),
				SourceLabel: data.ResourceName,
				Label:       ship.Name,
				TargetID:    id,
				TargetLabel: ship.ResourceName,
			}
			lists = append(lists, relationship)
		}

		if strings.HasPrefix(ship.Name, "Contains") {
			relationship := entity.AwsConfigRelationshipEntity{
				ID:          ship.ResourceID + "-" + data.ResourceID,
				SourceID:    id,
				SourceLabel: ship.ResourceName,
				Label:       ship.Name,
				TargetID:    getId(data.ARN, data.ResourceID),
				TargetLabel: data.ResourceName,
			}
			lists = append(lists, relationship)
		}

	}

	return lists
}

func getVpcInfo(resourceType string, configuration string) (vpcid string, subnetId string, subnetIds []string, securityGroups []string) {
	vpcid = ""
	subnetId = ""
	subnetIds = make([]string, 0)
	securityGroups = make([]string, 0)
	switch resourceType {
	case "AWS::EC2::VPCEndpoint":
		var config VPCEndpointConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			vpcid = config.VpcID
			subnetIds = config.SubnetIDS
			for _, v := range config.Groups {
				securityGroups = append(securityGroups, v.GroupID)
			}
		}
	case "AWS::EC2::VPC":
		var config VPCConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			vpcid = config.VpcID
		}
	case "AWS::EC2::Subnet":
		var config SubnetConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			vpcid = config.VpcID
			subnetId = config.SubnetID
		}
	case "AWS::AmazonMQ::Broker":
		var config AmazonMQConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			subnetIds = config.SubnetIDS
			securityGroups = config.SecurityGroups
		}
	case "AWS::EC2::NatGateway":
		var config NatGatewayConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			subnetId = config.SubnetID
			vpcid = config.VpcID
		}
	case "AWS::RDS::DBInstance":
		var config DBInstanceConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			for _, v := range config.VpcSecurityGroups {
				securityGroups = append(securityGroups, v.VpcSecurityGroupID)
			}
			vpcid = config.DBSubnetGroup.VpcID
			for _, v := range config.DBSubnetGroup.Subnets {
				subnetIds = append(subnetIds, v.SubnetIdentifier)
			}
		}
	case "AWS::EC2::SecurityGroup":
		var config SecurityGroupConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			vpcid = config.VpcID
			securityGroups = []string{config.GroupID}
		}
	case "AWS::EC2::NetworkInterface":
		var config NetworkInterfaceConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			vpcid = config.VpcID
			subnetId = config.SubnetID
			for _, v := range config.Groups {
				securityGroups = append(securityGroups, v.GroupID)
			}
		}
	case "AWS::Redshift::ClusterSubnetGroup":
		var config RedshiftConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			vpcid = config.VpcID

			for _, v := range config.Subnets {
				subnetIds = append(subnetIds, v.SubnetIdentifier)
			}
		}
	case "AWS::ElasticLoadBalancingV2::LoadBalancer":
		var config LoadBalancerConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			vpcid = config.VpcID
			securityGroups = config.SecurityGroups
			for _, v := range config.AvailabilityZones {
				subnetIds = append(subnetIds, v.SubnetID)
			}
		}
	case "AWS::ECS::Service":
		var config ECSServiceConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			subnetIds = config.NetworkConfiguration.AwsvpcConfiguration.Subnets
			securityGroups = config.NetworkConfiguration.AwsvpcConfiguration.SecurityGroups
		}
	case "AWS::EC2::NetworkAcl":
		var config NetworkAclConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			vpcid = config.VpcID
			for _, v := range config.Associations {
				subnetIds = append(subnetIds, v.SubnetID)
			}
		}
	case "AWS::Lambda::Function":
		var config FunctionConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			subnetIds = config.VpcConfig.SubnetIDS
			securityGroups = config.VpcConfig.SecurityGroupIDS
		}
	case "AWS::EC2::RouteTable":
		var config RouteTableConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			vpcid = config.VpcID
		}
	}
	return
}

type Relationship struct {
	ResourceID   string `json:"resourceId"`
	ResourceName string `json:"resourceName"`
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

type AmazonMQConfiguration struct {
	SecurityGroups          []string `json:"SecurityGroups"`
	SubnetIDS               []string `json:"SubnetIds"`
	DeploymentMode          string   `json:"DeploymentMode"`
	EngineType              string   `json:"EngineType"`
	Tags                    []Tag    `json:"Tags"`
	ConfigurationRevision   int64    `json:"ConfigurationRevision"`
	StorageType             string   `json:"StorageType"`
	EngineVersion           string   `json:"EngineVersion"`
	HostInstanceType        string   `json:"HostInstanceType"`
	AutoMinorVersionUpgrade bool     `json:"AutoMinorVersionUpgrade"`
}
type NatGatewayConfiguration struct {
	CreateTime          int64               `json:"createTime"`
	NatGatewayAddresses []NatGatewayAddress `json:"natGatewayAddresses"`
	NatGatewayID        string              `json:"natGatewayId"`
	State               string              `json:"state"`
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

type InternetGatewayConfiguration struct {
	Attachments       []Attachment `json:"attachments"`
	InternetGatewayID string       `json:"internetGatewayId"`
	OwnerID           string       `json:"ownerId"`
	Tags              []Tag        `json:"tags"`
}

type Attachment struct {
	State string `json:"state"`
	VpcID string `json:"vpcId"`
}

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

type DBInstanceConfiguration struct {
	DBInstanceIdentifier                   string                  `json:"dBInstanceIdentifier"`
	DBInstanceClass                        string                  `json:"dBInstanceClass"`
	Engine                                 string                  `json:"engine"`
	DBInstanceStatus                       string                  `json:"dBInstanceStatus"`
	MasterUsername                         string                  `json:"masterUsername"`
	Endpoint                               Endpoint                `json:"endpoint"`
	AllocatedStorage                       int64                   `json:"allocatedStorage"`
	InstanceCreateTime                     string                  `json:"instanceCreateTime"`
	PreferredBackupWindow                  string                  `json:"preferredBackupWindow"`
	BackupRetentionPeriod                  int64                   `json:"backupRetentionPeriod"`
	DBSecurityGroups                       []interface{}           `json:"dBSecurityGroups"`
	VpcSecurityGroups                      []VpcSecurityGroup      `json:"vpcSecurityGroups"`
	AvailabilityZone                       string                  `json:"availabilityZone"`
	DBSubnetGroup                          DBSubnetGroup           `json:"dBSubnetGroup"`
	PreferredMaintenanceWindow             string                  `json:"preferredMaintenanceWindow"`
	PendingModifiedValues                  PendingModifiedValues   `json:"pendingModifiedValues"`
	LatestRestorableTime                   string                  `json:"latestRestorableTime"`
	MultiAZ                                bool                    `json:"multiAZ"`
	EngineVersion                          string                  `json:"engineVersion"`
	AutoMinorVersionUpgrade                bool                    `json:"autoMinorVersionUpgrade"`
	ReadReplicaDBInstanceIdentifiers       []interface{}           `json:"readReplicaDBInstanceIdentifiers"`
	ReadReplicaDBClusterIdentifiers        []interface{}           `json:"readReplicaDBClusterIdentifiers"`
	LicenseModel                           string                  `json:"licenseModel"`
	OptionGroupMemberships                 []OptionGroupMembership `json:"optionGroupMemberships"`
	PubliclyAccessible                     bool                    `json:"publiclyAccessible"`
	StatusInfos                            []interface{}           `json:"statusInfos"`
	StorageType                            string                  `json:"storageType"`
	DBInstancePort                         int64                   `json:"dbInstancePort"`
	StorageEncrypted                       bool                    `json:"storageEncrypted"`
	KmsKeyID                               string                  `json:"kmsKeyId"`
	DbiResourceID                          string                  `json:"dbiResourceId"`
	CACertificateIdentifier                string                  `json:"cACertificateIdentifier"`
	DomainMemberships                      []interface{}           `json:"domainMemberships"`
	CopyTagsToSnapshot                     bool                    `json:"copyTagsToSnapshot"`
	MonitoringInterval                     int64                   `json:"monitoringInterval"`
	EnhancedMonitoringResourceArn          string                  `json:"enhancedMonitoringResourceArn"`
	MonitoringRoleArn                      string                  `json:"monitoringRoleArn"`
	DBInstanceArn                          string                  `json:"dBInstanceArn"`
	IAMDatabaseAuthenticationEnabled       bool                    `json:"iAMDatabaseAuthenticationEnabled"`
	PerformanceInsightsEnabled             bool                    `json:"performanceInsightsEnabled"`
	EnabledCloudwatchLogsExports           []interface{}           `json:"enabledCloudwatchLogsExports"`
	ProcessorFeatures                      []interface{}           `json:"processorFeatures"`
	DeletionProtection                     bool                    `json:"deletionProtection"`
	AssociatedRoles                        []interface{}           `json:"associatedRoles"`
	MaxAllocatedStorage                    int64                   `json:"maxAllocatedStorage"`
	TagList                                []interface{}           `json:"tagList"`
	DBInstanceAutomatedBackupsReplications []interface{}           `json:"dBInstanceAutomatedBackupsReplications"`
	CustomerOwnedIPEnabled                 bool                    `json:"customerOwnedIpEnabled"`
}

type DBSubnetGroup struct {
	DBSubnetGroupName        string   `json:"dBSubnetGroupName"`
	DBSubnetGroupDescription string   `json:"dBSubnetGroupDescription"`
	VpcID                    string   `json:"vpcId"`
	SubnetGroupStatus        string   `json:"subnetGroupStatus"`
	Subnets                  []Subnet `json:"subnets"`
}

type Subnet struct {
	SubnetIdentifier       string                 `json:"subnetIdentifier"`
	SubnetAvailabilityZone SubnetAvailabilityZone `json:"subnetAvailabilityZone"`
	SubnetStatus           string                 `json:"subnetStatus"`
}

type SubnetAvailabilityZone struct {
	Name string `json:"name"`
}

type Endpoint struct {
	Address      string `json:"address"`
	Port         int64  `json:"port"`
	HostedZoneID string `json:"hostedZoneId"`
}

type OptionGroupMembership struct {
	OptionGroupName string `json:"optionGroupName"`
	Status          string `json:"status"`
}

type PendingModifiedValues struct {
	ProcessorFeatures []interface{} `json:"processorFeatures"`
}

type VpcSecurityGroup struct {
	VpcSecurityGroupID string `json:"vpcSecurityGroupId"`
	Status             string `json:"status"`
}

type SecurityGroupConfiguration struct {
	Description         string         `json:"description"`
	GroupName           string         `json:"groupName"`
	IPPermissions       []IPPermission `json:"ipPermissions"`
	OwnerID             string         `json:"ownerId"`
	GroupID             string         `json:"groupId"`
	IPPermissionsEgress []IPPermission `json:"ipPermissionsEgress"`
	Tags                []interface{}  `json:"tags"`
	VpcID               string         `json:"vpcId"`
}

type IPPermission struct {
	IPProtocol       string            `json:"ipProtocol"`
	Ipv6Ranges       []interface{}     `json:"ipv6Ranges"`
	PrefixListIDS    []interface{}     `json:"prefixListIds"`
	UserIDGroupPairs []UserIDGroupPair `json:"userIdGroupPairs"`
	Ipv4Ranges       []Ipv4Range       `json:"ipv4Ranges"`
	IPRanges         []string          `json:"ipRanges"`
}

type Ipv4Range struct {
	CIDRIP string `json:"cidrIp"`
}

type UserIDGroupPair struct {
	GroupID string `json:"groupId"`
	UserID  string `json:"userId"`
}

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

type Association struct {
	IPOwnerID     string `json:"ipOwnerId"`
	PublicDNSName string `json:"publicDnsName"`
	PublicIP      string `json:"publicIp"`
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

type PrivateIPAddress struct {
	Association      Association `json:"association"`
	Primary          bool        `json:"primary"`
	PrivateDNSName   string      `json:"privateDnsName"`
	PrivateIPAddress string      `json:"privateIpAddress"`
}

type RedshiftConfiguration struct {
	ClusterSubnetGroupName string   `json:"clusterSubnetGroupName"`
	Description            string   `json:"description"`
	VpcID                  string   `json:"vpcId"`
	SubnetGroupStatus      string   `json:"subnetGroupStatus"`
	Subnets                []Subnet `json:"subnets"`
	Tags                   []Tag    `json:"tags"`
}

type LoadBalancerConfiguration struct {
	LoadBalancerArn       string             `json:"loadBalancerArn"`
	DNSName               string             `json:"dNSName"`
	CanonicalHostedZoneID string             `json:"canonicalHostedZoneId"`
	CreatedTime           string             `json:"createdTime"`
	LoadBalancerName      string             `json:"loadBalancerName"`
	Scheme                string             `json:"scheme"`
	VpcID                 string             `json:"vpcId"`
	Type                  string             `json:"type"`
	AvailabilityZones     []AvailabilityZone `json:"availabilityZones"`
	SecurityGroups        []string           `json:"securityGroups"`
	IPAddressType         string             `json:"ipAddressType"`
}

type AvailabilityZone struct {
	ZoneName              string        `json:"zoneName"`
	SubnetID              string        `json:"subnetId"`
	LoadBalancerAddresses []interface{} `json:"loadBalancerAddresses"`
}

type ECSServiceConfiguration struct {
	ServiceArn               string                  `json:"ServiceArn"`
	CapacityProviderStrategy []interface{}           `json:"CapacityProviderStrategy"`
	Cluster                  string                  `json:"Cluster"`
	DeploymentConfiguration  DeploymentConfiguration `json:"DeploymentConfiguration"`
	DeploymentController     DeploymentController    `json:"DeploymentController"`
	DesiredCount             int64                   `json:"DesiredCount"`
	EnableECSManagedTags     bool                    `json:"EnableECSManagedTags"`
	LaunchType               string                  `json:"LaunchType"`
	LoadBalancers            []interface{}           `json:"LoadBalancers"`
	Name                     string                  `json:"Name"`
	NetworkConfiguration     NetworkConfiguration    `json:"NetworkConfiguration"`
	PlacementConstraints     []interface{}           `json:"PlacementConstraints"`
	PlacementStrategies      []interface{}           `json:"PlacementStrategies"`
	PlatformVersion          string                  `json:"PlatformVersion"`
	Role                     string                  `json:"Role"`
	SchedulingStrategy       string                  `json:"SchedulingStrategy"`
	ServiceName              string                  `json:"ServiceName"`
	ServiceRegistries        []ServiceRegistry       `json:"ServiceRegistries"`
	Tags                     []interface{}           `json:"Tags"`
	TaskDefinition           string                  `json:"TaskDefinition"`
}

type DeploymentConfiguration struct {
	DeploymentCircuitBreaker DeploymentCircuitBreaker `json:"DeploymentCircuitBreaker"`
	MaximumPercent           int64                    `json:"MaximumPercent"`
	MinimumHealthyPercent    int64                    `json:"MinimumHealthyPercent"`
}

type DeploymentCircuitBreaker struct {
	Enable   bool `json:"Enable"`
	Rollback bool `json:"Rollback"`
}

type DeploymentController struct {
	Type string `json:"Type"`
}

type NetworkConfiguration struct {
	AwsvpcConfiguration AwsvpcConfiguration `json:"AwsvpcConfiguration"`
}

type AwsvpcConfiguration struct {
	Subnets        []string `json:"Subnets"`
	SecurityGroups []string `json:"SecurityGroups"`
	AssignPublicIP string   `json:"AssignPublicIp"`
}

type ServiceRegistry struct {
	RegistryArn string `json:"RegistryArn"`
}

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

type FunctionConfiguration struct {
	FunctionName         string               `json:"functionName"`
	FunctionArn          string               `json:"functionArn"`
	Runtime              string               `json:"runtime"`
	Role                 string               `json:"role"`
	Handler              string               `json:"handler"`
	CodeSize             int64                `json:"codeSize"`
	Description          string               `json:"description"`
	Timeout              int64                `json:"timeout"`
	MemorySize           int64                `json:"memorySize"`
	LastModified         string               `json:"lastModified"`
	CodeSha256           string               `json:"codeSha256"`
	Version              string               `json:"version"`
	VpcConfig            VpcConfig            `json:"vpcConfig"`
	TracingConfig        TracingConfig        `json:"tracingConfig"`
	RevisionID           string               `json:"revisionId"`
	Layers               []interface{}        `json:"layers"`
	State                string               `json:"state"`
	LastUpdateStatus     string               `json:"lastUpdateStatus"`
	FileSystemConfigs    []interface{}        `json:"fileSystemConfigs"`
	PackageType          string               `json:"packageType"`
	Architectures        []string             `json:"architectures"`
	EphemeralStorage     EphemeralStorage     `json:"ephemeralStorage"`
	SnapStart            SnapStart            `json:"snapStart"`
	RuntimeVersionConfig RuntimeVersionConfig `json:"runtimeVersionConfig"`
}

type EphemeralStorage struct {
	Size int64 `json:"size"`
}

type RuntimeVersionConfig struct {
	RuntimeVersionArn string `json:"runtimeVersionArn"`
}

type SnapStart struct {
	ApplyOn            string `json:"applyOn"`
	OptimizationStatus string `json:"optimizationStatus"`
}

type TracingConfig struct {
	Mode string `json:"mode"`
}

type VpcConfig struct {
	SubnetIDS        []string `json:"subnetIds"`
	SecurityGroupIDS []string `json:"securityGroupIds"`
}

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
