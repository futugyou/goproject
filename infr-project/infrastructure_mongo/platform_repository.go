package infrastructure_mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/futugyou/infr-project/extensions"
	"github.com/futugyou/infr-project/platform"
)

type PlatformRepository struct {
	BaseRepository[platform.Platform]
}

func NewPlatformRepository(client *mongo.Client, config DBConfig) *PlatformRepository {
	return &PlatformRepository{
		BaseRepository: *NewBaseRepository[platform.Platform](client, config),
	}
}

func (s *PlatformRepository) GetPlatformByName(ctx context.Context, name string) (*platform.Platform, error) {
	condition := extensions.NewSearch(nil, nil, nil, map[string]interface{}{"name": name})
	ent, err := s.BaseRepository.GetWithCondition(ctx, condition)
	if err != nil {
		return nil, err
	}
	if len(ent) == 0 {
		return nil, fmt.Errorf("no data found with name %s", name)
	}
	return &ent[0], nil
}

func (s *PlatformRepository) GetAllPlatform(ctx context.Context) ([]platform.Platform, error) {
	condition := extensions.NewSearch(nil, nil, nil, nil)
	return s.BaseRepository.GetWithCondition(ctx, condition)
}
