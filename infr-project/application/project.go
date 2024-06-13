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

func (s *ProjectService) UpdateProject(id string, data models.UpdateProjectRequest) (*project.Project, error) {
	res, err := s.repository.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}

	proj := *res
	if len(*data.Name) > 0 {
		proj.ChangeName(*data.Name)
	}

	if len(*data.Description) > 0 {
		proj.ChangeDescription(*data.Description)
	}

	if len(*data.ProjectState) > 0 {
		s := project.GetProjectState(*data.ProjectState)
		proj.ChangeProjectState(s)
	}

	if data.StartTime != nil {
		proj.ChangeStartDate(*data.StartTime)
	}

	if data.EndTime != nil {
		proj.ChangeEndDate(data.EndTime)
	}

	err = s.repository.Update(context.Background(), proj)
	if err != nil {
		return nil, err
	}
	return proj, nil
}

func (s *ProjectService) UpdateProjectPlatform(id string, datas []models.UpdateProjectPlatformRequest) (*project.Project, error) {
	res, err := s.repository.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}

	proj := *res
	platforms := make([]project.ProjectPlatform, 0)
	for _, data := range datas {
		platforms = append(platforms, project.ProjectPlatform{
			Name:        data.Name,
			Description: data.Description,
			PlatformId:  data.PlatformId,
		})
	}
	proj.UpdatePlatform(platforms)
	err = s.repository.Update(context.Background(), proj)
	if err != nil {
		return nil, err
	}
	return proj, nil
}
