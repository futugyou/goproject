package infrastructure_mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/futugyou/infr-project/extensions"
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
	condition := extensions.NewSearch(nil, nil, nil, map[string]interface{}{"name": name})
	ent, err := s.BaseRepository.GetWithCondition(ctx, condition)
	if err != nil {
		return nil, err
	}
	if len(ent) != 1 {
		return nil, fmt.Errorf("something wrong in 'GetProjectByName' sreach with name %s", name)
	}
	return ent[0], nil
}

func (s *ProjectRepository) GetAllProject(ctx context.Context) ([]project.Project, error) {
	condition := extensions.NewSearch(nil, nil, nil, nil)
	ent, err := s.BaseRepository.GetWithCondition(ctx, condition)
	if err != nil {
		return nil, err
	}

	result := make([]project.Project, len(ent))
	for i := 0; i < len(ent); i++ {
		result[i] = *ent[i]
	}
	return result, nil
}
