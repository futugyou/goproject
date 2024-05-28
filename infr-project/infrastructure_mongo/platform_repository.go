package infrastructure_mongo

import (
	"github.com/futugyou/infr-project/platform"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlatformRepository struct {
	BaseRepository[*platform.Platform]
}

func NewPlatformRepository(client *mongo.Client, config DBConfig) *PlatformRepository {
	return &PlatformRepository{
		BaseRepository: *NewBaseRepository[*platform.Platform](client, config),
	}
}
