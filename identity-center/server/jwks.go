package server

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

type JwksStore interface {
	GetJwksList(ctx context.Context) (string, error)
	CreateJwks(ctx context.Context, signed_key string) error
}

type JwkModel struct {
	ID      string `bson:"_id"`
	Payload string `bson:"payload"`
}

type MongoJwksStore struct {
	DBName         string
	CollectionName string
	client         *mongo.Client
}

func NewJwksStore() *MongoJwksStore {
	db := os.Getenv("db_name")
	u_name := "oauth2_jwks"
	url := os.Getenv("mongodb_url")
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))

	return &MongoJwksStore{
		DBName:         db,
		CollectionName: u_name,
		client:         client,
	}
}

func (u *MongoJwksStore) GetJwksList(ctx context.Context) (string, error) {
	coll := u.client.Database(u.DBName).Collection(u.CollectionName)

	result := make([]JwkModel, 0)
	filter := bson.D{}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	// end find

	if err = cursor.All(context.TODO(), &result); err != nil {
		panic(err)
	}

	for _, data := range result {
		cursor.Decode(&data)
	}

	s := make([]string, len(result))
	for i, v := range result {
		s[i] = v.Payload
	}

	return "{ \"keys\": [ " + strings.Join(s, ",") + " ] }", err
}

func (u *MongoJwksStore) CreateJwks(ctx context.Context, signed_key string) error {
	coll := u.client.Database(u.DBName).Collection(u.CollectionName)
	var jwkModel JwkModel
	err := coll.FindOne(ctx, bson.D{{Key: "_id", Value: signed_key}}).Decode(&jwkModel)
	if err == nil {
		return nil
	}

	raw, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate new rsa private key: %s", err)
	}

	key, err := jwk.FromRaw(raw)
	if err != nil {
		return fmt.Errorf("failed to create RSA key: %s", err)
	}
	if _, ok := key.(jwk.RSAPrivateKey); !ok {
		return fmt.Errorf("expected jwk.RSAPrivateKey, got %T", err)
	}

	key.Set(jwk.KeyIDKey, signed_key)

	buf, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal key into JSON: %s", err)
	}

	jwkModel = JwkModel{
		ID:      signed_key,
		Payload: string(buf),
	}

	_, err = coll.InsertOne(ctx, jwkModel)
	return err
}
