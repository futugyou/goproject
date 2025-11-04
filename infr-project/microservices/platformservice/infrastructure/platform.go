package infrastructure

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	coredomain "github.com/futugyou/domaincore/domain"
	"github.com/futugyou/domaincore/mongoimpl"

	"github.com/futugyou/platformservice/domain"
	"github.com/futugyou/platformservice/infrastructure/internal/dao"
	"github.com/futugyou/platformservice/infrastructure/internal/entity"
)

type PlatformRepository struct {
	*dao.PlatformDao
	*dao.ProjectDao
}

func (p *PlatformRepository) Delete(ctx context.Context, id string) error {
	err := p.PlatformDao.Delete(ctx, id)
	if err != nil {
		return err
	}

	return p.ProjectDao.DeleteByPlatformID(ctx, id)
}

func (p *PlatformRepository) Find(ctx context.Context, options *coredomain.QueryOptions) ([]domain.Platform, error) {
	datas, err := p.PlatformDao.Find(ctx, options)
	if err != nil {
		return nil, err
	}

	platforms := make([]domain.Platform, 0, len(datas))
	mapper := &entity.PlatformMapper{}
	for i := range datas {
		platforms[i] = *mapper.ToDomain(&datas[i])
	}

	return platforms, nil
}

// may by use $lookup
func (p *PlatformRepository) FindByID(ctx context.Context, id string) (*domain.Platform, error) {
	platdata, err := p.PlatformDao.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	projectDatas, err := p.ProjectDao.GetPlatformProjects(ctx, id)
	if err != nil {
		return nil, err
	}

	mapper := &entity.PlatformMapper{}
	result := *mapper.ToDomain(platdata)
	projectMapper := &entity.ProjectMapper{}
	for _, proj := range projectDatas {
		pro := projectMapper.ToDomain(&proj)
		result.Projects[pro.ID] = *pro
	}

	return &result, err
}

func (p *PlatformRepository) GetPlatformByIdOrName(ctx context.Context, name string) (*domain.Platform, error) {
	datas, err := p.PlatformDao.Find(ctx, coredomain.NewQueryOptions(nil, nil, nil, coredomain.NewQuery().Eq("name", name).Or(coredomain.NewQuery().Eq("id", name)).Build()))
	if err != nil {
		return nil, err
	}

	if len(datas) == 0 {
		return nil, fmt.Errorf("platform %s not found", name)
	}

	mapper := &entity.PlatformMapper{}

	return mapper.ToDomain(&datas[0]), nil
}

func (p *PlatformRepository) GetPlatformByName(ctx context.Context, name string) (*domain.Platform, error) {
	datas, err := p.PlatformDao.Find(ctx, coredomain.NewQueryOptions(nil, nil, nil, coredomain.NewQuery().Eq("name", name).Build()))
	if err != nil {
		return nil, err
	}

	if len(datas) == 0 {
		return nil, fmt.Errorf("platform %s not found", name)
	}

	mapper := &entity.PlatformMapper{}

	return mapper.ToDomain(&datas[0]), nil
}

func (p *PlatformRepository) GetPlatformProjectByProjectID(ctx context.Context, platformID string, projectID string) (*domain.PlatformProject, error) {
	data, err := p.ProjectDao.GetPlatformProjectByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if data.PlatformID != platformID {
		return nil, fmt.Errorf("project %s not found in platform %s", projectID, platformID)
	}
	mapper := &entity.ProjectMapper{}
	return mapper.ToDomain(data), nil
}

func (p *PlatformRepository) GetPlatformProjects(ctx context.Context, platformID string) ([]domain.PlatformProject, error) {
	datas, err := p.ProjectDao.GetPlatformProjects(ctx, platformID)
	if err != nil {
		return nil, err
	}

	result := make([]domain.PlatformProject, 0, len(datas))
	mapper := &entity.ProjectMapper{}
	for i := range datas {
		result[i] = *mapper.ToDomain(&datas[i])
	}

	return result, nil
}

func (p *PlatformRepository) Insert(ctx context.Context, aggregate domain.Platform) error {
	platMapper := &entity.PlatformMapper{}
	platData := platMapper.ToEntity(&aggregate)
	err := p.PlatformDao.Insert(ctx, *platData)
	if err != nil {
		return err
	}

	projMapper := &entity.ProjectMapper{}
	projects := []entity.ProjectEntity{}
	for _, v := range aggregate.Projects {
		projects = append(projects, *projMapper.ToEntity(aggregate.ID, &v))
	}

	return p.ProjectDao.MultipleInsert(ctx, projects)
}

func (p *PlatformRepository) SearchPlatforms(ctx context.Context, filter domain.PlatformSearch) ([]domain.Platform, error) {
	datas, err := p.PlatformDao.SearchPlatforms(ctx, filter)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Platform, 0, len(datas))
	mapper := &entity.PlatformMapper{}
	for i := range datas {
		result[i] = *mapper.ToDomain(&datas[i])
	}

	return result, nil
}

func (p *PlatformRepository) Update(ctx context.Context, aggregate domain.Platform) error {
	platMapper := &entity.PlatformMapper{}
	platData := platMapper.ToEntity(&aggregate)
	return p.PlatformDao.Update(ctx, *platData)
}

func (p *PlatformRepository) UpdateProject(ctx context.Context, platformID string, project domain.PlatformProject) error {
	projectMapper := &entity.ProjectMapper{}
	projectData := projectMapper.ToEntity(platformID, &project)
	return p.ProjectDao.Update(ctx, *projectData)
}

func NewPlatformRepository(client *mongo.Client, config mongoimpl.DBConfig) *PlatformRepository {
	plat := dao.NewPlatformDao(client, config)
	proj := dao.NewProjectDao(client, config)

	return &PlatformRepository{
		PlatformDao: plat,
		ProjectDao:  proj,
	}
}

var _ domain.PlatformRepository = (*PlatformRepository)(nil)
