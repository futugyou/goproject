package services

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/futugyousuzu/goproject/awsgolang/entity"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"
	c "github.com/futugyousuzu/goproject/awsgolang/viewmodel/awsconfigConfiguration"

	"golang.org/x/exp/slices"
)

func RemoveDuplicateRelationShip(ships []entity.AwsConfigRelationshipEntity) []entity.AwsConfigRelationshipEntity {
	result := make([]entity.AwsConfigRelationshipEntity, 0)
	for _, ship := range ships {
		if !slices.ContainsFunc(result, func(s entity.AwsConfigRelationshipEntity) bool {
			return (s.SourceID == ship.SourceID && s.TargetID == ship.TargetID) ||
				(s.SourceID == ship.TargetID && s.TargetID == ship.SourceID)
		}) {
			result = append(result, ship)
		}
	}
	return result
}

func CreateAwsConfigRelationshipEntity(data model.AwsConfigRawData, configs []entity.AwsConfigEntity) []entity.AwsConfigRelationshipEntity {
	lists := make([]entity.AwsConfigRelationshipEntity, 0)

	for _, ship := range data.Relationships {
		index := -1
		if len(ship.ResourceID) > 0 {
			index = slices.IndexFunc(configs, func(sd entity.AwsConfigEntity) bool {
				return ship.ResourceID == sd.ResourceID && sd.ResourceType == ship.ResourceType
			})
		} else {
			index = slices.IndexFunc(configs, func(sd entity.AwsConfigEntity) bool {
				return ship.ResourceName == sd.ResourceName && sd.ResourceType == ship.ResourceType
			})
		}

		if index == -1 {
			continue
		}

		target := configs[index]
		relationship := entity.AwsConfigRelationshipEntity{
			ID:                 data.ResourceID + "-" + target.ResourceID,
			SourceID:           getId(data.ARN, data.ResourceID),
			SourceLabel:        data.ResourceName,
			SourceResourceType: data.ResourceType,
			Label:              ship.Name,
			TargetID:           target.ID,
			TargetLabel:        target.ResourceName,
			TargetResourceType: target.ResourceType,
		}

		lists = append(lists, relationship)
	}

	return lists
}

func AddIndividualRelationShip(configs []entity.AwsConfigEntity) []entity.AwsConfigRelationshipEntity {
	ships := make([]entity.AwsConfigRelationshipEntity, 0)

	sgs := make([]entity.AwsConfigEntity, 0)
	sds := make([]entity.AwsConfigEntity, 0)
	eventsRules := make([]entity.AwsConfigEntity, 0)
	files := make([]entity.AwsConfigEntity, 0)
	kms := make([]entity.AwsConfigEntity, 0)
	lbs := make([]entity.AwsConfigEntity, 0)
	ecsClusters := make([]entity.AwsConfigEntity, 0)
	roles := make([]entity.AwsConfigEntity, 0)
	targetGroups := make([]entity.AwsConfigEntity, 0)
	ecsTaskDefinitions := make([]entity.AwsConfigEntity, 0)
	networkInterfaces := make([]entity.AwsConfigEntity, 0)
	accessPoints := make([]entity.AwsConfigEntity, 0)
	managedPolicies := make([]entity.AwsConfigEntity, 0)

	for _, config := range configs {
		switch config.ResourceType {
		case "AWS::EC2::SecurityGroup":
			sgs = append(sgs, config)
		case "AWS::ServiceDiscovery::Service":
			sds = append(sds, config)
		case "AWS::Events::Rule":
			eventsRules = append(eventsRules, config)
		case "AWS::EFS::FileSystem":
			files = append(files, config)
		case "AWS::KMS::Key":
			kms = append(kms, config)
		case "AWS::ElasticLoadBalancingV2::LoadBalancer":
			lbs = append(lbs, config)
		case "AWS::ECS::Cluster":
			ecsClusters = append(ecsClusters, config)
		case "AWS::IAM::Role":
			roles = append(roles, config)
		case "AWS::ElasticLoadBalancingV2::TargetGroup":
			targetGroups = append(targetGroups, config)
		case "AWS::ECS::TaskDefinition":
			ecsTaskDefinitions = append(ecsTaskDefinitions, config)
		case "AWS::EC2::NetworkInterface":
			networkInterfaces = append(networkInterfaces, config)
		case "AWS::EFS::AccessPoint":
			accessPoints = append(accessPoints, config)
		case "AWS::IAM::AWSManagedPolicy":
			managedPolicies = append(managedPolicies, config)
		}
	}

	for _, config := range configs {
		ss := make([]entity.AwsConfigRelationshipEntity, 0)
		switch config.ResourceType {
		case "AWS::ECS::Service":
			ss = CreateECSServiceIndividualRelationShip(config, sds,
				sgs, roles, ecsClusters, ecsTaskDefinitions, targetGroups)
		case "AWS::EC2::SecurityGroup":
			ss = CreateSecurityGroupIndividualRelationShip(config, sgs)
		case "AWS::EFS::AccessPoint":
			ss = CreateAccessPointIndividualRelationShip(config, files)
		case "AWS::EFS::FileSystem":
			ss = CreateFileSystemIndividualRelationShip(config, kms)
		case "AWS::ElasticLoadBalancingV2::Listener":
			ss = CreateLoadBalancingListenerIndividualRelationShip(config, lbs, targetGroups)
		case "AWS::Events::EventBus":
			ss = CreateEventBusIndividualRelationShip(config, eventsRules)
		case "AWS::Events::Rule":
			ss = CreateEventRuleIndividualRelationShip(config, configs)
		case "AWS::Lambda::Function":
			ss = CreateFunctionIndividualRelationShip(config, configs)
		case "AWS::EC2::NetworkInterface":
			ss = CreateNetworkInterfaceIndividualRelationShip(config, configs)
		case "AWS::EC2::RouteTable":
			ss = CreateRouteTableIndividualRelationShip(config, configs)
		case "AWS::AmazonMQ::Broker":
			ss = CreateMQBrokerIndividualRelationShip(config, sgs)
		case "AWS::RDS::DBInstance":
			ss = CreateDBInstanceIndividualRelationShip(config, kms)
		case "AWS::ECS::Task":
			ss = CreateEcsTaskRelationShip(config, ecsTaskDefinitions, networkInterfaces)
		case "AWS::ECS::TaskDefinition":
			ss = CreateEcsTaskDefinitionRelationShip(config, roles, files, accessPoints)
		case "AWS::ElasticLoadBalancingV2::TargetGroup":
			ss = CreateTargetGroupIndividualRelationShip(config, lbs)
		case "AWS::IAM::Role":
			ss = CreateRoleIndividualRelationShip(config, managedPolicies)
		case "AWS::IAM::User":
			ss = CreateUserIndividualRelationShip(config, managedPolicies)
		}
		ships = append(ships, ss...)
	}

	return ships
}

