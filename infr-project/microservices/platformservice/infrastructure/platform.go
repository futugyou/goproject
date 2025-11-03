package infrastructure

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	coredomain "github.com/futugyou/domaincore/domain"
	"github.com/futugyou/domaincore/mongoimpl"

	"github.com/futugyou/platformservice/domain"
)

type PlatformRepository struct {
	mongoimpl.BaseCRUD[domain.Platform]
}

func NewPlatformRepository(client *mongo.Client, config mongoimpl.DBConfig) *PlatformRepository {
	if config.CollectionName == "" {
		config.CollectionName = "platforms"
	}

	getID := func(a domain.Platform) string { return a.AggregateId() }

	return &PlatformRepository{
		BaseCRUD: *mongoimpl.NewBaseCRUD(client, config, getID),
	}
}

func (s *PlatformRepository) GetPlatformByName(ctx context.Context, name string) (*domain.Platform, error) {
	var page, size int = 1, 1
	condition := coredomain.NewQueryOptions(&page, &size, nil, map[string]any{"name": name})
	ent, err := s.Find(ctx, condition)
	if err != nil {
		return nil, err
	}
	if len(ent) == 0 {
		return nil, fmt.Errorf("%s with name %s", coredomain.DATA_NOT_FOUND_MESSAGE, name)
	}
	return &ent[0], nil
}

func (s *PlatformRepository) SearchPlatforms(ctx context.Context, filter domain.PlatformSearch) ([]domain.Platform, error) {
	f := s.buildSearchFilter(filter)
	condition := coredomain.NewQueryOptions(&filter.Page, &filter.Size, nil, f)
	return s.Find(ctx, condition)
}

func (s *PlatformRepository) buildSearchFilter(search domain.PlatformSearch) map[string]any {
	filter := map[string]any{}

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

func (s *PlatformRepository) GetPlatformByIdOrName(ctx context.Context, idOrName string) (*domain.Platform, error) {
	src, err := s.FindByID(ctx, idOrName)
	if err != nil {
		if !strings.HasPrefix(err.Error(), coredomain.DATA_NOT_FOUND_MESSAGE) {
			return nil, err
		}

		return s.GetPlatformByName(ctx, idOrName)
	}

	return src, nil
}
