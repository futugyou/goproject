package server

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RequestAuthorizeInfo struct {
	ID                   string `bson:"_id"`
	ClientId             string `bson:"client_id"`
	CodeChallenge        string `bson:"code_challenge"`
	CodeChallenge_method string `bson:"code_challenge_method"`
	RedirectUri          string `bson:"redirect_uri"`
	ResponseType         string `bson:"response_type"`
}

type RequestAuthorizeInfoStore interface {
	Create(ctx context.Context, info RequestAuthorizeInfo) error
	Get(ctx context.Context, code_challenge string) (RequestAuthorizeInfo, error)
}

type MongoRequestAuthorizeInfoStore struct {
	DBName         string
	CollectionName string
	client         *mongo.Client
}

func NewMongoRequestAuthorizeInfoStore() *MongoRequestAuthorizeInfoStore {
	db := os.Getenv("db_name")
	c_name := "oauth2_authorize_infos"
	url := os.Getenv("mongodb_url")
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))

	return &MongoRequestAuthorizeInfoStore{
		DBName:         db,
		CollectionName: c_name,
		client:         client,
	}
}

func (s *MongoRequestAuthorizeInfoStore) Create(ctx context.Context, info RequestAuthorizeInfo) error {
	c := s.client.Database(s.DBName).Collection(s.CollectionName)
	_, err := c.InsertOne(ctx, info)
	return err
}

func (s *MongoRequestAuthorizeInfoStore) Get(ctx context.Context, code_challenge string) (RequestAuthorizeInfo, error) {
	c := s.client.Database(s.DBName).Collection(s.CollectionName)
	entity := new(RequestAuthorizeInfo)
	filter := bson.D{{Key: "code_challenge", Value: code_challenge}}
	err := c.FindOne(ctx, filter).Decode(&entity)
	return *entity, err
}
