package mongostore

import (
	"context"
	"log"

	"github.com/chidiwilliams/flatbson"
	"github.com/futugyousuzu/identity-server/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UpdateStore[E core.IEntity, K any] struct {
	*MongoStore
}

func NewUpdateStore[E core.IEntity, K any](baseStore *MongoStore) *UpdateStore[E, K] {
	return &UpdateStore[E, K]{baseStore}
}

func (s *UpdateStore[E, K]) Update(ctx context.Context, obj E, id K) error {
	c := s.client.Database(s.DBName).Collection(obj.GetType())
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
