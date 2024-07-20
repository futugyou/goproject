package infrastructure_mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/extensions"
)

type BaseRepository[Aggregate domain.IAggregateRoot] struct {
	DBName string
	Client *mongo.Client
}

func NewBaseRepository[Aggregate domain.IAggregateRoot](client *mongo.Client, config DBConfig) *BaseRepository[Aggregate] {
	return &BaseRepository[Aggregate]{
		DBName: config.DBName,
		Client: client,
	}
}

func (s *BaseRepository[Aggregate]) Get(ctx context.Context, id string) (*Aggregate, error) {
	a := new(Aggregate)
	c := s.Client.Database(s.DBName).Collection((*a).AggregateName())

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

func (s *BaseRepository[Aggregate]) Delete(ctx context.Context, id string) error {
	a := new(Aggregate)
	c := s.Client.Database(s.DBName).Collection((*a).AggregateName())

	filter := bson.D{{Key: "id", Value: id}}
	opts := &options.DeleteOptions{}
	if _, err := c.DeleteOne(ctx, filter, opts); err != nil {
		return err
	}

	return nil
}

func (s *BaseRepository[Aggregate]) SoftDelete(ctx context.Context, id string) error {
	a := new(Aggregate)
	c := s.Client.Database(s.DBName).Collection((*a).AggregateName())

	filter := bson.D{{Key: "id", Value: id}}
	if _, err := c.UpdateOne(ctx, filter, bson.M{
		"$set": bson.D{{Key: "is_deleted", Value: true}},
	}); err != nil {
		return err
	}

	return nil
}

func (s *BaseRepository[Aggregate]) Insert(ctx context.Context, aggregate Aggregate) error {
	c := s.Client.Database(s.DBName).Collection(aggregate.AggregateName())
	_, err := c.InsertOne(ctx, aggregate)
	return err
}

func (s *BaseRepository[Aggregate]) Update(ctx context.Context, aggregate Aggregate) error {
	c := s.Client.Database(s.DBName).Collection(aggregate.AggregateName())
	opt := options.Update().SetUpsert(true)

	filter := bson.D{{Key: "id", Value: aggregate.AggregateId()}}
	if _, err := c.UpdateOne(ctx, filter, bson.M{
		"$set": aggregate,
	}, opt); err != nil {
		return err
	}

	return nil
}

func (s *BaseRepository[Aggregate]) GetWithCondition(ctx context.Context, condition *extensions.Search) ([]Aggregate, error) {
	a := new(Aggregate)
	c := s.Client.Database(s.DBName).Collection((*a).AggregateName())

	result := make([]Aggregate, 0)

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

	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	for _, data := range result {
		cursor.Decode(&data)
	}

	return result, nil
}

func (s *BaseRepository[Aggregate]) GetAsync(ctx context.Context, id string) (<-chan *Aggregate, <-chan error) {
	resultChan := make(chan *Aggregate, 1)
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

func (s *BaseRepository[Aggregate]) DeleteAsync(ctx context.Context, id string) <-chan error {
	errorChan := make(chan error, 1)

	go func() {
		defer close(errorChan)

		err := s.Delete(ctx, id)
		errorChan <- err
	}()

	return errorChan
}

func (s *BaseRepository[Aggregate]) SoftDeleteAsync(ctx context.Context, id string) <-chan error {
	errorChan := make(chan error, 1)

	go func() {
		defer close(errorChan)

		err := s.SoftDelete(ctx, id)
		errorChan <- err
	}()

	return errorChan
}

func (s *BaseRepository[Aggregate]) InsertAsync(ctx context.Context, aggregate Aggregate) <-chan error {
	errorChan := make(chan error, 1)

	go func() {
		defer close(errorChan)

		err := s.Insert(ctx, aggregate)
		errorChan <- err
	}()

	return errorChan
}

func (s *BaseRepository[Aggregate]) UpdateAsync(ctx context.Context, aggregate Aggregate) <-chan error {
	errorChan := make(chan error, 1)

	go func() {
		defer close(errorChan)

		err := s.Update(ctx, aggregate)
		errorChan <- err
	}()

	return errorChan
}

func (s *BaseRepository[Aggregate]) GetWithConditionAsync(ctx context.Context, condition *extensions.Search) (<-chan []Aggregate, <-chan error) {
	resultChan := make(chan []Aggregate, 1)
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
