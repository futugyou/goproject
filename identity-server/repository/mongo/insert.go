package mongoRepository

import (
	"context"
	"log"

	"github.com/futugyousuzu/identity-server/core"
)

type InsertRepository[E core.IEntity] struct {
	*MongoRepository
}

func NewInsertRepository[E core.IEntity](base *MongoRepository) *InsertRepository[E] {
	return &InsertRepository[E]{base}
}

func (s *InsertRepository[E]) Insert(ctx context.Context, obj E) error {
	c := s.Client.Database(s.DBName).Collection(obj.GetType())
	result, err := c.InsertOne(ctx, obj)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("insert id is: ", result.InsertedID)
	return nil
}
