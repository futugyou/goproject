package mongo2struct

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConfig struct {
	DBName        string
	ConnectString string
}

func (m MongoDBConfig) Check() error {
	if len(m.DBName) == 0 {
		return fmt.Errorf("mongodb name can not be nil")
	}
	if len(m.ConnectString) == 0 {
		return fmt.Errorf("mongodb url can not be nil")
	}
	return nil
}

func (m *MongoDBConfig) Generator() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(m.ConnectString))
	if err != nil {
		log.Println(err)
		return
	}
	db := client.Database(m.DBName)
	filter := bson.D{}
	tables, _ := db.ListCollectionSpecifications(context.Background(), filter)
	for _, v := range tables {
		fmt.Println(v.Name)
	}
}
