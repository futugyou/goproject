package services

import (
	"context"
	"os"

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

func (a *AwsConfigService) GetResourceGraph(ctx context.Context) model.ResourceGraph {
	configs, _ := a.repository.GetAll(ctx)
	ships, _ := a.relRepository.GetAll(ctx)
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