func CreateUserIndividualRelationShip(config entity.AwsConfigEntity, managedPolicies []entity.AwsConfigEntity) (ships []entity.AwsConfigRelationshipEntity) {
	ships = make([]entity.AwsConfigRelationshipEntity, 0)
	var conf c.UserConfiguration
	err := json.Unmarshal([]byte(config.Configuration), &conf)
	if err != nil {
		log.Println(err)
		return
	}
	for _, policy := range conf.AttachedManagedPolicies {
		index := slices.IndexFunc(managedPolicies, func(sd entity.AwsConfigEntity) bool {
			return policy.PolicyArn == sd.Arn && len(sd.Arn) > 0
		})
		if index != -1 {
			target := managedPolicies[index]
			ship := entity.AwsConfigRelationshipEntity{
				ID:                 config.ResourceID + "-" + target.ResourceID,
				SourceID:           config.ID,
				SourceLabel:        config.ResourceName,
				SourceResourceType: config.ResourceType,
				Label:              fmt.Sprintf("Is attached with %s", GetResourceTypeName(target.ResourceType)),
				TargetID:           target.ID,
				TargetLabel:        target.ResourceName,
				TargetResourceType: target.ResourceType,
			}
			ships = append(ships, ship)
		}
	}
	return ships
}

func CreateRoleIndividualRelationShip(config entity.AwsConfigEntity, managedPolicies []entity.AwsConfigEntity) (ships []entity.AwsConfigRelationshipEntity) {
	ships = make([]entity.AwsConfigRelationshipEntity, 0)
	var conf c.RoleConfiguration
	err := json.Unmarshal([]byte(config.Configuration), &conf)
	if err != nil {
		log.Println(err)
		return
	}
	for _, policy := range conf.AttachedManagedPolicies {
		index := slices.IndexFunc(managedPolicies, func(sd entity.AwsConfigEntity) bool {
			return policy.PolicyArn == sd.Arn && len(sd.Arn) > 0
		})
		if index != -1 {
			target := managedPolicies[index]
			ship := entity.AwsConfigRelationshipEntity{
				ID:                 config.ResourceID + "-" + target.ResourceID,
				SourceID:           config.ID,
				SourceLabel:        config.ResourceName,
				SourceResourceType: config.ResourceType,
				Label:              fmt.Sprintf("Is attached with %s", GetResourceTypeName(target.ResourceType)),
				TargetID:           target.ID,
				TargetLabel:        target.ResourceName,
				TargetResourceType: target.ResourceType,
			}
			ships = append(ships, ship)
		}
	}
	return ships
}

func CreateTargetGroupIndividualRelationShip(config entity.AwsConfigEntity, lbs []entity.AwsConfigEntity) (ships []entity.AwsConfigRelationshipEntity) {
	ships = make([]entity.AwsConfigRelationshipEntity, 0)
	var conf c.TargetGroupConfiguration
	err := json.Unmarshal([]byte(config.Configuration), &conf)
	if err != nil {
		log.Println(err)
		return
	}

	for _, lb := range conf.LoadBalancerArns {
		index := slices.IndexFunc(lbs, func(sd entity.AwsConfigEntity) bool {
			return lb == sd.Arn && len(lb) > 0
		})
		if index != -1 {
			target := lbs[index]
			ship := entity.AwsConfigRelationshipEntity{
				ID:                 config.ResourceID + "-" + target.ResourceID,
				SourceID:           config.ID,
				SourceLabel:        config.ResourceName,
				SourceResourceType: config.ResourceType,
				Label:              fmt.Sprintf("Is attached with %s", GetResourceTypeName(target.ResourceType)),
				TargetID:           target.ID,
				TargetLabel:        target.ResourceName,
				TargetResourceType: target.ResourceType,
			}
			ships = append(ships, ship)
		}
	}

	return ships
}

