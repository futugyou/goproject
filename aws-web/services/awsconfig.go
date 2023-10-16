package services

import (
	"context"
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

// path for aws config snapshot data (download from s3)
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
	ships := make([]entity.AwsConfigRelationshipEntity, 0)

	for _, data := range datas {
		config := data.CreateAwsConfigEntity()
		if len(config.Label) > 0 {
			configs = append(configs, config)
		}		
	}

	for _, data := range datas {
		ship := data.CreateAwsConfigRelationshipEntity(configs)
		if len(ship) > 0 {
			ships = append(ships, ship...)
		}
	}

	log.Println("configs count: ", len(configs))
	err = a.repository.BulkWrite(context.Background(), configs)
	log.Println("configs write finish: ", err)
	log.Println("relationships count: ", len(ships))
	err = a.relRepository.BulkWrite(context.Background(), ships)
	log.Println("relationships write finish: ", err)
}

func (a *AwsConfigService) GetResourceGraph() model.ResourceGraph {
	configs, _ := a.repository.GetAll(context.Background())
	ships, _ := a.relRepository.GetAll(context.Background())
	nodes := make([]model.Node, 0)
	edges := make([]model.Edge, 0)

	for _, config := range configs {
		node := model.Node{
			ID:    config.ID,
			Label: config.Label,
			Properties: model.Properties{
				AccountID:        config.AccountID,
				Arn:              config.Arn,
				AvailabilityZone: config.AvailabilityZone,
				AwsRegion:        config.AwsRegion,
				// Configuration:                config.Configuration,
				ConfigurationItemCaptureTime: config.ConfigurationItemCaptureTime,
				ConfigurationItemStatus:      config.ConfigurationItemStatus,
				ConfigurationStateID:         config.ConfigurationStateID,
				ResourceCreationTime:         config.ResourceCreationTime,
				ResourceID:                   config.ResourceID,
				ResourceName:                 config.ResourceName,
				ResourceType:                 config.ResourceType,
				Tags:                         config.Tags,
				Version:                      config.Version,
				VpcID:                        config.VpcID,
				SubnetID:                     config.SubnetID,
				SubnetIDS:                    config.SubnetIds,
				Title:                        config.Title,
				SecurityGroups:               config.SecurityGroups,
			},
		}
		nodes = append(nodes, node)
	}

	for _, ship := range ships {
		edge := model.Edge{
			ID:    ship.ID,
			Label: ship.Label,
			Source: model.EdgeItem{
				ID:    ship.SourceID,
				Label: ship.SourceLabel,
			},
			Target: model.EdgeItem{
				ID:    ship.TargetID,
				Label: ship.TargetLabel,
			},
		}
		edges = append(edges, edge)
	}

	return model.ResourceGraph{
		Nodes: nodes,
		Edges: edges,
	}
}

func filterResource(datas []model.AwsConfigFileData) []model.AwsConfigFileData {
	resuls := make([]model.AwsConfigFileData, 0)
	for _, d := range datas {
		if d.ResourceType == "AWS::EC2::VPCEndpoint" ||
			d.ResourceType == "AWS::EC2::VPC" ||
			d.ResourceType == "AWS::ServiceDiscovery::Service" ||
			// d.ResourceType == "AWS::Signer::SigningProfile" ||
			d.ResourceType == "AWS::EC2::Subnet" ||
			d.ResourceType == "AWS::AmazonMQ::Broker" ||
			// d.ResourceType == "AWS::CloudTrail::Trail" ||
			d.ResourceType == "AWS::EC2::NatGateway" ||
			d.ResourceType == "AWS::EC2::InternetGateway" ||
			d.ResourceType == "AWS::EC2::VPCPeeringConnection" ||
			d.ResourceType == "AWS::EFS::FileSystem" ||
			// d.ResourceType == "AWS::IAM::Role" ||
			d.ResourceType == "AWS::RDS::DBInstance" ||
			d.ResourceType == "AWS::SNS::Topic" ||
			d.ResourceType == "AWS::ECS::Cluster" ||
			d.ResourceType == "AWS::IAM::Group" ||
			// d.ResourceType == "AWS::ElasticLoadBalancingV2::Listener" ||
			d.ResourceType == "AWS::IAM::User" ||
			d.ResourceType == "AWS::EC2::SecurityGroup" ||
			d.ResourceType == "AWS::EFS::AccessPoint" ||
			// d.ResourceType == "AWS::IoT::ProvisioningTemplate" ||
			// d.ResourceType == "AWS::EC2::NetworkInterface" ||
			// d.ResourceType == "AWS::Route53Resolver::ResolverRuleAssociation" ||
			// d.ResourceType == "AWS::RDS::DBSubnetGroup" ||
			// d.ResourceType == "AWS::EC2::EIP" ||
			// d.ResourceType == "AWS::Redshift::ClusterSubnetGroup" ||
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
