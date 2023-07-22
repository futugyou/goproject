package mongoRepository

import (
	"context"
	"log"

	"github.com/chidiwilliams/flatbson"
	"github.com/futugyousuzu/identity-server/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UpdateRepository[E core.IEntity, K any] struct {
	*MongoRepository
}

func NewUpdateRepository[E core.IEntity, K any](base *MongoRepository) *UpdateRepository[E, K] {
	return &UpdateRepository[E, K]{base}
}

func (s *UpdateRepository[E, K]) Update(ctx context.Context, obj E, id K) error {
	c := s.Client.Database(s.DBName).Collection(obj.GetType())
	opt := options.Update().SetUpsert(true)
	doc, err := flatbson.Flatten(obj)
	if err != nil {
		log.Println(err)
		return err
	}

	filter := bson.D{{Key: "_id", Value: id}}
	result, err := c.UpdateOne(ctx, filter, bson.M{
		"$set": doc,
	}, opt)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("update count: ", result.UpsertedCount)
	return nil
}
