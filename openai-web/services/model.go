package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/beego/beego/v2/core/logs"
	openai "github.com/openai/openai-go"
	"github.com/redis/go-redis/v9"
)

type ModelListResponse struct {
	Name     string `json:"name"`
	Describe string `json:"describe"`
}

type ModelService struct {
	redisDb *redis.Client
	client  *openai.Client
}

func NewModelService(client *openai.Client, redisDb *redis.Client) *ModelService {
	if redisDb == nil {
		client, err := RedisClient(os.Getenv("REDIS_URL"))
		if err != nil {
			panic(err)
		}
		redisDb = client
	}

	return &ModelService{
		client:  client,
		redisDb: redisDb,
	}
}

const GetAllModelsKey string = "GetAllModelsKey"

func (s *ModelService) GetAllModels(ctx context.Context) []ModelListResponse {
	result := make([]ModelListResponse, 0)
	rmap, _ := s.redisDb.HGetAll(ctx, GetAllModelsKey).Result()

	if len(rmap) > 0 {
		for _, r := range rmap {
			m := openai.Model{}
			json.Unmarshal([]byte(r), &m)
			result = append(result, ModelListResponse{Name: m.ID})
		}

		return result
	}

	models, err := s.client.Models.List(ctx)
	if err != nil {
		return result
	}
	rset := make(map[string]interface{})
	if len(models.Data) > 0 {
		for _, model := range models.Data {
			result = append(result, ModelListResponse{Name: model.ID})
			modelstring, _ := json.Marshal(model)
			rset[model.ID] = string(modelstring)
		}
	}

	count, err := s.redisDb.HSet(ctx, GetAllModelsKey, rset).Result()
	if err != nil {
		logs.Error(err)
	} else {
		logs.Info(fmt.Sprintf("data count: %d", count))
	}

	return result
}

type MyHash struct {
	Key1 string `redis:"key1"`
	Key2 int    `redis:"key2"`
}
