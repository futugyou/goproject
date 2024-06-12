package infrastructure_mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

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
	ent, err := s.BaseRepository.GetAggregateByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return *ent, nil
}

func (s *ProjectRepository) GetAllProject(ctx context.Context) ([]project.Project, error) {
	ent, err := s.BaseRepository.GetAllAggregate(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]project.Project, len(ent))
	for i := 0; i < len(ent); i++ {
		result[i] = *ent[i]
	}
	return result, nil
}
