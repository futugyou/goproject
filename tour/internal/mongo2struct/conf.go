package mongo2struct

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConfig struct {
	EntityFolder    string
	RepoFolder      string
	CoreFoler       string
	MongoRepoFolder string
	PkgName         string
	DBName          string
	ConnectString   string
}

func (m *MongoDBConfig) Check() error {
	if len(m.DBName) == 0 {
		return fmt.Errorf("mongodb name can not be nil")
	}
	if len(m.ConnectString) == 0 {
		return fmt.Errorf("mongodb url can not be nil")
	}
	if len(m.EntityFolder) == 0 {
		m.EntityFolder = "entity"
	}
	if len(m.RepoFolder) == 0 {
		m.RepoFolder = "repository"
	}
	return nil
}

func (m *MongoDBConfig) ConnectDBDatabase() (*mongo.Database, error) {
	if err := m.Check(); err != nil {
		log.Println(err)
		return nil, err
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(m.ConnectString))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return client.Database(m.DBName), nil
}
