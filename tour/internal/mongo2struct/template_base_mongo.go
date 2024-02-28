package mongo2struct

type BaseMongoRepoConfig struct {
	PackageName     string
	BasePackageName string
	Folder          string
	FileName string
}

const base_mongorepo_TplString = `
package {{ .PackageName }}

import (
	"context"
	"log"

	"github.com/chidiwilliams/flatbson"
	"{{ .BasePackageName }}/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI(config.ConnectString))
	return &MongoRepository[E, K]{
		DBName: config.DBName,
		Client: client,
	}
}

func (s *MongoRepository[E, K]) Insert(ctx context.Context, obj E) error {
	tableNmae := obj.GetType()
	c := s.Client.Database(s.DBName).Collection(tableNmae)
	result, err := c.InsertOne(ctx, obj)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("table %s insert id %s \n", tableNmae, result.InsertedID)
	return nil
}

func (s *MongoRepository[E, K]) Delete(ctx context.Context, filter []core.DataFilterItem) error {
	obj := new(E)
	tableNmae := (*obj).GetType()
	c := s.Client.Database(s.DBName).Collection(tableNmae)
	ff := bson.D{}
	for _, val := range filter {
		ff = append(ff, bson.E(val))
	}
	result, err := c.DeleteOne(ctx, ff)
	if err != nil {
		return err
	}

	log.Printf("table %s deleted count %d \n", tableNmae, result.DeletedCount)
	return nil
}

func (s *MongoRepository[E, K]) GetOne(ctx context.Context, filter []core.DataFilterItem) (*E, error) {
	entity := new(E)
	c := s.Client.Database(s.DBName).Collection((*entity).GetType())

	ff := bson.D{}
	for _, val := range filter {
		ff = append(ff, bson.E(val))
	}
	err := c.FindOne(ctx, ff).Decode(&entity)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return entity, nil
}

func (s *MongoRepository[E, K]) Update(ctx context.Context, obj E, filter []core.DataFilterItem) error {
	tableNmae := obj.GetType()
	c := s.Client.Database(s.DBName).Collection(tableNmae)
	opt := options.Update().SetUpsert(true)
	doc, err := flatbson.Flatten(obj)
	if err != nil {
		log.Println(err)
		return err
	}

	ff := bson.D{}
	for _, val := range filter {
		ff = append(ff, bson.E(val))
	}
	result, err := c.UpdateOne(ctx, ff, bson.M{
		"$set": doc,
	}, opt)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("table %s update count %d \n", tableNmae, result.UpsertedCount)
	return nil
}

func (s *MongoRepository[E, K]) Paging(ctx context.Context, page core.Paging, filter []core.DataFilterItem) ([]E, error) {
	result := make([]E, 0)
	entity := new(E)
	c := s.Client.Database(s.DBName).Collection((*entity).GetType())

	ff := bson.D{}
	for _, val := range filter {
		ff = append(ff, bson.E(val))
	}
	var skip int64 = (page.Page - 1) * page.Limit
	op := options.Find().SetLimit(page.Limit).SetSkip(skip)
	if len(page.SortField) > 0 {
		sorted := bson.D{}
		if page.Direct == core.ASC {
			sorted = append(sorted, bson.E{Key: page.SortField, Value: 1})
		} else {
			sorted = append(sorted, bson.E{Key: page.SortField, Value: -1})
		}
		op.SetSort(sorted)
	}

	cursor, err := c.Find(ctx, ff, op)
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
 
`
