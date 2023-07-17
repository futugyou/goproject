package mongostore

import (
	"context"
	"log"

	"github.com/futugyousuzu/identity-server/core"
)

type InsertStore[E core.IEntity, K any] struct {
	*MongoStore
}

func NewInsertStore[E core.IEntity, K any](baseStore *MongoStore) *InsertStore[E, K] {
	return &InsertStore[E, K]{baseStore}
}

func (s *InsertStore[E, K]) Insert(ctx context.Context, obj E) error {
	c := s.client.Database(s.DBName).Collection(obj.GetType())
	result, err := c.InsertOne(ctx, obj)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("insert id is: ", result.InsertedID)
	return nil
}
