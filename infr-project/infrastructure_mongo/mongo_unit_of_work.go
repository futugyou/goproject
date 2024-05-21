package infrastructure_mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUnitOfWork struct {
	ClientSession mongo.Session
}

func NewMongoUnitOfWork(client *mongo.Client) (*MongoUnitOfWork, error) {
	session, err := client.StartSession()
	if err != nil {
		return nil, err
	}

	return &MongoUnitOfWork{ClientSession: session}, nil
}

func (u *MongoUnitOfWork) Start(ctx context.Context) (context.Context, error) {
	if u.ClientSession == nil {
		return nil, errors.New("session not initialized")
	}
	u.ClientSession.StartTransaction()
	return mongo.NewSessionContext(ctx, u.ClientSession), nil
}

func (u *MongoUnitOfWork) Commit(ctx context.Context) error {
	if u.ClientSession == nil {
		return errors.New("session not initialized")
	}
	return u.ClientSession.CommitTransaction(ctx)
}

func (u *MongoUnitOfWork) Rollback(ctx context.Context) error {
	if u.ClientSession == nil {
		return errors.New("session not initialized")
	}
	return u.ClientSession.AbortTransaction(ctx)
}

func (u *MongoUnitOfWork) End(ctx context.Context) {
	if u.ClientSession != nil {
		u.ClientSession.EndSession(ctx)
	}
}
