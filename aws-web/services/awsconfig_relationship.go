package services

import (
	"encoding/json"
	"strings"

	"github.com/futugyousuzu/goproject/awsgolang/entity"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"
	c "github.com/futugyousuzu/goproject/awsgolang/viewmodel/awsconfigConfiguration"

	"golang.org/x/exp/slices"
)

func CreateAwsConfigRelationshipEntity(data model.AwsConfigFileData, configs []entity.AwsConfigEntity) []entity.AwsConfigRelationshipEntity {
	lists := make([]entity.AwsConfigRelationshipEntity, 0)

	for _, ship := range data.Relationships {
		var id string
		for i := 0; i < len(configs); i++ {
			if configs[i].ResourceID == ship.ResourceID && configs[i].ResourceType == ship.ResourceType {
				id = configs[i].ID
				break
			}
		}

		if len(id) == 0 {
			continue
		}

		if strings.HasPrefix(ship.Name, "Is ") {
			relationship := entity.AwsConfigRelationshipEntity{
				ID:                 data.ResourceID + "-" + ship.ResourceID,
				SourceID:           getId(data.ARN, data.ResourceID),
				SourceLabel:        data.ResourceName,
				SourceResourceType: data.ResourceType,
				Label:              ship.Name,
				TargetID:           id,
				TargetLabel:        ship.ResourceName,
				TargetResourceType: ship.ResourceType,
			}
			lists = append(lists, relationship)
		}

		if strings.HasPrefix(ship.Name, "Contains ") {
			relationship := entity.AwsConfigRelationshipEntity{
				ID:                 ship.ResourceID + "-" + data.ResourceID,
				SourceID:           id,
				SourceLabel:        ship.ResourceName,
				SourceResourceType: ship.ResourceType,
				Label:              ship.Name,
				TargetID:           getId(data.ARN, data.ResourceID),
				TargetLabel:        data.ResourceName,
				TargetResourceType: data.ResourceType,
			}
			lists = append(lists, relationship)
		}

	}

	return lists
}

func AddIndividualRelationShip(configs []entity.AwsConfigEntity) []entity.AwsConfigRelationshipEntity {
	sgs := make([]entity.AwsConfigEntity, 0)
	sds := make([]entity.AwsConfigEntity, 0)
	ships := make([]entity.AwsConfigRelationshipEntity, 0)

	for _, config := range configs {
		switch config.ResourceType {
		case "AWS::ECS::Service":
			sgs = append(sgs, config)
		case "AWS::ServiceDiscovery::Service":
			sds = append(sds, config)
		}
	}

	for _, config := range configs {
		ss := make([]entity.AwsConfigRelationshipEntity, 0)
		switch config.ResourceType {
		case "AWS::ECS::Service":
			ss = CreateECSServiceIndividualRelationShip(config, sds, sgs)
		case "AWS::EC2::SecurityGroup":
			ss = CreateSecurityGroupIndividualRelationShip(config, sgs)
		case "AWS::EFS::AccessPoint":
			ss = CreateAccessPointIndividualRelationShip(config, configs)
		case "AWS::EFS::FileSystem":
			ss = CreateFileSystemIndividualRelationShip(config, configs)
		}
		ships = append(ships, ss...)
	}

	return ships
}

func CreateFileSystemIndividualRelationShip(config entity.AwsConfigEntity, configs []entity.AwsConfigEntity) (ships []entity.AwsConfigRelationshipEntity) {
	ships = make([]entity.AwsConfigRelationshipEntity, 0)
	var conf c.FileSystemConfiguration
	err := json.Unmarshal([]byte(config.Configuration), &conf)
	if err != nil {
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
		return
	}

	permissions := securityGroupIPPermissions(config, sgconfig.IPPermissions, sgs)
	ships = append(ships, permissions...)
	permissions = securityGroupIPPermissions(config, sgconfig.IPPermissionsEgress, sgs)
	ships = append(ships, permissions...)
	return
}

func CreateECSServiceIndividualRelationShip(config entity.AwsConfigEntity, sds []entity.AwsConfigEntity, sgs []entity.AwsConfigEntity) (ships []entity.AwsConfigRelationshipEntity) {
	ships = make([]entity.AwsConfigRelationshipEntity, 0)
	var ecsconfig c.ECSServiceConfiguration
	err := json.Unmarshal([]byte(config.Configuration), &ecsconfig)
	if err != nil {
		return
	}

	// ServiceDiscovery Relationship
	for _, sr := range ecsconfig.ServiceRegistries {
		index := slices.IndexFunc(sds, func(sd entity.AwsConfigEntity) bool {
			return sr.RegistryArn == sd.ID
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
