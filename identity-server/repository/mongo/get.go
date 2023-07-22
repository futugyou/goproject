package mongoRepository

import (
	"context"
	"log"

	"github.com/futugyousuzu/identity-server/core"
	"go.mongodb.org/mongo-driver/bson"
)

type GetRepository[E core.IEntity, K any] struct {
	*MongoRepository
}

func NewGetRepository[E core.IEntity, K any](base *MongoRepository) *GetRepository[E, K] {
	return &GetRepository[E, K]{base}
}

func (s *GetRepository[E, K]) Get(ctx context.Context, id K) (*E, error) {
	entity := new(E)
	c := s.Client.Database(s.DBName).Collection((*entity).GetType())

	filter := bson.D{{Key: "_id", Value: id}}
	err := c.FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return entity, nil
}