func CreateEcsTaskDefinitionRelationShip(config entity.AwsConfigEntity,
	roles []entity.AwsConfigEntity,
	files []entity.AwsConfigEntity,
	accessPoints []entity.AwsConfigEntity,
) (ships []entity.AwsConfigRelationshipEntity) {
	ships = make([]entity.AwsConfigRelationshipEntity, 0)
	var conf c.ECSTaskDefinitionConfiguration
	err := json.Unmarshal([]byte(config.Configuration), &conf)
	if err != nil {
		log.Println(err)
		return
	}
	// 1. ExecutionRoleArn
	index := slices.IndexFunc(roles, func(sd entity.AwsConfigEntity) bool {
		return conf.ExecutionRoleArn == sd.Arn && len(sd.Arn) > 0
	})
	if index != -1 {
		target := roles[index]
		ship := entity.AwsConfigRelationshipEntity{
			ID:                 config.ResourceID + "-" + target.ResourceID,
			SourceID:           config.ID,
			SourceLabel:        config.ResourceName,
			SourceResourceType: config.ResourceType,
			Label:              fmt.Sprintf("Is associated with %s", GetResourceTypeName(target.ResourceType)),
			TargetID:           target.ID,
			TargetLabel:        target.ResourceName,
			TargetResourceType: target.ResourceType,
		}
		ships = append(ships, ship)
	}

	// 2. TaskRoleArn
	index = slices.IndexFunc(roles, func(sd entity.AwsConfigEntity) bool {
		return conf.TaskRoleArn == sd.Arn && len(sd.Arn) > 0
	})
	if index != -1 {
		target := roles[index]
		ship := entity.AwsConfigRelationshipEntity{
			ID:                 config.ResourceID + "-" + target.ResourceID,
			SourceID:           config.ID,
			SourceLabel:        config.ResourceName,
			SourceResourceType: config.ResourceType,
			Label:              fmt.Sprintf("Is associated with %s", GetResourceTypeName(target.ResourceType)),
			TargetID:           target.ID,
			TargetLabel:        target.ResourceName,
			TargetResourceType: target.ResourceType,
		}
		ships = append(ships, ship)
	}

	// 3. volumes
	for _, volume := range conf.Volumes {
		if len(volume.EFSVolumeConfiguration.AuthorizationConfig.AccessPointID) > 0 {
			index = slices.IndexFunc(accessPoints, func(sd entity.AwsConfigEntity) bool {
				return volume.EFSVolumeConfiguration.AuthorizationConfig.AccessPointID == sd.ResourceID
			})
			if index != -1 {
				target := accessPoints[index]
				ship := entity.AwsConfigRelationshipEntity{
					ID:                 config.ResourceID + "-" + target.ResourceID,
					SourceID:           config.ID,
					SourceLabel:        config.ResourceName,
					SourceResourceType: config.ResourceType,
					Label:              fmt.Sprintf("Is associated with %s", GetResourceTypeName(target.ResourceType)),
					TargetID:           target.ID,
					TargetLabel:        target.ResourceName,
					TargetResourceType: target.ResourceType,
				}
				ships = append(ships, ship)
			}
		} else if len(volume.EFSVolumeConfiguration.FileSystemID) > 0 {
			index = slices.IndexFunc(files, func(sd entity.AwsConfigEntity) bool {
				return volume.EFSVolumeConfiguration.FileSystemID == sd.ResourceID
			})
			if index != -1 {
				target := files[index]
				ship := entity.AwsConfigRelationshipEntity{
					ID:                 config.ResourceID + "-" + target.ResourceID,
					SourceID:           config.ID,
					SourceLabel:        config.ResourceName,
					SourceResourceType: config.ResourceType,
					Label:              fmt.Sprintf("Is associated with %s", GetResourceTypeName(target.ResourceType)),
					TargetID:           target.ID,
					TargetLabel:        target.ResourceName,
					TargetResourceType: target.ResourceType,
				}
				ships = append(ships, ship)
			}
		}
	}
	return ships
}

func CreateEcsTaskRelationShip(config entity.AwsConfigEntity,
	ecsTaskDefinitions []entity.AwsConfigEntity,
	networkInterfaces []entity.AwsConfigEntity,
) (ships []entity.AwsConfigRelationshipEntity) {
	ships = make([]entity.AwsConfigRelationshipEntity, 0)
	var conf c.ECSTaskConfiguration
	err := json.Unmarshal([]byte(config.Configuration), &conf)
	if err != nil {
		log.Println(err)
		return
	}
	index := slices.IndexFunc(ecsTaskDefinitions, func(sd entity.AwsConfigEntity) bool {
		return conf.TaskDefinitionArn == sd.Arn && len(sd.Arn) > 0
	})
	if index != -1 {
		target := ecsTaskDefinitions[index]
		ship := entity.AwsConfigRelationshipEntity{
			ID:                 config.ResourceID + "-" + target.ResourceID,
			SourceID:           config.ID,
			SourceLabel:        config.ResourceName,
			SourceResourceType: config.ResourceType,
			Label:              fmt.Sprintf("Is attached with %s", GetResourceTypeName(target.ResourceType)),
			TargetID:           target.ID,
			TargetLabel:        target.ResourceName,
			TargetResourceType: target.ResourceType,
		}
		ships = append(ships, ship)
	}

	for _, att := range conf.Attachments {
		if att.Type == "ElasticNetworkInterface" {
			for _, d := range att.Details {
				if d.Name == "networkInterfaceId" && len(d.Value) > 0 {
					index = slices.IndexFunc(networkInterfaces, func(sd entity.AwsConfigEntity) bool {
						return d.Value == sd.ResourceID
					})
					if index != -1 {
						target := networkInterfaces[index]
						ship := entity.AwsConfigRelationshipEntity{
							ID:                 config.ResourceID + "-" + target.ResourceID,
							SourceID:           config.ID,
							SourceLabel:        config.ResourceName,
							SourceResourceType: config.ResourceType,
							Label:              fmt.Sprintf("Is associated with %s", GetResourceTypeName(target.ResourceType)),
							TargetID:           target.ID,
							TargetLabel:        target.ResourceName,
							TargetResourceType: target.ResourceType,
						}
						ships = append(ships, ship)
					}
				}
			}
		}
	}
	return ships
}

