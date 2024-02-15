package core

import (
	"context"
	"fmt"
	"log"

	"github.com/chidiwilliams/flatbson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InsertManyResult struct {
	InsertedCount int64
	MatchedCount  int64
	ModifiedCount int64
	DeletedCount  int64
	UpsertedCount int64
}

type IRepository[E IEntity, K any] interface {
	GetAll(ctx context.Context) ([]E, error)
	Paging(ctx context.Context, page Paging) ([]E, error)
	InsertMany(ctx context.Context, items []E, filter DataFilter[E]) (*InsertManyResult, error)
}

type DBConfig struct {
	DBName        string
	ConnectString string
}

type MongoRepository[E IEntity, K any] struct {
	DBName string
	Client *mongo.Client
}

func NewMongoRepository[E IEntity, K any](config DBConfig) *MongoRepository[E, K] {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.ConnectString))
	return &MongoRepository[E, K]{
		DBName: config.DBName,
		Client: client,
	}
}

func (s *MongoRepository[E, K]) GetAll(ctx context.Context) ([]E, error) {
	result := make([]E, 0)
	entity := new(E)
	c := s.Client.Database(s.DBName).Collection((*entity).GetType())

	filter := bson.D{}
	cursor, err := c.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		log.Println(err)
		return nil, err
	}

	for _, data := range result {
		cursor.Decode(&data)
	}

	return result, nil
}

type DataFilter[E IEntity] func(e E) primitive.D

func (s *MongoRepository[E, K]) InsertMany(ctx context.Context, items []E, filter DataFilter[E]) (*InsertManyResult, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("no data need to operate")
	}

	models := make([]mongo.WriteModel, len(items))
	for i := 0; i < len(items); i++ {
		e := items[i]
		doc, err := flatbson.Flatten(e)
		if err != nil {
			log.Println("BulkWrite: ", i, err)
			continue
		}
		model := mongo.NewUpdateOneModel().
			SetFilter(filter(e)).
			SetUpsert(true).
			SetUpdate(bson.M{
				"$set": doc,
			})
		models[i] = model
	}

	item := items[0]
	tableName := item.GetType()
	c := s.Client.Database(s.DBName).Collection(tableName)
	opts := options.BulkWrite().SetOrdered(false)
	results, err := c.BulkWrite(ctx, models, opts)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := &InsertManyResult{
		InsertedCount: results.InsertedCount,
		MatchedCount:  results.DeletedCount,
		ModifiedCount: results.MatchedCount,
		DeletedCount:  results.UpsertedCount,
		UpsertedCount: results.ModifiedCount,
	}
	return result, nil
}

func (s *MongoRepository[E, K]) Paging(ctx context.Context, page Paging) ([]*E, error) {
	result := make([]*E, 0)
	entity := new(E)
	c := s.Client.Database(s.DBName).Collection((*entity).GetType())

	filter := bson.D{}
	var skip int64 = (page.Page - 1) * page.Limit
	op := options.Find().SetLimit(page.Limit).SetSkip(skip)
	if len(page.SortField) > 0 {
		if page.Direct == ASC {
			op.SetSort(bson.D{{Key: page.SortField, Value: 1}})
		} else {
			op.SetSort(bson.D{{Key: page.SortField, Value: -1}})
		}
	}

	cursor, err := c.Find(ctx, filter, op)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		log.Println(err)
		return nil, err
	}

	for _, data := range result {
		cursor.Decode(&data)
	}

	return result, nil
}
