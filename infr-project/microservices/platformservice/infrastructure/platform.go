package infrastructure

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
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
	client *mongo.Client
	config mongoimpl.DBConfig
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

	platforms := make([]domain.Platform, len(datas))
	mapper := &entity.PlatformMapper{}
	for i := range datas {
		platforms[i] = *mapper.ToDomain(&datas[i])
	}

	return platforms, nil
}

type PlatformWithProjects struct {
	entity.PlatformEntity `bson:",inline"`
	Projects              []entity.ProjectEntity `bson:"projects"`
}

// may by use $lookup
func (p *PlatformRepository) FindByID(ctx context.Context, id string) (*domain.Platform, error) {
	filter := bson.D{{Key: "id", Value: id}}

	return p.useAggregateSearch(ctx, filter, id)
}

func (p *PlatformRepository) useAggregateSearch(ctx context.Context, filter bson.D, id string) (*domain.Platform, error) {
	coll := p.client.Database(p.config.DBName).Collection("platforms")
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "platform_projects"},
			{Key: "localField", Value: "id"},
			{Key: "foreignField", Value: "platform_id"},
			{Key: "as", Value: "projects"},
		}}},
		{{Key: "$limit", Value: 1}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate operate error: %w", err)
	}
	defer cursor.Close(ctx)

	var out PlatformWithProjects
	if cursor.Next(ctx) {
		if err := cursor.Decode(&out); err != nil {
			return nil, fmt.Errorf("decode error: %w", err)
		}
	}

	if len(out.ID) == 0 {
		return nil, fmt.Errorf("platform %s not found", id)
	}

	platdata := &entity.PlatformEntity{
		ID:          out.ID,
		Name:        out.Name,
		Description: out.Description,
		Url:         out.Url,
		Provider:    out.Provider,
		Properties:  out.Properties,
		Secrets:     out.Secrets,
		Tags:        out.Tags,
		IsDeleted:   out.IsDeleted,
	}

	mapper := &entity.PlatformMapper{}
	result := *mapper.ToDomain(platdata)
	projectMapper := &entity.ProjectMapper{}
	for _, proj := range out.Projects {
		pro := projectMapper.ToDomain(&proj)
		result.Projects[pro.ID] = *pro
	}

	return &result, err
}

func (p *PlatformRepository) GetPlatformByIdOrName(ctx context.Context, idOrName string) (*domain.Platform, error) {
	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "id", Value: idOrName}},
			bson.D{{Key: "name", Value: idOrName}},
		}},
	}

	return p.useAggregateSearch(ctx, filter, idOrName)
}

func (p *PlatformRepository) GetPlatformByIdOrNameWithoutProjects(ctx context.Context, idOrName string) (*domain.Platform, error) {
	datas, err := p.PlatformDao.Find(ctx, coredomain.NewQueryOptions(nil, nil, nil, coredomain.NewQuery().Eq("name", idOrName).Or(coredomain.NewQuery().Eq("id", idOrName)).Build()))
	if err != nil {
		return nil, err
	}

	if len(datas) == 0 {
		return nil, fmt.Errorf("platform %s not found", idOrName)
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

func (p *PlatformRepository) GetPlatformProjectByIDOrName(ctx context.Context, platformID string, projectID string) (*domain.PlatformProject, error) {
	data, err := p.ProjectDao.GetPlatformProjectByIDOrName(ctx, projectID)
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

	result := make([]domain.PlatformProject, len(datas))
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

	result := make([]domain.Platform, len(datas))
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

func (p *PlatformRepository) DeleteProject(ctx context.Context, platformID string, projectID string) error {
	return p.ProjectDao.Delete(ctx, projectID)
}

func (p *PlatformRepository) SyncProjects(ctx context.Context, platformID string, projects []domain.PlatformProject) error {
	projMapper := &entity.ProjectMapper{}
	projectEntities := []entity.ProjectEntity{}
	for _, v := range projects {
		projectEntities = append(projectEntities, *projMapper.ToEntity(platformID, &v))
	}

	return p.ProjectDao.SyncProjects(ctx, platformID, projectEntities)
}

func NewPlatformRepository(client *mongo.Client, config mongoimpl.DBConfig) *PlatformRepository {
	plat := dao.NewPlatformDao(client, config)
	proj := dao.NewProjectDao(client, config)

	return &PlatformRepository{
		PlatformDao: plat,
		ProjectDao:  proj,
		client:      client,
		config:      config,
	}
}

var _ domain.PlatformRepository = (*PlatformRepository)(nil)
