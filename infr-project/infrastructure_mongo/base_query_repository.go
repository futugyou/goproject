package infrastructure_mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/futugyou/infr-project/extensions"
	models "github.com/futugyou/infr-project/view_models"
)

type QueryDBConfig struct {
	DBName        string
	ConnectString string
}

type BaseQueryRepository[Query models.IQuery] struct {
	DBName string
	Client *mongo.Client
}

func NewBaseQueryRepository[Query models.IQuery](client *mongo.Client, config QueryDBConfig) *BaseQueryRepository[Query] {
	return &BaseQueryRepository[Query]{
		DBName: config.DBName,
		Client: client,
	}
}

func (s *BaseQueryRepository[Query]) Get(ctx context.Context, id string) (*Query, error) {
	a := new(Query)
	c := s.Client.Database(s.DBName).Collection((*a).GetTable())

	filter := bson.D{{Key: "id", Value: id}}
	opts := &options.FindOneOptions{}
	if err := c.FindOne(ctx, filter, opts).Decode(&a); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, fmt.Errorf("data not found with id: %s", id)
		} else {
			return nil, err
		}
	}

	return a, nil
}

func (s *BaseQueryRepository[Query]) GetAll(ctx context.Context) ([]Query, error) {
	a := new(Query)
	c := s.Client.Database(s.DBName).Collection((*a).GetTable())

	filter := bson.D{}
	cursor, err := c.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var querys []Query
	if err = cursor.All(ctx, &querys); err != nil {
		return nil, err
	}

	return querys, nil
}

func (s *BaseQueryRepository[Query]) GetWithSearch(ctx context.Context, condition extensions.Search) ([]Query, error) {
	a := new(Query)
	c := s.Client.Database(s.DBName).Collection((*a).GetTable())

	filter := bson.D{}
	for key, val := range condition.Filter {
		filter = append(filter, bson.E{Key: key, Value: val})
	}

	var skip int64 = (int64)((condition.Page - 1) * condition.Size)
	op := options.Find().SetLimit((int64)(condition.Size)).SetSkip(skip)

	sorter := bson.D{}
	for s, v := range condition.Sort {
		sorter = append(sorter, bson.E{Key: s, Value: v})
	}
	if len(sorter) > 0 {
		op.SetSort(sorter)
	}

	cursor, err := c.Find(ctx, filter, op)
	if err != nil {
		return nil, err
	}

	var querys []Query
	if err = cursor.All(ctx, &querys); err != nil {
		return nil, err
	}

	return querys, nil
}
