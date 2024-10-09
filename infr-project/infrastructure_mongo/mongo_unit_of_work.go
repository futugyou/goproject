package infrastructure_mongo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUnitOfWork struct {
	Client        *mongo.Client
	ClientSession mongo.Session
	nestedLevel   int
}

func NewMongoUnitOfWork(client *mongo.Client) (*MongoUnitOfWork, error) {
	return &MongoUnitOfWork{Client: client, nestedLevel: 0}, nil
	// session, err := client.StartSession()
	// if err != nil {
	// 	return nil, err
	// }

	// return &MongoUnitOfWork{ClientSession: session}, nil
}

func (u *MongoUnitOfWork) Start(ctx context.Context) (context.Context, error) {
	if u.Client == nil {
		return nil, errors.New("client not initialized")
	}

	if session := mongo.SessionFromContext(ctx); session != nil {
		u.ClientSession = session
		u.nestedLevel = u.nestedLevel + 2
		return ctx, nil
	}

	if session, err := u.Client.StartSession(); err != nil {
		return nil, err
	} else {
		u.ClientSession = session
		u.nestedLevel = 2
		u.ClientSession.StartTransaction()
		return mongo.NewSessionContext(ctx, u.ClientSession), nil
	}
}

func (u *MongoUnitOfWork) Commit(ctx context.Context) error {
	if u.ClientSession == nil {
		return errors.New("session not initialized")
	}

	if u.nestedLevel > 2 {
		u.nestedLevel--
		return nil
	}

	if err := u.ClientSession.CommitTransaction(ctx); err != nil {
		u.nestedLevel = 1
		return err
	}

	u.nestedLevel--
	return nil
}

func (u *MongoUnitOfWork) Rollback(ctx context.Context) error {
	if u.ClientSession == nil {
		return errors.New("session not initialized")
	}
	u.nestedLevel = 1
	return u.ClientSession.AbortTransaction(ctx)
}

func (u *MongoUnitOfWork) End(ctx context.Context) {
	if u.nestedLevel > 1 {
		u.nestedLevel--
		return
	}

	u.nestedLevel = 0
	if u.ClientSession != nil {
		u.ClientSession.EndSession(ctx)
	}
}

func (u *MongoUnitOfWork) StartAsync(ctx context.Context) (<-chan context.Context, <-chan error) {
	resultChan := make(chan context.Context, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		result, err := u.Start(ctx)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()

	return resultChan, errorChan
}

func (u *MongoUnitOfWork) CommitAsync(ctx context.Context) <-chan error {
	errorChan := make(chan error, 1)

	go func() {
		defer close(errorChan)

		err := u.Commit(ctx)
		errorChan <- err
	}()

	return errorChan
}

func (u *MongoUnitOfWork) RollbackAsync(ctx context.Context) <-chan error {
	errorChan := make(chan error, 1)

	go func() {
		defer close(errorChan)

		err := u.Rollback(ctx)
		errorChan <- err
	}()

	return errorChan
}

func (u *MongoUnitOfWork) EndAsync(ctx context.Context) <-chan error {
	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)

		if u.ClientSession != nil {
			// Handle any unexpected panics during EndSession call
			defer func() {
				if r := recover(); r != nil {
					errChan <- fmt.Errorf("panic: %v", r)
				}
			}()

			// Call EndSession asynchronously
			u.ClientSession.EndSession(ctx)
		}
	}()

	return errChan
}
