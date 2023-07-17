package mongostore

import (
	"context"
	"log"

	"github.com/futugyousuzu/identity-server/core"
	"go.mongodb.org/mongo-driver/bson"
)

type GetStore[E core.IEntity, K any] struct {
	*MongoStore
}

func NewGetStore[E core.IEntity, K any](baseStore *MongoStore) *GetStore[E, K] {
	return &GetStore[E, K]{baseStore}
}

func (s *GetStore[E, K]) Get(ctx context.Context, id K) (*E, error) {
	entity := new(E)
	c := s.client.Database(s.DBName).Collection((*entity).GetType())

	filter := bson.D{{Key: "_id", Value: id}}
	err := c.FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return entity, nil
}
