package infrastructure

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	coredomain "github.com/futugyou/domaincore/domain"
	"github.com/futugyou/domaincore/mongoimpl"

	"github.com/futugyou/projectservice/domain"
)

type ProjectRepository struct {
	mongoimpl.BaseCRUD[domain.Project]
}

func NewProjectRepository(client *mongo.Client, config mongoimpl.DBConfig) *ProjectRepository {
	if config.CollectionName == "" {
		config.CollectionName = "projects"
	}

	getID := func(a domain.Project) string { return a.AggregateId() }

	return &ProjectRepository{
		BaseCRUD: *mongoimpl.NewBaseCRUD(client, config, getID),
	}
}

func (s *ProjectRepository) GetProjectByName(ctx context.Context, name string) (*domain.Project, error) {
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

func (s *ProjectRepository) GetAllProject(ctx context.Context, page *int, size *int) ([]domain.Project, error) {
	condition := coredomain.NewQueryOptions(page, size, nil, nil)
	return s.Find(ctx, condition)
}
