package viewmodel

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/futugyousuzu/goproject/awsgolang/entity"

	"github.com/futugyousuzu/goproject/awsgolang/sdk/route53"
	"github.com/futugyousuzu/goproject/awsgolang/sdk/servicediscovery"
	"golang.org/x/exp/slices"
)

func AddIndividualData(configs []entity.AwsConfigEntity) []entity.AwsConfigEntity {
	namespaceEntityList := make([]entity.AwsConfigEntity, 0)
	namespaceList := make([]string, 0)

	// 1. get namespace data
	for i := 0; i < len(configs); i++ {
		if configs[i].ResourceType == "AWS::ServiceDiscovery::Service" {
			var configuration ServiceDiscoveryConfiguration
			err := json.Unmarshal([]byte(configs[i].Configuration), &configuration)
			if err != nil {
				continue
			}

			namespaceid := configuration.NamespaceID
			if slices.Contains(namespaceList, namespaceid) {
				for _, namespace := range namespaceEntityList {
					if namespace.ResourceID == namespaceid {
						configs[i].VpcID = namespace.VpcID
						break
					}
				}
				continue
			}

			namespaceList = append(namespaceList, namespaceid)
			namespace := servicediscovery.GetNamespaceDetail(namespaceid)

			if namespace == nil {
				continue
			}

			namespaceEntity := entity.AwsConfigEntity{
				ID:                           *namespace.Arn,
				Label:                        *namespace.Name,
				AccountID:                    configs[i].AccountID,
				Arn:                          *namespace.Arn,
				AvailabilityZone:             configs[i].AvailabilityZone,
				AwsRegion:                    configs[i].AwsRegion,
				Configuration:                "{}",
				ConfigurationItemCaptureTime: *namespace.CreateDate,
				ConfigurationItemStatus:      "",
				ConfigurationStateID:         0,
				ResourceCreationTime:         *namespace.CreateDate,
				ResourceID:                   *namespace.Id,
				ResourceName:                 *namespace.Name,
				ResourceType:                 "AWS::ServiceDiscovery::Namespace",
				Tags:                         "",
				Version:                      "",
				VpcID:                        "",
				SubnetID:                     "",
				SubnetIds:                    []string{},
				Title:                        *namespace.Name,
				SecurityGroups:               []string{},
				LoginURL:                     "",
				LoggedInURL:                  "",
			}

			if namespace.Properties != nil &&
				namespace.Properties.DnsProperties != nil &&
				namespace.Properties.DnsProperties.HostedZoneId != nil {
				vpcid := route53.GetHostedZoneVpcId(*namespace.Properties.DnsProperties.HostedZoneId)
				namespaceEntity.VpcID = vpcid
				configs[i].VpcID = vpcid
			}

			namespaceEntityList = append(namespaceEntityList, namespaceEntity)
		}
	}

	return append(configs, namespaceEntityList...)
}

func GetAllVpcInfos(datas []AwsConfigFileData) []VpcInfo {
	vpcInfos := make([]VpcInfo, 0)
	for _, data := range datas {
		if data.ResourceType == "AWS::EC2::Subnet" {
			var config SubnetConfiguration
			json.Unmarshal([]byte(getDataString(data.Configuration)), &config)

			vpcid := config.VpcID
			subnetId := config.SubnetID
			availabilityZone := config.AvailabilityZone

			found := false
			for i := 0; i < len(vpcInfos); i++ {
				if vpcid == vpcInfos[i].VpcId {
					found = true
					vpcInfos[i].Subnets = append(vpcInfos[i].Subnets, SubnetInfo{
						Subnet:           subnetId,
						AvailabilityZone: availabilityZone,
					})
				}
			}

			if !found {
				vpc := VpcInfo{
					VpcId: vpcid,
					Subnets: []SubnetInfo{{
						Subnet:           subnetId,
						AvailabilityZone: availabilityZone,
					}},
				}
				vpcInfos = append(vpcInfos, vpc)
			}
		}

	}
	return vpcInfos
}

func getId(arn string, resourceID string) string {
	if len(arn) == 0 {
		return resourceID
	}
	return arn
}

func getName(resourceID string, resourceName string, tags map[string]string) string {
	if len(tags) > 0 {
		if n, ok := tags["Name"]; ok && len(n) > 0 {
			return n
		}
	}

	if len(resourceName) != 0 {
		return resourceName
	}

	return resourceID
}

func getDataString(con interface{}) string {
	if con == nil {
		return "{}"
	} else {
		d, _ := json.Marshal(con)
		return string(d)
	}
}

func getVpcInfo(resourceType string, configuration string, vpcinfos []VpcInfo) (vpcid string, subnetId string, subnetIds []string, securityGroups []string, availabilityZone string) {
	vpcid = ""
	subnetId = ""
	subnetIds = make([]string, 0)
	securityGroups = make([]string, 0)
	availabilityZone = ""

	foundVpn := func(ids []string) string {
		for _, vpn := range vpcinfos {
			for _, net := range vpn.Subnets {
				if slices.Contains(ids, net.Subnet) {
					return vpn.VpcId
				}
			}
		}

		return ""
	}

	foundAvailabilityZone := func(ids []string) string {
		azs := make([]string, 0)
		for _, vpn := range vpcinfos {
			for _, net := range vpn.Subnets {
				if !slices.Contains(ids, net.Subnet) ||
					slices.Contains(azs, net.AvailabilityZone) {
					continue
				}

				azs = append(azs, net.AvailabilityZone)
			}
		}

		if len(azs) > 0 {
			slices.Sort(azs)
			return strings.Join(azs, ",")
		}

		return ""
	}
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
	case "AWS::EC2::Instance":
		var config Ec2Configuration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			vpcid = config.VpcID
			subnetId = config.SubnetID
			for _, sg := range config.SecurityGroups {
				securityGroups = append(securityGroups, sg.GroupID)
			}
		}
	}

	slices.Sort(subnetIds)

	if vpcid == "" {
		if len(subnetIds) > 0 {
			vpcid = foundVpn(subnetIds)
		} else if len(subnetId) > 0 {
			vpcid = foundVpn([]string{subnetId})
		}
	}

	if len(subnetId) == 0 && len(subnetIds) == 1 {
		subnetId = subnetIds[0]
	}

	if len(subnetIds) > 1 {
		availabilityZone = foundAvailabilityZone(subnetIds)
		subnetId = ""
	}

	return
}

