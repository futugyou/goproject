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
	// 1. read data from file
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}

	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)

	var rawDatas []model.AwsConfigFileData

	json.Unmarshal(byteValue, &rawDatas)

	if len(rawDatas) == 0 {
		return
	}

	// 2. filter data
	rawDatas = FilterResource(rawDatas)

	// 3. get all vpc info
	vpcinfos := GetAllVpcInfos(rawDatas)

	resources := make([]entity.AwsConfigEntity, 0)
	ships := make([]entity.AwsConfigRelationshipEntity, 0)

	// 4. create AwsConfigEntity list
	for _, data := range rawDatas {
		resource := CreateAwsConfigEntity(data, vpcinfos)
		resources = append(resources, resource)
	}

	// 4.1 add individual resource
	resources = AddIndividualResource(resources)

	// 5. create AwsConfigRelationshipEntity list
	for _, data := range rawDatas {
		ship := CreateAwsConfigRelationshipEntity(data, resources)
		ships = append(ships, ship...)
	}

	// 5.1 individual relation ship
	individualShips := AddIndividualRelationShip(resources)
	ships = append(ships, individualShips...)

	// 5.2 remove duplicate
	ships = RemoveDuplicateRelationShip(ships)

	// 6. BulkWrite data to db
	log.Println("resources count: ", len(resources))
	err = a.repository.BulkWrite(context.Background(), resources)
	log.Println("resources write finish: ", err)
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
				AccountID:                    config.AccountID,
				Arn:                          config.Arn,
				AvailabilityZone:             config.AvailabilityZone,
				AwsRegion:                    config.AwsRegion,
				Configuration:                config.Configuration,
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
				LoginURL:                     config.LoginURL,
				LoggedInURL:                  config.LoggedInURL,
			},
		}
		nodes = append(nodes, node)
	}

	for _, ship := range ships {
		edge := model.Edge{
			ID:    ship.ID,
			Label: ship.Label,
			Source: model.EdgeItem{
				ID:           ship.SourceID,
				Label:        ship.SourceLabel,
				ResourceType: ship.SourceResourceType,
			},
			Target: model.EdgeItem{
				ID:           ship.TargetID,
				Label:        ship.TargetLabel,
				ResourceType: ship.TargetResourceType,
			},
		}
		edges = append(edges, edge)
	}

	return model.ResourceGraph{
		Nodes: nodes,
		Edges: edges,
	}
}