func CreateDBInstanceIndividualRelationShip(config entity.AwsConfigEntity, configs []entity.AwsConfigEntity) (ships []entity.AwsConfigRelationshipEntity) {
	ships = make([]entity.AwsConfigRelationshipEntity, 0)
	var conf c.DBInstanceConfiguration
	err := json.Unmarshal([]byte(config.Configuration), &conf)
	if err != nil {
		log.Println(err)
		return
	}
	index := slices.IndexFunc(configs, func(sd entity.AwsConfigEntity) bool {
		return conf.KmsKeyID == sd.Arn && sd.ResourceType == "AWS::KMS::Key"
	})
	if index != -1 {
		target := configs[index]
		ship := entity.AwsConfigRelationshipEntity{
			ID:                 config.ResourceID + "-" + target.ResourceID,
			SourceID:           config.ID,
			SourceLabel:        config.ResourceName,
			SourceResourceType: config.ResourceType,
			Label:              "Is attached to KMS",
			TargetID:           target.ID,
			TargetLabel:        target.ResourceName,
			TargetResourceType: target.ResourceType,
		}
		ships = append(ships, ship)
	}
	return
}

func CreateMQBrokerIndividualRelationShip(config entity.AwsConfigEntity, configs []entity.AwsConfigEntity) (ships []entity.AwsConfigRelationshipEntity) {
	ships = make([]entity.AwsConfigRelationshipEntity, 0)
	var conf c.AmazonMQConfiguration
	err := json.Unmarshal([]byte(config.Configuration), &conf)
	if err != nil {
		log.Println(err)
		return
	}
	for _, target := range conf.SecurityGroups {
		index := slices.IndexFunc(configs, func(sd entity.AwsConfigEntity) bool {
			return target == sd.ResourceID
		})
		if index != -1 {
			target := configs[index]
			ship := entity.AwsConfigRelationshipEntity{
				ID:                 config.ResourceID + "-" + target.ResourceID,
				SourceID:           config.ID,
				SourceLabel:        config.ResourceName,
				SourceResourceType: config.ResourceType,
				Label:              fmt.Sprintf("Is associated with %s", GetResourceTypeName(target.ResourceType)),
				TargetID:           target.ID,
				TargetLabel:        target.ResourceName,
				TargetResourceType: target.ResourceType,
			}
			ships = append(ships, ship)
		}
	}

	return
}
func CreateRouteTableIndividualRelationShip(config entity.AwsConfigEntity, configs []entity.AwsConfigEntity) (ships []entity.AwsConfigRelationshipEntity) {
	ships = make([]entity.AwsConfigRelationshipEntity, 0)
	var conf c.RouteTableConfiguration
	err := json.Unmarshal([]byte(config.Configuration), &conf)
	if err != nil {
		log.Println(err)
		return
	}
	for _, route := range conf.Routes {
		index := slices.IndexFunc(configs, func(sd entity.AwsConfigEntity) bool {
			if len(route.NatGatewayId) > 0 {
				return sd.ResourceID == route.NatGatewayId && sd.ResourceType == "AWS::EC2::NatGateway"
			} else if strings.HasPrefix(route.GatewayID, "igw-") {
				return sd.ResourceID == route.GatewayID && sd.ResourceType == "AWS::EC2::InternetGateway"
			} else if strings.HasPrefix(route.GatewayID, "vpce-") {
				return sd.ResourceID == route.GatewayID && sd.ResourceType == "AWS::EC2::VPCEndpoint"
			}

			return false
		})
		if index != -1 {
			target := configs[index]
			ship := entity.AwsConfigRelationshipEntity{
				ID:                 config.ResourceID + "-" + target.ResourceID,
				SourceID:           config.ID,
				SourceLabel:        config.ResourceName,
				SourceResourceType: config.ResourceType,
				Label:              fmt.Sprintf("Contains %s", GetResourceTypeName(target.ResourceType)),
				TargetID:           target.ID,
				TargetLabel:        target.ResourceName,
				TargetResourceType: target.ResourceType,
			}
			ships = append(ships, ship)
		}
	}

	return
}

