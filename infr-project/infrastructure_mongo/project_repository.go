package infrastructure_mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/futugyou/infr-project/project"
)

type ProjectRepository struct {
	BaseRepository[*project.Project]
}

func NewProjectRepository(client *mongo.Client, config DBConfig) *ProjectRepository {
	return &ProjectRepository{
		BaseRepository: *NewBaseRepository[*project.Project](client, config),
	}
}

func (s *ProjectRepository) GetProjectByName(ctx context.Context, name string) (*project.Project, error) {
	var a project.Project
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

func (s *ProjectRepository) GetAllProject(ctx context.Context) ([]project.Project, error) {
	var a project.Project
	c := s.Client.Database(s.DBName).Collection(a.AggregateName())
	result := make([]project.Project, 0)

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
