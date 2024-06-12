package infrastructure_mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/futugyou/infr-project/platform"
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
	ent, err := s.BaseRepository.GetAggregateByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return *ent, nil
}

func (s *PlatformRepository) GetAllPlatform(ctx context.Context) ([]platform.Platform, error) {
	ent, err := s.BaseRepository.GetAllAggregate(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]platform.Platform, len(ent))
	for i := 0; i < len(ent); i++ {
		result[i] = *ent[i]
	}
	return result, nil
}
