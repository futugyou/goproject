package infrastructure_mongo

import (
	"context"
	"fmt"

	"github.com/futugyou/infr-project/platform"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PlatformRepository struct {
	BaseRepository[*platform.Platform]
}

func NewPlatformRepository(client *mongo.Client, config DBConfig) *PlatformRepository {
	return &PlatformRepository{
		BaseRepository: *NewBaseRepository[*platform.Platform](client, config),
	}
}

func (s *PlatformRepository) GetPlatformByName(ctx context.Context, name string) (*platform.Platform, error) {
	var a platform.Platform
	c := s.Client.Database(s.DBName).Collection(a.AggregateName())

	filter := bson.D{{Key: "name", Value: name}}
	opts := &options.FindOneOptions{}
	if err := c.FindOne(ctx, filter, opts).Decode(&a); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, fmt.Errorf("data not found with name: %s", name)
		} else {
			return nil, err
		}
	}

	return &a, nil
}

func (s *PlatformRepository) GetAllPlatform(ctx context.Context) ([]platform.Platform, error) {
	var a platform.Platform
	c := s.Client.Database(s.DBName).Collection(a.AggregateName())
	result := make([]platform.Platform, 0)

	filter := bson.D{}
	cursor, err := c.Find(ctx, filter)
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
