package mongo

import (
	"encoding/json"
	"time"

	"context"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TokenConfig token configuration parameters
type TokenConfig struct {
	// store txn collection name(The default is oauth2)
	TxnCName string
	// store token based data collection name(The default is oauth2_basic)
	BasicCName string
	// store access token data collection name(The default is oauth2_access)
	AccessCName string
	// store refresh token data collection name(The default is oauth2_refresh)
	RefreshCName string
}

// NewDefaultTokenConfig create a default token configuration
func NewDefaultTokenConfig() *TokenConfig {
	return &TokenConfig{
		TxnCName:     "oauth2_txn",
		BasicCName:   "oauth2_basic",
		AccessCName:  "oauth2_access",
		RefreshCName: "oauth2_refresh",
	}
}

// NewTokenStore create a token store instance based on mongodb
func NewTokenStore(cfg *Config, tcfgs ...*TokenConfig) (store *TokenStore) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.URL))
	if err != nil {
		panic(err)
	}

	return NewTokenStoreWithclient(client, cfg.DB, tcfgs...)
}

// NewTokenStoreWithclient create a token store instance based on mongodb
func NewTokenStoreWithclient(client *mongo.Client, dbName string, tcfgs ...*TokenConfig) (store *TokenStore) {
	ts := &TokenStore{
		dbName: dbName,
		client: client,
		tcfg:   NewDefaultTokenConfig(),
	}
	if len(tcfgs) > 0 {
		ts.tcfg = tcfgs[0]
	}

	store = ts
	return
}

// TokenStore MongoDB storage for OAuth 2.0
type TokenStore struct {
	tcfg   *TokenConfig
	dbName string
	client *mongo.Client
}

// Close close the mongo client
func (ts *TokenStore) Close() {
	ts.client.Disconnect(context.TODO())
}

func (ts *TokenStore) c(name string) *mongo.Collection {
	return ts.client.Database(ts.dbName).Collection(name)
}

func (ts *TokenStore) cHandler(name string, handler func(c *mongo.Collection)) {
	defer ts.Close()
	handler(ts.c(name))
	return
}

// Create create and store the new token information
func (ts *TokenStore) Create(info oauth2.TokenInfo) (err error) {
	jv, err := json.Marshal(info)
	if err != nil {
		return
	}

	if code := info.GetCode(); code != "" {
		ts.cHandler(ts.tcfg.BasicCName, func(c *mongo.Collection) {
			entity := basicData{
				ID:        code,
				Data:      jv,
				ExpiredAt: info.GetCodeCreateAt().Add(info.GetCodeExpiresIn()),
			}

			_, err = c.InsertOne(context.TODO(), entity)
		})

		return
	}

	aexp := info.GetAccessCreateAt().Add(info.GetAccessExpiresIn())
	rexp := aexp
	if refresh := info.GetRefresh(); refresh != "" {
		rexp = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn())
		if aexp.Second() > rexp.Second() {
			aexp = rexp
		}
	}

	var bgCtx = context.Background()
	session, err := ts.client.StartSession()
	defer func() {
		session.EndSession(bgCtx)
	}()

	id := primitive.NewObjectID().Hex()
	base := ts.client.Database(ts.dbName).Collection(ts.tcfg.BasicCName)
	base.InsertOne(bgCtx, basicData{
		ID:        id,
		Data:      jv,
		ExpiredAt: rexp,
	})

	access := ts.client.Database(ts.dbName).Collection(ts.tcfg.AccessCName)
	access.InsertOne(bgCtx, tokenData{
		ID:        info.GetAccess(),
		BasicID:   id,
		ExpiredAt: aexp,
	})

	if refresh := info.GetRefresh(); refresh != "" {
		fresh := ts.client.Database(ts.dbName).Collection(ts.tcfg.RefreshCName)
		fresh.InsertOne(bgCtx, tokenData{
			ID:        refresh,
			BasicID:   id,
			ExpiredAt: rexp,
		})
	}

	return session.CommitTransaction(bgCtx)
}

// RemoveByCode use the authorization code to delete the token information
func (ts *TokenStore) RemoveByCode(code string) (err error) {
	ts.cHandler(ts.tcfg.BasicCName, func(c *mongo.Collection) {
		filter := bson.D{{Key: "ID", Value: code}}
		if _, err = c.DeleteOne(context.TODO(), filter); err != nil {
			return
		}
	})
	return
}

// RemoveByAccess use the access token to delete the token information
func (ts *TokenStore) RemoveByAccess(access string) (err error) {
	ts.cHandler(ts.tcfg.AccessCName, func(c *mongo.Collection) {
		filter := bson.D{{Key: "ID", Value: access}}
		if _, err = c.DeleteOne(context.TODO(), filter); err != nil {
			return
		}
	})
	return
}

// RemoveByRefresh use the refresh token to delete the token information
func (ts *TokenStore) RemoveByRefresh(refresh string) (err error) {
	ts.cHandler(ts.tcfg.RefreshCName, func(c *mongo.Collection) {
		filter := bson.D{{Key: "ID", Value: refresh}}
		if _, err = c.DeleteOne(context.TODO(), filter); err != nil {
			return
		}
	})
	return
}

func (ts *TokenStore) getData(basicID string) (ti oauth2.TokenInfo, err error) {
	ts.cHandler(ts.tcfg.BasicCName, func(c *mongo.Collection) {
		var bd basicData
		filter := bson.D{{Key: "ID", Value: basicID}}
		if err = c.FindOne(context.TODO(), filter).Decode(&bd); err != nil {
			return
		}
		var tm models.Token
		err = json.Unmarshal(bd.Data, &tm)
		if err != nil {
			return
		}
		ti = &tm
	})
	return
}

func (ts *TokenStore) getBasicID(cname, token string) (basicID string, err error) {
	ts.cHandler(cname, func(c *mongo.Collection) {
		var td tokenData
		filter := bson.D{{Key: "ID", Value: token}}
		if err = c.FindOne(context.TODO(), filter).Decode(&td); err != nil {
			return
		}
		basicID = td.BasicID
	})
	return
}

// GetByCode use the authorization code for token information data
func (ts *TokenStore) GetByCode(code string) (ti oauth2.TokenInfo, err error) {
	ti, err = ts.getData(code)
	return
}

// GetByAccess use the access token for token information data
func (ts *TokenStore) GetByAccess(access string) (ti oauth2.TokenInfo, err error) {
	basicID, err := ts.getBasicID(ts.tcfg.AccessCName, access)
	if err != nil && basicID == "" {
		return
	}
	ti, err = ts.getData(basicID)
	return
}

// GetByRefresh use the refresh token for token information data
func (ts *TokenStore) GetByRefresh(refresh string) (ti oauth2.TokenInfo, err error) {
	basicID, err := ts.getBasicID(ts.tcfg.RefreshCName, refresh)
	if err != nil && basicID == "" {
		return
	}
	ti, err = ts.getData(basicID)
	return
}

type basicData struct {
	ID        string    `bson:"_id"`
	Data      []byte    `bson:"Data"`
	ExpiredAt time.Time `bson:"ExpiredAt"`
}

type tokenData struct {
	ID        string    `bson:"_id"`
	BasicID   string    `bson:"BasicID"`
	ExpiredAt time.Time `bson:"ExpiredAt"`
}
