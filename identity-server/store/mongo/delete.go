package mongostore

import (
	"context"
	"log"

	"github.com/futugyousuzu/identity-server/core"
	"go.mongodb.org/mongo-driver/bson"
)

type DeleteStore[E core.IEntity, K any] struct {
	*MongoStore
}

func NewDeleteStore[E core.IEntity, K any](baseStore *MongoStore) *DeleteStore[E, K] {
	return &DeleteStore[E, K]{baseStore}
}

func (s *DeleteStore[E, K]) Delete(ctx context.Context, obj E, id K) error {
	c := s.Client.Database(s.DBName).Collection(obj.GetType())
	filter := bson.D{{Key: "_id", Value: id}}
	result, err := c.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	log.Println("deleted count : ", result.DeletedCount)
	return nil
}
