package mongostore

import (
	"context"
	"log"

	"github.com/futugyousuzu/identity-server/core"
	"go.mongodb.org/mongo-driver/bson"
)

type GetAllStore[E core.IEntity, K any] struct {
	*MongoStore
}

func NewGetAllStore[E core.IEntity, K any](baseStore *MongoStore) *GetAllStore[E, K] {
	return &GetAllStore[E, K]{baseStore}
}

func (s *GetAllStore[E, K]) GetAll(ctx context.Context) ([]*E, error) {
	result := make([]*E, 0)
	entity := new(E)
	c := s.client.Database(s.DBName).Collection((*entity).GetType())

	filter := bson.D{}
	cursor, err := c.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err = cursor.All(context.TODO(), &result); err != nil {
		log.Println(err)
		return nil, err
	}

	for _, data := range result {
		cursor.Decode(&data)
	}

	return result, nil
}
