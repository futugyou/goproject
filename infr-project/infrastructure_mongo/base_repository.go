package infrastructure_mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/chidiwilliams/flatbson"
	"github.com/futugyou/infr-project/domain"
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
	opts.SetSort(bson.D{{Key: "version", Value: -1}})
	if err := c.FindOne(ctx, filter, opts).Decode(&a); err != nil {
		return nil, err
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

func (s *BaseRepository[Aggregate]) Insert(ctx context.Context, aggregate Aggregate) error {
	c := s.Client.Database(s.DBName).Collection(aggregate.AggregateName())
	_, err := c.InsertOne(ctx, aggregate)
	return err
}

func (s *BaseRepository[Aggregate]) Update(ctx context.Context, aggregate Aggregate) error {
	c := s.Client.Database(s.DBName).Collection(aggregate.AggregateName())
	opt := options.Update().SetUpsert(true)
	doc, err := flatbson.Flatten(aggregate)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "id", Value: aggregate.AggregateId()}}
	_, err = c.UpdateOne(ctx, filter, bson.M{
		"$set": doc,
	}, opt)
	if err != nil {
		return err
	}

	return nil
}
