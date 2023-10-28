package mongorepo

import (
	"context"
	"log"

	"github.com/chidiwilliams/flatbson"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AwsConfigRepository struct {
	*MongoRepository[entity.AwsConfigEntity, string]
}

func NewAwsConfigRepository(config DBConfig) *AwsConfigRepository {
	baseRepo := NewMongoRepository[entity.AwsConfigEntity, string](config)
	return &AwsConfigRepository{baseRepo}
}

type AwsConfigRelationshipRepository struct {
	*MongoRepository[entity.AwsConfigRelationshipEntity, string]
}

func NewAwsConfigRelationshipRepository(config DBConfig) *AwsConfigRelationshipRepository {
	baseRepo := NewMongoRepository[entity.AwsConfigRelationshipEntity, string](config)
	return &AwsConfigRelationshipRepository{baseRepo}
}

func (a *AwsConfigRepository) BulkWrite(ctx context.Context, entities []entity.AwsConfigEntity) error {
	models := make([]mongo.WriteModel, len(entities))
	for i := 0; i < len(entities); i++ {
		e := entities[i]
		doc, err := flatbson.Flatten(e)
		if err != nil {
			log.Println("BulkWrite: ", i, err)
			continue
		}

		filter := bson.D{{Key: "resourceType", Value: e.ResourceType}, {Key: "resourceId", Value: e.ResourceID}}
		model := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpsert(true).
			SetUpdate(bson.M{
				"$set": doc,
			})
		models[i] = model
	}

	return a.BulkOperate(ctx, models)
}

func (a *AwsConfigRelationshipRepository) BulkWrite(ctx context.Context, entities []entity.AwsConfigRelationshipEntity) error {
	models := make([]mongo.WriteModel, len(entities))
	for i := 0; i < len(entities); i++ {
		e := entities[i]
		doc, err := flatbson.Flatten(e)
		if err != nil {
			log.Println("BulkWrite: ", i, err)
			continue
		}

		filter := bson.D{{Key: "sourceId", Value: e.SourceID}, {Key: "targetId", Value: e.TargetID}}
		model := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpsert(true).
			SetUpdate(bson.M{
				"$set": doc,
			})
		models[i] = model
	}

	return a.BulkOperate(ctx, models)
}

func (a *AwsConfigRepository) DeleteAll(ctx context.Context) error {
	resource := new(entity.AwsConfigEntity)
	c := a.Client.Database(a.DBName).Collection(resource.GetType())
	filter := bson.D{}
	result, err := c.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	log.Println("DeletedAwsConfigEntityCount: ", result.DeletedCount)
	return nil
}

func (a *AwsConfigRelationshipRepository) DeleteAll(ctx context.Context) error {
	ship := new(entity.AwsConfigRelationshipEntity)
	c := a.Client.Database(a.DBName).Collection(ship.GetType())
	filter := bson.D{}
	result, err := c.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	log.Println("DeletedAwsConfigRelationshipEntityCount: ", result.DeletedCount)
	return nil
}