func CreateNetworkInterfaceIndividualRelationShip(config entity.AwsConfigEntity, configs []entity.AwsConfigEntity) (ships []entity.AwsConfigRelationshipEntity) {
	ships = make([]entity.AwsConfigRelationshipEntity, 0)
	var conf c.NetworkInterfaceConfiguration
	err := json.Unmarshal([]byte(config.Configuration), &conf)
	if err != nil {
		log.Println(err)
		return
	}

	index := slices.IndexFunc(configs, func(sd entity.AwsConfigEntity) bool {
		if strings.HasPrefix(conf.Description, "Interface for NAT Gateway ") {
			rid := strings.ReplaceAll(conf.Description, "Interface for NAT Gateway ", "")
			return rid == sd.ResourceID && sd.ResourceType == "AWS::EC2::NatGateway"
		} else if strings.HasPrefix(conf.Description, "ELB app") {
			arn := fmt.Sprintf("arn:aws:elasticloadbalancing:%s:%s:loadbalancer/%s", config.AwsRegion, config.AccountID, strings.Replace(conf.Description, "ELB ", "", 1))
			return arn == sd.Arn && sd.ResourceType == "AWS::ElasticLoadBalancingV2::LoadBalancer"
		} else if conf.InterfaceType == "vpc_endpoint" {
			rid := strings.ReplaceAll(conf.Description, "VPC Endpoint Interface ", "")
			return rid == sd.ResourceID && sd.ResourceType == "AWS::EC2::VPCEndpoint"
		} else if conf.RequesterID == "amazon-elasticsearch" {
			rid := fmt.Sprintf("%s/%s", config.AccountID, strings.Replace(conf.Description, "ES ", "", 1))
			return rid == sd.ResourceID && sd.ResourceType == "AWS::Elasticsearch::Domain"
		} else if conf.RequesterID == "lambda" {
			rid := strings.Replace(conf.Description, "AWS Lambda VPC ENI-", "", 1)
			if len(rid) > 37 {
				rid = rid[0 : len(rid)-37]
			}
			return rid == sd.ResourceID && sd.ResourceType == "AWS::Lambda::Function"
		} else if conf.RequesterID == "efs" {
			rid := strings.Split(strings.Replace(conf.Description, "EFS mount target for ", "", 1), " (fsmt-")
			if len(rid) == 0 {
				return false
			}
			return rid[0] == sd.ResourceID && sd.ResourceType == "AWS::EFS::FileSystem"
		}

		return false
	})

	if index != -1 {
		target := configs[index]
		ship := entity.AwsConfigRelationshipEntity{
			ID:                 config.ResourceID + "-" + target.ResourceID,
			SourceID:           config.ID,
			SourceLabel:        config.ResourceName,
			SourceResourceType: config.ResourceType,
			Label:              fmt.Sprintf("Is attached to %s", GetResourceTypeName(target.ResourceType)),
			TargetID:           target.ID,
			TargetLabel:        target.ResourceName,
			TargetResourceType: target.ResourceType,
		}
		ships = append(ships, ship)
	}
	return
}

func CreateFunctionIndividualRelationShip(config entity.AwsConfigEntity, configs []entity.AwsConfigEntity) (ships []entity.AwsConfigRelationshipEntity) {
	ships = make([]entity.AwsConfigRelationshipEntity, 0)
	var conf c.FunctionConfiguration
	err := json.Unmarshal([]byte(config.Configuration), &conf)
	if err != nil {
		log.Println(err)
		return
	}
	for _, target := range conf.FileSystemConfigs {
		index := slices.IndexFunc(configs, func(sd entity.AwsConfigEntity) bool {
			return target.Arn == sd.Arn && len(sd.Arn) > 0
		})
		if index != -1 {
			target := configs[index]
			ship := entity.AwsConfigRelationshipEntity{
				ID:                 config.ResourceID + "-" + target.ResourceID,
				SourceID:           config.ID,
				SourceLabel:        config.ResourceName,
				SourceResourceType: config.ResourceType,
				Label:              fmt.Sprintf("Is associated with %s", GetResourceTypeName(target.ResourceType)),
				TargetID:           target.ID,
				TargetLabel:        target.ResourceName,
				TargetResourceType: target.ResourceType,
			}
			ships = append(ships, ship)
		}
	}

	return
}

