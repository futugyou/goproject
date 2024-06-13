package application

import (
	"context"
	"fmt"
	"strings"

	domain "github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/project"
	models "github.com/futugyou/infr-project/view_models"
)

type ProjectService struct {
	innerService *AppService
	repository   project.IProjectRepository
}

func NewProjectService(
	unitOfWork domain.IUnitOfWork,
	repository project.IProjectRepository,
) *ProjectService {
	return &ProjectService{
		innerService: NewAppService(unitOfWork),
		repository:   repository,
	}
}

func (s *ProjectService) CreateProject(request models.CreateProjectRequest) (*project.Project, error) {
	var res *project.Project
	ctx := context.Background()
	res, err := s.repository.GetProjectByName(ctx, request.Name)
	if err != nil && !strings.HasPrefix(err.Error(), "data not found") {
		return nil, err
	}

	if res != nil && res.Name == request.Name {
		return nil, fmt.Errorf("name: %s is existed", request.Name)
	}

	err = s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		res = project.NewProject(request.Name, request.Description,
			project.GetProjectState(*request.ProjectState), request.StartTime, request.EndTime)
		return s.repository.Insert(ctx, res)
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *ProjectService) GetAllProject() ([]project.Project, error) {
	return s.repository.GetAllProject(context.Background())
}

func (s *ProjectService) GetProject(id string) (*project.Project, error) {
	res, err := s.repository.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return *res, nil
}
