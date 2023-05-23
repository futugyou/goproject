package mongo

import (
	"context"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ClientConfig client configuration parameters
type ClientConfig struct {
	// store clients data collection name(The default is oauth2_clients)
	ClientsCName string
}

// NewDefaultClientConfig create a default client configuration
func NewDefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		ClientsCName: "oauth2_clients",
	}
}

// NewClientStore create a client store instance based on mongodb
func NewClientStore(cfg *Config, ccfgs ...*ClientConfig) *ClientStore {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URL))
	if err != nil {
		panic(err)
	}

	return NewClientStoreWithclient(client, cfg.DB, ccfgs...)
}

// NewClientStoreWithclient create a client store instance based on mongodb
func NewClientStoreWithclient(client *mongo.Client, dbName string, ccfgs ...*ClientConfig) *ClientStore {
	cs := &ClientStore{
		dbName: dbName,
		client: client,
		ccfg:   NewDefaultClientConfig(),
	}
	if len(ccfgs) > 0 {
		cs.ccfg = ccfgs[0]
	}

	return cs
}

// ClientStore MongoDB storage for OAuth 2.0
type ClientStore struct {
	ccfg   *ClientConfig
	dbName string
	client *mongo.Client
}

// Close close the mongo client
func (cs *ClientStore) Close(ctx context.Context) {
	cs.client.Disconnect(ctx)
}

func (cs *ClientStore) c(name string) *mongo.Collection {
	return cs.client.Database(cs.dbName).Collection(name)
}

func (cs *ClientStore) cHandler(name string, ctx context.Context, handler func(c *mongo.Collection)) {
	handler(cs.c(name))
}

// Set set client information
func (cs *ClientStore) Set(ctx context.Context, info oauth2.ClientInfo) (err error) {
	cs.cHandler(cs.ccfg.ClientsCName, ctx, func(c *mongo.Collection) {
		entity := &client{
			ID:     info.GetID(),
			Secret: info.GetSecret(),
			Domain: info.GetDomain(),
			UserID: info.GetUserID(),
		}

		_, err = c.InsertOne(ctx, entity)
	})

	return
}

// GetByID according to the ID for the client information
func (cs *ClientStore) GetByID(ctx context.Context, id string) (info oauth2.ClientInfo, err error) {
	cs.cHandler(cs.ccfg.ClientsCName, ctx, func(c *mongo.Collection) {
		entity := new(client)
		filter := bson.D{{Key: "ID", Value: id}}

		if err = c.FindOne(ctx, filter).Decode(&entity); err != nil {
			return
		}

		info = &models.Client{
			ID:     entity.ID,
			Secret: entity.Secret,
			Domain: entity.Domain,
			UserID: entity.UserID,
		}
	})

	return
}

// RemoveByID use the client id to delete the client information
func (cs *ClientStore) RemoveByID(ctx context.Context, id string) (err error) {
	cs.cHandler(cs.ccfg.ClientsCName, ctx, func(c *mongo.Collection) {
		filter := bson.D{{Key: "ID", Value: id}}
		if _, err = c.DeleteOne(ctx, filter); err != nil {
			return
		}
	})

	return
}

type client struct {
	ID     string `bson:"_id"`
	Secret string `bson:"secret"`
	Domain string `bson:"domain"`
	UserID string `bson:"userid"`
}
