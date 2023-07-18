package mongostore

import (
	"context"
	"log"

	"github.com/futugyousuzu/identity-server/core"
	"go.mongodb.org/mongo-driver/bson"
)

type GetAllStore[E core.IEntity] struct {
	*MongoStore
}

func NewGetAllStore[E core.IEntity](baseStore *MongoStore) *GetAllStore[E] {
	return &GetAllStore[E]{baseStore}
}

func (s *GetAllStore[E]) GetAll(ctx context.Context) ([]*E, error) {
	result := make([]*E, 0)
	entity := new(E)
	c := s.Client.Database(s.DBName).Collection((*entity).GetType())

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
