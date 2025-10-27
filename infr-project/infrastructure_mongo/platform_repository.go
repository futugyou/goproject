package infrastructure_mongo

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
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
	var page, size int = 1, 1
	condition := extensions.NewSearch(&page, &size, nil, map[string]interface{}{"name": name})
	ent, err := s.BaseRepository.GetWithCondition(ctx, condition)
	if err != nil {
		return nil, err
	}
	if len(ent) == 0 {
		return nil, fmt.Errorf("%s with name %s", extensions.Data_Not_Found_Message, name)
	}
	return &ent[0], nil
}

func (s *PlatformRepository) SearchPlatforms(ctx context.Context, filter platform.PlatformSearch) ([]platform.Platform, error) {
	f := s.buildSearchFilter(filter)
	condition := extensions.NewSearch(&filter.Page, &filter.Size, nil, f)
	return s.BaseRepository.GetWithCondition(ctx, condition)
}

func (s *PlatformRepository) buildSearchFilter(search platform.PlatformSearch) map[string]interface{} {
	filter := map[string]interface{}{}

	if search.Name != "" {
		if search.NameFuzzy {
			filter["name"] = bson.D{{Key: "$regex", Value: search.Name}, {Key: "$options", Value: "i"}}
		} else {
			filter["name"] = search.Name
		}
	}

	if search.Activate != nil {
		filter["activate"] = &search.Activate
	}

	if len(search.Tags) > 0 {
		filter["tags"] = bson.D{{Key: "$in", Value: search.Tags}}
	}

	return filter
}

func (s *PlatformRepository) GetPlatformByIdOrName(ctx context.Context, idOrName string) (*platform.Platform, error) {
	src, err := s.Get(ctx, idOrName)
	if err != nil {
		if !strings.HasPrefix(err.Error(), extensions.Data_Not_Found_Message) {
			return nil, err
		}

		return s.GetPlatformByName(ctx, idOrName)
	}

	return src, nil
}
