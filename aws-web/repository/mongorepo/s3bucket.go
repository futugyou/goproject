package mongorepo

import (
	"context"
	"log"

	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type S3bucketRepository struct {
	*MongoRepository[entity.S3bucketEntity, string]
}

func NewS3bucketRepository(config DBConfig) *S3bucketRepository {
	baseRepo := NewMongoRepository[entity.S3bucketEntity, string](config)
	return &S3bucketRepository{baseRepo}
}

type S3bucketItemRepository struct {
	*MongoRepository[entity.S3bucketItemEntity, string]
}

func NewS3bucketItemRepository(config DBConfig) *S3bucketItemRepository {
	baseRepo := NewMongoRepository[entity.S3bucketItemEntity, string](config)
	return &S3bucketItemRepository{baseRepo}
}

func (a *S3bucketRepository) FilterPaging(ctx context.Context, page core.Paging, filter entity.S3bucketSearchFilter) ([]*entity.S3bucketEntity, error) {
	result := make([]*entity.S3bucketEntity, 0)
	entity := new(entity.S3bucketEntity)
	c := a.Client.Database(a.DBName).Collection((*entity).GetType())

	bsonfilter := bson.M{}
	if len(filter.BucketName) > 0 {
		bsonfilter = bson.M{"name": bson.M{"$regex": filter.BucketName, "$options": "im"}}
	}

	var skip int64 = (page.Page - 1) * page.Limit
	op := options.Find().SetLimit(page.Limit).SetSkip(skip).SetSort(bson.D{{Key: "operate_at", Value: -1}})
	cursor, err := c.Find(ctx, bsonfilter, op)
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

func (a *S3bucketRepository) DeleteAll(ctx context.Context) error {
	resource := new(entity.S3bucketEntity)
	c := a.Client.Database(a.DBName).Collection(resource.GetType())
	filter := bson.D{}
	result, err := c.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	log.Println("DeletedS3bucketEntityCount: ", result.DeletedCount)
	return nil
}

func (a *S3bucketItemRepository) DeleteByBucketName(ctx context.Context, bucketName string) error {
	entity := new(entity.S3bucketItemEntity)
	c := a.Client.Database(a.DBName).Collection((*entity).GetType())
	filter := bson.D{{Key: "name", Value: bucketName}}
	_, error := c.DeleteMany(ctx, filter)
	return error
}

func (a *S3bucketItemRepository) FilterPaging(ctx context.Context, page core.Paging, filter entity.S3bucketSearchFilter) ([]*entity.S3bucketItemEntity, error) {
	result := make([]*entity.S3bucketItemEntity, 0)
	entity := new(entity.S3bucketItemEntity)
	c := a.Client.Database(a.DBName).Collection((*entity).GetType())

	bsonfilter := bson.M{}
	if len(filter.BucketName) > 0 {
		bsonfilter = bson.M{"bucketName": bson.M{"$regex": filter.BucketName, "$options": "im"}}
	}

	var skip int64 = (page.Page - 1) * page.Limit
	op := options.Find().SetLimit(page.Limit).SetSkip(skip).SetSort(bson.D{{Key: "operate_at", Value: -1}})
	cursor, err := c.Find(ctx, bsonfilter, op)
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
