package application

import (
	"context"
	"fmt"
	"strings"

	coreapp "github.com/futugyou/domaincore/application"
	coredomain "github.com/futugyou/domaincore/domain"

	"github.com/futugyou/projectservice/domain"

	"github.com/futugyou/projectservice/viewmodel"
)

type ProjectService struct {
	innerService *coreapp.AppService
	repository   domain.ProjectRepository
}

func NewProjectService(
	unitOfWork coredomain.UnitOfWork,
	repository domain.ProjectRepository,
) *ProjectService {
	return &ProjectService{
		innerService: coreapp.NewAppService(unitOfWork),
		repository:   repository,
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, request viewmodel.CreateProjectRequest) (*viewmodel.ProjectView, error) {
	res, err := s.repository.GetProjectByName(ctx, request.Name)
	if err != nil && !strings.HasPrefix(err.Error(), coredomain.DATA_NOT_FOUND_MESSAGE) {
		return nil, err
	}

	if res != nil && res.Name == request.Name {
		return nil, fmt.Errorf("name: %s is existed", request.Name)
	}

	if err = s.innerService.WithUnitOfWork(ctx, func(ctx context.Context) error {
		res = domain.NewProject(request.Name, request.Description,
			domain.GetProjectState(request.ProjectState), request.StartTime, request.EndTime, request.Tags)
		return s.repository.Insert(ctx, *res)
	}); err != nil {
		return nil, err
	}

	return convertProjectEntityToViewModel(res), nil
}

func (s *ProjectService) GetAllProject(ctx context.Context, page *int, size *int) ([]viewmodel.ProjectView, error) {
	res, err := s.repository.GetAllProject(ctx, page, size)
	if err != nil {
		return nil, err
	}
	result := make([]viewmodel.ProjectView, len(res))
	for i := 0; i < len(res); i++ {
		result[i] = *convertProjectEntityToViewModel(&res[i])
	}
	return result, nil
}

func (s *ProjectService) GetProject(ctx context.Context, id string) (*viewmodel.ProjectView, error) {
	res, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return convertProjectEntityToViewModel(res), nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, id string, data viewmodel.UpdateProjectRequest) (*viewmodel.ProjectView, error) {
	proj, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	proj.ChangeName(data.Name)
	proj.ChangeDescription(data.Description)
	sta := domain.GetProjectState(data.ProjectState)
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

func (s *ProjectService) UpdateProjectPlatform(ctx context.Context, id string, datas []viewmodel.UpdateProjectPlatformRequest) (*viewmodel.ProjectView, error) {
	proj, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	platforms := make([]domain.ProjectPlatform, 0)
	for _, data := range datas {
		platforms = append(platforms, domain.ProjectPlatform{
			Name:        data.Name,
			Description: data.Description,
			PlatformID:  data.PlatformID,
			ProjectID:   data.ProjectID,
		})
	}

	proj.UpdatePlatform(platforms)
	if err = s.repository.Update(ctx, *proj); err != nil {
		return nil, err
	}

	return convertProjectEntityToViewModel(proj), nil
}

func (s *ProjectService) UpdateProjectDesign(ctx context.Context, id string, datas []viewmodel.UpdateProjectDesignRequest) (*viewmodel.ProjectView, error) {
	proj, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	designes := make([]domain.ProjectDesign, 0)
	for _, data := range datas {

		designes = append(designes, domain.ProjectDesign{
			Name:            data.Name,
			Description:     data.Description,
			ResourceID:      data.ResourceID,
			ResourceVersion: data.ResourceVersion,
		})
	}

	proj.UpdateDesign(designes)
	if err = s.repository.Update(ctx, *proj); err != nil {
		return nil, err
	}

	return convertProjectEntityToViewModel(proj), nil
}

func convertProjectEntityToViewModel(src *domain.Project) *viewmodel.ProjectView {
	if src == nil {
		return nil
	}

	platforms := make([]viewmodel.ProjectPlatform, len(src.Platforms))
	for i := 0; i < len(src.Platforms); i++ {
		platforms[i] = viewmodel.ProjectPlatform(src.Platforms[i])
	}

	design := make([]viewmodel.ProjectDesign, len(src.Designs))
	for i := 0; i < len(src.Designs); i++ {
		design[i] = viewmodel.ProjectDesign{
			Name:            src.Designs[i].Name,
			Description:     src.Designs[i].Description,
			ResourceID:      src.Designs[i].ResourceID,
			ResourceVersion: src.Designs[i].ResourceVersion,
		}
	}

	return &viewmodel.ProjectView{
		ID:          src.ID,
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
