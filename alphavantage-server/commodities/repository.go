package commodities

import (
	"context"
	"fmt"
	"log"

	"github.com/futugyou/alphavantage-server/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ICommoditiesRepository interface {
	core.IRepository[CommoditiesEntity, string]
	CreateIndex(ctx context.Context) error
}

type CommoditiesRepository struct {
	*core.MongoRepository[CommoditiesEntity, string]
}

func NewCommoditiesRepository(config core.DBConfig) *CommoditiesRepository {
	baseRepo := core.NewMongoRepository[CommoditiesEntity, string](config)
	return &CommoditiesRepository{baseRepo}
}

func (a *CommoditiesRepository) CreateIndex(ctx context.Context) error {
	resource := CommoditiesEntity{}
	c := a.Client.Database(a.DBName).Collection(resource.GetType())
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "date", Value: 1}, {Key: "type", Value: 1}},
	}
	name, err := c.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}
	fmt.Println("Name of Index Created: " + name)
	return nil
}

func (s *CommoditiesRepository) GetCommoditiesByType(ctx context.Context, dataType string) ([]CommoditiesEntity, error) {
	result := make([]CommoditiesEntity, 0)
	entity := new(CommoditiesEntity)
	c := s.Client.Database(s.DBName).Collection((*entity).GetType())

	filter := bson.D{{Key: "type", Value: dataType}}
	cursor, err := c.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return result, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		log.Println(err)
		return result, err
	}

	for _, data := range result {
		cursor.Decode(&data)
	}

	return result, nil
}