func CreateEventRuleIndividualRelationShip(config entity.AwsConfigEntity, configs []entity.AwsConfigEntity) (ships []entity.AwsConfigRelationshipEntity) {
	ships = make([]entity.AwsConfigRelationshipEntity, 0)
	var conf c.EventRoleConfiguration
	err := json.Unmarshal([]byte(config.Configuration), &conf)
	if err != nil {
		log.Println(err)
		return
	}
	for _, target := range conf.Targets {
		// 1. target.Arn
		index := slices.IndexFunc(configs, func(sd entity.AwsConfigEntity) bool {
			return target.Arn == sd.Arn && len(sd.Arn) > 0
		})
		if index != -1 {
			target := configs[index]
			ship := entity.AwsConfigRelationshipEntity{
				ID:                 config.ResourceID + "-" + target.ResourceID,
				SourceID:           config.ID,
				SourceLabel:        config.ResourceName,
				SourceResourceType: config.ResourceType,
				Label:              fmt.Sprintf("Is associated with %s", GetResourceTypeName(target.ResourceType)),
				TargetID:           target.ID,
				TargetLabel:        target.ResourceName,
				TargetResourceType: target.ResourceType,
			}
			ships = append(ships, ship)
		}
		// 2. target.RoleArn
		index = slices.IndexFunc(configs, func(sd entity.AwsConfigEntity) bool {
			return target.RoleArn == sd.Arn && len(sd.Arn) > 0
		})
		if index != -1 {
			target := configs[index]
			ship := entity.AwsConfigRelationshipEntity{
				ID:                 config.ResourceID + "-" + target.ResourceID,
				SourceID:           config.ID,
				SourceLabel:        config.ResourceName,
				SourceResourceType: config.ResourceType,
				Label:              fmt.Sprintf("Is associated with %s", GetResourceTypeName(target.ResourceType)),
				TargetID:           target.ID,
				TargetLabel:        target.ResourceName,
				TargetResourceType: target.ResourceType,
			}
			ships = append(ships, ship)
		}
		// 3. target.EcsParameters.TaskDefinitionArn
		index = slices.IndexFunc(configs, func(sd entity.AwsConfigEntity) bool {
			return target.EcsParameters.TaskDefinitionArn == sd.Arn && len(sd.Arn) > 0
		})
		if index != -1 {
			target := configs[index]
			ship := entity.AwsConfigRelationshipEntity{
				ID:                 config.ResourceID + "-" + target.ResourceID,
				SourceID:           config.ID,
				SourceLabel:        config.ResourceName,
				SourceResourceType: config.ResourceType,
				Label:              fmt.Sprintf("Is associated with %s", GetResourceTypeName(target.ResourceType)),
				TargetID:           target.ID,
				TargetLabel:        target.ResourceName,
				TargetResourceType: target.ResourceType,
			}
			ships = append(ships, ship)
		}
	}

	return
}

func GetResourceTypeName(name string) string {
	typeName := ""
	ls := strings.Split(name, "::")

	if len(ls) > 0 {
		typeName = ls[len(ls)-1]
	}
	return typeName
}

func CreateEventBusIndividualRelationShip(config entity.AwsConfigEntity, configs []entity.AwsConfigEntity) (ships []entity.AwsConfigRelationshipEntity) {
	ships = make([]entity.AwsConfigRelationshipEntity, 0)
	index := slices.IndexFunc(configs, func(sd entity.AwsConfigEntity) bool {
		var conf c.EventRoleConfiguration
		err := json.Unmarshal([]byte(sd.Configuration), &conf)
		if err != nil {
			log.Println(err)
			return false
		}

		eventBusName := conf.EventBusName
		if !strings.HasPrefix(eventBusName, "arn:") {
			eventBusName = fmt.Sprintf("arn:aws:events:%s:%s:event-bus/%s", sd.AwsRegion, sd.AccountID, eventBusName)
		}

		return config.Arn == eventBusName && len(eventBusName) > 0
	})

	if index != -1 {
		target := configs[index]
		ship := entity.AwsConfigRelationshipEntity{
			ID:                 config.ResourceID + "-" + target.ResourceID,
			SourceID:           config.ID,
			SourceLabel:        config.ResourceName,
			SourceResourceType: config.ResourceType,
			Label:              "Is associated with EventsRule",
			TargetID:           target.ID,
			TargetLabel:        target.ResourceName,
			TargetResourceType: target.ResourceType,
		}
		ships = append(ships, ship)
	}
	return
}

func CreateLoadBalancingListenerIndividualRelationShip(config entity.AwsConfigEntity, lbs, targetGroups []entity.AwsConfigEntity) (ships []entity.AwsConfigRelationshipEntity) {
	ships = make([]entity.AwsConfigRelationshipEntity, 0)
	var conf c.LoadBalancerListenerConfiguration
	err := json.Unmarshal([]byte(config.Configuration), &conf)
	if err != nil {
		log.Println(err)
		return
	}

	// 1. conf.LoadBalancerArn
	index := slices.IndexFunc(lbs, func(sd entity.AwsConfigEntity) bool {
		return conf.LoadBalancerArn == sd.Arn && len(sd.Arn) > 0
	})
	if index != -1 {
		target := lbs[index]
		ship := entity.AwsConfigRelationshipEntity{
			ID:                 config.ResourceID + "-" + target.ResourceID,
			SourceID:           config.ID,
			SourceLabel:        config.ResourceName,
			SourceResourceType: config.ResourceType,
			Label:              fmt.Sprintf("Is associated with %s", GetResourceTypeName(target.ResourceType)),
			TargetID:           target.ID,
			TargetLabel:        target.ResourceName,
			TargetResourceType: target.ResourceType,
		}
		ships = append(ships, ship)
	}

	// 2. conf.DefaultActions[].TargetGroupArn
	for _, act := range conf.DefaultActions {
		index = slices.IndexFunc(targetGroups, func(sd entity.AwsConfigEntity) bool {
			return act.TargetGroupArn == sd.Arn && len(sd.Arn) > 0
		})
		if index != -1 {
			target := targetGroups[index]
			ship := entity.AwsConfigRelationshipEntity{
				ID:                 config.ResourceID + "-" + target.ResourceID,
				SourceID:           config.ID,
				SourceLabel:        config.ResourceName,
				SourceResourceType: config.ResourceType,
				Label:              fmt.Sprintf("Is associated with %s", GetResourceTypeName(target.ResourceType)),
				TargetID:           target.ID,
				TargetLabel:        target.ResourceName,
				TargetResourceType: target.ResourceType,
			}
			ships = append(ships, ship)
		}

		//3. conf.DefaultActions[].ForwardConfig.TargetGroups[].TargetGroupArn
		for _, tg := range act.ForwardConfig.TargetGroups {
			index = slices.IndexFunc(targetGroups, func(sd entity.AwsConfigEntity) bool {
				return tg.TargetGroupArn == sd.Arn && len(sd.Arn) > 0
			})
			if index != -1 {
				target := targetGroups[index]
				ship := entity.AwsConfigRelationshipEntity{
					ID:                 config.ResourceID + "-" + target.ResourceID,
					SourceID:           config.ID,
					SourceLabel:        config.ResourceName,
					SourceResourceType: config.ResourceType,
					Label:              fmt.Sprintf("Is associated with %s", GetResourceTypeName(target.ResourceType)),
					TargetID:           target.ID,
					TargetLabel:        target.ResourceName,
					TargetResourceType: target.ResourceType,
				}
				ships = append(ships, ship)
			}
		}
	}

	return
}

