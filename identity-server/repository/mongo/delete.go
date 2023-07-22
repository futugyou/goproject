package mongoRepository

import (
	"context"
	"log"

	"github.com/futugyousuzu/identity-server/core"
	"go.mongodb.org/mongo-driver/bson"
)

type DeleteRepository[E core.IEntity, K any] struct {
	*MongoRepository
}

func NewDeleteRepository[E core.IEntity, K any](base *MongoRepository) *DeleteRepository[E, K] {
	return &DeleteRepository[E, K]{base}
}

func (s *DeleteRepository[E, K]) Delete(ctx context.Context, obj E, id K) error {
	c := s.Client.Database(s.DBName).Collection(obj.GetType())
	filter := bson.D{{Key: "_id", Value: id}}
	result, err := c.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	log.Println("deleted count : ", result.DeletedCount)
	return nil
}
