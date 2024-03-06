package dynamo2struct

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBConfig struct {
	EntityFolder     string
	RepoFolder       string
	CoreFoler        string
	DynamoRepoFolder string
	PkgName          string
	AccessKey        string
	AccessSecret     string
	Region           string
}

func (m *DynamoDBConfig) Check() error {
	if len(m.AccessKey) == 0 {
		return fmt.Errorf("dynamodb AccessKey can not be nil")
	}
	if len(m.AccessSecret) == 0 {
		return fmt.Errorf("dynamodb AccessSecret can not be nil")
	}
	if len(m.Region) == 0 {
		return fmt.Errorf("dynamodb Region can not be nil")
	}
	return nil
}

func (m *DynamoDBConfig) ConnectDBDatabase() (*dynamodb.Client, error) {
	if err := m.Check(); err != nil {
		log.Println(err)
		return nil, err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(m.AccessKey, m.AccessSecret, ""),
		),
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	cfg.Region = m.Region
	svc := dynamodb.NewFromConfig(cfg)

	return svc, nil
}
