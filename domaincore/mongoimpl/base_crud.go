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

type BaseCRUD[Aggregate domain.AggregateRoot] struct {
	DBName         string
	CollectionName string
	Client         *mongo.Client
}

func NewBaseCRUD[Aggregate domain.AggregateRoot](client *mongo.Client, config DBConfig) *BaseCRUD[Aggregate] {
	return &BaseCRUD[Aggregate]{
		DBName:         config.DBName,
		CollectionName: config.CollectionName,
		Client:         client,
	}
}

func (s *BaseCRUD[Aggregate]) FindByID(ctx context.Context, id string) (*Aggregate, error) {
	c := s.Client.Database(s.DBName).Collection(s.CollectionName)

	filter := bson.D{{Key: "id", Value: id}}
	opts := &options.FindOneOptions{}
	agg := new(Aggregate)
	if err := c.FindOne(ctx, filter, opts).Decode(&agg); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, fmt.Errorf("%s with id: %s", domain.DATA_NOT_FOUND_MESSAGE, id)
		} else {
			return nil, err
		}
	}

	return agg, nil
}

func (s *BaseCRUD[Aggregate]) Delete(ctx context.Context, id string) error {
	c := s.Client.Database(s.DBName).Collection(s.CollectionName)

	filter := bson.D{{Key: "id", Value: id}}
	opts := &options.DeleteOptions{}
	if _, err := c.DeleteOne(ctx, filter, opts); err != nil {
		return err
	}

	return nil
}

func (s *BaseCRUD[Aggregate]) SoftDelete(ctx context.Context, id string) error {
	c := s.Client.Database(s.DBName).Collection(s.CollectionName)

	filter := bson.D{{Key: "id", Value: id}}
	if _, err := c.UpdateOne(ctx, filter, bson.M{
		"$set": bson.D{{Key: "is_deleted", Value: true}},
	}); err != nil {
		return err
	}

	return nil
}

func (s *BaseCRUD[Aggregate]) Insert(ctx context.Context, aggregate Aggregate) error {
	c := s.Client.Database(s.DBName).Collection(s.CollectionName)
	_, err := c.InsertOne(ctx, aggregate)
	return err
}

func (s *BaseCRUD[Aggregate]) Update(ctx context.Context, aggregate Aggregate) error {
	c := s.Client.Database(s.DBName).Collection(s.CollectionName)
	opt := options.Update().SetUpsert(true)

	filter := bson.D{{Key: "id", Value: aggregate.AggregateId()}}
	if _, err := c.UpdateOne(ctx, filter, bson.M{
		"$set": aggregate,
	}, opt); err != nil {
		return err
	}

	return nil
}

func buildMongoFilter(input any) bson.D {
	switch val := input.(type) {
	case map[string]any:
		if len(val) == 1 {
			for k, v := range val {
				switch k {
				case "$or", "$and", "$nor":
					arr := bson.A{}
					if subList, ok := v.([]any); ok {
						for _, sub := range subList {
							arr = append(arr, buildMongoFilter(sub))
						}
					}
					return bson.D{{Key: k, Value: arr}}
				}
			}
		}

		doc := bson.D{}
		for k, v := range val {
			doc = append(doc, bson.E{Key: k, Value: v})
		}
		return doc

	default:
		return bson.D{}
	}
}

func (s *BaseCRUD[Aggregate]) Find(ctx context.Context, condition *domain.QueryOptions) ([]Aggregate, error) {
	c := s.Client.Database(s.DBName).Collection(s.CollectionName)
	result := make([]Aggregate, 0)

	filter := bson.D{}
	op := options.Find()

	if condition != nil {
		if condition.LogicFilter != nil {
			filter = buildMongoFilter(condition.LogicFilter)
		} else if len(condition.Filters) > 0 {
			for key, val := range condition.Filters {
				filter = append(filter, bson.E{Key: key, Value: val})
			}
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
