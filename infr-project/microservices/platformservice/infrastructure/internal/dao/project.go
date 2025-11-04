package dao

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	coredomain "github.com/futugyou/domaincore/domain"
	"github.com/futugyou/domaincore/mongoimpl"

	"github.com/futugyou/platformservice/infrastructure/internal/entity"
)

type ProjectDao struct {
	mongoimpl.BaseCRUD[entity.ProjectEntity]
}

func NewProjectDao(client *mongo.Client, config mongoimpl.DBConfig) *ProjectDao {
	if config.CollectionName == "" {
		config.CollectionName = "platform_projects"
	}

	getID := func(a entity.ProjectEntity) string { return a.ID }

	return &ProjectDao{
		BaseCRUD: *mongoimpl.NewBaseCRUD(client, config, getID),
	}
}

func (p *ProjectDao) GetPlatformProjects(ctx context.Context, platformID string) ([]entity.ProjectEntity, error) {
	query := coredomain.NewQuery().
		Eq("platform_id", platformID).
		Build()
	condition := coredomain.NewQueryOptions(nil, nil, nil, query)
	return p.Find(ctx, condition)
}

func (p *ProjectDao) GetPlatformProjectByProjectID(ctx context.Context, projectID string) (*entity.ProjectEntity, error) {
	query := coredomain.NewQuery().
		Eq("id", projectID).
		Build()
	condition := coredomain.NewQueryOptions(nil, nil, nil, query)
	projects, err := p.Find(ctx, condition)
	if err != nil {
		return nil, err
	}

	if len(projects) == 0 {
		return nil, fmt.Errorf("%s with ID %s", coredomain.DATA_NOT_FOUND_MESSAGE, projectID)
	}

	return &projects[0], nil
}

func (s *ProjectDao) DeleteByPlatformID(ctx context.Context, platformID string) error {
	c := s.Client.Database(s.DBName).Collection(s.CollectionName)

	filter := bson.D{{Key: "platform_id", Value: platformID}}
	opts := &options.DeleteOptions{}
	if _, err := c.DeleteMany(ctx, filter, opts); err != nil {
		return err
	}

	return nil
}

func (s *ProjectDao) MultipleInsert(ctx context.Context, datas []entity.ProjectEntity) error {
	c := s.Client.Database(s.DBName).Collection(s.CollectionName)
	docs := make([]any, len(datas))
	for i := range datas {
		docs[i] = datas[i]
	}

	if _, err := c.InsertMany(ctx, docs); err != nil {
		return err
	}

	return nil
}