func createSignInHostname(accountId string, service string) string {
	return fmt.Sprintf("https://%s.signin.aws.amazon.com/console/%s", accountId, service)
}

func createLoggedInHostname(awsRegion string, service string) string {
	return fmt.Sprintf(`https://%s.console.aws.amazon.com%s/home`, awsRegion, service)
}

func createConsoleUrls(resource AwsConfigFileData) (loginURL string, loggedInURL string) {
	resourceType := resource.ResourceType
	resourceName := resource.ResourceName
	accountId := resource.AwsAccountID
	awsRegion := resource.AwsRegion
	loginURL = ""
	loggedInURL = ""
	switch resourceType {
	case "AWS::Lambda::Function":
		loginURL =
			fmt.Sprintf(`%s?region=%s#/functions/%s?tab=graph`,
				createSignInHostname(accountId, "lambda"),
				awsRegion,
				resourceName,
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#/functions/%s?tab=graph`,
				createLoggedInHostname(awsRegion, "lambda"),
				awsRegion,
				resourceName,
			)
	case "AWS::IAM::Policy":
		loginURL =
			fmt.Sprintf(`%s?home?#%s`,
				createSignInHostname(accountId, "iam"),
				"/policies",
			)
		loggedInURL =
			fmt.Sprintf(`https://console.aws.amazon.com/%s/home?#%s`,
				"iam",
				"/policies",
			)
	case "AWS::S3::Bucket":
		loginURL =
			fmt.Sprintf(`%s?bucket=%s`,
				createSignInHostname(accountId, "s3"),
				resourceName,
			)
		loggedInURL =
			fmt.Sprintf(`https://s3.console.aws.amazon.com/s3/buckets/%s/?region=%s`,
				resourceName,
				awsRegion,
			)
	case "AWS::EC2::VPC":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "vpc"),
				awsRegion,
				"vpcs:sort=VpcId",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "vpc/v2"),
				awsRegion,
				"vpcs:sort=VpcId",
			)
	case "AWS::EC2::NetworkInterface":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "ec2"),
				awsRegion,
				"NIC:sort=description",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "ec2/v2"),
				awsRegion,
				"NIC:sort=description",
			)
	case "AWS::EC2::Instance":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "ec2"),
				awsRegion,
				"Instances:sort=instanceId",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "ec2/v2"),
				awsRegion,
				"Instances:sort=instanceId",
			)
	case "AWS::EC2::Volume":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "ec2"),
				awsRegion,
				"Volumes:sort=desc:name",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "ec2/v2"),
				awsRegion,
				"Volumes:sort=desc:name",
			)
	case "AWS::EC2::Subnet":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "vpc"),
				awsRegion,
				"subnets:sort=SubnetId",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "vpc/v2"),
				awsRegion,
				"subnets:sort=SubnetId",
			)
	case "AWS::EC2::SecurityGroup":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "ec2"),
				awsRegion,
				"SecurityGroups:sort=groupId",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "ec2/v2"),
				awsRegion,
				"SecurityGroups:sort=groupId",
			)
	case "AWS::EC2::RouteTable":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "vpc"),
				awsRegion,
				"RouteTables:sort=routeTableId",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "vpc/v2"),
				awsRegion,
				"RouteTables:sort=routeTableId",
			)
	case "AWS::EC2::InternetGateway":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "vpc"),
				awsRegion,
				"igws:sort=internetGatewayId",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "vpc/v2"),
				awsRegion,
				"igws:sort=internetGatewayId",
			)
	case "AWS::EC2::NetworkAcl":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "vpc"),
				awsRegion,
				"acls:sort=networkAclId",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "vpc/v2"),
				awsRegion,
				"acls:sort=networkAclId",
			)
	case "AWS::ElasticLoadBalancingV2::LoadBalancer":
	case "AWS::ElasticLoadBalancingV2::Listener":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "ec2"),
				awsRegion,
				"LoadBalancers:",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "ec2/v2"),
				awsRegion,
				"LoadBalancers:",
			)

	case "AWS::ElasticLoadBalancingV2::TargetGroup":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "ec2"),
				awsRegion,
				"TargetGroups:",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "ec2/v2"),
				awsRegion,
				"TargetGroups:",
			)
	case "AWS::EC2::EIP":
		loginURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createSignInHostname(accountId, "ec2"),
				awsRegion,
				"Addresses:sort=PublicIp",
			)
		loggedInURL =
			fmt.Sprintf(`%s?region=%s#%s`,
				createLoggedInHostname(awsRegion, "ec2/v2"),
				awsRegion,
				"Addresses:sort=PublicIp",
			)
	}

	return
}

type VpcInfo struct {
	VpcId   string
	Subnets []SubnetInfo
}

type SubnetInfo struct {
	Subnet           string
	AvailabilityZone string
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
	Code int64  `json:"code"`
	Name string `json:"name"`
}
