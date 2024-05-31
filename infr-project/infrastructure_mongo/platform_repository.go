package infrastructure_mongo

import (
	"context"

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
	a := new(platform.Platform)
	c := s.Client.Database(s.DBName).Collection((*a).AggregateName())

	filter := bson.D{{Key: "name", Value: name}}
	opts := &options.FindOneOptions{}
	if err := c.FindOne(ctx, filter, opts).Decode(&a); err != nil {
		return nil, err
	}

	return a, nil
}
