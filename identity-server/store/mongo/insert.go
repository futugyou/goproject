package mongostore

import (
	"context"
	"log"

	"github.com/futugyousuzu/identity-server/core"
)

type InsertStore[E core.IEntity] struct {
	*MongoStore
}

func NewInsertStore[E core.IEntity](baseStore *MongoStore) *InsertStore[E] {
	return &InsertStore[E]{baseStore}
}

func (s *InsertStore[E]) Insert(ctx context.Context, obj E) error {
	c := s.Client.Database(s.DBName).Collection(obj.GetType())
	result, err := c.InsertOne(ctx, obj)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("insert id is: ", result.InsertedID)
	return nil
}
