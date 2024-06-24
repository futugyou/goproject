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
	if err != nil && !strings.HasPrefix(err.Error(), "no data found") {
		return nil, err
	}

	if res != nil && res.Name == request.Name {
		return nil, fmt.Errorf("name: %s is existed", request.Name)
	}

	err = s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		res = project.NewProject(request.Name, request.Description,
			project.GetProjectState(request.ProjectState), request.StartTime, request.EndTime, request.Tags)
		return s.repository.Insert(ctx, *res)
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
	return s.repository.Get(context.Background(), id)
}

func (s *ProjectService) UpdateProject(id string, data models.UpdateProjectRequest) (*project.Project, error) {
	proj, err := s.repository.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}

	proj.ChangeName(data.Name)
	proj.ChangeDescription(data.Description)
	sta := project.GetProjectState(data.ProjectState)
	proj.ChangeProjectState(sta)

	if data.StartTime != nil {
		proj.ChangeStartDate(*data.StartTime)
	}

	if data.EndTime != nil {
		proj.ChangeEndDate(data.EndTime)
	}

	proj.ChangeTags(data.Tags)
	err = s.repository.Update(context.Background(), *proj)
	if err != nil {
		return nil, err
	}
	return proj, nil
}

func (s *ProjectService) UpdateProjectPlatform(id string, datas []models.UpdateProjectPlatformRequest) (*project.Project, error) {
	proj, err := s.repository.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}

	platforms := make([]project.ProjectPlatform, 0)
	for _, data := range datas {
		platforms = append(platforms, project.ProjectPlatform{
			Name:        data.Name,
			Description: data.Description,
			ProjectId:   data.ProjectId,
		})
	}
	proj.UpdatePlatform(platforms)
	err = s.repository.Update(context.Background(), *proj)
	if err != nil {
		return nil, err
	}
	return proj, nil
}

func (s *ProjectService) UpdateProjectDesign(id string, datas []models.UpdateProjectDesignRequest) (*project.Project, error) {
	proj, err := s.repository.Get(context.Background(), id)
	if err != nil {
		return nil, err
	}

	designes := make([]project.ProjectDesign, 0)
	for _, data := range datas {
		designes = append(designes, project.ProjectDesign{
			Name:        data.Name,
			Description: data.Description,
			Resources:   data.Resources,
		})
	}
	proj.UpdateDesign(designes)
	err = s.repository.Update(context.Background(), *proj)
	if err != nil {
		return nil, err
	}
	return proj, nil
}
