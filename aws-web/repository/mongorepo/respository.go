package mongorepo

import (
	"context"
	"fmt"
	"log"

	"github.com/chidiwilliams/flatbson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/futugyousuzu/goproject/awsgolang/core"
)

type DBConfig struct {
	DBName        string
	ConnectString string
}

type MongoRepository[E core.IEntity, K any] struct {
	DBName string
	Client *mongo.Client
}

func NewMongoRepository[E core.IEntity, K any](config DBConfig) *MongoRepository[E, K] {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.ConnectString))
	return &MongoRepository[E, K]{
		DBName: config.DBName,
		Client: client,
	}
}

func (s *MongoRepository[E, K]) Delete(ctx context.Context, id K) error {
	obj := new(E)
	c := s.Client.Database(s.DBName).Collection((*obj).GetType())
	filter := bson.D{{Key: "_id", Value: id}}
	result, err := c.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	log.Println("deleted count : ", result.DeletedCount)
	return nil
}

func (s *MongoRepository[E, K]) Get(ctx context.Context, id K) (*E, error) {
	entity := new(E)
	c := s.Client.Database(s.DBName).Collection((*entity).GetType())

	filter := bson.D{{Key: "_id", Value: id}}
	err := c.FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return entity, nil
}

func (s *MongoRepository[E, K]) GetByObjectId(ctx context.Context, id K) (*E, error) {
	stringid, ok := any(id).(string)
	if !ok {
		return nil, fmt.Errorf("%v can not convert to string", id)
	}

	entity := new(E)
	c := s.Client.Database(s.DBName).Collection((*entity).GetType())
	objID, _ := primitive.ObjectIDFromHex(stringid)

	filter := bson.D{{Key: "_id", Value: objID}}
	err := c.FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return entity, nil
}

func (s *MongoRepository[E, K]) GetAll(ctx context.Context) ([]*E, error) {
	result := make([]*E, 0)
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

func (s *MongoRepository[E, K]) Insert(ctx context.Context, obj E) error {
	c := s.Client.Database(s.DBName).Collection(obj.GetType())
	result, err := c.InsertOne(ctx, obj)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("insert id is: ", result.InsertedID)
	return nil
}

func (s *MongoRepository[E, K]) Update(ctx context.Context, obj E, id K) error {
	c := s.Client.Database(s.DBName).Collection(obj.GetType())
	opt := options.Update().SetUpsert(true)
	doc, err := flatbson.Flatten(obj)
	if err != nil {
		log.Println(err)
		return err
	}

	filter := bson.D{{Key: "_id", Value: id}}
	result, err := c.UpdateOne(ctx, filter, bson.M{
		"$set": doc,
	}, opt)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("update count: ", result.UpsertedCount)
	return nil
}

func (s *MongoRepository[E, K]) InsertMany(ctx context.Context, items []E) error {
	if len(items) == 0 {
		return nil
	}

	item := items[0]
	c := s.Client.Database(s.DBName).Collection(item.GetType())
	entitys := make([]interface{}, len(items))
	for i := 0; i < len(items); i++ {
		entitys[i] = items[i]
	}

	result, err := c.InsertMany(ctx, entitys)
	if err != nil {
		return err
	}

	log.Println("Inserted Count: ", len(result.InsertedIDs))
	return nil
}

func (s *MongoRepository[E, K]) Paging(ctx context.Context, page core.Paging) ([]*E, error) {
	result := make([]*E, 0)
	entity := new(E)
	c := s.Client.Database(s.DBName).Collection((*entity).GetType())

	filter := bson.D{}
	var skip int64 = (page.Page - 1) * page.Limit
	op := options.Find().SetLimit(page.Limit).SetSkip(skip)
	if len(page.SortField) > 0 {
		if page.Direct == core.ASC {
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

func (s *MongoRepository[E, K]) BulkOperate(ctx context.Context, models []mongo.WriteModel) error {
	if len(models) == 0 {
		return nil
	}

	entity := new(E)
	tableName := (*entity).GetType()
	c := s.Client.Database(s.DBName).Collection(tableName)
	opts := options.BulkWrite().SetOrdered(false)
	results, err := c.BulkWrite(ctx, models, opts)
	if err != nil {
		return err
	}

	log.Printf("%s, the number of documents inserted: %d\n", tableName, results.InsertedCount)
	log.Printf("%s, the number of documents deleted: %d\n", tableName, results.DeletedCount)
	log.Printf("%s, the number of documents matched by filters in update and replace operations: %d\n", tableName, results.MatchedCount)
	log.Printf("%s, the number of documents upserted by update and replace operations: %d\n", tableName, results.UpsertedCount)
	log.Printf("%s, the number of documents modified by update and replace operations: %d\n", tableName, results.ModifiedCount)

	return nil
}