func CreateFileSystemIndividualRelationShip(config entity.AwsConfigEntity, configs []entity.AwsConfigEntity) (ships []entity.AwsConfigRelationshipEntity) {
	ships = make([]entity.AwsConfigRelationshipEntity, 0)
	var conf c.FileSystemConfiguration
	err := json.Unmarshal([]byte(config.Configuration), &conf)
	if err != nil {
		log.Println(err)
		return
	}
	index := slices.IndexFunc(configs, func(sd entity.AwsConfigEntity) bool {
		return conf.KmsKeyID == sd.Arn && sd.ResourceType == "AWS::KMS::Key"
	})
	if index != -1 {
		target := configs[index]
		ship := entity.AwsConfigRelationshipEntity{
			ID:                 config.ResourceID + "-" + target.ResourceID,
			SourceID:           config.ID,
			SourceLabel:        config.ResourceName,
			SourceResourceType: config.ResourceType,
			Label:              "Is attached to KMS",
			TargetID:           target.ID,
			TargetLabel:        target.ResourceName,
			TargetResourceType: target.ResourceType,
		}
		ships = append(ships, ship)
	}
	return
}

func CreateAccessPointIndividualRelationShip(config entity.AwsConfigEntity, configs []entity.AwsConfigEntity) (ships []entity.AwsConfigRelationshipEntity) {
	ships = make([]entity.AwsConfigRelationshipEntity, 0)
	var conf c.AccessPointConfiguration
	err := json.Unmarshal([]byte(config.Configuration), &conf)
	if err != nil {
		log.Println(err)
		return
	}
	index := slices.IndexFunc(configs, func(sd entity.AwsConfigEntity) bool {
		return conf.FileSystemID == sd.ResourceID && sd.ResourceType == "AWS::EFS::FileSystem"
	})
	if index != -1 {
		target := configs[index]
		ship := entity.AwsConfigRelationshipEntity{
			ID:                 config.ResourceID + "-" + target.ResourceID,
			SourceID:           config.ID,
			SourceLabel:        config.ResourceName,
			SourceResourceType: config.ResourceType,
			Label:              "Is attached to FileSystem",
			TargetID:           target.ID,
			TargetLabel:        target.ResourceName,
			TargetResourceType: target.ResourceType,
		}
		ships = append(ships, ship)
	}
	return
}

func CreateSecurityGroupIndividualRelationShip(config entity.AwsConfigEntity, sgs []entity.AwsConfigEntity) (ships []entity.AwsConfigRelationshipEntity) {
	ships = make([]entity.AwsConfigRelationshipEntity, 0)
	var sgconfig c.SecurityGroupConfiguration
	err := json.Unmarshal([]byte(config.Configuration), &sgconfig)
	if err != nil {
		log.Println(err)
		return
	}

	permissions := securityGroupIPPermissions(config, sgconfig.IPPermissions, sgs)
	ships = append(ships, permissions...)
	permissions = securityGroupIPPermissions(config, sgconfig.IPPermissionsEgress, sgs)
	ships = append(ships, permissions...)
	return
}

