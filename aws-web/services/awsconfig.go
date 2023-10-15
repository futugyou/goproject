package services

import (
	"encoding/json"
	"io"
	"os"

	"log"

	"github.com/futugyousuzu/goproject/awsgolang/entity"
	"github.com/futugyousuzu/goproject/awsgolang/repository"
	"github.com/futugyousuzu/goproject/awsgolang/repository/mongorepo"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"
)

type AwsConfigService struct {
	repository    repository.IAwsConfigRepository
	relRepository repository.IAwsConfigRelationshipRepository
}

func NewAwsConfigService() *AwsConfigService {
	config := mongorepo.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	return &AwsConfigService{
		repository:    mongorepo.NewAwsConfigRepository(config),
		relRepository: mongorepo.NewAwsConfigRelationshipRepository(config),
	}
}

func (a *AwsConfigService) SyncFileResources(path string) {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)

	var datas []model.AwsConfigFileData

	json.Unmarshal(byteValue, &datas)

	if len(datas) == 0 {
		return
	}
	datas = filterResource(datas)
	configs := make([]entity.AwsConfigEntity, 0)

	for _, data := range datas {
		config := createAwsConfigEntity(data)
		configs = append(configs, config)
	}

	for _, v := range configs {
		log.Println(v.VpcID, v.SubnetID, v.SubnetIds, v.SecurityGroups)
	}
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

func createAwsConfigEntity(data model.AwsConfigFileData) entity.AwsConfigEntity {
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

func getVpcInfo(resourceType string, configuration string) (vpcid string, subnetId string, subnetIds []string, securityGroups []string) {
	vpcid = ""
	subnetId = ""
	subnetIds = make([]string, 0)
	securityGroups = make([]string, 0)
	switch resourceType {
	case "AWS::EC2::VPCEndpoint":
		var config model.VPCEndpointConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			vpcid = config.VpcID
			subnetIds = config.SubnetIDS
			for _, v := range config.Groups {
				securityGroups = append(securityGroups, v.GroupID)
			}
		}
	case "AWS::EC2::VPC":
		var config model.VPCConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			vpcid = config.VpcID
		}
	case "AWS::EC2::Subnet":
		var config model.SubnetConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			vpcid = config.VpcID
			subnetId = config.SubnetID
		}
	case "AWS::AmazonMQ::Broker":
		var config model.AmazonMQConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			subnetIds = config.SubnetIDS
			securityGroups = config.SecurityGroups
		}
	case "AWS::EC2::NatGateway":
		var config model.NatGatewayConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			subnetId = config.SubnetID
			vpcid = config.VpcID
		}
	case "AWS::EC2::InternetGateway":
		var config model.InternetGatewayConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			for range config.Attachments {
				// TODO: add to Relationships
			}
		}
	case "AWS::EC2::VPCPeeringConnection":
		var config model.VPCPeeringConnectionConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			// TODO: add to Relationships
		}
	case "AWS::RDS::DBInstance":
		var config model.DBInstanceConfiguration
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
		var config model.SecurityGroupConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			vpcid = config.VpcID
			securityGroups = []string{config.GroupID}
		}
	case "AWS::EC2::NetworkInterface":
		var config model.NetworkInterfaceConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			vpcid = config.VpcID
			subnetId = config.SubnetID
			for _, v := range config.Groups {
				securityGroups = append(securityGroups, v.GroupID)
			}
		}
	case "AWS::Redshift::ClusterSubnetGroup":
		var config model.RedshiftConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			vpcid = config.VpcID

			for _, v := range config.Subnets {
				subnetIds = append(subnetIds, v.SubnetIdentifier)
			}
		}
	case "AWS::ElasticLoadBalancingV2::LoadBalancer":
		var config model.LoadBalancerConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			vpcid = config.VpcID
			securityGroups = config.SecurityGroups
			for _, v := range config.AvailabilityZones {
				subnetIds = append(subnetIds, v.SubnetID)
			}
		}
	case "AWS::ECS::Service":
		var config model.ECSServiceConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			subnetIds = config.NetworkConfiguration.AwsvpcConfiguration.Subnets
			securityGroups = config.NetworkConfiguration.AwsvpcConfiguration.SecurityGroups
		}
	case "AWS::EC2::NetworkAcl":
		var config model.NetworkAclConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			vpcid = config.VpcID
			for _, v := range config.Associations {
				subnetIds = append(subnetIds, v.SubnetID)
			}
		}
	case "AWS::Lambda::Function":
		var config model.FunctionConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			subnetIds = config.VpcConfig.SubnetIDS
			securityGroups = config.VpcConfig.SecurityGroupIDS
		}
	case "AWS::EC2::RouteTable":
		var config model.RouteTableConfiguration
		err := json.Unmarshal([]byte(configuration), &config)
		if err == nil {
			vpcid = config.VpcID
		}
	}
	return
}

func filterResource(datas []model.AwsConfigFileData) []model.AwsConfigFileData {
	resuls := make([]model.AwsConfigFileData, 0)
	for _, d := range datas {
		if d.ResourceType == "AWS::EC2::VPCEndpoint" ||
			d.ResourceType == "AWS::EC2::VPC" ||
			d.ResourceType == "AWS::ServiceDiscovery::Service" ||
			d.ResourceType == "AWS::Signer::SigningProfile" ||
			d.ResourceType == "AWS::EC2::Subnet" ||
			d.ResourceType == "AWS::AmazonMQ::Broker" ||
			d.ResourceType == "AWS::CloudTrail::Trail" ||
			d.ResourceType == "AWS::EC2::NatGateway" ||
			d.ResourceType == "AWS::EC2::InternetGateway" ||
			d.ResourceType == "AWS::EC2::VPCPeeringConnection" ||
			d.ResourceType == "AWS::EFS::FileSystem" ||
			d.ResourceType == "AWS::IAM::Role" ||
			d.ResourceType == "AWS::RDS::DBInstance" ||
			d.ResourceType == "AWS::SNS::Topic" ||
			d.ResourceType == "AWS::ECS::Cluster" ||
			d.ResourceType == "AWS::IAM::Group" ||
			// d.ResourceType == "AWS::ElasticLoadBalancingV2::Listener" ||
			d.ResourceType == "AWS::IAM::User" ||
			d.ResourceType == "AWS::EC2::SecurityGroup" ||
			d.ResourceType == "AWS::EFS::AccessPoint" ||
			// d.ResourceType == "AWS::IoT::ProvisioningTemplate" ||
			d.ResourceType == "AWS::EC2::NetworkInterface" ||
			// d.ResourceType == "AWS::Route53Resolver::ResolverRuleAssociation" ||
			// d.ResourceType == "AWS::RDS::DBSubnetGroup" ||
			d.ResourceType == "AWS::EC2::EIP" ||
			d.ResourceType == "AWS::Redshift::ClusterSubnetGroup" ||
			d.ResourceType == "AWS::ElasticLoadBalancingV2::LoadBalancer" ||
			d.ResourceType == "AWS::ECS::Service" ||
			d.ResourceType == "AWS::EC2::NetworkAcl" ||
			d.ResourceType == "AWS::Lambda::Function" ||
			d.ResourceType == "AWS::S3::Bucket" ||
			d.ResourceType == "AWS::DynamoDB::Table" ||
			d.ResourceType == "AWS::EC2::RouteTable" {
			resuls = append(resuls, d)
		}
	}
	return resuls
}
