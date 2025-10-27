package infrastructure_mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUnitOfWork struct {
	client        *mongo.Client
	clientSession mongo.Session
	nestedLevel   int
}

func NewMongoUnitOfWork(client *mongo.Client) (*MongoUnitOfWork, error) {
	return &MongoUnitOfWork{client: client, nestedLevel: 0}, nil
}

func (u *MongoUnitOfWork) BeginTransaction(ctx context.Context) (context.Context, error) {
	if u.client == nil {
		return nil, errors.New("client not initialized")
	}

	if session := mongo.SessionFromContext(ctx); session != nil {
		u.clientSession = session
		u.nestedLevel = u.nestedLevel + 1
		return ctx, nil
	}

	if session, err := u.client.StartSession(); err != nil {
		return nil, err
	} else {
		u.clientSession = session
		u.nestedLevel = 1
		u.clientSession.StartTransaction()
		return mongo.NewSessionContext(ctx, u.clientSession), nil
	}
}

func (u *MongoUnitOfWork) Commit(ctx context.Context) error {
	if u.clientSession == nil {
		return errors.New("session not initialized or already end")
	}

	if u.nestedLevel > 1 {
		u.nestedLevel--
		return nil
	}

	err := u.clientSession.CommitTransaction(ctx)
	u.clientSession.EndSession(ctx)
	return err
}

func (u *MongoUnitOfWork) Rollback(ctx context.Context) error {
	if u.clientSession == nil {
		return errors.New("session not initialized or already end")
	}

	u.nestedLevel = 0
	err := u.clientSession.AbortTransaction(ctx)
	u.clientSession.EndSession(ctx)
	return err
}
