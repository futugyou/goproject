package infrastructure_mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/futugyou/infr-project/extensions"
	models "github.com/futugyou/infr-project/view_models"
)

type ResourceQueryRepository struct {
	BaseQueryRepository[models.ResourceView]
}

func NewResourceQueryRepository(client *mongo.Client, config QueryDBConfig) *ResourceQueryRepository {
	return &ResourceQueryRepository{
		BaseQueryRepository: *NewBaseQueryRepository[models.ResourceView](client, config),
	}
}

func (r *ResourceQueryRepository) GetResourceByName(ctx context.Context, name string) (*models.ResourceView, error) {
	condition := extensions.NewSearch(nil, nil, nil, map[string]interface{}{"name": name})
	ent, err := r.BaseQueryRepository.GetWithCondition(ctx, condition)
	if err != nil {
		return nil, err
	}
	if len(ent) == 0 {
		return nil, fmt.Errorf("%s with name %s", extensions.Data_Not_Found_Message, name)
	}
	return &ent[0], nil
}

func (r *ResourceQueryRepository) GetResourceByNameAsync(ctx context.Context, name string) (<-chan *models.ResourceView, <-chan error) {
	condition := extensions.NewSearch(nil, nil, nil, map[string]interface{}{"name": name})

	resultChan := make(chan *models.ResourceView, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		result, err := r.BaseQueryRepository.GetWithConditionAsync(ctx, condition)
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
