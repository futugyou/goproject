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
			return nil, fmt.Errorf("%s with id: %s", extensions.Data_Not_Found_Message, id)
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

func (s *BaseQueryRepository[Query]) GetWithCondition(ctx context.Context, condition *extensions.Search) ([]Query, error) {
	a := new(Query)
	c := s.Client.Database(s.DBName).Collection((*a).GetTable())

	filter := bson.D{}
	op := options.Find()
	if condition != nil {
		for key, val := range condition.Filter {
			filter = append(filter, bson.E{Key: key, Value: val})
		}

		var skip int64 = (int64)((condition.Page - 1) * condition.Size)
		op.SetLimit((int64)(condition.Size)).SetSkip(skip)

		sorter := bson.D{}
		for s, v := range condition.Sort {
			sorter = append(sorter, bson.E{Key: s, Value: v})
		}
		if len(sorter) > 0 {
			op.SetSort(sorter)
		}
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

func (s *BaseQueryRepository[Query]) GetAsync(ctx context.Context, id string) (<-chan *Query, <-chan error) {
	resultChan := make(chan *Query, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		result, err := s.Get(ctx, id)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()

	return resultChan, errorChan
}

func (s *BaseQueryRepository[Query]) GetAllAsync(ctx context.Context) (<-chan []Query, <-chan error) {
	resultChan := make(chan []Query, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		result, err := s.GetAll(ctx)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()

	return resultChan, errorChan
}

func (s *BaseQueryRepository[Query]) GetWithConditionAsync(ctx context.Context, condition *extensions.Search) (<-chan []Query, <-chan error) {
	resultChan := make(chan []Query, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		result, err := s.GetWithCondition(ctx, condition)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()

	return resultChan, errorChan
}
