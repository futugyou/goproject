package mongoimpl

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/futugyou/domaincore/domain"
)

type DBConfig struct {
	DBName         string
	CollectionName string
}

type BaseCRUD[T any] struct {
	DBName         string
	CollectionName string
	Client         *mongo.Client
	GetID          func(t T) string
}

func NewBaseCRUD[T any](client *mongo.Client, config DBConfig, getID func(t T) string) *BaseCRUD[T] {
	return &BaseCRUD[T]{
		DBName:         config.DBName,
		CollectionName: config.CollectionName,
		Client:         client,
		GetID:          getID,
	}
}

func (s *BaseCRUD[T]) FindByID(ctx context.Context, id string) (*T, error) {
	c := s.Client.Database(s.DBName).Collection(s.CollectionName)

	filter := bson.D{{Key: "id", Value: id}}
	opts := &options.FindOneOptions{}
	agg := new(T)
	if err := c.FindOne(ctx, filter, opts).Decode(&agg); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, fmt.Errorf("%s with id: %s", domain.DATA_NOT_FOUND_MESSAGE, id)
		} else {
			return nil, err
		}
	}

	return agg, nil
}

func (s *BaseCRUD[T]) Delete(ctx context.Context, id string) error {
	c := s.Client.Database(s.DBName).Collection(s.CollectionName)

	filter := bson.D{{Key: "id", Value: id}}
	opts := &options.DeleteOptions{}
	if _, err := c.DeleteOne(ctx, filter, opts); err != nil {
		return err
	}

	return nil
}

func (s *BaseCRUD[T]) SoftDelete(ctx context.Context, id string) error {
	c := s.Client.Database(s.DBName).Collection(s.CollectionName)

	filter := bson.D{{Key: "id", Value: id}}
	if _, err := c.UpdateOne(ctx, filter, bson.M{
		"$set": bson.D{{Key: "is_deleted", Value: true}},
	}); err != nil {
		return err
	}

	return nil
}

func (s *BaseCRUD[T]) Insert(ctx context.Context, data T) error {
	c := s.Client.Database(s.DBName).Collection(s.CollectionName)
	_, err := c.InsertOne(ctx, data)
	return err
}

func (s *BaseCRUD[T]) Update(ctx context.Context, data T) error {
	c := s.Client.Database(s.DBName).Collection(s.CollectionName)
	opt := options.Update().SetUpsert(true)
	filter := bson.D{{Key: "id", Value: s.GetID(data)}}
	if _, err := c.UpdateOne(ctx, filter, bson.M{
		"$set": data,
	}, opt); err != nil {
		return err
	}

	return nil
}

func (s *BaseCRUD[T]) Find(ctx context.Context, condition *domain.QueryOptions) ([]T, error) {
	c := s.Client.Database(s.DBName).Collection(s.CollectionName)
	result := make([]T, 0)

	filter := bson.D{}
	op := options.Find()

	if condition != nil {
		if condition.Filter != nil {
			filter = buildMongoFilter(condition.Filter)
		}

		var skip int64 = (int64)((condition.Page - 1) * condition.Limit)
		op.SetLimit((int64)(condition.Limit)).SetSkip(skip)

		sorter := bson.D{}
		for s, v := range condition.OrderBy {
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
