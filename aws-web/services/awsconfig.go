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

	var datas []model.AwsConfigFileData

	json.Unmarshal(byteValue, &datas)

	if len(datas) == 0 {
		return
	}

	// 2. filter data
	datas = filterResource(datas)

	// 3. get all vpc info
	vpcinfos := GetAllVpcInfos(datas)

	configs := make([]entity.AwsConfigEntity, 0)
	ships := make([]entity.AwsConfigRelationshipEntity, 0)

	// 4. create AwsConfigEntity list
	for _, data := range datas {
		config := CreateAwsConfigEntity(data, vpcinfos)
		configs = append(configs, config)
	}

	// 4.1 add cloud map data
	configs = AddIndividualData(configs)

	// 5. create AwsConfigRelationshipEntity list
	for _, data := range datas {
		ship := CreateAwsConfigRelationshipEntity(data, configs)
		ships = append(ships, ship...)
	}

	// 5.1 individual relation ship
	individualShips := AddIndividualRelationShip(configs)
	ships = append(ships, individualShips...)

	// 6. BulkWrite data to db
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
