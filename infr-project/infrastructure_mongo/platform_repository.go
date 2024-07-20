package infrastructure_mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/futugyou/infr-project/extensions"
	"github.com/futugyou/infr-project/platform"
)

type PlatformRepository struct {
	BaseRepository[platform.Platform]
}

func NewPlatformRepository(client *mongo.Client, config DBConfig) *PlatformRepository {
	return &PlatformRepository{
		BaseRepository: *NewBaseRepository[platform.Platform](client, config),
	}
}

func (s *PlatformRepository) GetPlatformByName(ctx context.Context, name string) (*platform.Platform, error) {
	condition := extensions.NewSearch(nil, nil, nil, map[string]interface{}{"name": name})
	ent, err := s.BaseRepository.GetWithCondition(ctx, condition)
	if err != nil {
		return nil, err
	}
	if len(ent) == 0 {
		return nil, fmt.Errorf("%s with name %s", extensions.Data_Not_Found_Message, name)
	}
	return &ent[0], nil
}

func (s *PlatformRepository) GetAllPlatform(ctx context.Context) ([]platform.Platform, error) {
	condition := extensions.NewSearch(nil, nil, nil, nil)
	return s.BaseRepository.GetWithCondition(ctx, condition)
}

func (s *PlatformRepository) GetPlatformByNameAsync(ctx context.Context, name string) (<-chan *platform.Platform, <-chan error) {
	condition := extensions.NewSearch(nil, nil, nil, map[string]interface{}{"name": name})
	resultChan := make(chan *platform.Platform, 1)
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
			errorChan <- fmt.Errorf("GetPlatformByNameAsync timeout, name %s", name)
		}
	}()

	return resultChan, errorChan
}

func (s *PlatformRepository) GetAllPlatformAsync(ctx context.Context) (<-chan []platform.Platform, <-chan error) {
	condition := extensions.NewSearch(nil, nil, nil, nil)
	return s.BaseRepository.GetWithConditionAsync(ctx, condition)
}