func CreateECSServiceIndividualRelationShip(
	config entity.AwsConfigEntity,
	sds []entity.AwsConfigEntity,
	sgs []entity.AwsConfigEntity,
	roles []entity.AwsConfigEntity,
	ecsClusters []entity.AwsConfigEntity,
	ecsTaskDefinitions []entity.AwsConfigEntity,
	targetGroups []entity.AwsConfigEntity,
) (ships []entity.AwsConfigRelationshipEntity) {
	ships = make([]entity.AwsConfigRelationshipEntity, 0)
	var ecsconfig c.ECSServiceConfiguration
	err := json.Unmarshal([]byte(config.Configuration), &ecsconfig)
	if err != nil {
		log.Println(err)
		return
	}

	// ServiceDiscovery Relationship
	for _, sr := range ecsconfig.ServiceRegistries {
		index := slices.IndexFunc(sds, func(sd entity.AwsConfigEntity) bool {
			return sr.RegistryArn == sd.ID && len(sr.RegistryArn) > 0
		})
		if index != -1 {
			target := sds[index]
			ship := entity.AwsConfigRelationshipEntity{
				ID:                 config.ResourceID + "-" + target.ResourceID,
				SourceID:           config.ID,
				SourceLabel:        config.ResourceName,
				SourceResourceType: config.ResourceType,
				Label:              "Is associated with ServiceDiscovery",
				TargetID:           target.ID,
				TargetLabel:        target.ResourceName,
				TargetResourceType: target.ResourceType,
			}
			ships = append(ships, ship)
		}
	}

	// SecurityGroup Relationship
	for _, sgg := range ecsconfig.NetworkConfiguration.AwsvpcConfiguration.SecurityGroups {
		index := slices.IndexFunc(sgs, func(sd entity.AwsConfigEntity) bool {
			return sgg == sd.ResourceID
		})
		if index != -1 {
			target := sgs[index]
			ship := entity.AwsConfigRelationshipEntity{
				ID:                 config.ResourceID + "-" + target.ResourceID,
				SourceID:           config.ID,
				SourceLabel:        config.ResourceName,
				SourceResourceType: config.ResourceType,
				Label:              "Is associated with SecurityGroup",
				TargetID:           target.ID,
				TargetLabel:        target.ResourceName,
				TargetResourceType: target.ResourceType,
			}
			ships = append(ships, ship)
		}
	}

	// role
	{
		index := slices.IndexFunc(roles, func(role entity.AwsConfigEntity) bool {
			return ecsconfig.Role == role.Arn && len(role.Arn) > 0
		})

		if index != -1 {
			target := roles[index]
			ship := entity.AwsConfigRelationshipEntity{
				ID:                 config.ResourceID + "-" + target.ResourceID,
				SourceID:           config.ID,
				SourceLabel:        config.ResourceName,
				SourceResourceType: config.ResourceType,
				Label:              "Is associated with Role",
				TargetID:           target.ID,
				TargetLabel:        target.ResourceName,
				TargetResourceType: target.ResourceType,
			}
			ships = append(ships, ship)
		}
	}

	// cluster
	{
		index := slices.IndexFunc(ecsClusters, func(ecsCluster entity.AwsConfigEntity) bool {
			return ecsconfig.Cluster == ecsCluster.Arn && len(ecsCluster.Arn) > 0
		})

		if index != -1 {
			target := ecsClusters[index]
			ship := entity.AwsConfigRelationshipEntity{
				ID:                 config.ResourceID + "-" + target.ResourceID,
				SourceID:           config.ID,
				SourceLabel:        config.ResourceName,
				SourceResourceType: config.ResourceType,
				Label:              "Is contained in Cluster",
				TargetID:           target.ID,
				TargetLabel:        target.ResourceName,
				TargetResourceType: target.ResourceType,
			}
			ships = append(ships, ship)
		}
	}

	// ecsTaskDefinitions
	{
		index := slices.IndexFunc(ecsTaskDefinitions, func(taskDefinition entity.AwsConfigEntity) bool {
			return ecsconfig.TaskDefinition == taskDefinition.Arn && len(taskDefinition.Arn) > 0
		})
		if index != -1 {
			target := ecsTaskDefinitions[index]
			ship := entity.AwsConfigRelationshipEntity{
				ID:                 config.ResourceID + "-" + target.ResourceID,
				SourceID:           config.ID,
				SourceLabel:        config.ResourceName,
				SourceResourceType: config.ResourceType,
				Label:              fmt.Sprintf("Is associated with %s", GetResourceTypeName(target.ResourceType)),
				TargetID:           target.ID,
				TargetLabel:        target.ResourceName,
				TargetResourceType: target.ResourceType,
			}
			ships = append(ships, ship)
		}
	}

	// targetGroups
	{
		for _, lb := range ecsconfig.LoadBalancers {
			index := slices.IndexFunc(targetGroups, func(tg entity.AwsConfigEntity) bool {
				return lb.TargetGroupArn == tg.Arn && len(tg.Arn) > 0
			})
			if index != -1 {
				target := targetGroups[index]
				ship := entity.AwsConfigRelationshipEntity{
					ID:                 config.ResourceID + "-" + target.ResourceID,
					SourceID:           config.ID,
					SourceLabel:        config.ResourceName,
					SourceResourceType: config.ResourceType,
					Label:              fmt.Sprintf("Is associated with %s", GetResourceTypeName(target.ResourceType)),
					TargetID:           target.ID,
					TargetLabel:        target.ResourceName,
					TargetResourceType: target.ResourceType,
				}
				ships = append(ships, ship)
			}
		}
	}
	return
}

func securityGroupIPPermissions(config entity.AwsConfigEntity, permissions []c.IPPermission, sgs []entity.AwsConfigEntity) []entity.AwsConfigRelationshipEntity {
	ships := make([]entity.AwsConfigRelationshipEntity, 0)
	for _, permission := range permissions {
		for _, pair := range permission.UserIDGroupPairs {
			index := slices.IndexFunc(sgs, func(sd entity.AwsConfigEntity) bool {
				return pair.GroupID == sd.ResourceID
			})
			if index != -1 {
				target := sgs[index]
				ship := entity.AwsConfigRelationshipEntity{
					ID:                 config.ResourceID + "-" + target.ResourceID,
					SourceID:           config.ID,
					SourceLabel:        config.ResourceName,
					SourceResourceType: config.ResourceType,
					Label:              "Is associated with SecurityGroup",
					TargetID:           target.ID,
					TargetLabel:        target.ResourceName,
					TargetResourceType: target.ResourceType,
				}
				ships = append(ships, ship)
			}
		}
	}
	return ships
}
