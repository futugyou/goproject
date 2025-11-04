package dao

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"

	coredomain "github.com/futugyou/domaincore/domain"
	"github.com/futugyou/domaincore/mongoimpl"

	"github.com/futugyou/platformservice/domain"
	"github.com/futugyou/platformservice/infrastructure/internal/entity"
)

type PlatformDao struct {
	mongoimpl.BaseCRUD[entity.PlatformEntity]
}

func NewPlatformDao(client *mongo.Client, config mongoimpl.DBConfig) *PlatformDao {
	if config.CollectionName == "" {
		config.CollectionName = "platforms"
	}

	getID := func(a entity.PlatformEntity) string { return a.ID }

	return &PlatformDao{
		BaseCRUD: *mongoimpl.NewBaseCRUD(client, config, getID),
	}
}

func (s *PlatformDao) GetPlatformByName(ctx context.Context, name string) (*entity.PlatformEntity, error) {
	var page, size int = 1, 1
	query := coredomain.NewQuery().
		Eq("name", name).
		Build()
	condition := coredomain.NewQueryOptions(&page, &size, nil, query)
	ent, err := s.Find(ctx, condition)
	if err != nil {
		return nil, err
	}
	if len(ent) == 0 {
		return nil, fmt.Errorf("%s with name %s", coredomain.DATA_NOT_FOUND_MESSAGE, name)
	}
	return &ent[0], nil
}

func (s *PlatformDao) SearchPlatforms(ctx context.Context, filter domain.PlatformSearch) ([]entity.PlatformEntity, error) {
	f := s.buildSearchFilter(filter)
	condition := coredomain.NewQueryOptions(&filter.Page, &filter.Size, nil, f)
	return s.Find(ctx, condition)
}

func (s *PlatformDao) buildSearchFilter(search domain.PlatformSearch) coredomain.FilterExpr {
	var filters []coredomain.FilterExpr

	if search.Name != "" {
		if search.NameFuzzy {
			filters = append(filters, coredomain.Like{
				Field:           "name",
				Pattern:         search.Name,
				CaseInsensitive: true,
			})
		} else {
			filters = append(filters, coredomain.Eq{
				Field: "name",
				Value: search.Name,
			})
		}
	}

	if search.Activate != nil {
		filters = append(filters, coredomain.Eq{
			Field: "activate",
			Value: *search.Activate,
		})
	}

	if len(search.Tags) > 0 {
		filters = append(filters, coredomain.In{
			Field:  "tags",
			Values: anySlice(search.Tags),
		})
	}

	if len(filters) == 0 {
		return nil
	}

	if len(filters) == 1 {
		return filters[0]
	}
	return coredomain.And(filters)
}

func anySlice(slice []string) []any {
	res := make([]any, len(slice))
	for i, v := range slice {
		res[i] = v
	}
	return res
}

func (s *PlatformDao) GetPlatformByIdOrName(ctx context.Context, idOrName string) (*entity.PlatformEntity, error) {
	src, err := s.FindByID(ctx, idOrName)
	if err != nil {
		if !strings.HasPrefix(err.Error(), coredomain.DATA_NOT_FOUND_MESSAGE) {
			return nil, err
		}

		return s.GetPlatformByName(ctx, idOrName)
	}

	return src, nil
}
