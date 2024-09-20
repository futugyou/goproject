package infrastructure_mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/futugyou/infr-project/extensions"
	"github.com/futugyou/infr-project/project"
)

type ProjectRepository struct {
	BaseRepository[project.Project]
}

func NewProjectRepository(client *mongo.Client, config DBConfig) *ProjectRepository {
	return &ProjectRepository{
		BaseRepository: *NewBaseRepository[project.Project](client, config),
	}
}

func (s *ProjectRepository) GetProjectByName(ctx context.Context, name string) (*project.Project, error) {
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

func (s *ProjectRepository) GetAllProject(ctx context.Context, page *int, size *int) ([]project.Project, error) {
	condition := extensions.NewSearch(page, size, nil, nil)
	return s.BaseRepository.GetWithCondition(ctx, condition)
}

func (s *ProjectRepository) GetProjectByNameAsync(ctx context.Context, name string) (<-chan *project.Project, <-chan error) {
	var page, size int = 1, 1
	condition := extensions.NewSearch(&page, &size, nil, map[string]interface{}{"name": name})

	resultChan := make(chan *project.Project, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		result, err := s.BaseRepository.GetWithConditionAsync(ctx, condition)
		select {
		case datas := <-result:
			if len(datas) == 0 {
				errorChan <- fmt.Errorf("%s with name %s", extensions.Data_Not_Found_Message, name)
			} else {
				resultChan <- (&datas[0])
			}
		case errM := <-err:
			errorChan <- errM
		case <-ctx.Done():
			errorChan <- fmt.Errorf("GetResourceByNameAsync timeout, name %s", name)
		}
	}()

	return resultChan, errorChan
}

func (s *ProjectRepository) GetAllProjectAsync(ctx context.Context, page *int, size *int) (<-chan []project.Project, <-chan error) {
	condition := extensions.NewSearch(page, size, nil, nil)
	return s.BaseRepository.GetWithConditionAsync(ctx, condition)
}
