package mongostore

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBConfig struct {
	DBName        string
	ConnectString string
}

type MongoStore struct {
	DBName string
	client *mongo.Client
}

func NewMongoStore(config DBConfig) *MongoStore {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.ConnectString))
	return &MongoStore{
		DBName: config.DBName,
		client: client,
	}
}
