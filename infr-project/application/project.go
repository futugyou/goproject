package application

import (
	"context"
	"fmt"
	"strings"

	domain "github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/extensions"
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

func (s *ProjectService) CreateProject(ctx context.Context, request models.CreateProjectRequest) (*models.ProjectView, error) {
	res, err := s.repository.GetProjectByName(ctx, request.Name)
	if err != nil && !strings.HasPrefix(err.Error(), extensions.Data_Not_Found_Message) {
		return nil, err
	}

	if res != nil && res.Name == request.Name {
		return nil, fmt.Errorf("name: %s is existed", request.Name)
	}

	if err = s.innerService.withUnitOfWork(ctx, func(ctx context.Context) error {
		res = project.NewProject(request.Name, request.Description,
			project.GetProjectState(request.ProjectState), request.StartTime, request.EndTime, request.Tags)
		return s.repository.Insert(ctx, *res)
	}); err != nil {
		return nil, err
	}

	return convertProjectEntityToViewModel(res), nil
}

func (s *ProjectService) GetAllProject(ctx context.Context, page *int, size *int) ([]models.ProjectView, error) {
	res, err := s.repository.GetAllProject(ctx, page, size)
	if err != nil {
		return nil, err
	}
	result := make([]models.ProjectView, len(res))
	for i := 0; i < len(res); i++ {
		result[i] = *convertProjectEntityToViewModel(&res[i])
	}
	return result, nil
}

func (s *ProjectService) GetProject(ctx context.Context, id string) (*models.ProjectView, error) {
	res, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return convertProjectEntityToViewModel(res), nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, id string, data models.UpdateProjectRequest) (*models.ProjectView, error) {
	proj, err := s.repository.Get(ctx, id)
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
	if err = s.repository.Update(ctx, *proj); err != nil {
		return nil, err
	}

	return convertProjectEntityToViewModel(proj), nil
}

func (s *ProjectService) UpdateProjectPlatform(ctx context.Context, id string, datas []models.UpdateProjectPlatformRequest) (*models.ProjectView, error) {
	proj, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	platforms := make([]project.ProjectPlatform, 0)
	for _, data := range datas {
		platforms = append(platforms, project.ProjectPlatform{
			Name:        data.Name,
			Description: data.Description,
			PlatformId:  data.PlatformId,
			ProjectId:   data.ProjectId,
		})
	}

	proj.UpdatePlatform(platforms)
	if err = s.repository.Update(ctx, *proj); err != nil {
		return nil, err
	}

	return convertProjectEntityToViewModel(proj), nil
}

func (s *ProjectService) UpdateProjectDesign(ctx context.Context, id string, datas []models.UpdateProjectDesignRequest) (*models.ProjectView, error) {
	proj, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	designes := make([]project.ProjectDesign, 0)
	for _, data := range datas {

		designes = append(designes, project.ProjectDesign{
			Name:            data.Name,
			Description:     data.Description,
			ResourceId:      data.ResourceId,
			ResourceVersion: data.ResourceVersion,
		})
	}

	proj.UpdateDesign(designes)
	if err = s.repository.Update(ctx, *proj); err != nil {
		return nil, err
	}

	return convertProjectEntityToViewModel(proj), nil
}

func convertProjectEntityToViewModel(src *project.Project) *models.ProjectView {
	if src == nil {
		return nil
	}

	platforms := make([]models.ProjectPlatform, len(src.Platforms))
	for i := 0; i < len(src.Platforms); i++ {
		platforms[i] = models.ProjectPlatform(src.Platforms[i])
	}

	design := make([]models.ProjectDesign, len(src.Designs))
	for i := 0; i < len(src.Designs); i++ {
		design[i] = models.ProjectDesign{
			Name:            src.Designs[i].Name,
			Description:     src.Designs[i].Description,
			ResourceId:      src.Designs[i].ResourceId,
			ResourceVersion: src.Designs[i].ResourceVersion,
		}
	}

	return &models.ProjectView{
		Id:          src.Id,
		Name:        src.Name,
		Description: src.Description,
		State:       src.State.String(),
		StartDate:   src.StartDate,
		EndDate:     src.EndDate,
		Platforms:   platforms,
		Designs:     design,
		Tags:        src.Tags,
	}
}
