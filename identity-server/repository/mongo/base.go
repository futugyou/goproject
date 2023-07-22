package mongoRepository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBConfig struct {
	DBName        string
	ConnectString string
}

type MongoRepository struct {
	DBName string
	Client *mongo.Client
}

func NewMongoRepository(config DBConfig) *MongoRepository {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.ConnectString))
	return &MongoRepository{
		DBName: config.DBName,
		Client: client,
	}
}
